package binpacking

// BestFitPack pack a single item using the best fit algorithm
func BestFitPack(binCollection *BinCollectionImpl, item Item) {

	// find the bin with the smallest non-negative remaining space after adding the item
	smallestRemainderIndex := 0
	var smallestRemainder Size
	for i := 0; i < int(binCollection.GetTotalBins()); i++ {
		// remainder = the amount of space left over after adding the item
		remainder := binCollection.GetBin(i).Remaining() - Size(item)
		if i == 0 {
			smallestRemainder = remainder
		} else if remainder >= 0 && (remainder < smallestRemainder || smallestRemainder < 0) {
			smallestRemainder = remainder
			smallestRemainderIndex = i
		}
	}
	// if we found a bin that the item will fit inside
	if smallestRemainder >= 0 {
		remainderBin := binCollection.GetBin(smallestRemainderIndex)
		remainderBin.Pack(item)
	} else {
		// create a new bin
		newBin := binCollection.NewBin()
		newBin.Pack(item)
	}

}

// BestFitDecreasingPack pack the next item using the best fit decreasing algorithm
func BestFitDecreasingPack(binCollection *BinCollectionImpl, item Item) {
	BestFitPack(binCollection, item)
}
