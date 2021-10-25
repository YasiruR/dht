package main

import (
	"bufio"
	"github.com/tryfix/log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

func main() {
	client := &http.Client{Timeout: 10 * time.Second}
	if len(os.Args) != 3 {
		panic(`incorrect arguments (eg: ./tester <join/leave/crash> <num_of_nodes>)`)
	}

	typ := os.Args[1]
	numOfNodes, err := strconv.ParseInt(os.Args[2], 10, 64)
	if err != nil {
		log.Fatal(err, `incorrect parameter for number of nodes (NaN)`)
	}

	file, err := os.Open(`started_nodes.txt`)
	if err != nil {
		log.Fatal(err, `opening file failed`)
	}
	defer file.Close()

	var nodes []string
	port := `:52520`
	scanner := bufio.NewScanner(file)
	counter := 0
	for scanner.Scan() {
		if counter == int(numOfNodes) {
			break
		}
		nodes = append(nodes, scanner.Text())
		counter++
	}

	wg := &sync.WaitGroup{}
	if typ == `join` {
		log.Debug(`starting the test with nodes for join test`, len(nodes))
		startTime := time.Now().UTC()
		for i, node := range nodes {
			if i == 0 {
				continue
			}

			wg.Add(1)
			go func(node string, wg *sync.WaitGroup) {
				_, err := client.Post(`http://` + node + port + `/join?nprime=` + nodes[0] + port, `text/plain`, nil)
				if err != nil {
					log.Error(err, `http://` + node + port + `/join?nprime=` + nodes[0] + port)
				}
				wg.Done()
			}(node, wg)
		}

		wg.Wait()
		log.Debug(`all requests sent within (ms)`, time.Since(startTime).Milliseconds())
	} else if typ == `crash` {
		log.Debug(`starting the test with nodes for crash test`, len(nodes))
		startTime := time.Now().UTC()
		for _, node := range nodes {
			wg.Add(1)
			go func(node string, wg *sync.WaitGroup) {
				_, err := client.Post(`http://` + node + port + `/sim-crash`, `text/plain`, nil)
				if err != nil {
					log.Error(err,`http://` + node + port + `/sim-crash`)
				}
				wg.Done()
			}(node, wg)
		}

		wg.Wait()
		log.Debug(`all requests sent within (ms)`, time.Since(startTime).Milliseconds())
	} else {
		log.Error(`invalid type or not yet implemented`)
	}
}


