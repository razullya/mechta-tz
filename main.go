package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
)

type Sum struct {
	A int `json:"a"`
	B int `json:"b"`
}

func main() {
	filename := flag.String("file", "file.json", "filname")
	numWorkers := flag.Int("workers", 10, "nums of workers")
	flag.Parse()

	file, err := os.Open(*filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	_, err = decoder.Token()
	if err != nil {
		log.Fatal(err)
	}

	results := make(chan int, *numWorkers)
	sums := make(chan Sum, *numWorkers)

	var wg sync.WaitGroup

	for i := 0; i < *numWorkers; i++ {
		wg.Add(1)
		go work(sums, results, &wg)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	go func() {
		for decoder.More() {
			var sum Sum
			if err := decoder.Decode(&sum); err != nil {
				log.Println(err)
				continue
			}
			sums <- sum
		}
		close(sums)
	}()

	result := 0
	for r := range results {
		result += r
	}
	fmt.Println(result)
}

func work(sums <-chan Sum, resulst chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for sum := range sums {
		s := sum.A + sum.B
		resulst <- s
	}
}
