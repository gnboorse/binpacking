package binpacking

// PackingList input for a bin packing problem
type PackingList struct {
	// BinSize the size of the bins being packed
	Size `json:"capacity"`
	// ItemCount the number of items being passed in
	Count `json:"count"`
	// Algorithm the algorithm being used to solve the problem
	Algorithm `json:"algorithm"`
	// Items the actual items being passed in
	Items       `json:"items"`
	Variability `json:"variability"`
	Center      int `json:"center"`
}
