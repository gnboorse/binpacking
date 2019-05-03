package binpacking

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/gnboorse/centipede"
)

// import ("github.com/gnboorse/centipede")

// ConstraintPackingImpl implementation of Bin Collection
// that uses constraint programming
type ConstraintPackingImpl struct {
	BinCapacity Size  `json:"capacity"`
	TotalBins   Count `json:"count"`
	Bins        `json:"bins"`
	Algorithm   `json:"algorithm"`
}

// PackAll solve the underlying bin packing problem using constraint programming
func (binCollection *ConstraintPackingImpl) PackAll(items Items) {
	// reverse sort items by size
	sort.Sort(sort.Reverse(items))

	// get total number of items
	itemCount := len(items)

	vars := make(centipede.Variables, 0)
	constraints := make(centipede.Constraints, 0)
	propagations := make(centipede.Propagations, 0)

	sum := 0 // get sum of all items to pack
	for _, item := range items {
		sum += int(item)
	}

	// domain of all size variables is always limited from {1...C}
	siVariableDomain := centipede.IntRange(1, int(binCollection.BinCapacity+1))
	xijVariableDomain := centipede.IntRange(0, 2) // only 0 or 1
	// domain of load variables is always {0...C} in case there is an empty bin
	// liVariableDomain := centipede.IntRange(0, int(binCollection.BinCapacity+1))

	for i := 0; i < itemCount; i++ {
		siVariableName := centipede.VariableName("S" + strconv.Itoa(i))

		// create size variables and assign to them
		vars = append(vars, centipede.NewVariable(siVariableName, siVariableDomain))
		vars.SetValue(centipede.VariableName("S"+strconv.Itoa(i)), int(items[i]))

		xijVariableNames := make(centipede.VariableNames, 0)
		for j := 0; j < int(binCollection.TotalBins); j++ {
			xijVariableName := centipede.VariableName("X" + strconv.Itoa(i) + "_" + strconv.Itoa(j))
			xijVariableNames = append(xijVariableNames, xijVariableName)
			vars = append(vars, centipede.NewVariable(xijVariableName, xijVariableDomain))
		}
		// constraint indicating that the sum of all xij variables for this i and all j must be 1
		constraints = append(constraints, IntSumConstraint(1, Equality(), xijVariableNames))

		// do propagation to trim domains
		propagation := centipede.Propagation{
			Vars:                xijVariableNames,
			PropagationFunction: GetXIJPropagationFunction(i, int(binCollection.TotalBins))}
		propagations = append(propagations, propagation)
	}

	// iterate through all bins and create load variables
	// liVariableNames := make(centipede.VariableNames, 0)
	totalVariableNames := make(centipede.VariableNames, 0)
	for j := 0; j < int(binCollection.TotalBins); j++ {
		// liVariableName := centipede.VariableName("L" + strconv.Itoa(j))
		// liVariableNames = append(liVariableNames, liVariableName)
		// vars = append(vars, centipede.NewVariable(liVariableName, liVariableDomain))

		// calculate the list of VariableNames we care about
		constraintVarNames := make(centipede.VariableNames, 0)
		// constraintVarNames = append(constraintVarNames, liVariableName)
		for i := 0; i < itemCount; i++ {
			xijVariableName := centipede.VariableName("X" + strconv.Itoa(i) + "_" + strconv.Itoa(j))
			siVariableName := centipede.VariableName("S" + strconv.Itoa(i))
			constraintVarNames = append(constraintVarNames, xijVariableName)
			constraintVarNames = append(constraintVarNames, siVariableName)
		}
		// fmt.Printf("%v\n", constraintVarNames)

		// add sum constraint
		constraints = append(constraints, centipede.Constraint{
			Vars:               constraintVarNames,
			ConstraintFunction: GetSumConstraintFunction(j, itemCount, int(binCollection.BinCapacity))})

		totalVariableNames = append(totalVariableNames, constraintVarNames...)
	}

	// get the sum propagation
	sumPropagation := centipede.Propagation{
		Vars:                totalVariableNames,
		PropagationFunction: GetXIJSumPropagationFunction(int(binCollection.TotalBins), itemCount, int(binCollection.BinCapacity))}
	propagations = append(propagations, sumPropagation)

	// constraint function used for redundant check that all sums are equal to the total
	redundantConstraintFunction := func(variables *centipede.Variables) bool {
		runningSum := 0
		for j := 0; j < int(binCollection.TotalBins); j++ {
			for i := 0; i < itemCount; i++ {
				siVariableName := centipede.VariableName("S" + strconv.Itoa(i))
				siVariable := variables.Find(siVariableName)
				if siVariable.Empty {
					return true
				}
				xijVariableName := centipede.VariableName("X" + strconv.Itoa(i) + "_" + strconv.Itoa(j))
				xijVariable := variables.Find(xijVariableName)
				if xijVariable.Empty {
					return true
				}
				// sum is sigma(xij * si)
				runningSum += (siVariable.Value.(int) * xijVariable.Value.(int))
			}
		}
		return runningSum == sum
	}

	// redundant constraint indicating that the sum of all load variables must
	// equal the total size of items to be packed
	constraints = append(constraints, centipede.Constraint{Vars: totalVariableNames, ConstraintFunction: redundantConstraintFunction})

	// create solver
	solver := centipede.NewBackTrackingCSPSolverWithPropagation(vars, constraints, propagations)

	fmt.Printf("Starting solve... \n")
	// solve for constraints
	begin := time.Now()
	success := solver.Solve()
	elapsed := time.Since(begin)

	if success {
		fmt.Printf("Found solution in %s\n", elapsed)
		fmt.Println("Printing loads...")

	} else {
		fmt.Printf("Could not find solution in %s\n", elapsed)
	}

	for i := 0; i < itemCount; i++ {
		siVariableName := centipede.VariableName("S" + strconv.Itoa(i))
		variable := solver.State.Vars.Find(siVariableName)
		fmt.Printf(" (%v %v) ", variable.Name, variable.Value)
	}

	// for _, n := range liVariableNames {
	// 	variable := solver.State.Vars.Find(n)
	// 	fmt.Printf(" (%v %v) ", variable.Name, variable.Value)
	// }

}

