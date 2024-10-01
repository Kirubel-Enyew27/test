Here’s a Go program that reads multiple `.txt` files concurrently from a given directory, processes the log data, and counts word occurrences in a thread-safe manner. It uses goroutines, worker pools, and a mutex for concurrency control. Additionally, it handles errors gracefully and includes basic error handling for file operations.

### Program Structure:
1. **Read and Process Files Concurrently**: Uses worker pools to limit concurrent file reading.
2. **Thread-Safe Word Counting**: A map protected by a `sync.Mutex` to ensure concurrent writes are safe.
3. **Error Handling**: Manages errors for file opening and processing.
4. **Case-Insensitive Word Counting**: Makes word comparison case-insensitive.
5. **Graceful Exit with Error Logging**: Logs any encountered errors and prints a summary after all work is complete.

### Code Implementation:

```go
package main

import (
    "bufio"
    "fmt"
    "io/ioutil"
    "log"
    "os"
    "path/filepath"
    "strings"
    "sync"
)

var (
    wordCount   = make(map[string]int)
    mutex       = &sync.Mutex{}
    maxWorkers  = 3
    errorLogger = log.New(os.Stderr, "ERROR: ", log.LstdFlags)
)

// Function to process a file and count word occurrences
func processFile(path string, wg *sync.WaitGroup) {
    defer wg.Done()

    file, err := os.Open(path)
    if err != nil {
        errorLogger.Printf("Failed to open file %s: %v\n", path, err)
        return
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    scanner.Split(bufio.ScanWords)

    for scanner.Scan() {
        word := strings.ToLower(scanner.Text()) // Case-insensitive count
        word = strings.Trim(word, `.,;:"'!?()[]`) // Remove punctuation

        mutex.Lock()
        wordCount[word]++
        mutex.Unlock()
    }

    if err := scanner.Err(); err != nil {
        errorLogger.Printf("Error reading file %s: %v\n", path, err)
    }
}

// Worker function that processes files from the jobs channel
func worker(jobs <-chan string, wg *sync.WaitGroup) {
    for path := range jobs {
        processFile(path, wg)
    }
}

// Function to read all .txt files from a directory and start worker pool
func readFilesConcurrently(dir string) {
    files, err := ioutil.ReadDir(dir)
    if err != nil {
        log.Fatalf("Failed to read directory %s: %v\n", dir, err)
    }

    var wg sync.WaitGroup
    jobs := make(chan string, len(files))

    // Start worker pool
    for i := 0; i < maxWorkers; i++ {
        go worker(jobs, &wg)
    }

    // Add files to job queue
    for _, file := range files {
        if filepath.Ext(file.Name()) == ".txt" {
            wg.Add(1)
            jobs <- filepath.Join(dir, file.Name())
        }
    }

    close(jobs)
    wg.Wait()
}

// Main function
func main() {
    if len(os.Args) < 2 {
        log.Fatal("Please provide a directory containing .txt files.")
    }

    dir := os.Args[1]
    readFilesConcurrently(dir)

    // Print word count summary
    fmt.Println("Word Count Summary:")
    for word, count := range wordCount {
        fmt.Printf("%s: %d\n", word, count)
    }
}
```

### Key Points:
- **Concurrency Control**: The `worker` function handles file processing, and the number of concurrent workers is controlled using the `maxWorkers` variable.
- **Mutex for Safe Access**: A `sync.Mutex` ensures that the shared `wordCount` map is accessed safely across multiple goroutines.
- **Case-Insensitive Word Counting**: Words are converted to lowercase and stripped of punctuation before counting.
- **Error Handling**: Errors encountered during file opening or reading are logged using an `errorLogger`.
- **Graceful Shutdown**: The program waits for all goroutines to complete using a `sync.WaitGroup`.

### Testable Components:
You can write unit tests for the following:
1. **File Reading Function** (`processFile`)
2. **Word Counting Logic** (checking case insensitivity and punctuation handling)
3. **Concurrency and Error Handling** (mocking file operations)

### Example Test:

```go
package main

import (
    "os"
    "testing"
)

func TestProcessFile(t *testing.T) {
    // Create a temp file for testing
    file, err := os.CreateTemp("", "testfile*.txt")
    if err != nil {
        t.Fatal(err)
    }
    defer os.Remove(file.Name())

    // Write some log data to the file
    content := "Go is great! Go is fast."
    file.WriteString(content)
    file.Close()

    var wg sync.WaitGroup
    wg.Add(1)

    // Process the file
    processFile(file.Name(), &wg)
    wg.Wait()

    // Validate word counts
    expectedCounts := map[string]int{
        "go":   2,
        "is":   2,
        "great": 1,
        "fast": 1,
    }

    for word, expectedCount := range expectedCounts {
        if count := wordCount[word]; count != expectedCount {
            t.Errorf("Expected %d occurrences of %q, got %d", expectedCount, word, count)
        }
    }
}
```

### Bonus:
- **Case-Insensitive Word Search**: The program already converts all words to lowercase for case-insensitive counting.
- **Handling Large Files**: The program reads files using `bufio.Scanner`, which is efficient for large files.

This program efficiently handles multiple `.txt` files with proper concurrency control and error handling while providing a summary of the word counts across all files.
