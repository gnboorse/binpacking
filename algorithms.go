package binpacking

// Algorithm types of algorithms supported by the library
type Algorithm int

const (
	// Unknown for default value initialization
	Unknown Algorithm = iota
	// NextFit is a naive approach that puts objects in the next bin that can fit
	NextFit
	// FirstFit puts objects in the first bin currently known to hold its capacity
	FirstFit
	// FirstFitDecreasing first sorts objects by size (decreasing) and then applies FirstFit
	FirstFitDecreasing
	// BestFit puts objects in the tighest spot in all bins currently known
	BestFit
	// BestFitDecreasing first sorts objects by size (decreasing) and then applies BestFit
	BestFitDecreasing
)

var names = []string{"Unknown", "NextFit", "FirstFit", "FirstFitDecreasing", "BestFit", "BestFitDecreasing"}

func (algorithm Algorithm) String() string {
	return names[algorithm]
}

// GetAlgorithm get an algorithm from string
func GetAlgorithm(s string) Algorithm {
	for i, name := range names {
		if name == s {
			return Algorithm(i)
		}
	}
	return Unknown
}
