package parse

import (
	"CallsForService/CFS-Parsing/pkg/config"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html"
)

type CFS struct {
	EventTime    time.Time `json:"eventtime"`
	RawEventTime string    `json:"rawEventTime"`
	EventDate    string    `json:"eventDay"`
	Id           string    `json:"Id"`
	Location     string    `json:"location"`
	Description  string    `json:"description"`
}

type Location struct {
	Location           string  `json:"location"`
	NormalizedLocation string  `json:"normalizedLocation"`
	Lat                float64 `json:"lat"`
	Lng                float64 `json:"lng"`
	Ward               string  `json:"ward"`
	Neighborhood       string  `json:"neighborhood"`
	Zipcode            string  `json:"zipcode"`
	HasIssue           bool    `json:"hasIssue"`
}

const timeLayout = "2006-01-02 15:04:05"

func CallCFS(config *config.Config) ([]CFS, error) {

	response, err := config.Client.Get(config.CFSWebsite)

	if err != nil {
		log.Fatalf("Error received from http client: %v", err)
	}

	if response.StatusCode != http.StatusOK {
		log.Printf("Non 200 code received: %v", response.StatusCode)
		return nil, errors.New(fmt.Sprintf("Non 200 status code received: %v", response.StatusCode))
	}
	defer response.Body.Close()

	doc := html.NewTokenizer(response.Body)
	cfs, err := ParseCFS(doc)

	io.Copy(ioutil.Discard, response.Body)
	return cfs, err
}

func ParseCFS(doc *html.Tokenizer) ([]CFS, error) {
	loc, err := time.LoadLocation("America/Chicago")
	if err != nil {
		log.Fatal("Error loading time location")
	}

	var CFSs []CFS

	end := true
	for end {

		t := doc.Next()
		switch {
		case t == html.ErrorToken:
			end = false
		case t == html.StartTagToken:
			tt := doc.Token()

			isTR := tt.Data == "tr"

			if isTR {

				i := 0
				var cfs CFS
				for i < 5 {
					n := doc.Next()
					if n == html.TextToken {
						temp := string(doc.Text())
						switch i {
						case 1:
							rawTime, err := time.ParseInLocation(timeLayout, temp, loc)

							if err != nil {

								cfs.RawEventTime = temp
								fixedTime, err := buildValidTimeStamp(loc, cfs.RawEventTime)

								if err != nil {
									log.Printf("Error building fixed timestamp: %v", err)
								}

								cfs.EventTime = fixedTime
							} else {
								cfs.EventTime = rawTime.UTC()
							}

							eventDate := strings.Split(temp, " ")
							cfs.EventDate = strings.Trim(eventDate[0], " ")

						case 2:
							cfs.Id = temp
						case 3:
							cfs.Location = temp
						case 4:
							cfs.Description = temp
							CFSs = append(CFSs, cfs)
						}
						i += 1
					}
				}

			}
		}

	}
	return CFSs, nil
}

func buildValidTimeStamp(loc *time.Location, badEventTime string) (time.Time, error) {

	seperateDateAndTime := strings.Split(badEventTime, " ")

	if len(seperateDateAndTime) != 2 {
		errMessage := fmt.Sprintf("Expected the dateAndTime to have length two, instead it had length: %v", len(seperateDateAndTime))
		log.Printf(errMessage)
		return time.Time{}, errors.New(errMessage)
	}

	dateOnly := strings.Split(seperateDateAndTime[0], "-")

	timeOnly := strings.Split(seperateDateAndTime[1], ":")

	if len(timeOnly) != 3 {
		errMessage := fmt.Sprintf("Expected the dateAndTime to have length three, instead it had length: %v", len(timeOnly))
		log.Printf(errMessage)
		return time.Time{}, errors.New(errMessage)
	}

	year, err := strconv.Atoi(dateOnly[0])

	if err != nil {
		return time.Time{}, err
	}

	month, err := strconv.Atoi(dateOnly[1])

	if err != nil {
		return time.Time{}, err
	}

	day, err := strconv.Atoi(dateOnly[2])

	if err != nil {
		return time.Time{}, err
	}

	if timeOnly[0] == "" {
		timeOnly[0] = "0"
	}

	hour, err := strconv.Atoi(timeOnly[0])

	if err != nil {
		return time.Time{}, err
	}

	if timeOnly[1] == "" {
		timeOnly[1] = "0"
	}

	minute, err := strconv.Atoi(timeOnly[1])

	if err != nil {
		return time.Time{}, err
	}

	if timeOnly[2] == "" {
		timeOnly[2] = "0"
	}

	second, err := strconv.Atoi(timeOnly[2])

	if err != nil {
		return time.Time{}, err
	}

	validDate := time.Date(year, time.Month(month), day, hour, minute, second, 0, loc)

	return validDate.UTC(), nil

}
