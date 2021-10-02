package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetBirdsHandler(t *testing.T) {
	mockStore := InitMockStore()

	mockStore.On("GetBirds").Return([]*Bird{
		{"sparrow", "A small harmless bird"},
	}, nil).Once()

	req, err := http.NewRequest("GET", "", nil)

	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()

	hf := http.HandlerFunc(getBirdHandler)

	hf.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	b := []Bird{}
	err = json.NewDecoder(recorder.Body).Decode(&b)
	if err != nil {
		t.Fatal(err)
	}
	actual := b[0]
	assert.Equal(t, Bird{"sparrow", "A small harmless bird"}, actual)
}

func TestCreateBirdsHandler(t *testing.T) {
	mockStore := InitMockStore()
	mockStore.On("CreateBird", &Bird{"eagle", "A bird of prey"}).Return(nil)
	form := newCreateBirdForm()
	req, err := http.NewRequest("POST", "", bytes.NewBufferString(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(form.Encode())))
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()

	hf := http.HandlerFunc(createBirdHandler)

	hf.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusFound, recorder.Code)
	mockStore.AssertExpectations(t)
}

func newCreateBirdForm() *url.Values {
	form := url.Values{}
	form.Set("species", "eagle")
	form.Set("description", "A bird of prey")
	return &form
}
