package binpacking

import (
	"strconv"

	"github.com/gnboorse/centipede"
)

// PackAllConstraint pack all items using constraints
func (binCollection *BinCollectionImpl) PackAllConstraint(items Items) {
	itemCount := len(items)
	vars := make(centipede.Variables, 0)
	constraints := make(centipede.Constraints, 0)
	propagations := make(centipede.Propagations, 0)

	// placement range can be any index in the range of bins
	itemPlacementVariableNames := make(centipede.VariableNames, 0)
	itemPlacementVariableDomain := centipede.IntRange(0, int(binCollection.GetTotalBins()))
	for i := 0; i < itemCount; i++ {
		itemPlacementVariableName := centipede.VariableName("ItemPlacement" + strconv.Itoa(i))
		vars = append(vars, centipede.NewVariable(itemPlacementVariableName, itemPlacementVariableDomain))
		itemPlacementVariableNames = append(itemPlacementVariableNames, itemPlacementVariableName)
	}

	sumConstraint := centipede.Constraint{
		Vars: itemPlacementVariableNames,
		ConstraintFunction: func(variables *centipede.Variables) bool {
			sums := make([]int, binCollection.GetTotalBins(), binCollection.GetTotalBins())
			for i := 0; i < itemCount; i++ {
				itemPositionVar := variables.Find(itemPlacementVariableNames[i])
				if !itemPositionVar.Empty {
					itemPosition := itemPositionVar.Value.(int)
					itemSize := int(items[i])
					if sums[itemPosition]+itemSize > int(binCollection.GetBinCapacity()) {
						return false
					}
					sums[itemPosition] += itemSize
				}
			}
			return true
		},
	}

	constraints = append(constraints, sumConstraint)

	// adding this propagation improves performance by factors of 5
	sumPropagation := centipede.Propagation{
		Vars: itemPlacementVariableNames,
		PropagationFunction: func(assignment centipede.VariableAssignment, variables *centipede.Variables) []centipede.DomainRemoval {
			binIndexAssigned := assignment.Value.(int)
			// calculate runningSum to be the total sum of all items placed in the bin just assigned to
			runningSum := 0
			potentialDomainRemovals := make(centipede.DomainRemovals, 0)
			// iterate over items
			for i := 0; i < itemCount; i++ {
				// find the position variable
				itemPositionVar := variables.Find(itemPlacementVariableNames[i])
				if !itemPositionVar.Empty {
					itemPosition := itemPositionVar.Value.(int)
					// check if this item is in the same bin that we just assigned to
					if itemPosition == binIndexAssigned {
						itemSize := int(items[i])
						runningSum += itemSize
					}
				} else {
					// pre-calculate what our domain removals would be if we have maxed out the sum
					potentialDomainRemovals = append(potentialDomainRemovals, centipede.DomainRemoval{
						VariableName: itemPlacementVariableNames[i],
						Value:        assignment.Value})
				}
			}
			// return domain removals if necessary
			if runningSum > int(binCollection.BinCapacity) {
				return potentialDomainRemovals
			}
			return []centipede.DomainRemoval{}
		},
	}

	propagations = append(propagations, sumPropagation)

	// create solver
	solver := centipede.NewBackTrackingCSPSolverWithPropagation(vars, constraints, propagations)

	// solve for constraints
	solver.Solve()

	for i := 0; i < itemCount; i++ {
		itemPlacementVariableName := itemPlacementVariableNames[i]
		variableValue := solver.State.Vars.Find(itemPlacementVariableName)
		if !variableValue.Empty {
			bin := binCollection.GetBin(variableValue.Value.(int))
			bin.Pack(items[i])
		}
	}
}
