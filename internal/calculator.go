// Package calculator conatains the actual propability calculation
package calculator

import (
	"iter"
	"maps"
	"math/rand/v2"
)

type Calculator struct {
	DiceSides             int
	AmountOfDice          int
	MinimumRollForSuccess int
	MinimumRollToExplode  int
	MaximumExplodingRolls int
	AmountOfRolls         int
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
	var results []calculatorResult
	successCount := 0

	for i := 0; i < c.AmountOfRolls; i++ {
		rolledDice := 0
		explodedDice := 0

		for rolledDice < c.AmountOfDice {
			roll := roll(1, c.DiceSides)

			if roll >= c.MinimumRollForSuccess {
				successCount++
			}

			if c.MaximumExplodingRolls == explodedDice {
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
	return rand.IntN(max+1-min) + min
}
