package main

import (
	"context"
	"flag"
	"fmt"
	stresstest "fullcycle-desafio-tecnico-2/internal/stress_test"
)

func main() {
	// variables declaration
	var argUrl string
	var argRequests int
	var argConcurrency int

	// flags declaration using flag package
	flag.StringVar(&argUrl, "url", "", "URL do serviço a ser testado.")
	flag.IntVar(&argRequests, "requests", 0, "Número total de requests")
	flag.IntVar(&argConcurrency, "concurrency", 0, "Número de chamadas simultâneas")

	flag.Usage = func() {
		fmt.Printf("Argumentos obrigatórios: \n")
		fmt.Printf("--url: 		URL do serviço a ser testado.\n")
		fmt.Printf("--requests: 	Número total de requests.\n")
		fmt.Printf("--concurrency: 	Número de chamadas simultâneas.\n")
	}

	flag.Parse() // after declaring flags we need to call it

	if argUrl == "" || argRequests == 0 || argConcurrency == 0 {
		flag.Usage()
		return
	}

	//st := stresstest.NewStressTest("https://www.google.com", 10, 1)
	st := stresstest.NewStressTest(argUrl, argRequests, argConcurrency)

	report, err := st.Run(context.Background())

	if err != nil {
		println("Error")
	}

	report.Print()
}
