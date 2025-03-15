package service

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMexcClient_GetPrice(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"symbol":"KASUSDT","price":"1.0000"}`))
	}))
	defer server.Close()

	client := &MexcClient{
		httpClient: server.Client(),
		baseURL:    server.URL,
	}

	price, err := client.GetPrice("KASUSDT")
	assert.NoError(t, err)
	assert.Equal(t, 1.0, price)
}

func TestMexcClient_GetPrice_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := &MexcClient{
		httpClient: server.Client(),
		baseURL:    server.URL,
	}

	_, err := client.GetPrice("KASUSDT")
	assert.Error(t, err)
}
