package main

import (
	"fmt"
	"math"
	"sync"
	"sync/atomic"
	"time"
)

var totalPrimeNumbers int32 = 0
var currentNum int32 = 2
var MAX_INT int32 = 100000

var numThreads int = 10

func checkPrime(x int) {
	if x&1 == 0 {
		return
	}

	for i := 2; i <= int(math.Sqrt(float64(x))); i++ {
		if x%i == 0 {
			return
		}
	}
	atomic.AddInt32(&totalPrimeNumbers, 1)
}

func doWork(threadNum int, wg *sync.WaitGroup) {
	for {
		startTime := time.Now()
		atomic.AddInt32(&currentNum, 1)
		if currentNum > MAX_INT {
			break
		}
		checkPrime(int(currentNum))
		fmt.Println("Thread", threadNum, "is working on number", currentNum, "time taken: ", time.Since(startTime))
	}
	wg.Done()
}

func main() {
	startTime := time.Now()
	wg := sync.WaitGroup{}
	for i := 0; i < numThreads; i++ {
		wg.Add(1)
		go doWork(i, &wg)
	}

	wg.Wait()
	fmt.Println("Total prime numbers: ", totalPrimeNumbers+1)
	fmt.Println("Total time taken: ", time.Since(startTime))
}
