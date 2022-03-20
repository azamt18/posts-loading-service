package gorest_api_service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"posts.loader/pkg/models"
	"time"
)

type GoRestApiService interface {
	GetPosts() (*[]models.Post, error)
}

type goRestApiService struct {
	host string
}

func (g goRestApiService) GetPosts() (*[]models.Post, error) {
	var (
		client               = http.Client{Timeout: 30 * time.Second}
		req    *http.Request = nil
		resp                 = &http.Response{}
		err    error         = nil
	)

	// prepare request
	{
		req, _ = http.NewRequest(http.MethodGet, fmt.Sprintf("%s/public/v2/posts	", g.host), nil)
		req.Header.Add("Accept", "application/json")

		// make a request
		if resp, err = client.Do(req); err != nil {
			return nil, err
		}

		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return nil, errors.New("unexpected response")
		}

	}

	body, err := ioutil.ReadAll(resp.Body)

	// Unmarshal result
	var posts []models.Post
	err = json.Unmarshal(body, &posts)
	if err != nil {
		log.Printf("Reading body failed: %s", err)
	}

	//err = json.NewDecoder(resp.Body).Decode(&post)
	//if err != nil {
	//	return nil, err
	//}

	return &posts, nil
}

func NewGoRestApiService(host string) GoRestApiService {
	return &goRestApiService{
		host: host,
	}
}
