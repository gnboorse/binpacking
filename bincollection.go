package binpacking

import (
	"encoding/json"
	"fmt"
	"math"
	"sort"
)

// BinCollection an interface representing an instance
// of the bin packing problem.
type BinCollection interface {
	PackAll(items Items)
	GetTotalBins() Count
	GetBinCapacity() Size
	SetTime(nanoseconds int64)
	String() string
}

// NewBinCollection create an instance of the bin packing problem
// from a PackingList object
func NewBinCollection(pList *PackingList) BinCollection {

	collection := &BinCollectionImpl{
		BinCapacity: pList.Size,
		TotalBins:   0,
		Bins:        make(Bins, 0), // pre-allocate memory for a reasonably large capacity
		Algorithm:   pList.Algorithm}

	// init bins for packing constraint
	if pList.Algorithm == PackingConstraint {
		roundedLowerBound := int(math.Round(float64(pList.LowerBound) * 1.2))
		for i := 0; i < roundedLowerBound; i++ {
			collection.NewBin() // add new bins for all of them
		}
	} else if pList.Algorithm != ModifiedFirstFitDecreasing {
		collection.NewBin() // always create first bin if not MFFD or constraint
	}
	return collection

}

// BinCollectionImpl default implementation of an
// instance of the bin packing problem
type BinCollectionImpl struct {
	BinCapacity  Size  `json:"capacity"`
	TotalBins    Count `json:"count"`
	Bins         `json:"bins"`
	Algorithm    `json:"algorithm"`
	SolutionTime int64 `json:"solution_time"`
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
		binCollection.Algorithm == BestFitDecreasing ||
		binCollection.Algorithm == ModifiedFirstFitDecreasing ||
		binCollection.Algorithm == PackingConstraint {
		sort.Sort(sort.Reverse(items))
	}
	if binCollection.Algorithm == ModifiedFirstFitDecreasing {
		binCollection.PackAllMFFD(items)
	} else if binCollection.Algorithm == PackingConstraint {
		binCollection.PackAllConstraint(items)
	} else {
		for _, item := range items {
			binCollection.PackItem(item)
		}
	}
	binCollection.cleanupBins()
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
		panic(fmt.Errorf("unsupported algorithm for PackItem: %v", binCollection.Algorithm))
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

// SetTime set the execution time for a single run
func (binCollection *BinCollectionImpl) SetTime(nanoseconds int64) {
	binCollection.SolutionTime = nanoseconds
}

// cleanupBins convenience method used to clean out unused bins from the list
func (binCollection *BinCollectionImpl) cleanupBins() {

	for {
		emptyIndex := -1
		for i := 0; i < int(binCollection.GetTotalBins()); i++ {
			bin := binCollection.GetBin(i)
			if bin.Usage == 0 && len(bin.Items) == 0 {
				// empty bin.
				emptyIndex = i
			}
		}
		if emptyIndex >= 0 {
			binCollection.Bins[emptyIndex] = binCollection.Bins[len(binCollection.Bins)-1] // Copy last element to index i.
			binCollection.Bins[len(binCollection.Bins)-1] = Bin{}                          // Erase last element (write zero value).
			binCollection.Bins = binCollection.Bins[:len(binCollection.Bins)-1]
			binCollection.TotalBins--
		} else {
			break
		}
	}
}
