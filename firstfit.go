package binpacking

// FirstFitPack pack the next item using the first fit algorithm
func FirstFitPack(binCollection BinCollection, item Item) {
	found := binCollection.Find(func(bin *Bin) bool { return bin.CanFit(item) })
	if found != nil {
		found.Pack(item)
	} else {
		newBin := binCollection.NewBin()
		newBin.Pack(item)
	}
}

// FirstFitDecreasingPack pack the next item using the first fit decreasing algorithm
func FirstFitDecreasingPack(binCollection BinCollection, item Item) {
	FirstFitPack(binCollection, item)
}
