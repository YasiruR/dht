package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-ping/ping"
	"github.com/tryfix/log"
	"io/ioutil"
	"net/http"
)

type Client struct {
	successorGetKeyURL   string
	successorStoreKeyURL string
	*http.Client
}

var client Client

func initClient() {

}

func testPeerConn(ctx context.Context) {
	if !config.NeighbourCheck {
		log.DebugContext(ctx, `checking for neighbours is disabled`)
		return
	}

	p1, err := ping.NewPinger(config.Predecessor + `:` + config.PredecessorPort)
	if err != nil {
		log.Fatal(`creating new pinger failed`, config.Predecessor+`:`+config.PredecessorPort)
	}
	p1.Count = 3

	err = p1.Run()
	if err != nil {
		log.Fatal(`pinging to predecessor failed`, config.Predecessor+`:`+config.PredecessorPort)
	}

	p2, err := ping.NewPinger(config.Successor + `:` + config.SuccessorPort)
	if err != nil {
		log.Fatal(`creating new pinger failed`, config.Successor+`:`+config.SuccessorPort)
	}
	p2.Count = 3

	err = p2.Run()
	if err != nil {
		log.Fatal(`pinging to successor failed`, config.Successor+`:`+config.SuccessorPort)
	}
}

func (c *Client) proceedGetKey(req *http.Request) (string, error) {
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