package console

import (
	"testing"
)

func TestColorShow(t *testing.T) {
	t.Log("this is the color test ---")
	str := "this is the test String"

	t.Logf(" %v", Red(str))
	t.Logf(" %v", Yellow(str))
	t.Logf(" %v", Green(str))
	t.Logf(" %v", Black(str))
	t.Logf(" %v", Blue(str))
	t.Logf(" %v", Magenta(str))
	t.Logf(" %v", Cyan(str))
	t.Logf(" %v", White(str))
}
