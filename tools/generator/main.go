package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/gnboorse/binpacking"
)

// entrypoint for generating bin packing problems
func main() {
	itemCount := flag.Int("count", 100, "The number of items to pack")
	itemMaxSize := flag.Int("max", 100, "The maximum size for a bin")
	itemCenter := flag.Int("center", 50, "Center of concentrated values")
	itemVariability := flag.Int("variability", int(binpacking.LowVariability), "Measure of the variability of values (should be 1, 2, or 3)")
	algorithm := flag.String("algorithm", "NextFit", "The name of the algorithm to use when solving the problem")
	duplicates := flag.Int("dups", 1, "How many of this kind of problem to generate")
	flag.Parse()
	for i := 0; i < *duplicates; i++ {
		// randomly generate items based on params provided
		items := binpacking.GenerateItems(*itemCount, *itemMaxSize, *itemCenter, binpacking.Variability(*itemVariability))

		// calculate lower bound for most optimal solution
		tmp := make(binpacking.Items, len(items))
		copy(tmp, items)
		lowerBound := binpacking.CalculateLowerBound(tmp, binpacking.Size(*itemMaxSize))

		// create packing list
		packingList := binpacking.PackingList{
			Size:        binpacking.Size(*itemMaxSize),
			Count:       binpacking.Count(*itemCount),
			Center:      *itemCenter,
			Variability: binpacking.Variability(*itemVariability),
			Algorithm:   binpacking.GetAlgorithm(*algorithm),
			Items:       items,
			LowerBound:  lowerBound}

		jsonValue, err := json.MarshalIndent(packingList, "", "  ")
		if err != nil {
			panic(err)
		}
		filename := fmt.Sprintf("binpacking%v_%vcount_%vmax_%vcenter_%vvariability_%s.json",
			i, *itemCount, *itemMaxSize, *itemCenter, *itemVariability, packingList.Algorithm)
		err = ioutil.WriteFile(filename, jsonValue, 0644)
		if err != nil {
			panic(err)
		}

	}
}
