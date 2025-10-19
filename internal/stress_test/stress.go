package stresstest

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type StressTest struct {
	Url                       string
	Requests                  int
	Concurrency               int
	responseStatusCodeChannel chan responseStatusCode
	report                    *StressTestReport
}

type responseStatusCode struct {
	NameLoop   string
	StatusCode int
}

func NewStressTest(url string, requests int, concurrency int) *StressTest {
	return &StressTest{
		Url:         url,
		Requests:    requests,
		Concurrency: concurrency,
	}
}

func (st *StressTest) Run(ctx context.Context) (StressTestReport, error) {

	st.responseStatusCodeChannel = make(chan responseStatusCode, 10)
	st.report = &StressTestReport{
		TotalResquest: 0,
		StatusCodes:   map[int]int{},
	}

	var totalLastLoop int = 0

	totalPerLoops := st.Requests / st.Concurrency
	loopExtra := st.Requests % st.Concurrency

	if loopExtra > 0 {
		totalPerLoops++
		totalLastLoop = st.Requests - (totalPerLoops * (st.Concurrency - 1))
	} else {
		totalLastLoop = totalPerLoops
	}

	startTime := time.Now()

	var wgLoop sync.WaitGroup
	var wgResp sync.WaitGroup

	go st.registerStatusCode(&wgResp)

	for i := 0; i < (st.Concurrency - 1); i++ {
		wgLoop.Add(1)
		go st.callResquests(ctx, &wgLoop, &wgResp, strconv.Itoa(i+1), totalPerLoops)
	}

	if loopExtra > 0 {
		wgLoop.Add(1)
		go st.callResquests(ctx, &wgLoop, &wgResp, strconv.Itoa(st.Concurrency), totalLastLoop)
	}

	wgLoop.Wait()
	wgResp.Wait()

	elapsed := time.Since(startTime)

	return StressTestReport{
		TotalTime:     elapsed,
		TotalResquest: st.report.TotalResquest,
		StatusCodes:   st.report.StatusCodes,
	}, nil
}

func (st *StressTest) registerStatusCode(wgResp *sync.WaitGroup) {
	defer close(st.responseStatusCodeChannel)

	for {
		resp, ok := <-st.responseStatusCodeChannel
		if !ok {
			return
		}

		fmt.Printf("L[%v] Status Code: %v\n", resp.NameLoop, resp.StatusCode)

		st.report.StatusCodes[resp.StatusCode]++
		st.report.TotalResquest++

		wgResp.Done()
	}
}

func (st *StressTest) callResquests(ctx context.Context, wgLoop *sync.WaitGroup, wgResp *sync.WaitGroup, nameLoop string, numRequest int) {
	for i := 0; i < numRequest; i++ {
		status, err := st.callRequest(ctx)

		wgResp.Add(1)

		if err != nil {
			st.responseStatusCodeChannel <- responseStatusCode{
				NameLoop:   nameLoop,
				StatusCode: 0,
			}
			continue
		}

		st.responseStatusCodeChannel <- responseStatusCode{
			NameLoop:   nameLoop,
			StatusCode: status,
		}
	}
	wgLoop.Done()
}

func (st *StressTest) callRequest(ctx context.Context) (int, error) {

	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		st.Url,
		nil)

	if err != nil {
		return 0, err
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Do(req)

	if err != nil {
		return 0, err
	}

	return resp.StatusCode, nil
}
