package main

import (
	"os"

	"github.com/MuhammadOrabi/ews"
	"github.com/gin-gonic/gin"
)

func getUserAvailability(c *gin.Context) {
	username := os.Getenv("EWS_USERNAME")
	password := os.Getenv("EWS_PASSWORD")

	hostname := os.Getenv("EWS_HOSTNAME")
	if dn := hostname; len(dn) > 0 {
		username = dn + "\\" + username
	}

	client := ews.NewClient(
		os.Getenv("EWS_URL"),
		username,
		password,
		&ews.Config{
			Dump:    false,
			NTLM:    true,
			SkipTLS: false,
		},
	)

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
				StartTime: "2020-08-02T00:00:00Z",
				EndTime:   "2020-08-02T01:00:00Z",
			},
			MergedFreeBusyIntervalInMinutes: 60,
			RequestedView:                   "Detailed",
		},
	}

	resp, err := ews.GetUserAvailability(client, r)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(200, gin.H{
		"result": resp,
	})
}