// GetXIJPropagationFunction get a propagation function for elminating xij variables such that we can only assign 1 to a single xij for an item
func GetXIJPropagationFunction(i, totalBins int) func(assignment centipede.VariableAssignment, variables *centipede.Variables) []centipede.DomainRemoval {
	return func(assignment centipede.VariableAssignment, variables *centipede.Variables) []centipede.DomainRemoval {
		domainRemovals := make(centipede.DomainRemovals, 0)
		targetVariable := variables.Find(assignment.VariableName)
		if targetVariable.Value.(int) == 1 {
			// this xij variable has been set to 1, all others can only be 0
			// iterate through bins
			for j := 0; j < int(totalBins); j++ {
				xijVariableName := centipede.VariableName("X" + strconv.Itoa(i) + "_" + strconv.Itoa(j))
				foundXiJVariable := variables.Find(xijVariableName)
				if assignment.VariableName != xijVariableName &&
					foundXiJVariable.Empty &&
					foundXiJVariable.Domain.Contains(1) {
					domainRemovals = append(domainRemovals, centipede.DomainRemoval{
						VariableName: xijVariableName,
						Value:        1})
				}
			}
		}

		return domainRemovals
	}
}

// GetXIJSumPropagationFunction get
func GetXIJSumPropagationFunction(totalBins, itemCount, capacity int) func(assignment centipede.VariableAssignment, variables *centipede.Variables) []centipede.DomainRemoval {
	return func(assignment centipede.VariableAssignment, variables *centipede.Variables) []centipede.DomainRemoval {
		domainRemovals := make(centipede.DomainRemovals, 0)
		// iterate through all bins
		for j := 0; j < totalBins; j++ {
			// get sum of items already packed in this bin
			binSum := 0
			for i := 0; i < itemCount; i++ {
				xijVariableName := centipede.VariableName("X" + strconv.Itoa(i) + "_" + strconv.Itoa(j))
				foundXiJVariable := variables.Find(xijVariableName)
				if !foundXiJVariable.Empty {
					siVariableName := centipede.VariableName("S" + strconv.Itoa(i))
					foundSiVariable := variables.Find(siVariableName)
					if !foundSiVariable.Empty {
						binSum += (foundSiVariable.Value.(int) * foundXiJVariable.Value.(int))
					}
				}
			}
			// iterate through all items
			for i := 0; i < itemCount; i++ {
				xijVariableName := centipede.VariableName("X" + strconv.Itoa(i) + "_" + strconv.Itoa(j))
				foundXiJVariable := variables.Find(xijVariableName)
				// looking for unassigned xij variables
				if foundXiJVariable.Empty {
					siVariableName := centipede.VariableName("S" + strconv.Itoa(i))
					foundSiVariable := variables.Find(siVariableName)
					if !foundSiVariable.Empty {
						itemSize := foundSiVariable.Value.(int)
						if (binSum >= capacity || (binSum+itemSize) > capacity) && foundXiJVariable.Domain.Contains(1) {
							// if the bin is already overpacked, or if the current item would make the bin overpacked
							// then reject the possibility that the current item could be in that bin
							fmt.Printf("Removing item %v with size %v from consideration for bin %v with sum %v (capacity: %v)\n",
								i, itemSize, j, binSum, capacity)
							domainRemovals = append(domainRemovals, centipede.DomainRemoval{
								VariableName: xijVariableName,
								Value:        1})
						}
					}
				}
			}
		}
		return domainRemovals
	}
}

