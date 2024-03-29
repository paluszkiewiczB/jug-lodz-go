//go:build integration

package counter_test

import (
	"testing"
	"time"
)

func Test_Integration_VerySlow(t *testing.T) {
	time.Sleep(10 * time.Second)
}
