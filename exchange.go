package main

import (
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mhewedy/ews"
)

func getUserAvailability(c *gin.Context) {
	username := os.Getenv("EWS_USERNAME")
	password := os.Getenv("EWS_PASSWORD")
	url := os.Getenv("EWS_URL")
	// hostname := os.Getenv("EWS_HOSTNAME")

	client := ews.NewClient(url, username, password, &ews.Config{
		Dump: true, NTLM: true,
	})
	startTime, _ := time.Parse("2006-01-02", "2020-07-16")
	endTime, _ := time.Parse("2020-08-16", "2020-07-16")
	r := &ews.GetUserAvailabilityRequest{
		TimeZone: ews.TimeZone{
			Bias: 480,
			StandardTime: ews.TimeZoneTime{
				Bias:      0,
				Time:      "02:00:00",
				DayOrder:  5,
				Month:     10,
				DayOfWeek: "Sunday",
			},
			DaylightTime: ews.TimeZoneTime{
				Bias:      -60,
				Time:      "02:00:00",
				DayOrder:  1,
				Month:     4,
				DayOfWeek: "Sunday",
			},
		},
		MailboxDataArray: ews.MailboxDataArray{
			[]ews.MailboxData{
				{
					Email: ews.Email{
						Address: "itsvc_teamsbot_exch@EFG-HERMES.com",
					},
					AttendeeType:     "Required",
					ExcludeConflicts: false,
				},
			},
		},
		FreeBusyViewOptions: ews.FreeBusyViewOptions{
			TimeWindow: ews.TimeWindow{
				StartTime: startTime,
				EndTime:   endTime,
			},
			MergedFreeBusyIntervalInMinutes: 60,
			RequestedView:                   "Detailed",
		},
	}

	// xmlBytes, err := xml.MarshalIndent(r, "", "  ")
	// if err != nil {
	// 	log.Println("err", err)
	// 	c.JSON(400, err)
	// 	return
	// }

	resp, err := ews.GetUserAvailability(client, r)
	if err != nil {
		log.Println("Error on unmarshaling xml. ", err.Error())
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	// bb, err := client.SendAndReceive(xmlBytes)
	// log.Println("bb, err", bb, err)
	// if err != nil {
	// 	c.JSON(400, err)
	// 	return
	// }

	c.JSON(200, gin.H{
		"result": resp,
	})
}
