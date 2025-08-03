package main

import (
	"fmt"
	"sync"
	"time"
)

const (
	MAX_QUEUE_SIZE = 21
)

var wgE sync.WaitGroup

var wgD sync.WaitGroup

type BlockingQueue struct {
	queue          [MAX_QUEUE_SIZE]int
	tail           int
	head           int
	mutex          sync.Mutex
	spaceAvailable sync.Cond
	dataAvailable  sync.Cond
}

func (q *BlockingQueue) Init() {
	q.tail = 0
	q.head = 0
	q.spaceAvailable = sync.Cond{L: &q.mutex}
	q.dataAvailable = sync.Cond{L: &q.mutex}
}

func (q *BlockingQueue) Enqueue(value int) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	for (q.tail+1)%MAX_QUEUE_SIZE == q.head {
		fmt.Printf("Thread %d is waiting for space to be available\n", value)
		q.spaceAvailable.Wait()
	}
	fmt.Printf("🟡 Thread %d is enqueuing %d\n", value, value*100)
	q.queue[q.tail] = value * 100
	q.tail = (q.tail + 1) % MAX_QUEUE_SIZE
	q.dataAvailable.Signal()
}

func (q *BlockingQueue) Dequeue(value int) int {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	for q.head == q.tail {
		fmt.Printf("Thread %d is waiting for data to be available\n", value)
		q.dataAvailable.Wait()
	}

	fmt.Printf("🟢 Thread %d is dequeuing %d\n", value, q.queue[q.head])
	element := q.queue[q.head]
	q.head = (q.head + 1) % MAX_QUEUE_SIZE
	q.spaceAvailable.Signal()
	return element
}

func (q *BlockingQueue) Size() int {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	return (q.tail - q.head + MAX_QUEUE_SIZE) % MAX_QUEUE_SIZE
}

func main() {

	queue := BlockingQueue{}
	queue.Init()

	for i := range MAX_QUEUE_SIZE - 1 {
		queue.Enqueue(i)
	}

	for i := range 22 {
		wgE.Add(1)
		go func(i int) {
			queue.Enqueue(i)
			wgE.Done()
		}(i)
	}

	for i := range 22 {
		wgD.Add(1)
		go func(i int) {
			time.Sleep(1 * time.Second)
			queue.Dequeue(i)
			wgD.Done()
		}(i)
	}

	wgE.Wait()
	wgD.Wait()
	fmt.Println(queue.Size())
	fmt.Println("-------------------------------- remaining elements -------------------------------")
	for i := queue.head; i != queue.tail; i = (i + 1) % MAX_QUEUE_SIZE {
		fmt.Println(queue.queue[i])
	}

}
