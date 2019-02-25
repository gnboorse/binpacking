package binpacking

import (
	"math"
	"math/rand"
	"time"
)

// Variability basic type representing variability in data
type Variability float64

const (
	//HighVariability level of variability
	HighVariability Variability = iota + 1
	//MediumVariability level of variability
	MediumVariability
	//LowVariability level of variability
	LowVariability
)

// GenerateItems generate some number of items based on the number
// of items needed, the size of the items, and the desired center for the items.
func GenerateItems(itemCount, maxItemSize, itemCenter int, variability Variability) Items {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	half := maxItemSize / 2
	var sigma float64
	if half >= itemCenter {
		// most items will occupy less than half of the container
		sigma = float64(itemCenter) / float64(variability)
	} else {
		// most items will occupy more than half of the container
		sigma = float64(maxItemSize-itemCenter) / float64(variability)
	}
	items := make(Items, itemCount)
	for i := 0; i < itemCount; i++ {
		item := int(math.Round(NormalRandom(r, sigma, float64(itemCenter))))
		if item <= 0 { //minItemSize
			item = 1
		} else if item >= maxItemSize {
			item = maxItemSize - 1
		}
		items[i] = Item(item)
	}
	return items
}

// NormalRandom get a random number in the normal distribution described
// by the standard deviation and mean provided
func NormalRandom(r *rand.Rand, standardDeviation, mean float64) float64 {
	return r.NormFloat64()*standardDeviation + mean
}

// UnityBasedNormalization given x on the range of xMin to xMax,
// scale x on the range a to b
func UnityBasedNormalization(x, xMin, xMax, a, b float64) float64 {
	return a + ((x-xMin)*(b-a))/(xMax-xMin)
}
