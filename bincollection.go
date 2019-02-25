package binpacking

import (
	"encoding/json"
	"sort"
)

// BinCollection an interface representing an instance
// of the bin packing problem.
type BinCollection interface {
	NewBin() *Bin
	PackItem(item Item)
	PackAll(items Items)
	GetBin(index int) *Bin
	GetFirstBin() *Bin
	GetLastBin() *Bin
	Find(predicate func(*Bin) bool) *Bin
	GetTotalBins() Count
	GetBinCapacity() Size
	String() string
}

// NewBinCollection create an instance of the bin packing problem
// from a PackingList object
func NewBinCollection(pList *PackingList) BinCollection {
	collection := &BinCollectionImpl{
		BinCapacity: pList.Size,
		TotalBins:   0,
		Bins:        make(Bins, 0, (pList.Count/2)+1), // pre-allocate memory for a reasonably large capacity
		Algorithm:   pList.Algorithm}

	collection.NewBin() // always create first bin

	return collection
}

// BinCollectionImpl default implementation of an
// instance of the bin packing problem
type BinCollectionImpl struct {
	BinCapacity Size  `json:"capacity"`
	TotalBins   Count `json:"count"`
	Bins        `json:"bins"`
	Algorithm   `json:"algorithm"`
}

// GetTotalBins getter for the total number of bins
func (binCollection *BinCollectionImpl) GetTotalBins() Count {
	return binCollection.TotalBins
}

// GetBinCapacity getter for the individual bin capacities
func (binCollection *BinCollectionImpl) GetBinCapacity() Size {
	return binCollection.BinCapacity
}

// NewBin method used for allocating a new bin when necessary.
// returns the new bin just created
func (binCollection *BinCollectionImpl) NewBin() *Bin {
	binCollection.Bins = append(binCollection.Bins, NewBin(binCollection.BinCapacity))
	binCollection.TotalBins++ // update our number of bins used
	return binCollection.GetLastBin()
}

// PackAll solve the underlying bin packing problem
func (binCollection *BinCollectionImpl) PackAll(items Items) {
	// reverse sort for algorithms that require it
	if binCollection.Algorithm == FirstFitDecreasing ||
		binCollection.Algorithm == BestFitDecreasing {
		sort.Sort(sort.Reverse(items))
	}
	for _, item := range items {
		binCollection.PackItem(item)
	}
}

// GetFirstBin getter for the first bin created
func (binCollection *BinCollectionImpl) GetFirstBin() *Bin {
	return &binCollection.Bins[0]
}

// GetLastBin getter for the bin most recently created in the problem
func (binCollection *BinCollectionImpl) GetLastBin() *Bin {
	return &binCollection.Bins[len(binCollection.Bins)-1]
}

// GetBin get an element at the given index in our list of bins
func (binCollection *BinCollectionImpl) GetBin(index int) *Bin {
	return &binCollection.Bins[index]
}

// PackItem method used for packing each item in succession
func (binCollection *BinCollectionImpl) PackItem(item Item) {
	switch binCollection.Algorithm {
	case NextFit:
		NextFitPack(binCollection, item)
	case FirstFit:
		FirstFitPack(binCollection, item)
	case FirstFitDecreasing:
		FirstFitDecreasingPack(binCollection, item)
	case BestFit:
		BestFitPack(binCollection, item)
	case BestFitDecreasing:
		BestFitDecreasingPack(binCollection, item)
	default:
		// do nothing here. Algorithm not supported
	}
}

// Find find an item based on the given predicate, nil if not found
func (binCollection *BinCollectionImpl) Find(predicate func(*Bin) bool) *Bin {
	var found *Bin
	for i := 0; i < int(binCollection.GetTotalBins()); i++ {
		if predicate(binCollection.GetBin(i)) {
			found = binCollection.GetBin(i)
		}
	}
	return found
}

// String return representation of this object as a string
func (binCollection *BinCollectionImpl) String() string {
	jsonString, _ := json.MarshalIndent(binCollection, "", "  ")
	return string(jsonString)
}
