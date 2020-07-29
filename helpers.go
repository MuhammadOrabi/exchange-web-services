package main

import (
	"log"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Numbers ...
var Numbers = []map[string]string{
	{"en": "0", "ar": "٠"}, {"en": "1", "ar": "١"}, {"en": "2", "ar": "٢"}, {"en": "3", "ar": "٣"}, {"en": "4", "ar": "٤"},
	{"en": "5", "ar": "٥"}, {"en": "6", "ar": "٦"}, {"en": "7", "ar": "٧"}, {"en": "8", "ar": "٨"}, {"en": "9", "ar": "٩"},
}

const charset = "abcdefghijklmnopqrstuvwxyz" + "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789" + "0123456789" + "-@*_"

// ReplaceArabicNumbers ...
func ReplaceArabicNumbers(str string) string {
	var result string
	split := strings.Split(str, "")
	for _, c := range split {
		var found string
		for _, number := range Numbers {
			if c == number["ar"] {
				found = number["en"]
			}
		}
		if found != "" {
			result += found
		} else {
			result += c
		}
	}
	return result
}

// CheckValidDate ...
func CheckValidDate(str string) string {
	var result string
	var dateSeprator string
	split := strings.Split(str, "")
	for _, c := range split {
		_, err := strconv.ParseInt(c, 10, 32)
		if err != nil {
			dateSeprator = c
		}
	}
	split = strings.Split(str, dateSeprator)
	if len(split) > 0 {
		day, _ := strconv.ParseInt(split[0], 10, 64)
		if day == 0 || day > 31 {
			return ""
		} else if day < 10 {
			result += "0" + strconv.Itoa(int(day))
		} else {
			result += strconv.Itoa(int(day))
		}
	}

	if len(split) > 1 {
		var month int64
		month, _ = strconv.ParseInt(split[1], 10, 64)
		if month == 0 || month > 12 {
			return ""
		} else if month < 10 {
			result += "-0" + strconv.Itoa(int(month))
		} else {
			result += "-" + strconv.Itoa(int(month))
		}
	}

	var year int64
	if len(split) > 2 {
		year, _ = strconv.ParseInt(split[2], 10, 64)
	}
	if year == 0 {
		year, _, _ := time.Now().Date()
		result += "-" + strconv.Itoa(year)
	} else {
		result += "-" + strconv.Itoa(int(year))
	}

	return result
}

// ConvertTime12to24 ...
func ConvertTime12to24(time string) string {
	var result, time12h string
	modifier := "AM"
	time = strings.Replace(strings.ToLower(time), " ", "", -1)
	if len(strings.Split(time, "pm")) == 2 {
		modifier = "PM"
		time12h = strings.Split(time, "pm")[0]
	} else if len(strings.Split(time, "am")) == 2 {
		time12h = strings.Split(time, "am")[0]
	}

	var timeSeprator string
	split := strings.Split(time12h, "")
	for _, c := range split {
		_, err := strconv.ParseInt(c, 10, 32)
		if err != nil {
			timeSeprator = c
		}
	}

	split = strings.Split(time12h, timeSeprator)

	if len(split) == 2 {
		hours := split[0]
		if hours == "12" {
			hours = "00"
		}

		if modifier == "PM" {
			h, _ := strconv.ParseInt(hours, 10, 64)
			hours = strconv.Itoa(int(h + 12))
		}
		result = hours + ":" + split[1]

	}
	return result
}

// ValidatePassword ...
func ValidatePassword(password string) string {
	if len(password) > 10 {
		return ""
	}
	m := regexp.MustCompile(`\d+|[a-z|A-Z]|[-@*_]`)
	strped := m.ReplaceAllString(password, "")
	log.Println("password", password)
	log.Println("strped", strped)
	if strped != "" {
		return ""
	}
	return password
}

// GeneratePassword ...
func GeneratePassword() string {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, 10)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
