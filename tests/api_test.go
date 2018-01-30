package main

import (
	"net/http"
	"strings"
	"testing"
)

const p1 = "world"
const p2 = "spiderman"

const baseapi = "http://127.0.0.1:9696"

// const baseapiHTTPS = "https://127.0.0.1:8443"

const baseapiHTTPS = "http://127.0.0.1:9696"

var GETAPI = []string{
	strings.Join([]string{baseapi, "hello", p1, "goblin", p2}, "/"),
	strings.Join([]string{baseapiHTTPS, "hello", p1, "goblin", p2}, "/"),
}

var POSTAPI = []string{
	strings.Join([]string{baseapi, "hello", p1, "goblin", p2}, "/"),
	strings.Join([]string{baseapiHTTPS, "hello", p1, "goblin", p2}, "/"),
}

var PUTAPI = []string{
	strings.Join([]string{baseapi, "hello", p1, "goblin", p2}, "/"),
	strings.Join([]string{baseapiHTTPS, "hello", p1, "goblin", p2}, "/"),
}

var DELETEAPI = []string{
	strings.Join([]string{baseapi, "hello", p1, "goblin", p2}, "/"),
	strings.Join([]string{baseapiHTTPS, "hello", p1, "goblin", p2}, "/"),
}

var PATCHAPI = []string{
	strings.Join([]string{baseapi, "hello", p1, "goblin", p2}, "/"),
	strings.Join([]string{baseapiHTTPS, "hello", p1, "goblin", p2}, "/"),
}

var OPTIONSAPI = []string{
	strings.Join([]string{baseapi, "hello", p1, "goblin", p2}, "/"),
	strings.Join([]string{baseapiHTTPS, "hello", p1, "goblin", p2}, "/"),
}

func TestGet(t *testing.T) {
	for _, url := range GETAPI {
		resp, err := Get(url)
		if err != nil {
			t.Log(err)
			t.Fail()
			return
		}

		if resp.Data["method"] != http.MethodGet {
			t.Log("URL:", url, "response method:", resp.Data["method"], " required method:", http.MethodGet)
			t.Fail()
		}

		if resp.Data["p1"] != p1 {
			t.Log("p1:", resp.Data["p1"])
			t.Fail()
		}

		if resp.Data["p2"] != p2 {
			t.Log("p2:", resp.Data["p2"])
			t.Fail()
		}
	}
}

func TestPost(t *testing.T) {
	var payload = []byte(`{"payload": "nothing"}`)

	for _, url := range POSTAPI {
		resp, err := Post(url, payload)
		if err != nil {
			t.Log(err)
			t.Fail()
			return
		}

		if resp.Data["method"] != http.MethodPost {
			t.Log("response method:", resp.Data["method"], " required method:", http.MethodPost)
			t.Fail()
		}

		if resp.Data["p1"] != p1 {
			t.Log("p1:", resp.Data["p1"])
			t.Fail()
		}

		if resp.Data["p2"] != p2 {
			t.Log("p2:", resp.Data["p2"])
			t.Fail()
		}

		if resp.Data["payload"] != string(payload) {
			t.Log("payload:", resp.Data["payload"])
			t.Fail()
		}
	}
}

func TestPut(t *testing.T) {
	var payload = []byte(`{"payload": "nothing"}`)
	for _, url := range PUTAPI {
		resp, err := Put(url, payload)
		if err != nil {
			t.Log(err)
			t.Fail()
			return
		}

		if resp.Data["method"] != http.MethodPut {
			t.Log("response method:", resp.Data["method"], " required method:", http.MethodPut)
			t.Fail()
		}

		if resp.Data["p1"] != p1 {
			t.Log("p1:", resp.Data["p1"])
			t.Fail()
		}

		if resp.Data["p2"] != p2 {
			t.Log("p2:", resp.Data["p2"])
			t.Fail()
		}

		if resp.Data["payload"] != string(payload) {
			t.Log("payload:", resp.Data["payload"])
			t.Fail()
		}
	}
}

func TestPatch(t *testing.T) {
	var payload = []byte(`{"payload": "nothing"}`)
	for _, url := range PATCHAPI {
		resp, err := Patch(url, payload)
		if err != nil {
			t.Log(err)
			t.Fail()
			return
		}

		if resp.Data["method"] != http.MethodPatch {
			t.Log("response method:", resp.Data["method"], " required method:", http.MethodPatch)
		}

		if resp.Data["p1"] != p1 {
			t.Log("p1:", resp.Data["p1"])
			t.Fail()
		}

		if resp.Data["p2"] != p2 {
			t.Log("p2:", resp.Data["p2"])
			t.Fail()
		}

		if resp.Data["payload"] != string(payload) {
			t.Log("payload:", resp.Data["payload"])
			t.Fail()
		}
	}
}

func TestDelete(t *testing.T) {
	var payload = []byte(`{"payload": "nothing"}`)
	for _, url := range DELETEAPI {
		resp, err := Delete(url, payload)
		if err != nil {
			t.Log(err)
			t.Fail()
			return
		}

		if resp.Data["method"] != http.MethodDelete {
			t.Log("response method:", resp.Data["method"], " required method:", http.MethodDelete)
			t.Fail()
		}

		if resp.Data["p1"] != p1 {
			t.Log("p1:", resp.Data["p1"])
			t.Fail()
		}

		if resp.Data["p2"] != p2 {
			t.Log("p2:", resp.Data["p2"])
			t.Fail()
		}

		if resp.Data["payload"] != string(payload) {
			t.Log("payload:", resp.Data["payload"])
			t.Fail()
		}
	}
}

func TestOptions(t *testing.T) {
	var payload = []byte(`{"payload": "nothing"}`)

	for _, url := range OPTIONSAPI {
		resp, err := Options(url, payload)
		if err != nil {
			t.Log(err)
			t.Fail()
			return
		}

		if resp.Data["method"] != http.MethodOptions {
			t.Log("response method:", resp.Data["method"], " required method:", http.MethodOptions)
			t.Fail()
		}

		if resp.Data["p1"] != p1 {
			t.Log("p1:", resp.Data["p1"])
			t.Fail()
		}

		if resp.Data["p2"] != p2 {
			t.Log("p2:", resp.Data["p2"])
			t.Fail()
		}

		if resp.Data["payload"] != string(payload) {
			t.Log("payload:", resp.Data["payload"])
			t.Fail()
		}
	}
}
