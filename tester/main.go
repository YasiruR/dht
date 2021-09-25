package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

func main() {
	client := &http.Client{}
	if len(os.Args) != 4 {
		panic(`incorrect arguments (eg: ./tester <GET/PUT> <host:port> <num_of_requests>)`)
	}

	typ := os.Args[1]
	host := os.Args[2]
	numOfReqs, err := strconv.ParseInt(os.Args[3], 10, 64)
	if err != nil {
		panic(`incorrect parameter for number of requests (NaN)`)
	}

	// keys from bucket id=0 to 15
	url0 := `http://` + host + `/storage/aaa`
	url1 := `http://` + host + `/storage/assignment`
	url2 := `http://` + host + `/storage/999`
	url3 := `http://` + host + `/storage/123`
	url4 := `http://` + host + `/storage/1234`
	url5 := `http://` + host + `/storage/12345`
	url6 := `http://` + host + `/storage/9876`
	url7 := `http://` + host + `/storage/333`
	url8 := `http://` + host + `/storage/klm`
	url9 := `http://` + host + `/storage/abcd`
	url10 := `http://` + host + `/storage/54321`
	url11 := `http://` + host + `/storage/098`
	url12 := `http://` + host + `/storage/klmn`
	url13 := `http://` + host + `/storage/abc`
	url14 := `http://` + host + `/storage/111`
	url15 := `http://` + host + `/storage/dht`

	setReq0, _ := http.NewRequest(http.MethodPut, url0, nil)
	setReq1, _ := http.NewRequest(http.MethodPut, url1, nil)
	setReq2, _ := http.NewRequest(http.MethodPut, url2, nil)
	setReq3, _ := http.NewRequest(http.MethodPut, url3, nil)
	setReq4, _ := http.NewRequest(http.MethodPut, url4, nil)
	setReq5, _ := http.NewRequest(http.MethodPut, url5, nil)
	setReq6, _ := http.NewRequest(http.MethodPut, url6, nil)
	setReq7, _ := http.NewRequest(http.MethodPut, url7, nil)
	setReq8, _ := http.NewRequest(http.MethodPut, url8, nil)
	setReq9, _ := http.NewRequest(http.MethodPut, url9, nil)
	setReq10, _ := http.NewRequest(http.MethodPut, url10, nil)
	setReq11, _ := http.NewRequest(http.MethodPut, url11, nil)
	setReq12, _ := http.NewRequest(http.MethodPut, url12, nil)
	setReq13, _ := http.NewRequest(http.MethodPut, url13, nil)
	setReq14, _ := http.NewRequest(http.MethodPut, url14, nil)
	setReq15, _ := http.NewRequest(http.MethodPut, url15, nil)

	setReqs := []*http.Request{setReq0, setReq1, setReq2, setReq3, setReq4, setReq5, setReq6, setReq7, setReq8, setReq9, setReq10, setReq11, setReq12, setReq13, setReq14, setReq15}

	getReq0, _ := http.NewRequest(http.MethodGet, url0, nil)
	getReq1, _ := http.NewRequest(http.MethodGet, url1, nil)
	getReq2, _ := http.NewRequest(http.MethodGet, url2, nil)
	getReq3, _ := http.NewRequest(http.MethodGet, url3, nil)
	getReq4, _ := http.NewRequest(http.MethodGet, url4, nil)
	getReq5, _ := http.NewRequest(http.MethodGet, url5, nil)
	getReq6, _ := http.NewRequest(http.MethodGet, url6, nil)
	getReq7, _ := http.NewRequest(http.MethodGet, url7, nil)
	getReq8, _ := http.NewRequest(http.MethodGet, url8, nil)
	getReq9, _ := http.NewRequest(http.MethodGet, url9, nil)
	getReq10, _ := http.NewRequest(http.MethodGet, url10, nil)
	getReq11, _ := http.NewRequest(http.MethodGet, url11, nil)
	getReq12, _ := http.NewRequest(http.MethodGet, url12, nil)
	getReq13, _ := http.NewRequest(http.MethodGet, url13, nil)
	getReq14, _ := http.NewRequest(http.MethodGet, url14, nil)
	getReq15, _ := http.NewRequest(http.MethodGet, url15, nil)

	getReqs := []*http.Request{getReq0, getReq1, getReq2, getReq3, getReq4, getReq5, getReq6, getReq7, getReq8, getReq9, getReq10, getReq11, getReq12, getReq13, getReq14, getReq15}

	var counter int64
	wg := &sync.WaitGroup{}
	var reqs []*http.Request
	if typ == `GET` {
		reqs = getReqs
	} else if typ == `PUT` {
		reqs = setReqs
	} else {
		panic(`incorrect method type (eg: ./tester <GET/PUT> <host:port> <num_of_requests>)`)
	}

	startTime := time.Now()
	getReqLoop:
		for _, req := range reqs {
			if counter == numOfReqs {
				break
			}
			wg.Add(1)
			go func(req *http.Request, wg *sync.WaitGroup) {
				_, _ = client.Do(req)
				wg.Done()
			}(req, wg)
			counter++
		}

		if counter < numOfReqs {
			goto getReqLoop
		}

	wg.Wait()
	elapsedTime := time.Since(startTime).Seconds()
	fmt.Println(`Type of request: `, typ)
	fmt.Println(`Number of requests: `, numOfReqs)
	fmt.Println(`Total time (s): `, elapsedTime)
	fmt.Println(`Average time per request (s): `, elapsedTime/float64(numOfReqs))
}
