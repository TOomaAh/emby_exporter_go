package emby

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Server struct {
	Url    string
	Token  string
	UserID string
	Port   int
}

var headers = map[string]string{
	"Content-Type": "application/json",
	"X-Emby-Token": "",
}

func NewServer(url, token, userID string, port int) *Server {

	server := &Server{
		Url:    url,
		Token:  token,
		UserID: userID,
		Port:   port,
	}

	return server
}

func (s *Server) GetServerInfo() (*SystemInfo, error) {
	resp, err := s.request("GET", "/System/Info", "")
	if err != nil {
		return nil, err
	}

	var systemInfo SystemInfo
	err = json.Unmarshal(resp, &systemInfo)
	if err != nil {
		return nil, err
	}

	return &systemInfo, nil
}

func (s *Server) GetLibrary() (*[]Library, error) {
	resp, err := s.request("GET", "/Items", "")
	if err != nil {
		return nil, err
	}

	var library []Library
	err = json.Unmarshal(resp, &library)
	if err != nil {
		return nil, err
	}

	return &library, nil
}

func (s *Server) GetSessions() (*[]Sessions, error) {
	resp, err := s.request("GET", "/Sessions", "")
	if err != nil {
		return nil, err
	}

	var sessions []Sessions
	err = json.Unmarshal(resp, &sessions)
	if err != nil {
		return nil, err
	}
	var sessionResult []Sessions
	for i, session := range sessions {
		if session.PlayState.PlayMethod != "" {
			sessionResult = append(sessionResult, sessions[i])
		}
	}

	return &sessionResult, nil
}

func (s *Server) request(method string, path string, body string) ([]byte, error) {
	req, _ := http.NewRequest(method, fmt.Sprintf("%s%s", s.Url, path), strings.NewReader(body))
	req.Header.Set("X-Emby-Token", s.Token)
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
