package pool

import (
	"log"
	"runtime"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/janithT/webpage-analyzer/analyzers"
)

// Run the analyzers job
// func ExecuteAnalyzers(anlz []analyzers.Analyzer, doc *goquery.Document, raw string) []analyzers.Result {
// 	numWorkers := runtime.NumCPU()
// 	log.Printf("Num of workers - %v ", numWorkers) // Number of concurrent goroutines

// 	jobs := make(chan int, len(anlz))
// 	results := make([]analyzers.Result, len(anlz))
// 	var wg sync.WaitGroup

// 	// Start worker goroutines
// 	for w := 0; w < numWorkers; w++ {
// 		wg.Add(1)
// 		go func() {
// 			for i := range jobs {
// 				results[i] = anlz[i].Analyze(doc, raw)
// 			}
// 			wg.Done()
// 		}()
// 	}

// 	// Send all job indices to jobs channel
// 	for i := range anlz {
// 		jobs <- i
// 	}
// 	close(jobs)

// 	wg.Wait()
// 	return results
// }

func ExecuteAnalyzers(anlz []analyzers.Analyzer, doc *goquery.Document, raw string) []analyzers.Result {
	numWorkers := runtime.NumCPU()
	log.Printf("Number of workers: %d", numWorkers)

	jobs := make(chan int, len(anlz))
	results := make([]analyzers.Result, 0, len(anlz)) // start empty slice with capacity
	var wg sync.WaitGroup
	var mu sync.Mutex

	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := range jobs {
				res := anlz[i].Analyze(doc, raw)

				// Lock before appending
				mu.Lock()
				results = append(results, res)
				mu.Unlock()
			}
		}()
	}

	for i := range anlz {
		jobs <- i
	}
	close(jobs)

	wg.Wait()
	return results
}
