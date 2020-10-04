package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const picsumURL = "https://picsum.photos/"

// NewPicSum returns a new picsum client
func NewPicSum(c *http.Client) *PicSum {
	if c != nil {
		return &PicSum{
			c: c,
		}
	}
	return &PicSum{
		c: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// PicSum client
type PicSum struct {
	c *http.Client
}

// Image basic image metadata returned by picsum
type Image struct {
	ID          string `json:"id"`
	Author      string `json:"author"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
	URL         string `json:"url"`
	DownloadURL string `json:"download_url"`
}

// MakeURL creates a new url given a width
func (i Image) MakeURL(width int) string {
	height := float64(width) * float64(i.Height) / float64(i.Width)
	return fmt.Sprintf("https://picsum.photos/id/%s/%d/%d", i.ID, width, int(height))
}

//ListImages lists images from picsum
func (p *PicSum) ListImages() ([]Image, error) {
	url := picsumURL + "v2/list"
	res, err := p.c.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to list images %w", err)
	}
	imgs := []Image{}
	d := json.NewDecoder(res.Body)
	err = d.Decode(&imgs)
	if err != nil {
		return nil, fmt.Errorf("failed to decode images %w", err)
	}
	return imgs, nil
}
