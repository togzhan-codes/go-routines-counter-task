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
	"time"
)

const filename = "sample_file_10k.json"

// make sure this file is generated or existing in your directory first
// const filename = "sample_file_1m.json"

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
	for {
		fmt.Fprint(os.Stderr, "How many go routines you want to use? this is mandatory value to provide, between 1 and 1,000,000\n")
		s, _ = r.ReadString('\n')
		if s != "" {
			break
		}
	}
	for {
		fmt.Fprint(os.Stderr, "Would you like to regenerate json file with numbers? proivde 'y' or 'n' answer. If no new file needed, the existing sample file will be used\n")
		regenerateFile, _ = r.ReadString('\n')
		if s != "" {
			break
		}
	}
	workers, err := strconv.Atoi(strings.TrimSpace(s))
	if err != nil {
		return workers, err
	}
	if strings.TrimSpace(regenerateFile) == "y" {
		err = generateJsonFile(-10, 10, 10000)
		if err != nil {
			return workers, err
		}
		log.Print("you asked to generate the file, so you need to wait a bit after generation...like 30s of sleep:)")
		time.Sleep(30 * time.Second)
	}
	return workers, nil
}

func main() {
	// this is code without any complications with go routines, uncomment it and see the result
	// or just set num of go routines to 1...
	// start := time.Now()
	// plan, _ := os.ReadFile(filename)
	// var data []obj
	// err := json.Unmarshal(plan, &data)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// var sum int
	// for _, element := range data {
	// 	sum += element.A + element.B
	// }
	// fmt.Println("\nThe sum is ", sum)
	// timeElapsed := time.Since(start)
	// fmt.Printf("The `for` loop took %s", timeElapsed)

	workers, err := readArguments()
	if err != nil {
		log.Fatal(err)
	}
	start := time.Now()

	dataFile, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	var data []obj
	err = json.Unmarshal(dataFile, &data)
	if err != nil {
		log.Fatal(err)
	}
	if workers > len(data) {
		workers = len(data)
	}
	processRounds := int(math.Ceil(float64(len(data)) / float64(workers)))

	var eachRoundWgs = make([]sync.WaitGroup, processRounds)
	for i := range processRounds {
		eachRoundWgs[i].Add(workers)
	}

	var wg sync.WaitGroup
	wg.Add(workers)
	var boundary1, boundary2, finalResSum int
	for i := range workers {
		go func(j int) {
			defer wg.Done()
			boundary1 = j * processRounds
			boundary2 = j*processRounds + processRounds
			if j == workers-1 {
				boundary2 = len(data)
			}
			if boundary1 > len(data) {
				boundary1 = len(data)
			}
			finalResSum += work(processRounds, eachRoundWgs, data[boundary1:boundary2])
		}(i)

	}
	wg.Wait()
	fmt.Println("\nThe sum is ", finalResSum)
	timeElapsed := time.Since(start)
	fmt.Printf("\nsthe time taken is %s", timeElapsed)
}

func work(processRounds int, eachRoundWgs []sync.WaitGroup, objects []obj) int {
	var overallSum int
	for r := range processRounds {
		var sum int
		difference := r - len(objects)
		if len(objects) > 0 && difference < 0 {
			sum = objects[r].A + objects[r].B
		} else if len(objects) <= 0 || difference >= 0 {
			sum = 0
		}
		overallSum += sum
		eachRoundWgs[r].Done()
		eachRoundWgs[r].Wait()
	}
	return overallSum
}
