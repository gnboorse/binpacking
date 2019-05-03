package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"

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

	problem.PackAll(packingList.Items)

	jsonValue, err := json.MarshalIndent(problem, "", "  ")
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(*outputFile, jsonValue, 0644)
	if err != nil {
		panic(err)
	}
}
