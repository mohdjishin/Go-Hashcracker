package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"runtime"
	"sync"
)

func generateHash(text string) string {
	hash := sha256.Sum256(([]byte(text)))
	return hex.EncodeToString(hash[:])
}

func crackHash(hash string, wordlistPath string, maxWorkers int) string {
	file, err := os.Open(wordlistPath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var wg sync.WaitGroup
	wg.Add(maxWorkers)

	jobs := make(chan string)

	for i := 0; i < maxWorkers; i++ {
		go func() {
			defer wg.Done()
			for word := range jobs {
				if generateHash(word) == hash {
					fmt.Printf("Password found: %s\n", word)
					return
				}
			}
		}()
	}

	for scanner.Scan() {
		jobs <- scanner.Text()
	}

	close(jobs)
	wg.Wait()

	return ""
}

func main() {
	hashToCrack := "5e884898da28047151d0e56f8dc6292773603d0d6aabbdd62a11ef721d1542d8"
	wordlistPath := "wordlist.txt"
	maxWorkers := runtime.NumCPU()

	result := crackHash(hashToCrack, wordlistPath, maxWorkers)
	fmt.Printf("Hash: %s\nPassword: %s\n", hashToCrack, result)
}
