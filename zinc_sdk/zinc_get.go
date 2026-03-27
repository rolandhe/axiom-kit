package zinc_sdk

import (
	"encoding/json"
	"github.com/rolandhe/go-base/commons"
	"time"
)

func GetDocumentById[T any](bc *commons.BaseContext, zincIndex *ZincIndexer, docId int64, indexName string) (*T, error) {
	body, err := zincIndex.GetDocumentById(bc, docId, indexName)
	if err != nil {
		return nil, err
	}
	var result GetDocumentResult[T]
	err = json.Unmarshal([]byte(body), &result)
	if err != nil {
		return nil, err
	}
	return result.Source, nil
}

type GetDocumentResult[T any] struct {
	Index     string    `json:"_index"`
	Type      string    `json:"_type"`
	Id        string    `json:"_id"`
	Score     int       `json:"_score"`
	Timestamp time.Time `json:"@timestamp"`
	Error     string    `json:"error"`
	Source    *T        `json:"_source"`
}