// GetSumConstraintFunction constraint function that makes sure that the sum of all
// items packed in a bin is less than or equal to the bin capacity
func GetSumConstraintFunction(j, itemCount, capacity int) func(variables *centipede.Variables) bool {
	return func(variables *centipede.Variables) bool {
		// for all items...
		runningSum := 0
		for i := 0; i < itemCount; i++ {
			siVariableName := centipede.VariableName("S" + strconv.Itoa(i))
			siVariable := variables.Find(siVariableName)
			if siVariable.Empty {
				return true
			}
			xijVariableName := centipede.VariableName("X" + strconv.Itoa(i) + "_" + strconv.Itoa(j))
			xijVariable := variables.Find(xijVariableName)
			if xijVariable.Empty {
				return true
			}
			// sum is sigma(xij * si)
			runningSum += (siVariable.Value.(int) * xijVariable.Value.(int))
		}
		// li = sum of loads
		return runningSum <= capacity
	}
}

// ComparisonFunction function type used to compare a sum with a comparisonValue
type ComparisonFunction func(c, s int) bool

// Equality comparison function
func Equality() ComparisonFunction {
	return func(c, s int) bool {
		return c == s
	}
}

// IntSumConstraint constraint generator used to compare the sum of a set of variables to a constant value
func IntSumConstraint(
	comparisonValue int,
	comparison ComparisonFunction,
	names centipede.VariableNames) centipede.Constraint {
	return centipede.Constraint{
		Vars: names,
		ConstraintFunction: func(variables *centipede.Variables) bool {
			runningSum := 0
			for _, name := range names {
				s := variables.Find(name)
				if s.Empty {
					return true
				} else {
					runningSum += s.Value.(int)
				}
			}
			return comparison(comparisonValue, runningSum)
		}}
}

// GetTotalBins getter for the total number of bins
func (binCollection *ConstraintPackingImpl) GetTotalBins() Count {
	return binCollection.TotalBins
}

// GetBinCapacity getter for the individual bin capacities
func (binCollection *ConstraintPackingImpl) GetBinCapacity() Size {
	return binCollection.BinCapacity
}

// String return representation of this object as a string
func (binCollection *ConstraintPackingImpl) String() string {
	jsonString, _ := json.MarshalIndent(binCollection, "", "  ")
	return string(jsonString)
}
