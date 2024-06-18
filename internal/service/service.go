package service

import (
	"Epic55/go_currency_app2/internal/models"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Service struct {
	Client *http.Client
}

func NewService() *Service {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	return &Service{
		Client: client,
	}
}

func (s *Service) GetData1(ctx context.Context, data string, APIURL string) *models.Rates {
	//start := time.Now()
	apiURL := fmt.Sprintf("%s?fdate=%s", APIURL, data)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		fmt.Println("Failed to create request with context", err)
		return nil
	}

	resp, err := s.Client.Do(req)
	if err != nil {
		fmt.Println("Failed to GET URL", err)
		return nil
	}
	defer resp.Body.Close()

	//

	xmlData, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Failed to Read response Body", err)
		return nil
	}

	var rates *models.Rates
	if err := xml.Unmarshal(xmlData, &rates); err != nil {
		fmt.Println("Failed to parse XML data", err)
		return nil
	}
	return rates
}
