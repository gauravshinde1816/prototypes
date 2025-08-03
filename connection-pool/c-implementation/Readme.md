## Blocking Queue Implementation:

## What `pthread_cond_wait(&queue->not_full, &queue->mutex)` does:

This line is the **heart of the blocking behavior**. Here's what happens:

### The Magic of `pthread_cond_wait`:
1. **🔓 Releases the mutex** - So other threads can access the queue
2. **😴 Puts thread to sleep** - Waiting for the `not_full` condition to be signaled  
3. **🔔 Wakes up when signaled** - Another thread calls `pthread_cond_signal(&not_full)`
4. **🔒 Re-acquires the mutex** - Before returning, so the thread can safely continue

### Why is this needed?

**The Problem:** When the queue is full (`tail == MAX_QUEUE_SIZE`), we can't add more items. But we don't want to just return an error - we want to **wait** until space becomes available.

**The Solution:** 
- **Release the lock** so other threads can dequeue items (making space)
- **Sleep** until someone signals that space is available
- **Wake up and try again**

### Simple analogy:
Think of it like a parking lot:
- 🚗 You want to park (enqueue) but lot is full
- 🔓 You step away from the entrance (release mutex) 
- 😴 You wait nearby for someone to leave
- 🔔 Security guard tells you "space available!" (signal)
- 🔒 You go back to entrance (re-acquire mutex)
- ✅ You park your car!

The `pthread_cond_wait` line ensures that enqueue **blocks** (waits) when the queue is full, rather than failing or spinning in a busy loop.