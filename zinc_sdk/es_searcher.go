package zinc_sdk

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rolandhe/go-base/commons"
	"github.com/rolandhe/go-base/https_sdks/http11"
	"github.com/rolandhe/go-base/logger"
	"net/http"
	"time"
)

type ESSearcher struct {
	host        string
	authToken   string
	httpTimeout time.Duration
}

type WithScore interface {
	UseScore(score float64)
}

func NewESSearcher(host string, httpTimeout time.Duration, authToken string) *ESSearcher {
	host = buildHostUrl(host, authToken)
	return &ESSearcher{
		host:        host,
		authToken:   authToken,
		httpTimeout: httpTimeout,
	}
}

func (es *ESSearcher) searchCore(bc *commons.BaseContext, indexName string, queryJson []byte) (string, error) {
	const urlFormat = "%s/es/%s/_search"
	targetUrl := fmt.Sprintf(urlFormat, es.host, indexName)

	logger.WithBaseContextInfof(bc)("search zinc %s body:%s", es.host, string(queryJson))

	//jsonBody, _ := json.Marshal(queryJson)
	ctx, cancel := context.WithTimeout(context.Background(), es.httpTimeout)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, targetUrl, bytes.NewReader(queryJson))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	setAuthHeader(req, es.authToken)

	result, _, err := http11.CallWithResult[string](bc, req, nil)
	if err != nil {
		return "", err
	}
	if result == nil {
		return "", nil
	}
	return *result, nil
}

func Search[T any](bc *commons.BaseContext, es *ESSearcher, indexName string, query *ESQueryReq) ([]*T, error) {
	jsonBody, _ := json.Marshal(query)

	return searchCore[T](bc, es, indexName, jsonBody)
}

func SearchAll[T any](bc *commons.BaseContext, es *ESSearcher, indexName string, limit int) ([]*T, error) {
	if limit <= 0 || limit > 10000 {
		limit = 10000
	}
	queryTpl := `{
			"query": {
				"match_all": {}
			},
			"size": %d
		}`
	queryStr := fmt.Sprintf(queryTpl, limit)
	return searchCore[T](bc, es, indexName, []byte(queryStr))
}

func SearchByIds[T any, I string | int64](bc *commons.BaseContext, es *ESSearcher, indexName string, ids []I, viewFields ...string) ([]*T, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	if len(ids) > 10000 {
		return nil, errors.New("too many ids,exceed 10000")
	}

	idsCondition := &TermsCondition[I]{
		Terms: map[string]any{
			"id": ids,
		},
	}
	req := &ESQueryReq{
		Query: idsCondition,
		Size:  len(ids),
	}
	if len(viewFields) > 0 {
		req.Source = viewFields
	}
	jsonBody, _ := json.Marshal(req)

	return searchCore[T](bc, es, indexName, jsonBody)
}

func searchCore[T any](bc *commons.BaseContext, esSearch *ESSearcher, indexName string, queryJson []byte) ([]*T, error) {
	body, err := esSearch.searchCore(bc, indexName, queryJson)
	if err != nil {
		return nil, err
	}
	var esResult ESResult[T]
	err = json.Unmarshal([]byte(body), &esResult)
	if err != nil {
		return nil, err
	}

	if esResult.Error != "" {
		return nil, commons.NewError(commons.CommonErr, esResult.Error)
	}
	if len(esResult.Hits.Hits) == 0 {
		return nil, nil
	}
	list := make([]*T, 0, len(esResult.Hits.Hits))
	for _, hit := range esResult.Hits.Hits {
		ws, ok := any(hit.Source).(WithScore)
		if ok {
			ws.UseScore(hit.Score)
		}
		list = append(list, hit.Source)
	}
	return list, nil
}

type ESResult[T any] struct {
	Took     int  `json:"took"`
	TimedOut bool `json:"timed_out"`
	Shards   struct {
		Total      int `json:"total"`
		Successful int `json:"successful"`
		Skipped    int `json:"skipped"`
		Failed     int `json:"failed"`
	} `json:"_shards"`
	Hits struct {
		Total struct {
			Value int `json:"value"`
		} `json:"total"`
		MaxScore float64           `json:"max_score"`
		Hits     []*HitDocument[T] `json:"hits"`
	} `json:"hits"`
	Error string `json:"error"`
}

type HitDocument[T any] struct {
	Index     string    `json:"_index"`
	Type      string    `json:"_type"`
	Id        string    `json:"_id"`
	Score     float64   `json:"_score"`
	Timestamp time.Time `json:"@timestamp"`
	Source    *T        `json:"_source"`
}
