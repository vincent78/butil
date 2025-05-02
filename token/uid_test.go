package token

import (
	"encoding/base32"
	"encoding/hex"
	"fmt"
	"github.com/rs/xid"
	"os"
	"strings"
	"testing"
)

func TestUId(t *testing.T) {

	for i := 0; i < 1000000; i++ {
		X.NewID()
	}

}

func TestXIDFromString(t *testing.T) {
	pgHex := "shetaghfdsjfh1827087"
	pgBin, err := hex.DecodeString(pgHex)
	if err != nil {
		fmt.Printf("Failed to decode, err: %v\n", err)
		os.Exit(1)
	}
	str32 := strings.ToLower(strings.TrimRight(base32.StdEncoding.EncodeToString(pgBin), "="))
	guid, err := xid.FromString(str32)
	if err != nil {
		fmt.Printf("Failed to parse: %s, err: %v\n", str32, err)
		os.Exit(1)
	}
	println(guid.String())
}
