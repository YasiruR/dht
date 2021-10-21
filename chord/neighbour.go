package chord

import (
	"context"
	"dht/logger"
	"errors"
	"fmt"
	"github.com/tryfix/log"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	predHostname string
	sucHostname  string
	*http.Client
}

var neighbors *Client

func InitClient(ctx context.Context) {
	neighbors = &Client{
		Client:   &http.Client{Timeout: time.Duration(Config.RequestTimeout) * time.Second},
	}
	logger.Log.InfoContext(ctx, `client initiated for neighbour requests`)
}

func (c *Client) updateSuccessor(hostname string) {
	c.sucHostname = hostname
	node.updateSucId(hostname)
}

func (c *Client) updatePredecessor(hostname string) {
	c.predHostname = hostname
	node.updatePredId(hostname)
}

func (c *Client) proceedGetKey(key string, req *http.Request) (string, int, error) {
	u, err := url.Parse(`http://` + c.sucHostname + `/storage/` + key)
	if err != nil {
		log.Error(err, `failed parsing successor storeUrl`)
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

	logger.Log.Debug(fmt.Sprintf(`proceeding %s for get key request from %s`, string(bytes), c.sucHostname))
	return string(bytes), res.StatusCode, nil
}

func (c *Client) proceedStoreKey(key string, req *http.Request) (int, error) {
	u, err := url.Parse(`http://` + c.sucHostname + `/storage/` + key)
	if err != nil {
		log.Error(err, `failed parsing successor storeUrl`)
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

	logger.Log.Debug(fmt.Sprintf(`store key:%s was proceeded successfully to %s`, key, c.sucHostname))
	return res.StatusCode, nil
}

func (c *Client) proceedJoin(key string, req *http.Request) (int, error) {
	u, err := url.Parse(`http://` + c.sucHostname + `/internal/join/` + key)
	if err != nil {
		log.Error(err, `failed parsing successor join url`)
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

	logger.Log.Debug(fmt.Sprintf(`internal join for key:%s was proceeded successfully to %s`, key, c.sucHostname))
	return res.StatusCode, nil
}

func (c *Client) notifyPredecessor(hostname string) error {
	res, err := c.Client.Post(`http://` + c.predHostname + `/internal/update-successor` + hostname, `application/json`, nil)
	if err != nil {
		log.Error(err, `failed sending inform predecessor request`)
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return extractError(res)
	}

	logger.Log.Debug(fmt.Sprintf(`internal update of predecssor was proceeded successfully to %s`, c.predHostname))
	return nil
}

func extractError(res *http.Response) error {
	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf(err.Error(), `reading err response body failed`)
	}

	return errors.New(string(bytes))
}
