package binpackingtests

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/gnboorse/binpacking"
)

// TestCalculateLowerBound unit test for calculating the lower bound (optimal solution)
func TestCalculateLowerBound(t *testing.T) {
	b, err := ioutil.ReadFile("test.json")
	if err != nil {
		panic(err)
	}

	var packingList binpacking.PackingList
	err = json.Unmarshal(b, &packingList)
	if err != nil {
		panic(err)
	}
	lowerBound := binpacking.CalculateLowerBound(packingList.Items, packingList.Size)
	if lowerBound != 57 { // very high lower bound for the test json
		t.Errorf("Calculated Lower Bound was: %v", lowerBound)
	}

}
