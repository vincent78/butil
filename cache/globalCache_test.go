package cache

import (
	"github.com/stretchr/testify/assert"
	"testing"
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
