package emby

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Emby struct {
	hostname string
	apiKey   string
	userId   string
}

func New(hostname string, apiKey string, userId string) Emby {
	if hostname[len(hostname)-1] == '/' {
		return Emby{
			hostname: hostname[:hostname[len(hostname)-2]],
			apiKey:   apiKey,
			userId:   userId,
		}
	}
	return Emby{
		hostname: hostname,
		apiKey:   apiKey,
		userId:   userId,
	}
}

// SystemInfo retrieves server information corresponding to the System/Info API endpoint
func (c *Emby) GetSystemInfo() (*SystemInfo, error) {
	raw, err := c.request("GET", "/System/Info", "")
	if err != nil {
		return nil, err
	}
	sysInfo := &SystemInfo{}
	err = json.Unmarshal(raw, sysInfo)
	if err != nil {
		return nil, err
	}
	return sysInfo, nil
}

func (c *Emby) GetSessions() (*[]Sessions, error) {
	raw, err := c.request("GET", "/Sessions", "")
	if err != nil {
		return nil, err
	}
	session := &[]Sessions{}
	err = json.Unmarshal(raw, session)
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (c *Emby) GetMediaItem() (*[]UserView, error) {

	raw, err := c.request("GET", fmt.Sprintf("/Users/%s/Views", c.userId), "")
	if err != nil {
		return nil, err
	}
	mediaItemList := &MediaItemList{}
	err = json.Unmarshal(raw, mediaItemList)
	if err != nil {
		return nil, err
	}
	return &mediaItemList.Items, nil
}

func (e *Emby) request(method string, url string, body string) ([]byte, error) {
	req, _ := http.NewRequest(method, fmt.Sprintf("%s%s", e.hostname, url), strings.NewReader(body))
	req.Header.Set("X-Emby-Token", e.apiKey)
	req.Header.Set("Application-Type", "application/json")

	if len(body) > 0 {
		bodybytes := []byte(body)
		buf := bytes.NewBuffer(bodybytes)
		req.Body = ioutil.NopCloser(buf)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respbody, _ := ioutil.ReadAll(resp.Body)
	return respbody, nil
}
