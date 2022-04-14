package e2e

import (
	"fmt"
	"github.com/GeneralKenobi/census/pkg/util"
	math_rand "math/rand"
)

func init() {
	seed, err := util.RngSeed()
	if err != nil {
		panic("error seeding random number generator: " + err.Error())
	}
	fmt.Printf("RNG seed: %d\n", seed)
	math_rand.Seed(seed)
}
