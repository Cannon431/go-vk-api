package vkapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	GET  = 0
	POST = 1
)

type VK struct {
	Token   string
	Version string
	Lang    string
}

type Error struct {
	Error struct {
		ErrorCode     int    `json:"error_code"`
		ErrorMsg      string `json:"error_msg"`
		RequestParams []struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		}
	} `json:"error"`
}

func (vk *VK) getUrl(method string) string {
	return "https://api.vk.com/method/" + method
}

func (vk *VK) getQuery(params map[string]string) url.Values {
	query := url.Values{}
	query.Add("access_token", vk.Token)
	query.Add("v", vk.Version)
	query.Add("lang", vk.Lang)

	for key, value := range params {
		query.Add(key, value)
	}

	return query
}

func (vk *VK) Request(_type int, method string, params map[string]string) ([]byte, error) {
	var resp *http.Response
	var err error
	var data []byte
	var errorResponse Error

	u := vk.getUrl(method)
	query := vk.getQuery(params)

	if _type == GET {
		resp, err = http.Get(u + "?" + query.Encode())
	} else if _type == POST {
		resp, err = http.PostForm(u, query)
	} else {
		return data, errors.New("undefined request type")
	}

	if err != nil {
		return data, err
	}

	defer resp.Body.Close()

	data, err = ioutil.ReadAll(resp.Body)

	if err != nil {
		return data, err
	}

	err = json.Unmarshal(data, &errorResponse)

	if err != nil {
		return data, err
	}

	if errorResponse.Error.ErrorMsg != "" {
		return data, errors.New(fmt.Sprintf(
			"request error. [%d] %s",
			errorResponse.Error.ErrorCode,
			errorResponse.Error.ErrorMsg,
		))
	}

	return data, nil
}
