package main

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"exchange-web-services/ntlm"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httputil"

	"github.com/Azure/go-ntlmssp"
)

const (
	soapStart = `<?xml version="1.0" encoding="utf-8" ?>
<soap:Envelope xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" 
		xmlns:m="http://schemas.microsoft.com/exchange/services/2006/messages" 
		xmlns:t="http://schemas.microsoft.com/exchange/services/2006/types" 
		xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  		<soap:Header>
    		<t:RequestServerVersion Version="Exchange2010" />
  		</soap:Header>
  		<soap:Body>
`
	soapEnd = `
</soap:Body></soap:Envelope>`
)

// Transport ...
type Transport struct {
	Domain   string
	User     string
	Password string
	conn     *httputil.ClientConn
	host     string
}

var encBase64 = base64.StdEncoding.EncodeToString
var decBase64 = base64.StdEncoding.DecodeString

// Config ...
type Config struct {
	Dump    bool
	NTLM    bool
	SkipTLS bool
}

// Client ...
type Client interface {
	SendAndReceive(body []byte) ([]byte, error)
	GetEWSAddr() string
	GetUsername() string
}

type client struct {
	EWSAddr  string
	Username string
	Password string
	HostName string
	config   *Config
}

func cloneRequest(req *http.Request) *http.Request {
	r2 := *req
	r2.Header = http.Header{}
	for k, v := range req.Header {
		r2.Header[k] = v
	}
	return &r2
}

func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.Host != t.host {
		if t.conn != nil {
			t.conn.Close()
		}

		t.host = req.Host

		host, port, err := net.SplitHostPort(t.host)
		if err != nil {
			port = "80"
		}

		sock, err := net.Dial("tcp", net.JoinHostPort(host, port))
		if err != nil {
			return nil, err
		}

		t.conn = httputil.NewClientConn(sock, nil)
		req = cloneRequest(req)
		req.Header.Set("Authorization", encBase64(ntlm.Negotiate()))
	}

	resp, err := t.conn.Do(req)

	if err == nil && resp.StatusCode == http.StatusUnauthorized {
		chlg, err := decBase64(resp.Header.Get("Www-Authenticate"))
		if err != nil {
			return nil, err
		}

		auth, err := ntlm.Authenticate(chlg, t.Domain, t.User, t.Password)
		if err != nil {
			return nil, err
		}

		req = cloneRequest(req)
		req.Header.Set("Authorization", encBase64(auth))

		resp, err = t.conn.Do(req)
	}

	return resp, err
}

func (c *client) GetEWSAddr() string {
	return c.EWSAddr
}

func (c *client) GetUsername() string {
	return c.Username
}

// NewClient ....
func NewClient(ewsAddr, username, password, hostName string, config *Config) Client {
	return &client{
		EWSAddr:  ewsAddr,
		Username: username,
		Password: password,
		HostName: hostName,
		config:   config,
	}
}

func (c *client) SendAndReceive(body []byte) ([]byte, error) {

	bb := []byte(soapStart)
	bb = append(bb, body...)
	bb = append(bb, soapEnd...)
	req, err := http.NewRequest("POST", c.EWSAddr, bytes.NewReader([]byte(bb)))
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()
	ts := &Transport{
		Domain:   "efg-hermes.com",
		User:     c.Username,
		Password: c.Password,
		host:     c.HostName,
	}
	req.Header.Set("Content-Type", "text/xml; charset=utf-8")
	req.Header.Set("Connection", "keep-alive")
	resp, err := ts.RoundTrip(req)
	// req.Header.Set("WWW-Authenticate", "NTLM")
	// req.Header.Set("SOAPAction", `"http://schemas.microsoft.com/exchange/services/2006/messages/GetUserAvailability"`)

	// req.SetBasicAuth(c.Username, c.Password)
	// auth := c.Username + ":" + c.Password
	// token := Ntlmgen(auth)
	// token := base64.StdEncoding.EncodeToString([]byte(auth))
	// req.Header.Set("Authorization", "NTLM "+token)
	// req.Header.Set("Authorization", "NTLM TlRMTVNTUAABAAAAB6IIogAAAAAoAAAAAAAAACgAAAAFASgKAAAADw==")
	// req.Header.Add("cookie", "exchangecookie=6e3dc47cb8a14895b17a1878a5cda997; Expires=Tue, 27 Jul 2021 23:40:35 GMT; Path=/; HttpOnly; Domain=mail.efg-hermes.com")
	// client := &http.Client{
	// 	// Transport: ntlmssp.Negotiator{
	// 	// 	RoundTripper: &http.Transport{},
	// 	// },
	// 	Transport: &httpntlm.NtlmTransport{
	// 		Domain:   c.HostName,
	// 		User:     c.Username,
	// 		Password: c.Password,
	// 		// Configure RoundTripper if necessary, otherwise DefaultTransport is used
	// 		RoundTripper: &http.Transport{
	// 			// provide tls config
	// 			TLSClientConfig: &tls.Config{},
	// 			// other properties RoundTripper, see http.DefaultTransport
	// 		},
	// 	},
	// }

	// applyConfig(c.config, client)
	logRequest(c, req)

	// resp, err := client.Do(req)
	// resp, err := ntlm.DoNTLMRequest(client, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	logResponse(c, resp)

	if resp.StatusCode != http.StatusOK {
		return nil, err
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return respBytes, err
}

func applyConfig(config *Config, client *http.Client) {
	if config.NTLM {
		client.Transport = ntlmssp.Negotiator{}
	}
	if config.SkipTLS {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}
}

func logRequest(c *client, req *http.Request) {
	if c.config != nil && c.config.Dump {
		dump, err := httputil.DumpRequestOut(req, true)
		if err != nil {
			fmt.Println(err)
		}
		log.Printf("Request:\n%v\n----\n", string(dump))
	}
}

func logResponse(c *client, resp *http.Response) {
	if c.config != nil && c.config.Dump {
		dump, err := httputil.DumpResponse(resp, true)
		if err != nil {
			log.Println(err)
		}
		log.Printf("Response:\n%v\n----\n", string(dump))
	}
}
