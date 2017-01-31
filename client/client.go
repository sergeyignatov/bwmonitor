package client

import (
	"fmt"
	"github.com/sergeyignatov/bwmonitor/common"
	"io/ioutil"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyz"

type Client struct {
	dest     string
	client   *http.Client
	context  *common.Context
	minBytes int
}

func genetateRandName(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func (c *Client) setHistory(delta time.Duration) {

	c.context.History[c.dest] = common.History{delta, c.minBytes}
}
func (c *Client) getHistory() (common.History, bool) {
	t, ok := c.context.History[c.dest]
	return t, ok
}

func NewClient(dest string, timeout int, context *common.Context) *Client {
	client := &http.Client{
		Timeout: time.Second * time.Duration(timeout),
	}

	return &Client{dest, client, context, 1 * 1024 * 1024}
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

	body, _ := ioutil.ReadAll(resp.Body)
	return len(body), time.Now().Sub(start).Seconds(), nil
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
	minBytes := 128 * 1024
	maxBytes := 100 * 1024 * 1024
	if t, ok := c.getHistory(); ok != false {
		var tmp int
		if t.Duration.Seconds() < 1 {
			tmp = int(float64(t.MinBytes) * (1.0 / t.Duration.Seconds()))
		} else {
			tmp = int(float64(t.MinBytes) / t.Duration.Seconds())
		}
		if tmp > minBytes && tmp < maxBytes {
			t.MinBytes = tmp
		} else if tmp < minBytes {
			t.MinBytes = minBytes
		} else if tmp >= maxBytes {
			t.MinBytes = maxBytes
		}
		c.minBytes = t.MinBytes
	}
	//fmt.Println(c.minBytes)
	for t := 0; t < connections; t++ {
		go func(tt int) {
			s, t, err := c.download(random(c.minBytes/2, 2*c.minBytes))
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
	}
	c.setHistory(time.Now().Sub(start))
	return int(tb / time.Now().Sub(start).Seconds()), nil
}
