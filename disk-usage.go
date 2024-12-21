package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// dirSize calculates the total size of a directory and its subdirectories.
func dirSize(path string) (int64, error) {
	var totalSize int64

	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Skip directories but count their contents
		if !info.IsDir() {
			totalSize += info.Size()
		}
		return nil
	})

	return totalSize, err
}

func main() {
	root := "/mnt/c/Users/indra/Downloads/" // Replace with your root directory

	var wg sync.WaitGroup
	results := make(chan string, 10) // Buffered channel to store results

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Error accessing path %s: %v\n", path, err)
			return err
		}

		if info.IsDir() {
			// Increment WaitGroup counter
			wg.Add(1)

			// Launch a goroutine to calculate directory size
			go func(dir string) {
				defer wg.Done() // Decrement counter when goroutine completes
				size, err := dirSize(dir)
				if err != nil {
					results <- fmt.Sprintf("Error calculating size for %s: %v", dir, err)
				} else {
					results <- fmt.Sprintf("Directory: %s, Size: %d bytes", dir, size)
				}
			}(path)
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error walking the path %s: %v\n", root, err)
	}

	// Close the results channel when all goroutines complete
	go func() {
		wg.Wait()
		close(results)
	}()

	// Print results as they come in
	for result := range results {
		fmt.Println(result)
	}
}
