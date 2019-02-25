package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/gnboorse/binpacking"
)

func main() {
	// main entrypoint for running bin packing problems
	inputFile := flag.String("file", "input.json", "File to run.")
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

	fmt.Printf("Problem solution was: %s\n", problem)
}
