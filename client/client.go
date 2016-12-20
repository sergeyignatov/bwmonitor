package client

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyz"

type Client struct {
	dest   string
	client *http.Client
}

func genetateRandName(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func NewClient(dest string) *Client {
	client := &http.Client{
		Timeout: time.Second * 15,
	}

	return &Client{dest, client}
}

func (c *Client) download(size int) (int, float64, error) {
	start := time.Now()
	req, err := http.NewRequest("GET", fmt.Sprintf("http://%s:5312/api/1.0/bw/%s?size=%d", c.dest, genetateRandName(8), size), nil)
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("User-Agent", "BWMonitor")
	resp, err := c.client.Do(req)
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, 0, err
	}
	bodyLen := len(body)
	finish := time.Now()
	return bodyLen, finish.Sub(start).Seconds(), nil
}

type R struct {
	bodylen int
	spend   float64
	err     error
}

func random(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

func (c *Client) DownloadSpeed() (int, error) {
	var wg sync.WaitGroup
	connections := 4
	wg.Add(connections)
	ch := make(chan R, 0)
	start := time.Now()
	for t := 0; t < connections; t++ {
		go func(tt int) {
			s, t, err := c.download(random(1*1024*1024, 5*1024*1024))
			ch <- R{s, t, err}
			wg.Done()
		}(t)
	}
	go func() {
		wg.Wait()
		close(ch)
	}()
	var tb, tt float64
	for t := range ch {
		tb += float64(t.bodylen)
		tt += t.spend
		if t.err != nil {
			return 0, t.err
		}
	}
	return int(tb / time.Now().Sub(start).Seconds()), nil
}
