
# 📝 File-Processor

**One-Line Description:**

> High-performance concurrent file processor in Go with live metrics, autoscaling workers, and SHA256 hashing.

**Developed by:**
Syed Shaheer Hussain © 2026

**Technologies / Language:**

* Go (Golang)
* Standard Library: `context`, `crypto/sha256`, `sync`, `runtime`, `os`, `filepath`

**Tags:**
#GoLang #Concurrency #SystemsProgramming #WorkerPool #Autoscaler #FileProcessing #SHA256 #Metrics

# 📖 Introduction

FileProcessor is a **production-ready Go application** that scans directories, calculates SHA256 hashes of files, and processes them concurrently using a **dynamic worker pool**.

It includes:

* **Live metrics reporting** (processed files, failed files, queue length, goroutines, memory usage)
* **Worker autoscaling** (adds workers automatically when the queue grows)
* **Graceful shutdown** via Ctrl+C or system signals
* **Error handling** and atomic counters for concurrency safety

This project demonstrates **real-world systems programming concepts** in Go.

# 🛠 What This Project Does

* Scans directories recursively for files
* Hashes files using SHA256
* Processes multiple files concurrently with **worker goroutines**
* Reports live metrics every second
* Dynamically adds workers if backlog grows
* Gracefully shuts down on **Ctrl+C** or termination signals

# 🧩 Architecture / Flow

**Flowchart / Process Flow:**

```
          +-----------------+
          |  Main Program   |
          +-----------------+
                   |
                   v
         +--------------------+
         | Walk Directory     |
         | Collect File Paths |
         +--------------------+
                   |
                   v
        +---------------------+
        | Jobs Channel (Chan) |
        +---------------------+
        /          |          \
       /           |           \
      v            v            v
+---------+   +---------+   +---------+
| Worker  |   | Worker  |   | Worker  |
| Goroutine|  | Goroutine|  | Goroutine|
+---------+   +---------+   +---------+
      \           |           /
       \          |          /
        v         v         v
   +-------------------------+
   | Process Files (SHA256)  |
   +-------------------------+
                   |
                   v
          +----------------+
          | Metrics Reporter|
          | Memory Usage    |
          | Queue Length    |
          | Goroutines      |
          +----------------+
                   |
                   v
          +----------------+
          | Autoscaler     |
          | Add/Remove     |
          | Workers        |
          +----------------+

```

# 🏗 Folder Structure

```
FileProcessor/
├── main.go               # Main application
├── go.mod                # Go modules file

```

# 💡 Features

* ✅ Concurrency with **worker pools**
* ✅ **Dynamic autoscaling** of workers
* ✅ SHA256 hashing of files
* ✅ Live **metrics reporting** (processed, failed, queue, goroutines, memory)
* ✅ Graceful shutdown with **Ctrl+C**
* ✅ Atomic counters for safe concurrent updates
* ✅ Error logging and collection

# ⚙ Functions Overview

| Function             | Purpose                                                                         |
| -------------------- | ------------------------------------------------------------------------------- |
| `main()`             | Initializes workers, metrics reporter, autoscaler, walks directories            |
| `worker()`           | Processes jobs from the channel, computes SHA256, updates metrics               |
| `processFile()`      | Opens file, computes SHA256, simulates processing delay                         |
| `metricsReporter()`  | Prints live metrics every second (processed, failed, queue, goroutines, memory) |
| `workerAutoscaler()` | Dynamically adds workers if backlog grows, tracks logical reduction             |

# 💾 Installation / Setup

**Requirements:**

* Go >= 1.25
* Windows, Linux, or macOS

**Steps:**

1. Clone repo (or create project folder):

```bash
git clone https://github.com/SyedShaheerHussain/Concurrency-FileProcessor-GO-lang

```
```
cd FileProcessor

```

2. Initialize Go modules (if not done):

```bash
go mod init fileprocessor

```

3. Build or run:

```bash
go run main.go -dir=C:\Users\YourUser\Documents -workers=4

```

**Optional build:**

```bash
go build -o fileprocessor main.go
./fileprocessor -dir=C:\Users\YourUser\Documents -workers=4

```

# 🏃 How to Use

* Run with `-dir` to specify directory
* Run with `-workers` to specify initial number of workers
* Monitor metrics printed every second
* Press **Ctrl+C** to gracefully stop

**Example:**

```bash
go run main.go -dir=C:\Windows -workers=4

```
```
go run main.go -dir=C:\Windows -workers=6

```
# ⚠️ Cautions & Warnings

>[!caution]
> * The program reads **all files in the directory recursively** — do not point it to extremely large directories without enough RAM.
> * SHA256 hashing can be CPU-intensive for very large files.
> * Autoscaler increases workers dynamically — too many workers can overwhelm CPU.
> * Only **files** are processed, directories are skipped.

# ✅ Advantages

