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

func dialTimeout(network, addr string) (net.Conn, error) {
	//time.Duration, is the duration the app will wait for opening a TCP connection to the respective host
	return net.DialTimeout(network, addr, time.Duration(1*time.Second))
}

var httpClient = &http.Client{
	Timeout: time.Second * time.Duration(10),
	Transport: &http.Transport{
		Dial: dialTimeout,
	},
}

// GetAnyJSON does an HTTP get request and unmarshal result to a generic map[string]interface{}
func GetAnyJSON(url string) (interface{}, error) {
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}

	result := make(map[string]interface{}, 2)
	err = json.NewDecoder(resp.Body).Decode(&result)
	return &result, err
}

// Get does an HTTP get request
func Get(url string) (*response, error) {
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}

	result := response{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	return &result, err
}

// Post does an HTTP post request
func Post(url string, payload []byte) (*response, error) {
	resp, err := httpClient.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	result := response{}
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

	result := response{}
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

	result := response{}
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

	result := response{}
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

	result := response{}
	err = json.NewDecoder(resp.Body).Decode(&result)

	return &result, err
}
