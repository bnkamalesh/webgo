package main

import (
	"bytes"
	"encoding/json"
	"net"
	"net/http"
	"time"
)

type response struct {
	Data   map[string]string `json:"data"`
	Status int               `json:"status"`
}

type newstruct struct {
	Hey   string
	World string
}

var httpClient *http.Client

func dialTimeout(network, addr string) (net.Conn, error) {
	//time.Duration, is the duration the app will wait for opening a TCP connection to the respective host
	return net.DialTimeout(network, addr, time.Duration(15*time.Second))
}

func init() {
	httpClient = &http.Client{
		Timeout: time.Second * time.Duration(10),
		Transport: &http.Transport{
			Dial: dialTimeout,
		},
	}
}

// Get does an HTTP get request
func Get(url string) (*response, error) {
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result = response{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	return &result, err
}

// Post does an HTTP post request
func Post(url string, payload []byte) (*response, error) {
	resp, err := httpClient.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result = response{}
	err = json.NewDecoder(resp.Body).Decode(&result)

	return &result, err
}

// Put does an HTTP put request
func Put(url string, payload []byte) (*response, error) {
	req, _ := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(payload))
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result = response{}
	err = json.NewDecoder(resp.Body).Decode(&result)

	return &result, err
}

// Patch does an HTTP patch request
func Patch(url string, payload []byte) (*response, error) {
	req, _ := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(payload))
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result = response{}
	err = json.NewDecoder(resp.Body).Decode(&result)

	return &result, err
}

// Delete does an HTTP delete request
func Delete(url string, payload []byte) (*response, error) {
	req, _ := http.NewRequest(http.MethodDelete, url, bytes.NewBuffer(payload))
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result = response{}
	err = json.NewDecoder(resp.Body).Decode(&result)

	return &result, err
}

// Options does an HTTP options request
func Options(url string, payload []byte) (*response, error) {
	req, _ := http.NewRequest(http.MethodOptions, url, bytes.NewBuffer(payload))
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result = response{}
	err = json.NewDecoder(resp.Body).Decode(&result)

	return &result, err
}
