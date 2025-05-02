package mapUtil

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLocation(t *testing.T) {

	m := map[int]int{
		0:   0,
		10:  10,
		100: 100,
		200: 200,
		300: 300,
	}

	testCases := []struct {
		Value  int
		Map    map[int]int
		Result int
	}{
		{Value: 0, Map: m, Result: 0},
		{Value: 10, Map: m, Result: 10},
		{Value: 11, Map: m, Result: 10},
		{Value: 99, Map: m, Result: 10},
		{Value: 100, Map: m, Result: 100},
		{Value: 101, Map: m, Result: 100},
		{Value: 400, Map: m, Result: 300},
	}

	for _, test := range testCases {
		t.Run(fmt.Sprintf("check-%v", test.Value), func(t *testing.T) {
			t.Parallel()
			if r, err := Location[int](test.Value, test.Map); err != nil {
				assert.Error(t, err)
			} else {
				assert.Equal(t, r, test.Result)
			}

		})
	}
}
