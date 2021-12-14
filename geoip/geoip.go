package geoip

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type GeoToIP struct {
	IP string
}

type GeoIPInformation struct {
	Query       string  `json:"query"`
	Status      string  `json:"status"`
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	Region      string  `json:"region"`
	RegionName  string  `json:"regionName"`
	City        string  `json:"city"`
	Zip         string  `json:"zip"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Timezone    string  `json:"timezone"`
	Isp         string  `json:"isp"`
	Org         string  `json:"org"`
	As          string  `json:"as"`
}

func New(ip string) *GeoToIP {
	return &GeoToIP{
		IP: ip,
	}
}

func (g *GeoToIP) GetInfo() (*GeoIPInformation, error) {
	url := "http://ip-api.com/json/%s"

	resp, err := http.Get(fmt.Sprintf(url, g.IP))
	if err != nil {
		log.Fatalln(err)
	}
	//We Read the response body on the line below.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()
	//Convert the body to type string

	geoInformation := &GeoIPInformation{}
	err = json.Unmarshal(body, geoInformation)
	if err != nil {
		return nil, err
	}
	return geoInformation, nil
}
