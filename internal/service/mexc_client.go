package service

import (
	"encoding/json"
	"finalyTask/internal/config"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

type PriceProvider interface {
	GetPrice(symbol string) (float64, error)
	GetUSDTPrice() (float64, error)
}

type MexcClient struct {
	httpClient *http.Client
	baseURL    string
}

func NewMexcClient(cfg *config.Config) *MexcClient {
	return &MexcClient{
		httpClient: &http.Client{Timeout: cfg.MexcTimeout},
		baseURL:    "https://api.mexc.com",
	}
}

func (c *MexcClient) GetPrice(symbol string) (float64, error) {
	resp, err := c.httpClient.Get(fmt.Sprintf("%s/api/v3/ticker/price?symbol=%s", c.baseURL, symbol))
	if err != nil {
		return 0, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read response: %w", err)
	}

	var result struct {
		Price string `json:"price"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return 0, fmt.Errorf("JSON parse error: %w", err)
	}

	price, err := strconv.ParseFloat(result.Price, 64)
	if err != nil {
		return 0, fmt.Errorf("price conversion error: %w", err)
	}

	return price, nil
}

func (c *MexcClient) GetUSDTPrice() (float64, error) {
	url := fmt.Sprintf("%s/api/v3/ticker/price?symbol=KASUSDT", c.baseURL)
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return 0, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	var result struct {
		Price string `json:"price"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return 0, fmt.Errorf("JSON unmarshal error: %w", err)
	}

	price, err := strconv.ParseFloat(result.Price, 64)
	if err != nil {
		return 0, fmt.Errorf("price parsing error: %w", err)
	}

	return price, nil
}
