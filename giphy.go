package main

import (
	"encoding/json"
	"errors"

	"github.com/go-resty/resty/v2"
)

type GiphyHandler struct {
	client *resty.Client
	token  string
}

type GiphyGif struct {
	Id  string `json:"id"`
	Url string `json:"url"`
}

type GiphySearchResponse struct {
	Data []*GiphyGif `json:"data"`
}

func (g *GiphyHandler) SearchFirst(q string) (*GiphyGif, error) {
	req := g.client.R()
	req.SetQueryParam("api_key", g.token)
	req.SetQueryParam("q", q)
	req.SetQueryParam("limit", "1")

	resp, err := req.Get("https://api.giphy.com/v1/gifs/search")
	if err != nil {
		return nil, err
	}

	response := GiphySearchResponse{}
	err = json.Unmarshal(resp.Body(), &response)
	if err != nil {
		return nil, err
	}

	gifs := response.Data
	if len(gifs) == 0 {
		return nil, errors.New("no gifs found")
	}

	return gifs[0], nil
}

func NewGiphyHandler(token string) *GiphyHandler {
	client := resty.New()
	giphyHandler := GiphyHandler{client, token}
	return &giphyHandler
}
