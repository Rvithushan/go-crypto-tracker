package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
)

// CryptoResponse maps the JSON response from CoinCap API

type CryptoResponse struct {
	Data struct {
		ID string `json:"id"`
		Symbol string `json:"symbol"`
		Name string `json:"name"`
		PriceUSD string `json:"priceUsd"`
	}`json:"data"`
}

// CoinGeckoResponse maps the JSON response from CoinGecko API
type CoinGeckoResponse map[string]map[string]float64

type PageData struct {
	Crypto string 
	Price string
	PodName string
	Error string
}

//fetchCryptoPrice queries the CoinCap API for a given crypto asset

func fetchCryptoPrice(crypto string) (string, float64, error) {
	apiURL := fmt.Sprintf("https://api.coingecko.com/api/v3/simple/price?ids=%s&vs_currencies=usd", crypto)

	resp, err := http.Get(apiURL)
	if err != nil {
		return "", 0, fmt.Errorf("network/DNS error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", 0, fmt.Errorf("API returned status code %d", resp.StatusCode)
	}

	var data CoinGeckoResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", 0, fmt.Errorf("failed to parse response: %w", err)
	}

	priceData, exists := data[crypto]
	if !exists || priceData["usd"] == 0 {
		return "", 0, fmt.Errorf("asset '%s' not found or invalid", crypto)
	}

	return crypto, priceData["usd"], nil
}


func cryptoHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		http.Error(w, "Unable to load HTML template", http.StatusInternalServerError)
		return
	}

	asset := r.URL.Query().Get("asset")
	if asset == "" {
		asset = "bitcoin"
	}

	// Read POD_NAME environment variable (set automatically by Kubernetes later)
	podName := os.Getenv("POD_NAME")
	if podName == "" {
		podName = "Local Machine / Standalone Docker"
	}

	data := PageData{PodName: podName}

	name, price, err := fetchCryptoPrice(asset)
	if err != nil {
		data.Error = err.Error()
	} else {
		// data.Crypto = fmt.Sprintf("%s (%s)", cryptoData.Data.Name, cryptoData.Data.Symbol)
		// Format price to 2 decimal places
		// var priceFloat float64
		// fmt.Sscanf(cryptoData.Data.PriceUSD, "%f", &priceFloat)
		//data.Price = fmt.Sprintf("$%.2f USD", priceFloat)
		data.Crypto = name
    		data.Price = fmt.Sprintf("$%.2f USD", price)
	}

	tmpl.Execute(w, data)
}

func main() {
	http.HandleFunc("/", cryptoHandler)
	fmt.Println("Server running on http://localhost:8080...")
	http.ListenAndServe(":8080", nil)
}
