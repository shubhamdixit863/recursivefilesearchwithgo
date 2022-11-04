package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
)

func processed(fileName string, processedDirectories []string) bool {
	for i := 0; i < len(processedDirectories); i++ {
		if processedDirectories[i] != fileName {
			continue
		}
		return true
	}
	return false
}

func ListAllTheFileInDirectory(path string, dirs []string, channel chan string) {

	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)

	}
	files, err := f.Readdir(0)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		var newPath string
		if path != "/" {
			newPath = fmt.Sprintf("%s/%s", path, f.Name())
		} else {
			newPath = fmt.Sprintf("%s%s", path, f.Name())
		}

		if f.IsDir() {
			if !processed(newPath, dirs) {
				dirs = append(dirs, newPath)
				ListAllTheFileInDirectory(newPath, dirs, channel)

			}
		} else {
			// Passing the string path to channel to be received at the other routine
			channel <- newPath

		}
	}
}

// This method searches for the given string in the file -------------

func SearchForContentInfile(filePath string, content string) {

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), content) {
			fmt.Println(fmt.Sprintf("%s contain : %s", filePath, content))
		}

	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

}

func main() {
	// Wait group such that main routine can wait for the other routines to be finished
	var wg sync.WaitGroup

	channel := make(chan string)
	wg.Add(1)
	go func() {

		defer wg.Done()

		// Passing the folder ,in which you want to search for the file with given string
		// Make sure the path is correct
		ListAllTheFileInDirectory("demo", []string{}, channel)
		close(channel)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for elem := range channel {
			SearchForContentInfile(elem, "hello")
		}

	}()

	wg.Wait()

}
