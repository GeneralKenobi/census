package e2eutil

import "github.com/GeneralKenobi/census/pkg/util"

// RandomEmail generates a random email address.
func RandomEmail() string {
	return util.RandomAlphanumericString(16) + "@test.com"
}
