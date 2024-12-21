// package main

// import (
// 	"fmt"
// 	"os"
// 	"path/filepath"
// )

// func iterateRoot(root string) map[string]int64 {
// 	dirs := make(map[string]int64)

// 	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
// 		if err != nil {
// 			return err
// 		}

// 		if !info.IsDir() {
// 			fileDir := filepath.Dir(path)
// 			dirs[fileDir] += info.Size()
// 		}

// 		return nil
// 	})

// 	return dirs
// }

// func propagateSizes(dirs map[string]int64) map[string]int64 {
// 	for dir := range dirs {
// 		parent := filepath.Dir(dir)
// 		if parent != dir { // Avoid propagating to root's parent
// 			dirs[parent] += dirs[dir]
// 		}
// 	}
// 	return dirs
// }

// func bytesToGB(bytes int64) float64 {
// 	return float64(bytes) / (1024 * 1024 * 1024)
// }

// func main() {
// 	root := "/mnt/c/Users/indra/Downloads/"
// 	res := iterateRoot(root)
// 	totalSizes := propagateSizes(res)

// 	fmt.Printf("Directory\t\t\tSize (GB):\n")
// 	for dir, size := range totalSizes {
// 		fmt.Printf("%s\t\t\t%.2f\n", dir, bytesToGB(size))
// 	}
// }

// package main

// import (
// 	"fmt"
// 	"os"
// 	"path/filepath"
// )

// // dirSize calculates the total size of a directory and its subdirectories.
// func dirSize(path string) (int64, error) {
// 	var totalSize int64

// 	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
// 		if err != nil {
// 			return err
// 		}
// 		// Skip directories but count their contents
// 		if !info.IsDir() {
// 			totalSize += info.Size()
// 		}
// 		return nil
// 	})

// 	return totalSize, err
// }

// func main() {
// 	root := "/mnt/c/Users/indra/Downloads/" // Replace with your root directory

// 	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
// 		if err != nil {
// 			fmt.Printf("Error accessing path %s: %v\n", path, err)
// 			return err
// 		}

// 		// Print size only for directories
// 		if info.IsDir() {
// 			size, err := dirSize(path)
// 			if err != nil {
// 				fmt.Printf("Error calculating size for %s: %v\n", path, err)
// 			} else {
// 				fmt.Printf("Directory: %s, Size: %d bytes\n", path, size)
// 			}
// 		}
// 		return nil
// 	})

// 	if err != nil {
// 		fmt.Printf("Error walking the path %s: %v\n", root, err)
// 	}
// }

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
