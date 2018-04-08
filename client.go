package main

import (
	"encoding/json"
	"errors"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

var (
	invalidJsonError = errors.New("invalid JSON")
)

type ErgoNodeClient struct {
	Logger  *zap.SugaredLogger
	URL     string
	Client  *http.Client
	Timeout int
}

type GetBlocksRequest struct {
	Limit  int
	Offset int
}

type GetBlocksSuccessResponse struct {
	BlocksIds []BlockId
}

func (enc *ErgoNodeClient) fetch(requestUrl string, body []byte) (*http.Response, error) {
	enc.Logger.Infof("fetch url %s", requestUrl)
	req, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		enc.Logger.Errorf("got error at request create %s", err)
		return nil, err
	}
	resp, err := enc.Client.Do(req)
	if err != nil {
		enc.Logger.Errorf("got error response %s", err)
		return nil, err
	}
	return resp, err
}

func (enc *ErgoNodeClient) GetBlocks(req *GetBlocksRequest) (*GetBlocksSuccessResponse, error) {
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
	switch resp.StatusCode {
	case http.StatusOK:
		blocks := make([]BlockId, 0, 0)
		err := json.Unmarshal(data, &blocks)
		if err != nil {
			enc.Logger.Errorf("got error at JSON unmarshal: err=%s", err)
			return nil, invalidJsonError
		}
		enc.Logger.Infof("got %d blocks", len(blocks))
		getBlocksResponse := GetBlocksSuccessResponse{BlocksIds: blocks}
		return &getBlocksResponse, nil
	default:
		enc.Logger.Infof("got err='%s', responseCode=%d", string(data), resp.StatusCode)
		return nil, errors.New(string(data))
	}
}

func (enc *ErgoNodeClient) GetBlock(headerId string) (*Block, error) {
	requestUrl := enc.URL + "/blocks/" + headerId
	resp, err := enc.fetch(requestUrl, nil)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	switch resp.StatusCode {
	case http.StatusOK:
		block := &Block{}
		err := json.Unmarshal(data, &block)
		if err != nil {
			enc.Logger.Errorf("got error at JSON unmarshal: err=%s", err)
			return nil, invalidJsonError
		}
		enc.Logger.Infof("got block %+v", block)
		return block, nil
	default:
		errString := string(data)
		enc.Logger.Infof("got err='%v', responseCode=%d", errString, resp.StatusCode)
		return nil, errors.New(errString)
	}
}
