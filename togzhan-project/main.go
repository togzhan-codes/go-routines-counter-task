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

// make sure this file is generated or existing in your directory first
// to be able to upload it to github, file with 10,000 objects was created, for your test regenerate from console new file with 1m objects
// const filename = "sample_file_10k.json"

const filename = "sample_file_1m.json"

type obj struct {
	A int `json:"a"`
	B int `json:"b"`
}

func generateJsonFile(valMin, valMax, objNum int) error {
	var data []obj
	for i := 0; i < objNum; i++ {
		val1 := rand.Intn(valMax-valMin+1) + valMin
		val2 := rand.Intn(valMax-valMin+1) + valMin
		data = append(data, obj{
			A: val1,
			B: val2,
		})
	}

	b, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		return err
	}
	err = os.WriteFile(filename, b, 0644)
	if err != nil {
		return err
	}
	return nil
}

func readArguments() (int, error) {
	var s, regenerateFile string
	r := bufio.NewReader(os.Stdin)
	fmt.Fprint(os.Stderr, "How many go routines you want to use? this is mandatory value to provide\n")
	s, _ = r.ReadString('\n')
	fmt.Fprint(os.Stderr, "Would you like to regenerate json file with numbers? provide 'y' or 'n' answer. If no new file needed, the existing sample file will be used\n")
	regenerateFile, _ = r.ReadString('\n')
	workers, err := strconv.Atoi(strings.TrimSpace(s))
	if err != nil {
		return workers, err
	}
	//if user wants to have new file, file generator function is called with hardcoded 1,000,000 number of objects
	if strings.TrimSpace(regenerateFile) == "y" {
		err = generateJsonFile(-10, 10, 1000000)
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
	var data []obj
	// unmarshaling bytes from file to the slice of objects of type obj
	err = json.Unmarshal(dataFileBytes, &data)
	if err != nil {
		log.Fatal(err)
	}
	// getting number of process rounds for each goroutine/ worker, so that objects are taken in slices of length=processRounds
	processRounds := int(math.Ceil(float64(len(data)) / float64(workers)))
	// introducing wait groups with the counter = number of goroutines/workers, so that we wait for ALL routines to be finished
	var wg sync.WaitGroup
	wg.Add(workers)
	var finalResSum int
	var mu sync.Mutex
	for i := 0; i < workers; i++ {
		go func(j int) {
			// decrease wait group counter when and only when goroutine finished its job
			defer wg.Done()
			// creating limits for each worker, according to which we get slice from data,
			// e.g. if num of goroutines=10000, first worker will get data[0:100], second - data [100:200], etc.
			boundary1 := j * processRounds
			boundary2 := j*processRounds + processRounds
			// boundaries might get > than length of data, therefore introducing 2 if statements to avoid panic
			if boundary2 > len(data) {
				boundary2 = len(data)
			}
			if boundary1 > len(data) {
				boundary1 = len(data)
			}
			// finally each goroutine is counting sum of its data slice
			sum := countSum(data[boundary1:boundary2])
			// using mutex lock/unlock functions to be sure that only one goroutine can access finalResSum variable at a time
			mu.Lock()
			finalResSum += sum
			mu.Unlock()
		}(i)
	}
	wg.Wait()
	fmt.Println("\nThe final sum is ", finalResSum)
}

func countSum(objects []obj) int {
	var sum int
	for _, object := range objects {
		sum += object.A + object.B
	}
	return sum
}
