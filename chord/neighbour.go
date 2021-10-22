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
	"strings"
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
		Client: &http.Client{Timeout: time.Duration(Config.RequestTimeout) * time.Second},
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

func (c *Client) clearNeighbors() {
	c.sucHostname = ""
	c.predHostname = ""
	node.leave()
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

func (c *Client) initJoin(networkHost string) (string, string, error) {
	res, err := c.Client.Post(`http://`+networkHost+`/internal/join/`+node.hostname, `application/json`, nil)
	if err != nil {
		log.Error(err, `failed initiating internal join request`)
		return ``, ``, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return ``, ``, extractError(res)
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Error(err, `failed reading json body of internal join response`)
		return ``, ``, err
	}

	// additional processing is required as response body is in raw string format - 'predecessor,successor'
	newNeighbors := strings.Split(string(data)[1:len(data)-2], ",")
	if len(newNeighbors) != 2 {
		return ``, ``, fmt.Errorf(`returned an invalid response for neighbors upon joining [res: %s]`, newNeighbors)
	}

	logger.Log.Debug(fmt.Sprintf(`internal join was initiated successfully to %s`, networkHost))
	return newNeighbors[0], newNeighbors[1], nil
}

func (c *Client) proceedJoin(hostname string, req *http.Request) ([]byte, error) {
	u, err := url.Parse(`http://` + c.sucHostname + `/internal/join/` + hostname)
	if err != nil {
		log.Error(err, `failed parsing successor join url`)
		return nil, err
	}

	req.RequestURI = ""
	req.URL = u
	res, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, extractError(res)
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logger.Log.Error(err, `reading internal join proceed response failed`)
		return nil, err
	}

	logger.Log.Debug(fmt.Sprintf(`internal join for hostname:%s was proceeded successfully to %s`, hostname, c.sucHostname))
	return data, nil
}

func (c *Client) notifyNeighbor(hostname string, predecessor bool) error {
	var u string
	if predecessor {
		u = `http://`+c.predHostname+`/internal/update-successor/`+hostname
	} else {
		u = `http://`+c.sucHostname+`/internal/update-predecessor/`+hostname
	}

	res, err := c.Client.Post(u, `application/json`, nil)
	if err != nil {
		log.Error(err, fmt.Sprintf(`failed sending inform neighbor request [predecessor: %t]`, predecessor))
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return extractError(res)
	}

	logger.Log.Debug(fmt.Sprintf(`internal update of neighbor was proceeded successfully [predecessor: %t]`, predecessor))
	return nil
}

func extractError(res *http.Response) error {
	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf(err.Error(), `reading err response body failed`)
	}

	return errors.New(string(bytes))
}
