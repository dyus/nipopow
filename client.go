package main

import (
	"net/http"
	"time"
	"errors"
	"strconv"
	"encoding/json"
	"net/url"
	"io/ioutil"
	"go.uber.org/zap"
)

var (
	client = &http.Client{Timeout: time.Second}
	invalidJsonError = errors.New("invalid JSON")
)

type ErgoNodeClient struct {
	Logger  *zap.SugaredLogger
	URL     string
	Timeout int
}

type Block string

type GetBlocksRequest struct {
	Limit  int
	Offset int
}

type GetBlocksSucessResponse struct {
	Blocks []Block
}

func (enc *ErgoNodeClient) fetch(requestUrl string, body []byte) (*http.Response, error) {
	enc.Logger.Infof("fetch url %s", requestUrl)
	req, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		enc.Logger.Errorf("got error at request create %s", err)
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		enc.Logger.Errorf("got error response %s", err)
		return nil, err
	}
	return resp, err
}

func (enc *ErgoNodeClient) GetBlocks(req *GetBlocksRequest) (*GetBlocksSucessResponse, error) {
	queryParams := url.Values{}
	if req.Limit != 0 {
		queryParams.Add("limit", strconv.Itoa(req.Limit))
	}
	if req.Offset != 0 {
		queryParams.Add("offset", strconv.Itoa(req.Offset))
	}
	requestUrl := enc.URL + "/blocks"

	queryParamsString := queryParams.Encode()
	if queryParamsString != "" {
		requestUrl += "?" + queryParamsString
	}
	resp, err := enc.fetch(requestUrl, nil)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	//fmt.Println(string(data))
	switch resp.StatusCode {
	case http.StatusOK:
		blocks := make([]Block, 0, 0)
		err := json.Unmarshal(data, &blocks)
		if err != nil {
			enc.Logger.Errorf("got error at JSON unmarshal: err=%s", err)
			return nil, invalidJsonError
		}
		enc.Logger.Infof("got %d blocks", len(blocks))
		getBlocksResponse := GetBlocksSucessResponse{Blocks: blocks}
		return &getBlocksResponse, nil
	default:
		enc.Logger.Infof("got err='%s', responseCode=%d", err, resp.StatusCode)
		return nil, errors.New(string(data))
	}
}
