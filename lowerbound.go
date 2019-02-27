package binpacking

import (
	"math"
	"sort"
)

// CalculateLowerBound calculate the estimated lower bound
// on the minimum (optimal) number of bins that should be
// used in a problem instance.
// This is the basic algorithm used by Richard E. Korf in
// his first Bin Completion paper.
func CalculateLowerBound(items Items, binSize Size) Count {
	sort.Sort(sort.Reverse(items)) // sort items in decreasing order
	waste := 0                     // total space wasted in the ideal solution
	j := 1                         // j = pointer to end of items list
	carry := 0                     // carry over for sums
	itemSum := 0
	for i := 0; i <= len(items)-j; i++ {
		x := items[i]              // iterate over every item x
		itemSum += int(x)          // add to item sum
		r := int(binSize) - int(x) // find remaining space in bin
		// find all elements <= r
		lessThanR := make(Items, 0)
		// iterate from end of array
		for k := len(items) - j; k > i; k-- {
			sItem := items[k]
			if int(sItem) <= r {
				// consider all items less than or equal to r
				lessThanR = append(lessThanR, sItem)
			}
			if int(sItem) > r {
				break // no need to consider anything larger, since items are sorted
			}
		}
		s := 0 // sum of all items less than r
		if r >= carry {
			// only remove items to consider if we have a carry < r
			for _, sItem := range lessThanR {
				s += int(sItem)
				j++ // move up "pointer" to end of list
				itemSum += int(sItem)
			}
		}
		s += carry // add space carried over from previous bin

		if r == s {
			// no wasted space, and no carry over to next bin
			carry = 0
		} else if s < r {
			// some wasted space, but no carry over
			waste += r - s
			carry = 0
		} else if r < s {
			// no wasted space, but carry over s-r to next bin
			carry = s - r
		}
	}
	// return (sum of items + waste) divided by the bin size, rounded up
	return Count(int(math.Round(float64(itemSum+waste) / float64(int(binSize)))))
}
