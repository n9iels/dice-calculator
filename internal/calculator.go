// Package calculator conatains the actual propability calculation
package calculator

import (
	"iter"
	"maps"
	"math/rand/v2"
	"slices"
	"strconv"
	"strings"
)

type Calculator struct {
	DiceSides             int
	AmountOfDice          int
	MinimumRollForSuccess int
	MinimumRollToExplode  int
	MaximumExplodingRolls int
	AmountOfRolls         int
	DiceSidesForFailure   string
}

type CalculatorOutput struct {
	AmountOfSuccess int
	AmountOfRolls   int
	Probability     int
}

type calculatorResult struct {
	AmountOfSuccess int
	AmountOfFailure int
}

func (c Calculator) Calculate() iter.Seq[CalculatorOutput] {
	var failureDiceSidesSlice = strings.Split(c.DiceSidesForFailure, ",")
	var results []calculatorResult
	successCount := 0

	for i := 0; i < c.AmountOfRolls; i++ {
		rolledDice := 0
		explodedDice := 0

		for rolledDice < c.AmountOfDice {
			roll := roll(1, c.DiceSides)

			if slices.Contains(failureDiceSidesSlice, strconv.Itoa(roll)) {
				successCount--
				rolledDice++
				continue
			}

			if roll >= c.MinimumRollForSuccess {
				successCount++
			}

			if c.MaximumExplodingRolls != 0 && c.MaximumExplodingRolls == explodedDice {
				rolledDice++
				continue
			}

			if c.MinimumRollToExplode == 0 || roll < c.MinimumRollToExplode {
				rolledDice++
			}
		}

		results = append(results, calculatorResult{
			AmountOfSuccess: successCount,
		})
		successCount = 0
		explodedDice = 0
	}

	var output = make(map[int]CalculatorOutput)
	for _, v := range results {
		val, ok := output[v.AmountOfSuccess]

		if ok {
			val.AmountOfRolls++
			output[v.AmountOfSuccess] = val
			continue
		}

		output[v.AmountOfSuccess] = CalculatorOutput{
			AmountOfSuccess: v.AmountOfSuccess,
			AmountOfRolls:   1,
		}
	}

	return maps.Values(output)
}

func roll(min int, max int) int {
	return rand.IntN(max-min+1) + min
}
