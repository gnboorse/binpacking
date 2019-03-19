package binpacking

// FirstFitPack pack the next item using the first fit algorithm
func FirstFitPack(binCollection *BinCollectionImpl, item Item) {
	found := binCollection.Find(func(bin *Bin) bool { return bin.CanFit(item) })
	if found != nil {
		found.Pack(item)
	} else {
		newBin := binCollection.NewBin()
		newBin.Pack(item)
	}
}

// FirstFitDecreasingPack pack the next item using the first fit decreasing algorithm
func FirstFitDecreasingPack(binCollection *BinCollectionImpl, item Item) {
	FirstFitPack(binCollection, item)
}
