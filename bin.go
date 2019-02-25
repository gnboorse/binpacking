package binpacking

// Bin representing an individual bin containing items
type Bin struct {
	Capacity Size `json:"capacity"`
	Items    `json:"items"`
	Usage    Size `json:"usage"`
}

// Bins collection type for Bin
type Bins []Bin

// NewBin create a new bin
func NewBin(size Size) Bin {
	return Bin{size, make(Items, 0), 0}
}

// Remaining get the amount of remaining space in this bin
func (bin *Bin) Remaining() Size {
	return bin.Capacity - bin.Usage
}

// CanFit check if the given bin can fit an item
func (bin *Bin) CanFit(item Item) bool {
	return bin.Remaining() >= Size(item)
}

// Pack adds an item to a Bin
func (bin *Bin) Pack(item Item) {
	bin.Items = append(bin.Items, item)
	bin.Usage += Size(item) // update bin capacity
}