* Handles large directories concurrently
* Real-time metrics for observability
* Dynamic adjustment of workers for performance
* Graceful shutdown prevents resource leaks
* Cross-platform (Windows/Linux/macOS)

# ❌ Disadvantages

* High memory usage if queue size is huge
* Autoscaler currently **cannot reduce active workers forcibly**; idle workers exit naturally
* Metrics printing may slightly slow down very high-throughput processing
* No persistence of processed files metadata yet

# 🚀 Future Enhancements

* True **worker scaling down** (idle workers terminate automatically)
* Add **throughput stats** (files/sec)
* **Prometheus metrics** endpoint for external monitoring
* **Terminal dashboard UI**
* **Retry mechanism** for failed files
* **Distributed processing** with multiple machines
* **Configurable thresholds** for autoscaling

# ⚙ How It Works (Step-by-Step)

1. Program starts, parses flags `-dir` and `-workers`
2. **Context** and graceful shutdown signal are initialized
3. Jobs channel with buffer 100 is created
4. Initial workers (`worker()` goroutines) start
5. **Worker autoscaler** starts monitoring the queue
6. **Metrics reporter** prints live metrics every second
7. Directory is walked recursively; files are sent to jobs channel
8. Workers read jobs, compute SHA256, update metrics
9. Autoscaler adds workers if backlog grows
10. Ctrl+C triggers context cancellation
11. Workers and metrics reporter exit gracefully
12. Final summary is printed

# 🏷 Market Value / Use Cases

* File indexing and backup systems
* Antivirus or file integrity scanning
* Log aggregation or crawler pipelines
* Educational tool for **Go concurrency and systems programming**

# 🛠 Developed By

Syed Shaheer Hussain © 2026

# ⚖️ Disclaimer

>[!warning]
> * Use responsibly; scanning system directories may require admin permissions
> * Designed for learning, testing, and real-world file processing scenarios

# ⚡ What This Project Can Do

* Scan directories concurrently
* Compute SHA256 of files
* Autoscale workers
* Live metrics monitoring (including memory usage)
* Handle graceful shutdown

# ⚡ What This Project Cannot Do

* Process files beyond memory/disk constraints
* Scale workers across multiple machines (currently single-machine)
* Persist results to database (requires extension)

# 📦 Summary of Current Features

| Feature            | Status |
| ------------------ | ------ |
| Concurrency        | ✅      |
| Worker Autoscaling | ✅      |
| Live Metrics       | ✅      |
| Memory Usage Stats | ✅      |
| SHA256 hashing     | ✅      |
| Graceful Shutdown  | ✅      |
| Error Logging      | ✅      |

# 💻 Languages

* Go (Golang) 1.25+

# 📈 Pros

* Highly concurrent
* Dynamic scaling
* Real-time observability
* Cross-platform

# 📉 Cons

* CPU & memory usage grows with large directories
* Autoscaler downscaling is **conceptual only**

# 🎯 When to Use

* Large file directories
* Systems programming practice in Go
* Learning concurrency patterns, worker pools, atomic counters, context handling

# 🧪 Notes

>[!note]
> * Sleep time in `processFile()` simulates CPU-bound work (50ms default)
> * Metrics are printed every 1 second
> * Autoscaler ticks every 2 seconds

# 🏗 How You Made This

* Designed worker pool with channel for job distribution
* Added context cancellation for graceful shutdown
* Used `atomic` counters for processed/failed files
* Added metrics reporter for live stats and memory usage
* Added worker autoscaler for dynamic concurrency

# ⚡ Step-By-Step Installation

1. Install Go >= 1.25
2. Clone repository
3. Open terminal and navigate to project folder
4. Run: `go run main.go -dir=<directory> -workers=<number>`
5. Observe metrics and logs in real-time
6. Press **Ctrl+C** to gracefully stop

# 📂 Folder / File Structure

```
FileProcessor/
├── main.go
├── go.mod

```

# 📝 Summary

FileProcessor is a **real-time, concurrent, scalable file processor** written in Go. It’s suitable for **systems programming, educational purposes, and real-world concurrent file processing**.

It demonstrates **worker pools, context cancellation, atomic counters, live metrics, autoscaling, SHA256 hashing, and graceful shutdown** in a single project.

## ⭐ Support & Engagement

If you find this repository useful or insightful, please consider:

- ⭐ Starring the repository
- 🔁 Sharing it within your network
- 👤 Following my GitHub profile for future projects and updates

Your support helps drive continued innovation and open-source contributions.

— Syed Shaheer Hussain

[![GitHub followers](https://img.shields.io/github/followers/SyedShaheerHussain?label=Follow&style=social)](https://github.com/SyedShaheerHussain)

![Followers](https://img.shields.io/github/followers/SyedShaheerHussain?label=Followers&color=blue)

![Stars](https://img.shields.io/github/stars/SyedShaheerHussain/Concurrency-FileProcessor-GO-lang?label=Stars&color=yellow)
