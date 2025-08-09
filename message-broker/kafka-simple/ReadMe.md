# Kafka POC in Go – Learning Notes

This document summarizes the key concepts, issues, and fixes from building a **simple Producer & Consumer in Go** with Apache Kafka. It’s based on real debugging steps and lessons learned.

---

## 1. What is Kafka?

Apache Kafka is a **distributed, fault-tolerant event streaming platform**.  
Think of it as a durable, append-only log that producers write to and consumers read from.

**Core concepts:**
- **Topic** – A named stream of events.
- **Partition** – A subset of a topic. Ordering is guaranteed **within** a partition.
- **Offset** – A sequential ID for each message in a partition.
- **Producer** – Publishes messages to topics.
- **Consumer** – Reads messages from topics.
- **Broker** – Kafka server storing partitions.

---

## 2. Why Kafka is Popular

- **High throughput & scalability** – Millions of messages/sec.
- **Durable storage** – Messages replicated across brokers.
- **Decoupling** – Producers don’t know about consumers.
- **Replay capability** – Consumers can re-read messages by offset.
- **Rich ecosystem** – Kafka Connect, Streams, ksqlDB.



## 3. How Kafka Works

Below is a simplified diagram of the Producer → Kafka → Consumer flow with partitions.

```

```
    +-----------+            +---------------------+            +-----------+
    |  Producer |   writes   |     Kafka Broker     |   reads    | Consumer  |
    |  (Go app) +----------->|  my-topic            +----------->|  (Go app) |
    +-----------+            |                     |            +-----------+
                              |  Partition 0        |
                              |  0: msg             |
                              |  1: msg             |
                              |                     |
                              |  Partition 1        |
                              |  0: msg             |
                              |  1: msg             |
                              |                     |
                              |  Partition 2        |
                              |  0: msg             |
                              |  1: msg             |
                              +---------------------+
```

```

- **Producers** send messages to a topic. Kafka decides which partition to store them in (unless you target explicitly).
- **Partitions** hold ordered, immutable sequences of messages.
- **Consumers** can read from one or more partitions and keep track of their offset.

---

## 4. Common Issues & Fixes

### **A. DNS / Connection Issues**
**Symptom:** `lookup localhost: i/o timeout`  
**Cause:** Broker advertises `localhost`, client is outside container.  
**Fix:** Set:
```

KAFKA\_CFG\_ADVERTISED\_LISTENERS=PLAINTEXT://127.0.0.1:9092

```

---

### **B. Duplicate Messages**
- Kafka is **at-least-once** → retries can cause duplicates.
- Reading from earliest offset replays all history.

**Fixes:**
- Set `MaxAttempts: 1` in producer for dev.
- Use `StartOffset: kafka.LastOffset` or a `GroupID` to read only new data.
- Delete/recreate topic to clear history.

---

### **C. Messages on Wrong Partitions**
- `kafka.Writer` ignores `Message.Partition` field.
- Use a **balancer** (e.g., `Hash`) with keys or `DialLeader` for explicit partition targeting.

---

### **D. Out-of-Order Messages**
- Order is guaranteed only **within a partition**.
- Balancer can reorder per partition in mixed batches.
- To preserve order, send sequential writes per partition or per key.

---

## 5. Clearing Existing Data

**Option 1 – Delete & Recreate Topic**
```

kafka-topics.sh --delete --topic my-topic
kafka-topics.sh --create --topic my-topic --partitions 3 --replication-factor 1

```

**Option 2 – Purge via Retention**
```

kafka-configs.sh --alter --topic my-topic --add-config retention.ms=0

```
(Wait a few seconds, then restore retention.)

---

## 6. Key Takeaways

1. **Advertised listeners** must match client address.
2. **StartOffset** or **GroupID** prevents replay of old messages.
3. `Writer` ignores `Message.Partition`; use `DialLeader` for explicit targeting.
4. Kafka is **at-least-once** by default → duplicates are possible.
5. Ordering is only guaranteed **within a partition**.
6. Old messages stay until retention deletes them or the topic is deleted.

---

## 7. References
- [Kafka Official Docs](https://kafka.apache.org/documentation/)
- [kafka-go Documentation](https://pkg.go.dev/github.com/segmentio/kafka-go)
- [Bitnami Kafka Docker Image](https://hub.docker.com/r/bitnami/kafka)
```

