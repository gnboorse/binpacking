package binpacking

// Item representation of an item being packed into a bin
type Item int

// Items collection type for Item
type Items []Item

// Len used to implement sort.Interface
func (items Items) Len() int {
	return len(items)
}

// Less used to implement sort.Interface
func (items Items) Less(i, j int) bool {
	return items[i] < items[j]
}

func (items Items) Swap(i, j int) {
	items[i], items[j] = items[j], items[i]
}

// Size scalar indicating the size of something
type Size int

// Count scalar indicating a number of items
type Count int
