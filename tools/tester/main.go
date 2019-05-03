package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"time"

	"github.com/gnboorse/binpacking"
)

func main() {
	// main entrypoint for running bin packing problems
	inputFile := flag.String("file", "input.json", "File to run.")
	outputFile := flag.String("output", "output.json", "Output file for results.")
	flag.Parse()
	b, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		panic(err)
	}

	var packingList binpacking.PackingList
	err = json.Unmarshal(b, &packingList)
	if err != nil {
		panic(err)
	}

	problem := binpacking.NewBinCollection(&packingList)

	start := time.Now()
	// time how long it takes to pack
	problem.PackAll(packingList.Items)
	elapsed := time.Since(start)
	// set duration of run
	problem.SetTime(elapsed.Nanoseconds())

	jsonValue, err := json.MarshalIndent(problem, "", "  ")
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(*outputFile, jsonValue, 0644)
	if err != nil {
		panic(err)
	}
}
