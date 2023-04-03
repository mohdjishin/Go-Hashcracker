package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
)

func generateHash(text string, hashType string) string {
	var hash [32]byte
	if hashType == "sha256" {
		hash = sha256.Sum256(([]byte(text)))
	} else {
		panic("Unsupported hash type")
	}
	return hex.EncodeToString(hash[:])
}

func crackHash(hash string, wordlistPath string, maxWorkers int, hashType string) string {
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
				if generateHash(word, hashType) == hash {
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
	var maxWorkers *int

	hashType := flag.String("type", "sha256", "Hash type (sha256 only supported)")
	hashValue := flag.String("hash", "", "Hash value to crack")
	wordlistPath := flag.String("wordlist", "wordlist.txt", "Path to wordlist")
	maxWorkers = flag.Int("workers", runtime.NumCPU(), "Number of workers to use")
	flag.Parse()

	if *hashValue == "" {
		fmt.Println("Please provide a hash value to crack using the -hash flag")
		return
	}

	result := crackHash(strings.TrimSpace(*hashValue), *wordlistPath, *maxWorkers, *hashType)
	fmt.Printf("Hash: %s\nPassword: %s\n", *hashValue, result)
}

// useage:= go run main.go -type=sha256 -hash=5e884898da28047151d0e56f8dc6292773603d0d6aabbdd62a11ef721d1542d8 -wordlist=wordlist.txt -workers=10
