package consumer

import (
	"encoding/json"

	"github.com/Wundagor/high-throughput-data-consumer/models"
)

func decodeMessage(body []byte) (*models.SourceData, error) {
	var data models.SourceData
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}
	return &data, nil
}
