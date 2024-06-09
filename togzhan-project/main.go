package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
)

const filename = "sample_file_1m.json"

type obj struct {
	A int `json:"a"`
	B int `json:"b"`
}

func generateJsonFile(filename string, valMin, valMax, objNum int) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString("[")
	if err != nil {
		return err
	}

	encoder := json.NewEncoder(file)
	for i := 0; i < objNum; i++ {
		val1 := rand.Intn(valMax-valMin+1) + valMin
		val2 := rand.Intn(valMax-valMin+1) + valMin
		if err := encoder.Encode(obj{A: val1, B: val2}); err != nil {
			return err
		}
		if i < objNum-1 {
			_, err = file.WriteString(",")
			if err != nil {
				return err
			}
		}
	}

	_, err = file.WriteString("]")
	if err != nil {
		return err
	}

	return nil
}

func readArguments() (int, error) {
	var workersInput, regenerateFile string
	r := bufio.NewReader(os.Stdin)
	var workers int
	var err error
	for {
		fmt.Fprint(os.Stderr, "How many go routines you want to use? this is mandatory value to provide\n")
		workersInput, _ = r.ReadString('\n')
		workers, err = strconv.Atoi(strings.TrimSpace(workersInput))
		if err != nil {
			return workers, err
		}
		if workers > 0 {
			break
		}
		log.Print("invalid value")
	}
	for {
		fmt.Fprint(os.Stderr, "Would you like to regenerate json file with numbers? provide 'y' or 'n' answer. If no new file needed, the existing sample file will be used\n")
		regenerateFile, _ = r.ReadString('\n')
		regenerateFile = strings.TrimSpace(regenerateFile)
		regenerateFile = strings.ToLower(regenerateFile)
		if regenerateFile == "y" || regenerateFile == "n" {
			break
		}
		log.Print("invalid value")
	}

	//if user wants to have new file, file generator function is called with hardcoded 1,000,000 number of objects
	if regenerateFile == "y" {
		err := generateJsonFile(filename, -10, 10, 1000000)
		if err != nil {
			return workers, err
		}
	}
	return workers, nil
}

func main() {
	// reading arguments from console
	workers, err := readArguments()
	if err != nil {
		log.Fatal(err)
	}
	// reading bytes from the given file
	dataFileBytes, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	// unmarshaling bytes from file to the slice of objects of type obj
	var data []obj
	if err = json.Unmarshal(dataFileBytes, &data); err != nil {
		log.Fatal(err)
	}
	// getting number of process rounds for each goroutine/ worker, so that objects are taken in slices of length=processRounds
	processRounds := int(math.Ceil(float64(len(data)) / float64(workers)))
	// initializing wait group var to be used to wait for each goroutine to be finished indeed
	var wg sync.WaitGroup
	// creating channel for saving result of each goroutine task parallely with the buffer=number of goroutines
	results := make(chan int, workers)
	for i := 0; i < workers; i++ {
		// increasing waitGroup counter for each and every routine
		wg.Add(1)
		go func(j int) {
			// decrease wait group counter when and only when goroutine finished its job
			defer wg.Done()
			// creating limits for each worker, according to which we get slice from data,
			// e.g. if num of goroutines=10000, first worker will get data[0:100], second - data [100:200], etc.
			boundary1 := j * processRounds
			boundary2 := (j + 1) * processRounds
			// boundaries might get > than length of data, therefore introducing 2 if statements to avoid panic
			if boundary2 > len(data) {
				boundary2 = len(data)
			}
			if boundary1 > len(data) {
				boundary1 = len(data)
			}
			// finally each goroutine is counting sum of its data slice
			sum := countSum(data[boundary1:boundary2])
			// sending sum result of each routine to the channel
			results <- sum
		}(i)
	}

	go func() {
		// waitGroup waits for all workers to be done
		wg.Wait()
		// when they are done channel is closed and no more values are expected to be sent there
		close(results)
	}()

	var finalResSum int
	// counting final sum based on values kept in the result channel
	for sum := range results {
		finalResSum += sum
	}

	fmt.Println("Final Result Sum:", finalResSum)
}

func countSum(objects []obj) int {
	sum := 0
	for _, object := range objects {
		sum += object.A + object.B
	}
	return sum
}
