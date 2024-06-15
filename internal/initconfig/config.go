package initconfig

import (
	"Epic55/go_currency_app2/internal/models"
	"encoding/json"
	"os"
)

func InitConfig(filename string) (*models.Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var Config models.Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&Config); err != nil {
		return nil, err
	}
	return &Config, nil
}
