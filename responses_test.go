package webgo

import (
	"encoding/json"
	"html/template"
	"net/http"
	"testing"
)

func TestResponses(t *testing.T) {
	_, respRec := setup()
	R200(respRec, nil)
	if respRec.Code != http.StatusOK {
		t.Error("Expected response status 200, got", respRec.Code)
	}

	_, respRec = setup()
	R201(respRec, nil)
	if respRec.Code != http.StatusCreated {
		t.Error("Expected response status 201, got", respRec.Code)
	}

	_, respRec = setup()
	R204(respRec)
	if respRec.Code != http.StatusNoContent {
		t.Error("Expected response status 204, got", respRec.Code)
	}

	_, respRec = setup()
	R302(respRec, nil)
	if respRec.Code != http.StatusFound {
		t.Error("Expected response status 302, got", respRec.Code)
	}

	_, respRec = setup()
	R400(respRec, nil)
	if respRec.Code != http.StatusBadRequest {
		t.Error("Expected response status 400, got", respRec.Code)
	}

	_, respRec = setup()
	R403(respRec, nil)
	if respRec.Code != http.StatusForbidden {
		t.Error("Expected response status 403, got", respRec.Code)
	}

	_, respRec = setup()
	R404(respRec, nil)
	if respRec.Code != http.StatusNotFound {
		t.Error("Expected response status 404, got", respRec.Code)
	}

	_, respRec = setup()
	R406(respRec, nil)
	if respRec.Code != http.StatusNotAcceptable {
		t.Error("Expected response status 406, got", respRec.Code)
	}

	_, respRec = setup()
	R451(respRec, nil)
	if respRec.Code != http.StatusUnavailableForLegalReasons {
		t.Error("Expected response status 451, got", respRec.Code)
	}

	_, respRec = setup()
	R500(respRec, nil)
	if respRec.Code != http.StatusInternalServerError {
		t.Error("Expected response status 500, got", respRec.Code)
	}
}

func TestInvalidResponses(t *testing.T) {
	_, respRec := setup()

	R200(respRec, make(chan int))
	resp := response{}
	err := json.NewDecoder(respRec.Body).Decode(&resp)
	if err != nil {
		t.Error(err)
		return
	}

	if resp.Status != http.StatusInternalServerError {
		t.Error("Expected status 500, got:", resp.Status)
	}

	if respRec.Result().StatusCode != http.StatusInternalServerError {
		t.Error("Expected status 500, got:", respRec.Result().StatusCode)
	}

	_, respRec = setup()

	R400(respRec, make(chan int))
	resp = response{}
	err = json.NewDecoder(respRec.Body).Decode(&resp)
	if err != nil {
		t.Error(err)
		return
	}

	if resp.Status != http.StatusInternalServerError {
		t.Error("Expected status 500, got:", resp.Status)
	}

	if respRec.Result().StatusCode != http.StatusInternalServerError {
		t.Error("Expected status 500, got:", respRec.Result().StatusCode)
	}
}

func TestSend(t *testing.T) {
	_, respRec := setup()
	Send(respRec, "text/html", "hello", http.StatusOK)

	if respRec.Code != http.StatusOK {
		t.Error("Expected status '200', got", respRec.Code)
	}

	str := respRec.Body.String()
	if str != "hello" {
		t.Error("Expected 'hello', got", str)
	}
}

func TestRender(t *testing.T) {
	_, respResc := setup()

	tmpl, err := template.New("sample").Parse(`hello world`)
	if err != nil {
		t.Error(err)
		return
	}

	Render(respResc, nil, http.StatusOK, tmpl)
	str := respResc.Body.String()
	if str != `hello world` {
		t.Error(str)
	}
}
func TestRender404(t *testing.T) {
	_, respResc := setup()

	tmpl, err := template.New("sample").Parse(`{{.ErrDescription}}`)
	if err != nil {
		t.Error(err)
		return
	}

	Render404(respResc, tmpl)

	str := respResc.Body.String()
	if str != `Sorry, the URL you requested was not found on this server... Or you&#39;re lost :-/` {
		t.Error(str)
	}
}

func TestSendHeader(t *testing.T) {
	_, respResc := setup()
	SendHeader(respResc, http.StatusAccepted)
	if respResc.Code != http.StatusAccepted {
		t.Error("Expected response code 202, got:", respResc.Code)
	}
}
