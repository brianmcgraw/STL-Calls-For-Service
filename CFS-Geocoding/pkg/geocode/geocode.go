package geocode

import (
	"CallsForService/CFS-Geocoding/pkg/config"
	"CallsForService/CFS-Geocoding/pkg/pg"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

const (
	GMAPSAPIKEY = "Maps_API_Key"
	MAPSURL     = "https://maps.googleapis.com/maps/api/geocode/json"
)

type MapsResponse struct {
	Results []MapsResults `json:"results"`
	Status  string        `json:"string"`
}

type MapsResults struct {
	AddressComponents []AddressComponent `json:"address_components"`
	FormattedAddress  string             `json:"formatted_address"`
	Geometry          MapsGeometry       `json:"geometry"`
	PlaceId           string             `json:"place_id"`
	Types             []string           `json:"types"`
}

type AddressComponent struct {
	LongName  string   `json:"long_name"`
	ShortName string   `json:"short_name"`
	Types     []string `json:"types"`
}

type MapsGeometry struct {
	Location LatLng `json:"location"`
}

type LatLng struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

func CallMaps(config *config.Config, loc string) (pg.Location, error) {

	if config.Debug_Env == "debug" {
		return CallMapsDebug(config, loc)
	}

	normalizedAddress := NormalizeAddress(loc)
	url := buildURL(config.MapsConfig.URL, normalizedAddress, config.MapsConfig.API_Key)

	req, err := http.NewRequest(http.MethodGet, url.String(), nil)

	if err != nil {
		log.Printf("Err building http request for google maps: %v", err)
	}

	resp, err := config.Client.Do(req)

	if err != nil {
		log.Println("Error from google maps API")
		//TODO FIGURE out how tohandle this error
	}
	var mapsResponse MapsResponse

	// Turn these responses into errors
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("Error calling google maps api, status code %v received", resp.StatusCode)
		return pg.Location{}, err
	}

	err = json.NewDecoder(resp.Body).Decode(&mapsResponse)

	if err != nil {
		err = fmt.Errorf("Error decoding google maps api response: %v", err)
		return pg.Location{}, err
	}

	var location pg.Location
	if len(mapsResponse.Results) > 0 {

		location = pg.Location{
			Location:           loc,
			Lat:                mapsResponse.Results[0].Geometry.Location.Lat,
			Lng:                mapsResponse.Results[0].Geometry.Location.Lng,
			NormalizedLocation: normalizedAddress,
			HasIssue:           false,
		}

		for _, value := range mapsResponse.Results[0].AddressComponents {
			if checkContains(value.Types, "postal_code") {
				location.Zipcode = value.LongName
			}
		}
	} else {
		return pg.Location{Location: loc, NormalizedLocation: normalizedAddress, HasIssue: true}, errors.New("Error getting result from geocoding api")
	}

	return location, err
}

func NormalizeAddress(s string) (f string) {
	// f = strings.ReplaceAll(s, "X", "0")
	f = s
	for index, value := range f {
		if index > 0 && string(value) == "X" && StringValueIsNumber(string(f[index-1])) {
			f = strings.Replace(f, "X", "0", 1)
		}
	}
	return f
}

func StringValueIsNumber(s string) bool {
	if s == "0" || s == "1" || s == "2" || s == "3" || s == "4" || s == "5" || s == "6" || s == "7" || s == "8" || s == "9" {
		return true
	}
	return false
}

func buildURL(stringURL string, s string, k string) *url.URL {
	toURL, err := url.Parse(stringURL)
	if err != nil {
		log.Printf("Bad value for url env variable: %v", err)
	}
	q := toURL.Query()
	fixedString := s + "+SAINT+LOUIS+MO"
	q.Set("address", fixedString)
	q.Set("key", k)
	toURL.RawQuery = q.Encode()
	return toURL
}

func checkContains(arrayToCheck []string, valuetoCheck string) bool {
	for _, v := range arrayToCheck {
		if v == valuetoCheck {
			return true
		}
	}

	return false
}

func CallMapsDebug(mapsConfig *config.Config, loc string) (pg.Location, error) {

	return pg.Location{}, nil

}
