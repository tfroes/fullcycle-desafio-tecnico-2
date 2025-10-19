package stresstest

import (
	"fmt"
	"time"
)

type StressTestReport struct {
	TotalTime     time.Duration
	TotalResquest int
	StatusCodes   map[int]int
}

func (rep *StressTestReport) Print() {
	fmt.Printf("Tempo total gasto na execução: %s", rep.TotalTime.String())
	fmt.Println()
	fmt.Printf("Quantidade total de request realizados: %d", rep.TotalResquest)
	fmt.Println()
	fmt.Printf("Quantidade de requests com status HTTP 200: %d", rep.StatusCodes[200])
	fmt.Println()

	var errorTimeout int = 0

	for status, count := range rep.StatusCodes {

		if status == 200 {
			continue
		}

		if status == 0 {
			errorTimeout = count
			continue
		}

		fmt.Printf("Quantidade de requests com status HTTP %v: %v", status, count)
	}

	if errorTimeout > 0 {
		fmt.Printf("Quantidade de requests com timeout: %v", errorTimeout)
	}
}
