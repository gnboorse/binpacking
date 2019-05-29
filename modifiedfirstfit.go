package binpacking

// MFFDCategory type representing a category in MFFD
type MFFDCategory int

func categorize(item Item, capacity Size) MFFDCategory {
	if Size(item) > capacity/2 {
		return MFFDCategory(a)
	} else if Size(item) > capacity/3 {
		return MFFDCategory(b)
	} else if Size(item) > capacity/4 {
		return MFFDCategory(c)
	} else if Size(item) > capacity/5 {
		return MFFDCategory(d)
	} else if Size(item) > capacity/6 {
		return MFFDCategory(e)
	} else if Size(item) > (11/71)*capacity {
		return MFFDCategory(f)
	} else {
		return MFFDCategory(g)
	}
}

const (
	a = iota + 1
	b
	c
	d
	e
	f
	g
)

// PackAllMFFD pack all items for MFFD
func (binCollection *BinCollectionImpl) PackAllMFFD(items Items) {
	aCount := 0
	bItems := make(Items, 0)
	cdeItems := make(Items, 0)
	fgItems := make(Items, 0)

	// STEP 1

	for _, item := range items {
		category := categorize(item, binCollection.BinCapacity)
		if category == a {
			aCount++
			newBin := binCollection.NewBin()
			newBin.Pack(item)
		} else if category == b {
			// add to B items list, noting that items stay sorted in descending order
			bItems = append(bItems, item)
		} else if category == c || category == d || category == e {
			// add to a union C D E list
			cdeItems = append(cdeItems, item)
		} else if category == f || category == g {
			// add to a union F G list
			fgItems = append(fgItems, item)
		}
	}

	// STEP 2

	// maintain a list of all A containers containing a B value
	bPresenceList := make([]bool, aCount, aCount)
	// loop through B containers
	for i := 0; i < aCount; i++ {
		// indicates whether a B item was packed or not.
		bItemPacked := -1
		// iterate through B items
		for bIndex, bItem := range bItems {
			// if the bItem remains unpacked
			if bItem > 0 {
				// get the current A bin
				bin := binCollection.GetBin(i)
				if bin.CanFit(bItem) {
					// the unpacked item can fit, so pack it
					bin.Pack(bItem)
					bPresenceList[i] = true
					bItemPacked = bIndex
					break
				}
			}
		}
		// if we packed an item, set its value in bItems to -1 so we
		// don't pack it twice
		if bItemPacked > 0 {
			bItems[bItemPacked] = -1
		}
	}

	// STEP 3

	// iterate backwards through A containers
	for i := aCount - 1; i >= 0; i-- {
		// don't pack anything if we have already packed a B in this bin
		if bPresenceList[i] {
			continue
		}

		// this is the indices of the current two smallest unpacked items in C D or E
		twoSmallest := [2]int{-1, -1}
		// iterate through the C D E list from smallest to largest
		for j := len(cdeItems) - 1; j >= 0; j-- {
			// the first item in the twoSmallest array should be the smallest
			if twoSmallest[0] == -1 && cdeItems[j] != -1 {
				twoSmallest[0] = j
				// the second item in the twoSmallest array should be the second smallest
			} else if twoSmallest[0] != -1 && twoSmallest[1] == -1 && cdeItems[j] != -1 {
				twoSmallest[1] = j
			}
		}

		if twoSmallest[0] < 0 || twoSmallest[1] < 0 {
			// do nothing if we have fewer than 2 smallest items
			continue
		}

		// get the current A bin
		bin := binCollection.GetBin(i)

		// do nothing if the bin cannot fit the sum of the two smallest items in C D E
		if !bin.CanFit(Item(cdeItems[twoSmallest[0]] + cdeItems[twoSmallest[1]])) {
			continue
		}

		// pack item since we know it and at least one other C D E item will fit in the bin
		bin.Pack(cdeItems[twoSmallest[0]])
		// make sure we never pack it again
		cdeItems[twoSmallest[0]] = -1

		// loop through the C D E items again and pack the largest unpacked C D E item that will fit
		for j := 0; j < len(cdeItems); j++ {
			cdeItem := cdeItems[j]
			if cdeItem > 0 {
				if bin.CanFit(cdeItem) {
					bin.Pack(cdeItem)
					cdeItems[j] = -1
				}
			}
		}
	}

	// STEP 4

	unpackedCounter := 0
	// try to pack remaining items that fit in A bins
	for {
		madeAssignment := false
		// loop through A bins
		for i := 0; i < aCount; i++ {
			unpackedCounter = 0 // reset to zero for each bin to get an accurate count
			// get the current A bin
			bin := binCollection.GetBin(i)

			// attempt to pack a B item
			for j := 0; j < len(bItems); j++ {
				if bItems[j] > 0 {
					if bin.CanFit(bItems[j]) {
						bin.Pack(bItems[j])
						bItems[j] = -1
						madeAssignment = true
					} else {
						unpackedCounter++
					}
				}
			}

			// attempt to pack a C D E item
			for j := 0; j < len(cdeItems); j++ {
				if cdeItems[j] > 0 {
					if bin.CanFit(cdeItems[j]) {
						bin.Pack(cdeItems[j])
						cdeItems[j] = -1
						madeAssignment = true
					} else {
						unpackedCounter++
					}
				}
			}

			// attempt to pack an F G item
			for j := 0; j < len(fgItems); j++ {
				if fgItems[j] > 0 {
					if bin.CanFit(fgItems[j]) {
						bin.Pack(fgItems[j])
						fgItems[j] = -1
						madeAssignment = true
					} else {
						unpackedCounter++
					}
				}
			}
		}

		if !madeAssignment {
			break
		}
	}
	// STEP 5

	// if we actually have remaining items
	if unpackedCounter > 0 || binCollection.GetTotalBins() == 0 {
		// assume we have to create at least one new bin
		binCollection.NewBin()
		for {
			innerLoopAssignment := false
			unpackedCounter := 0
			// loop through bins beyond A bins
			for i := aCount; i < int(binCollection.GetTotalBins()); i++ {
				unpackedCounter = 0 // reset to zero for each bin to get an accurate count
				// get the current bin
				bin := binCollection.GetBin(i)

				// attempt to pack a B item
				for j := 0; j < len(bItems); j++ {
					if bItems[j] > 0 {
						if bin.CanFit(bItems[j]) {
							bin.Pack(bItems[j])
							bItems[j] = -1
							innerLoopAssignment = true
						} else {
							unpackedCounter++
						}
					}
				}

				// attempt to pack a C D E item
				for j := 0; j < len(cdeItems); j++ {
					if cdeItems[j] > 0 {
						if bin.CanFit(cdeItems[j]) {
							bin.Pack(cdeItems[j])
							cdeItems[j] = -1
							innerLoopAssignment = true
						} else {
							unpackedCounter++
						}
					}
				}

				// attempt to pack an F G item
				for j := 0; j < len(fgItems); j++ {
					if fgItems[j] > 0 {
						if bin.CanFit(fgItems[j]) {
							bin.Pack(fgItems[j])
							fgItems[j] = -1
							innerLoopAssignment = true
						} else {
							unpackedCounter++
						}
					}
				}
			}

			// if we haven't made an assignment to any of the bins and we still have more to pack
			if !innerLoopAssignment && unpackedCounter > 0 {
				binCollection.NewBin()
			}

			// break when there are no more bins to pack
			if unpackedCounter == 0 {
				break
			}
		}
	}

}
