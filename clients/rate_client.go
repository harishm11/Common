package clients

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	ratingmodels "github.com/harishm11/Common/models/RatingModels"
)


const RATING_SERVICE_URL = "http://localhost:8002/rate"

func CallRatingService(rateRequest *ratingmodels.RateRequest) (ratingmodels.RateResponse, error) {
	client := &http.Client{Timeout: 10 * time.Second}

	ratingServiceURL := RATING_SERVICE_URL
	// ratingServiceURL := os.Getenv("RATING_SERVICE_URL")
	requestBody, err := json.Marshal(rateRequest)
	var rateResponse ratingmodels.RateResponse
	if err != nil {
		return rateResponse, err
	}

	resp, err := client.Post(ratingServiceURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return rateResponse, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return rateResponse, errors.New("failed to get valid response from rating service")
	}

	if err := json.NewDecoder(resp.Body).Decode(&rateResponse); err != nil {
		return rateResponse, err
	}

	return rateResponse, nil
}
