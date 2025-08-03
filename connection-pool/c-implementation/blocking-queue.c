#include <stdio.h>
#include <stdlib.h>
#include <pthread.h>
#include <unistd.h>

#define MAX_QUEUE_SIZE 2

typedef struct {
    int data[MAX_QUEUE_SIZE];
    int head;
    int tail;
    pthread_mutex_t mutex;
    pthread_cond_t not_empty;
    pthread_cond_t not_full;
} BlockingQueue;

void init_queue(BlockingQueue *queue) {
    queue->head = 0;
    queue->tail = 0;
    pthread_mutex_init(&queue->mutex, NULL);
    pthread_cond_init(&queue->not_empty, NULL);
    pthread_cond_init(&queue->not_full, NULL);
}

void enqueue(BlockingQueue *queue, int value) {
    printf("🟡 Producer trying to enqueue %d\n", value);
    pthread_mutex_lock(&queue->mutex);
    while (queue->tail == MAX_QUEUE_SIZE) {
        printf("🔴 Producer BLOCKED! Queue is full, waiting...\n");
        pthread_cond_wait(&queue->not_full, &queue->mutex);
        printf("🟢 Producer WOKE UP! Checking queue again...\n");
    }
    queue->data[queue->tail] = value;
    queue->tail = (queue->tail + 1) % MAX_QUEUE_SIZE;
    printf("✅ Producer enqueued %d\n", value);
    pthread_cond_signal(&queue->not_empty);
    pthread_mutex_unlock(&queue->mutex);
}

int dequeue(BlockingQueue *queue) {
    printf("🟡 Consumer trying to dequeue\n");
    pthread_mutex_lock(&queue->mutex);
    while (queue->head == queue->tail) {
        printf("🔴 Consumer BLOCKED! Queue is empty, waiting...\n");
        pthread_cond_wait(&queue->not_empty, &queue->mutex);
        printf("🟢 Consumer WOKE UP! Checking queue again...\n");
    }
    int value = queue->data[queue->head];
    queue->head = (queue->head + 1) % MAX_QUEUE_SIZE;
    printf("✅ Consumer dequeued %d\n", value);
    pthread_cond_signal(&queue->not_full);
    pthread_mutex_unlock(&queue->mutex);
    return value;
}

void destroy_queue(BlockingQueue *queue) {
    pthread_mutex_destroy(&queue->mutex);
    pthread_cond_destroy(&queue->not_empty);
    pthread_cond_destroy(&queue->not_full);
}

// Producer function - renamed to avoid conflict
void* producer_func(void* arg) {
    BlockingQueue *queue = (BlockingQueue*)arg;
    
    for (int i = 0; i < 5; i++) {
        enqueue(queue, i + 100);  // Add 100 to make values unique
        usleep(500000); // Sleep 0.5 seconds between enqueues
    }
    printf("🏁 Producer finished!\n");
    return NULL;
}

// Consumer function - renamed to avoid conflict  
void* consumer_func(void* arg) {
    BlockingQueue *queue = (BlockingQueue*)arg;
    
    for (int i = 0; i < 5; i++) {
        int value = dequeue(queue);
        printf("📦 Final result: Consumer got %d\n", value);
        usleep(1000000); // Sleep 1 second between dequeues (slower than producer)
    }
    printf("🏁 Consumer finished!\n");
    return NULL;
}

int main() {
    printf("🚀 Starting Blocking Queue Demo!\n");
    printf("Queue size: %d\n\n", MAX_QUEUE_SIZE);
    
    BlockingQueue queue;
    init_queue(&queue);

    pthread_t producer_thread, consumer_thread;

    // Create producer and consumer threads
    pthread_create(&producer_thread, NULL, producer_func, &queue);
    pthread_create(&consumer_thread, NULL, consumer_func, &queue);
    
    // Wait for both threads to finish
    pthread_join(producer_thread, NULL);
    pthread_join(consumer_thread, NULL);

    printf("\n🎉 Demo completed!\n");
    destroy_queue(&queue);
    return 0;
}