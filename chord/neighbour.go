package chord

import (
	"context"
	"dht/logger"
	"errors"
	"fmt"
	"github.com/tryfix/log"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Client struct {
	predHostname string
	predLock     *sync.Mutex
	lostPredList []string
	sucHostname  string
	sucLock      *sync.Mutex
	crashClient  *http.Client
	*http.Client
}

var neighbors *Client

func InitClient(ctx context.Context) {
	neighbors = &Client{
		predLock:     &sync.Mutex{},
		sucLock:      &sync.Mutex{},
		predHostname: node.hostname,
		sucHostname:  node.hostname,
		crashClient:  &http.Client{Timeout: time.Duration(Config.DetectCrashTimeout*100) * time.Second},
		Client:       &http.Client{Timeout: time.Duration(Config.RequestTimeout) * time.Second},
	}
	neighbors.initProbe()
	logger.Log.InfoContext(ctx, `client initiated for neighbour requests`)
}

func (c *Client) initProbe() {
	go func() {
		ticker := time.NewTicker(time.Duration(Config.ProbeInterval) * time.Second)
	probeLoop:
		for {
			select {
			case <-ticker.C:
				if node.single {
					continue
				}

				if len(c.lostPredList) != 0 {
					for i, p := range c.lostPredList {
						_, err := c.Client.Get(`http://` + p + `/internal/probe`)
						if err != nil {
							continue
						}

						joined := false
						func(joined *bool) {
							req, err := http.NewRequest(http.MethodPost, `http://`+p+`/join`, nil)
							if err != nil {
								logger.Log.Error(err, `creating request for join of crashed node failed`)
								return
							}
							q := req.URL.Query()
							if node.single {
								q.Add(`nprime`, node.hostname)
							} else {
								// join request is forwarded to successor directly to reduce cost of http requests
								q.Add(`nprime`, c.sucHostname)
							}
							req.URL.RawQuery = q.Encode()

							res, err := c.Client.Do(req)
							if err != nil {
								logger.Log.Error(err, `join request to crashed node failed`)
								return
							}
							defer res.Body.Close()

							if res.StatusCode != http.StatusOK {
								logger.Log.Debug(`received non-2xx response for join of crashed node`, res.StatusCode)
								return
							}

							// updating lost list with only better predecessors
							c.lostPredList = c.lostPredList[:i]
							*joined = true
						}(&joined)

						if joined == true {
							logger.Log.Trace(`crashed node join was successful`, c.lostPredList)
							continue probeLoop
						}
					}
				}

				// enclosed in an unnamed function to handle closing body properly
				func() {
					res, err := c.Client.Get(`http://` + c.predHostname + `/internal/probe`)
					if err != nil {
						if netErr, ok := err.(net.Error); ok {
							if netErr.Timeout() {
								lastNode, err := c.proceedFixCrash(node.hostname)
								if err != nil {
									logger.Log.Error(err, `initiating fix crash failed`)
									return
								}

								if lastNode == noCrashResponse {
									logger.Log.Debug(`fix crash was initiated but could not find any defect`)
									return
								}

								// no lock required since it will be always this go-routine that updates the list
								c.lostPredList = append(c.lostPredList, c.predHostname)
								c.updatePredecessor(lastNode)
								logger.Log.Debug(`broken network ring was detected and fixed`, lastNode)
								return
							}
						}
						logger.Log.Error(err, `probing predecessor failed`)
						return
					}
					defer res.Body.Close()
				}()
			}
		}
	}()
}

func (c *Client) updateSuccessor(hostname string) {
	c.sucLock.Lock()
	defer c.sucLock.Unlock()
	c.sucHostname = hostname
	node.updateSucId(hostname)
}

func (c *Client) updatePredecessor(hostname string) {
	c.predLock.Lock()
	defer c.predLock.Unlock()
	c.predHostname = hostname
	node.updatePredId(hostname)
}

func (c *Client) clearNeighbors() {
	c.sucLock.Lock()
	defer c.sucLock.Unlock()
	c.sucHostname = node.hostname

	c.predLock.Lock()
	defer c.predLock.Unlock()
	c.predHostname = node.hostname
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
		u = `http://` + c.predHostname + `/internal/update-successor/` + hostname
	} else {
		u = `http://` + c.sucHostname + `/internal/update-predecessor/` + hostname
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

func (c *Client) proceedFixCrash(initiator string) (string, error) {
	ticker := time.NewTicker(time.Duration(Config.DetectCrashTimeout) * time.Second)
	resChan := make(chan string)
	errChan := make(chan error)
	go func() {
		res, err := c.crashClient.Post(`http://`+c.sucHostname+`/internal/fix-crash/`+initiator, `text/plain`, nil)
		if err != nil {
			logger.Log.Error(err, `proceeding fix crash request failed`, initiator)
			errChan <- err
			return
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			logger.Log.Error(`received non-2xx response code`, res.StatusCode)
			errChan <- errors.New(`request failed`)
			return
		}

		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			logger.Log.Error(err, `reading fix crash response body failed`, initiator)
			errChan <- err
			return
		}

		resChan <- string(data)
	}()

	for {
		select {
		case <-ticker.C:
			// if detect period is reached check if the successor is alive and if so wait longer for the response. If not,
			// current node is the last node.
			_, err := c.Client.Get(`http://` + c.sucHostname + `/internal/probe`)
			if err != nil {
				if netErr, ok := err.(net.Error); ok {
					if netErr.Timeout() {
						c.updateSuccessor(initiator)
						logger.Log.Trace(`timed out as successor is out of reach`)
						return node.hostname, nil
					}
				}
				logger.Log.Error(err, `reaching successor in fix crash failed`, initiator)
			}
		case lastNode := <-resChan:
			logger.Log.Trace(`immediate successor was reached successfully`, lastNode)
			return lastNode, nil
		case err := <-errChan:
			return "", err
		}
	}
}

func (c *Client) proceedGetClusterInfo(hostname string, req *http.Request) (string, error) {
	u, err := url.Parse(`http://` + c.sucHostname + `/internal/cluster-info/` + hostname)
	if err != nil {
		log.Error(err, `failed parsing successor get nodes url`)
		return "", err
	}

	req.RequestURI = ""
	req.URL = u
	res, err := c.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", extractError(res)
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logger.Log.Error(err, `reading get nodes body failed`)
		return "", err
	}

	num, err := strconv.Atoi(string(data))
	if err != nil {
		logger.Log.Error(err, `converting get nodes response to number failed`, string(data))
		return "", err
	}

	num++
	logger.Log.Debug(fmt.Sprintf(`get nodes was successfully proceeded to %s`, c.sucHostname))
	return strconv.Itoa(num), nil
}

func extractError(res *http.Response) error {
	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf(err.Error(), `reading err response body failed`)
	}

	return errors.New(string(bytes))
}
