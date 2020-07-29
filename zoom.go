package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func getMeetings(c *gin.Context) {
	claims := &jwt.StandardClaims{
		Issuer:    os.Getenv("ZOOM_KEY"),
		ExpiresAt: 1496091964000,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("ZOOM_SECRET")))
	if err != nil {
		log.Println("SignedString: ", err)
		panic(err)
	}

	req, err := http.NewRequest("GET", "https://api.zoom.us/v2/users/rdQj6tVKSAOErV6T9VUQlA/meetings?status=active&page_size=30&page_number=1", bytes.NewBuffer([]byte{}))
	if err != nil {
		log.Println("error: ", err)
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+tokenString)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("error: ", err)
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var data map[string]interface{}

	json.Unmarshal(body, &data)
	c.JSON(200, gin.H{
		"meetings": data["meetings"],
	})
}

// MeetingBody ...
type MeetingBody struct {
	Date      string `json:"date"`
	Duration  string `json:"duration"`
	Password  string `json:"password"`
	Time      string `json:"time"`
	StartTime string `json:"start_time"`
	Type      int16  `json:"type"`
	Topic     string `json:"topic"`
	Timezone  string `json:"timezone"`
}

func createMeeting(c *gin.Context) {
	var data MeetingBody
	body, _ := c.GetRawData()

	json.Unmarshal(body, &data)

	var trasnformed MeetingBody
	trasnformed.Date = CheckValidDate(ReplaceArabicNumbers(data.Date))
	trasnformed.Duration = ReplaceArabicNumbers(data.Duration)
	trasnformed.Time = ConvertTime12to24(ReplaceArabicNumbers(data.Time))
	trasnformed.Topic = ReplaceArabicNumbers(data.Topic)
	trasnformed.Type = data.Type
	trasnformed.Timezone = "Africa/Cairo"

	if data.Password != "" {
		trasnformed.Password = ValidatePassword(ReplaceArabicNumbers(data.Password))
		if trasnformed.Password == "" {
			c.JSON(200, gin.H{
				"status": 403, "field": "password", "message": "invalid password",
			})
			return
		}
	} else {
		trasnformed.Password = GeneratePassword()
	}

	claims := &jwt.StandardClaims{
		Issuer:    os.Getenv("ZOOM_KEY"),
		ExpiresAt: 1496091964000,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("ZOOM_SECRET")))
	if err != nil {
		log.Println("SignedString: ", err)
		panic(err)
	}

	reqBody, err := json.Marshal(trasnformed)
	if err != nil {
		log.Fatalln(err)
	}

	req, err := http.NewRequest("POST", "https://api.zoom.us/v2/users/rdQj6tVKSAOErV6T9VUQlA/meetings", bytes.NewBuffer(reqBody))
	if err != nil {
		log.Println("error: ", err)
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+tokenString)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("error: ", err)
		panic(err)
	}
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)
	var response map[string]interface{}

	json.Unmarshal(respBody, &response)

	c.JSON(200, gin.H{
		"response": response,
	})
}
