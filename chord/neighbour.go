package chord

import (
	"context"
	"dht/logger"
	"encoding/json"
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
	url    *url.URL
	*http.Client
}

var client *Client

func InitClient(ctx context.Context) {
	u, err := url.Parse(`http://` + Config.Successor + `:` + Config.SuccessorPort + `/storage/key`)
	if err != nil {
		log.FatalContext(ctx, err, `failed parsing successor url`)
	}

	client = &Client{
		url:    u,
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

func (c *Client) proceedGetKey(req *http.Request) (string, error) {
	req.RequestURI = ""
	req.URL = c.url
	res, err := c.Client.Do(req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", extractError(res)
	}

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func (c *Client) proceedStoreKey(req *http.Request) (int, error) {
	req.RequestURI = ""
	req.URL = c.url

	res, err := c.Client.Do(req)
	if err != nil {
		return 0, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return 0, extractError(res)
	}

	return res.StatusCode, nil
}

func extractError(res *http.Response) error {
	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf(err.Error(), `reading err response body failed`)
	}

	var msg response
	err = json.Unmarshal(bytes, &msg)
	if err != nil {
		return fmt.Errorf(err.Error(), `unmarshalling err response failed`)
	}

	return errors.New(msg.Error)
}
