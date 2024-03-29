package footest

import (
	"testing"
)

func ShouldReturnOne(t *testing.T, f func() int) {
	if res := f(); res != 1 {
		t.Errorf("expected 1, got: %d", res)
	}
}
