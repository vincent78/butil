package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	k = "testkey"
	v = "testvalue"
)

func TestNormal(t *testing.T) {
	Save(k, v)
	nv := Get(k)
	assert.Equal(t, v, nv)
}
