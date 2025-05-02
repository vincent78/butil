package token

import (
	"crypto/rand"
	"reflect"
	"regexp"
	"testing"
)

func TestUUID(t *testing.T) {
	uuid := NewUUID()
	t.Log("the uuid:", uuid.String())
}

func TestGetUUID(t *testing.T) {
	uuid := GetUUID()
	t.Log("the uuid:", uuid)
}

func TestNew(t *testing.T) {
	uuidStr := New()
	t.Log(uuidStr)
}

func TestNewShort(t *testing.T) {
	for i := 0; i < 20; i++ {
		uuidStr := NewShort()
		t.Log(uuidStr)
	}
}

func TestGenerateUUID(t *testing.T) {
	prev, err := GenerateUUID()
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 100; i++ {
		id, err := GenerateUUID()
		if err != nil {
			t.Fatal(err)
		}
		if prev == id {
			t.Fatalf("Should get a new ID!")
		}

		matched, err := regexp.MatchString(
			"[\\da-f]{8}-[\\da-f]{4}-[\\da-f]{4}-[\\da-f]{4}-[\\da-f]{12}", id)
		if !matched || err != nil {
			t.Fatalf("expected match %s %v %s", id, matched, err)
		}
	}
}

func TestParseUUID(t *testing.T) {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		t.Fatalf("failed to read random bytes: %v", err)
	}

	uuidStr, err := FormatUUID(buf)
	if err != nil {
		t.Fatal(err)
	}

	parsedStr, err := ParseUUID(uuidStr)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(parsedStr, buf) {
		t.Fatalf("mismatched buffers")
	}
}

func BenchmarkGenerateUUID(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_, _ = GenerateUUID()
	}
}
