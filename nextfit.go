package binpacking

// NextFitPack pack the next item using the next fit algorithm
func NextFitPack(binCollection *BinCollectionImpl, item Item) {
	mostRecentBin := binCollection.GetLastBin()
	if mostRecentBin.CanFit(item) {
		mostRecentBin.Pack(item)
	} else {
		newBin := binCollection.NewBin()
		newBin.Pack(item)
	}
}
