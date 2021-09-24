package chord

import (
	"context"
	"dht/logger"
	"errors"
	"fmt"
	"github.com/go-ping/ping"
	"github.com/tryfix/log"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	url string
	*http.Client
}

var client *Client

func InitClient(ctx context.Context) {
	client = &Client{
		url:    `http://` + Config.Successor + `:` + Config.SuccessorPort + `/storage/`,
		Client: &http.Client{Timeout: time.Duration(Config.RequestTimeout) * time.Second},
	}
	logger.Log.InfoContext(ctx, `client initiated for neighbour requests`)
}

func TestPeerConn(ctx context.Context) {
	if !Config.NeighbourCheck {
		log.DebugContext(ctx, `checking for neighbours is disabled`)
		return
	}

	p1, err := ping.NewPinger(Config.Predecessor)
	if err != nil {
		log.Fatal(`creating new pinger failed`, Config.Predecessor)
	}
	p1.Count = 3

	err = p1.Run()
	if err != nil {
		log.Fatal(`pinging to predecessor failed`, Config.Predecessor)
	}

	p2, err := ping.NewPinger(Config.Successor)
	if err != nil {
		log.Fatal(`creating new pinger failed`, Config.Successor)
	}
	p2.Count = 3

	err = p2.Run()
	if err != nil {
		log.Fatal(`pinging to successor failed`, Config.Successor)
	}
}

func (c *Client) proceedGetKey(key string, req *http.Request) (string, int, error) {
	u, err := url.Parse(c.url + key)
	if err != nil {
		log.Error(err, `failed parsing successor url`)
		return "", 0, err
	}

	req.RequestURI = ""
	req.URL = u
	res, err := c.Client.Do(req)
	if err != nil {
		return "", 0, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", res.StatusCode, extractError(res)
	}

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", res.StatusCode, err
	}

	logger.Log.Debug(fmt.Sprintf(`proceeding %s for get key request from %s`, string(bytes), Config.Successor))
	return string(bytes), res.StatusCode, nil
}

func (c *Client) proceedStoreKey(key string, req *http.Request) (int, error) {
	u, err := url.Parse(c.url + key)
	if err != nil {
		log.Error(err, `failed parsing successor url`)
		return 0, err
	}

	req.RequestURI = ""
	req.URL = u
	res, err := c.Client.Do(req)
	if err != nil {
		return 0, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return 0, extractError(res)
	}

	logger.Log.Debug(fmt.Sprintf(`store key:%s was proceeded successfully to %s`, key, Config.Successor))
	return res.StatusCode, nil
}

func extractError(res *http.Response) error {
	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf(err.Error(), `reading err response body failed`)
	}

	return errors.New(string(bytes))
}
