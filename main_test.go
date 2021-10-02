package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	r := newRouter()
	mockServer := httptest.NewServer(r)
	resp, err := http.Get(mockServer.URL + "/hello")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "Hello World!", string(b))
}

func TestRouterForNonExistentRoute(t *testing.T) {
	r := newRouter()
	mockServer := httptest.NewServer(r)
	resp, err := http.Post(mockServer.URL+"/hello", "", nil)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "", string(b))
}

func TestStaticFileServer(t *testing.T) {
	r := newRouter()
	mockServer := httptest.NewServer(r)
	resp, err := http.Get(mockServer.URL + "/assets/")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	contentType := resp.Header.Get("Content-Type")
	assert.Equal(t, "text/html; charset=utf-8", contentType)
}
