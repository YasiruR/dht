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
	if len(os.Args) != 4 {
		panic(`incorrect arguments (eg: ./tester <join/leave/crash> <num_of_nodes> <ticker_interval>)`)
	}

	typ := os.Args[1]
	numOfNodes, err := strconv.ParseInt(os.Args[2], 10, 64)
	if err != nil {
		log.Fatal(err, `incorrect parameter for number of nodes (NaN)`)
	}

	tickerInterval, err := strconv.ParseInt(os.Args[3], 10, 64)
	if err != nil {
		log.Fatal(err, `incorrect parameter for ticker intercal (NaN)`)
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
		startTime := time.Now().UTC()
		log.Debug(`starting the test with nodes`, len(nodes), startTime)
		for i, node := range nodes {
			if i == 0 {
				continue
			}

			wg.Add(1)
			go func(node string, wg *sync.WaitGroup) {
				_, err := client.Post(`http://` + node + port + `/join?nprime=` + nodes[0] + port, `application/json`, nil)
				if err != nil {
					log.Error(err, `http://` + node + port + `/join?nprime=` + nodes[0] + port)
				}
				wg.Done()
			}(node, wg)
		}

		wg.Wait()
		log.Debug(`all requests sent within (ms)`, time.Since(startTime).Milliseconds())
		time.Sleep(time.Duration(tickerInterval) * time.Millisecond)
		log.Debug(`after sleep`, tickerInterval)

		//ticker := time.NewTicker(time.Duration(tickerInterval) * time.Millisecond)
		//done := make(chan bool)
		//resChan := make(chan time.Time)
		//wg := &sync.WaitGroup{}
		//go func(numOfNodes int, wg *sync.WaitGroup) {
		//	for {
		//		select {
		//		case <- done:
		//			return
		//		case <- ticker.C:
		//			fmt.Println("REQQQQQQ")
		//			wg.Add(1)
		//			go func(numOfNodes int, wg *sync.WaitGroup) {
		//				reqStartTime := time.Now().UTC()
		//				res, err := client.Get(`http://` + nodes[1] + port + `/cluster-info`)
		//				if err != nil {
		//					log.Error(err)
		//				}
		//
		//				data, err := ioutil.ReadAll(res.Body)
		//				if err != nil {
		//					log.Error(err)
		//				}
		//
		//				if string(data) == strconv.Itoa(numOfNodes) {
		//					resChan <- reqStartTime
		//					done <- true
		//				}
		//				wg.Done()
		//			}(numOfNodes, wg)
		//		}
		//	}
		//}(int(numOfNodes), wg)
		//
		//reqStartTime := <- resChan
		//fmt.Println(`time to settle (ms): `, time.Since(startTime).Milliseconds())
		//fmt.Println(`time for cluster-info request (ms): `, time.Since(reqStartTime).Milliseconds())
		//wg.Wait()
	}
}


