package zinc_sdk

import (
	"context"
	"fmt"
	"github.com/rolandhe/go-base/commons"
	"github.com/rolandhe/go-base/https_sdks/http11"
	"net/http"
	"strings"
	"time"
)

type ZincIndexer struct {
	host        string
	authToken   string
	httpTimeout time.Duration
}

func NewZincIndexer(host string, httpTimeout time.Duration, authToken string) *ZincIndexer {
	host = buildHostUrl(host, authToken)
	return &ZincIndexer{
		host:        host,
		authToken:   authToken,
		httpTimeout: httpTimeout,
	}
}

func (zi *ZincIndexer) CreateDocument(bc *commons.BaseContext, bodyJsonFunc func() (string, int64, error), indexName string) error {
	jsonBody, docId, err := bodyJsonFunc()
	if err != nil {
		return err
	}
	if jsonBody == "" {
		return commons.NewError(commons.BadRequest, "CreateIndex failed. Body is empty")
	}
	const documentUrlFmt = "%s/api/%s/_doc/%d"
	targetUrl := fmt.Sprintf(documentUrlFmt, zi.host, indexName, docId)
	ctx, cancel := context.WithTimeout(context.Background(), zi.httpTimeout)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, targetUrl, strings.NewReader(jsonBody))
	if err != nil {
		return err
	}
	setAuthHeader(req, zi.authToken)
	req.Header.Set("Content-Type", "application/json")
	result, _, err := http11.CallWithResult[CreateDocumentResult](bc, req, nil)
	if err != nil {
		return err
	}
	if result.Message == "ok" {
		return nil
	}
	return commons.NewError(commons.CommonErr, result.Message)
}

func (zi *ZincIndexer) UpdateDocument(bc *commons.BaseContext, bodyJsonFunc func() (string, int64, error), indexName string) error {
	jsonBody, docId, err := bodyJsonFunc()
	if err != nil {
		return err
	}
	if jsonBody == "" {
		return commons.NewError(commons.BadRequest, "UpdateDocument failed. Body is empty")
	}
	const documentUrlFmt = "%s/api/%s/_update/%d"
	targetUrl := fmt.Sprintf(documentUrlFmt, zi.host, indexName, docId)
	ctx, cancel := context.WithTimeout(context.Background(), zi.httpTimeout)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, targetUrl, strings.NewReader(jsonBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	setAuthHeader(req, zi.authToken)
	result, _, err := http11.CallWithResult[UpdateDocumentResult](bc, req, nil)
	if err != nil {
		return err
	}
	if result.Id == "" {
		return ErrDocumentNoExist
	}
	return nil
}

func (zi *ZincIndexer) DeleteDocument(bc *commons.BaseContext, docId int64, indexName string) error {
	const documentUrlFmt = "%s/api/%s/_doc/%d"
	targetUrl := fmt.Sprintf(documentUrlFmt, zi.host, indexName, docId)
	ctx, cancel := context.WithTimeout(context.Background(), zi.httpTimeout)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, targetUrl, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	setAuthHeader(req, zi.authToken)
	result, _, err := http11.CallWithResult[DeleteDocumentResult](bc, req, idNotFoundRespCallback)
	if err != nil {
		return err
	}
	if result.Id == "" {
		return commons.NewError(commons.BadRequest, result.Message)
	}
	return nil
}

func (zi *ZincIndexer) GetDocumentById(bc *commons.BaseContext, docId int64, indexName string) (string, error) {
	const documentUrlFmt = "%s/api/%s/_doc/%d"
	targetUrl := fmt.Sprintf(documentUrlFmt, zi.host, indexName, docId)
	ctx, cancel := context.WithTimeout(context.Background(), zi.httpTimeout)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, targetUrl, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	setAuthHeader(req, zi.authToken)
	body, _, err := http11.CallWithResult[string](bc, req, idNotFoundRespCallback)
	if err != nil {
		return "", err
	}
	if body == nil {
		return "", ErrDocumentNoExist
	}
	return *body, nil
}

func (zi *ZincIndexer) DeleteIndex(bc *commons.BaseContext, indexName string) error {
	const deleteIndexUrlFmt = "%s/api/index/%s"
	targetUrl := fmt.Sprintf(deleteIndexUrlFmt, zi.host, indexName)
	ctx, cancel := context.WithTimeout(context.Background(), zi.httpTimeout)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, targetUrl, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	setAuthHeader(req, zi.authToken)
	result, _, err := http11.CallWithResult[DeleteIndexResult](bc, req, nil)
	if err != nil {
		return err
	}
	if result.Message == "deleted" {
		return nil
	}
	msg := result.Error
	if msg == fmt.Sprintf("index %s does not exists", indexName) {
		return nil
	}
	if msg == "" {
		msg = result.Message
	}
	return commons.NewError(commons.CommonErr, msg)
}
func (zi *ZincIndexer) CreateIndex(bc *commons.BaseContext, bodyJson string) error {
	const createIndexUrlFmt = "%s/api/index"
	targetUrl := fmt.Sprintf(createIndexUrlFmt, zi.host)
	ctx, cancel := context.WithTimeout(context.Background(), zi.httpTimeout)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, targetUrl, strings.NewReader(bodyJson))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	setAuthHeader(req, zi.authToken)
	result, _, err := http11.CallWithResult[CreateIndexResult](bc, req, nil)
	if err != nil {
		return err
	}
	if result.Message != "ok" {
		msg := result.Error
		if msg == "" {
			msg = result.Message
		}
		return commons.NewError(commons.CommonErr, msg)
	}
	return nil
}

func (zi *ZincIndexer) CreateEsSearcher() *ESSearcher {
	return NewESSearcher(zi.host, zi.httpTimeout, zi.authToken)
}

type CreateDocumentResult struct {
	Message     string `json:"message"`
	Id          string `json:"_id"`
	Index       string `json:"_index"`
	Version     int    `json:"_version"`
	SeqNo       int    `json:"_seq_no"`
	PrimaryTerm int    `json:"_primary_term"`
	Result      string `json:"result"`
}

type UpdateDocumentResult struct {
	Message string `json:"message"`
	Id      string `json:"id"`
}

type DeleteDocumentResult struct {
	Message string `json:"message"`
	Index   string `json:"index"`
	Id      string `json:"id"`
}

type DeleteIndexResult struct {
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

type CreateIndexResult struct {
	Message     string `json:"message,omitempty"`
	Index       string `json:"index"`
	StorageType string `json:"storage_type"`
	Error       string `json:"error,omitempty"`
}
