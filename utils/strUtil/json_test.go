package strUtil

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"
	"testing"
)

const td = `
[
{
	"site" : "npr",
	"link" : "http://www.npr.org/rss/rss.php?id=1001",
	"type" : "rss"
},
{
	"site" : "npr",
	"link" : "http://www.npr.org/rss/rss.php?id=1008",
	"type" : "rss"
},
{
	"site" : "npr",
	"link" : "http://www.npr.org/rss/rss.php?id=1006",
	"type" : "rss"
}
]
`

const t1 = `{
"site" : "npr",
"link" : "http://www.npr.org/rss/rss.php?id=1001",
"type" : "rss"
}
`

type Feed struct {
	Name string `json:"site"`
	URI  string `json:"link"`
	Type string `json:"type"`
}

func TestDecoder2Json(t *testing.T) {
	var fp []*Feed
	err := Decoder2Json(td, &fp)
	if err != nil {
		t.Errorf("解析出错: %v", err)
	} else {
		t.Logf("point:%T  value:%v", fp, fp)
	}

	var fo []Feed
	err = Decoder2Json(td, &fo)
	if err != nil {
		t.Errorf("解析出错: %v", err)
	} else {
		t.Logf("point:%T  value:%v", fo, fo)
	}

	var f1 Feed
	err = Decoder2Json(t1, &f1)
	if err != nil {
		t.Errorf("解析出错: %v", err)
	} else {
		t.Logf("point:%T  value:%v", f1, f1)
	}

	fp1 := &Feed{}
	err = Decoder2Json(t1, fp1)
	if err != nil {
		t.Errorf("解析出错: %v", err)
	} else {
		t.Logf("point:%T  value:%v", fp1, fp1)
	}

}

func TestReadFromFile(t *testing.T) {

	const jsonStream = `
	{"Name": "Ed", "Text": "Knock knock."}
	{"Name": "Sam", "Text": "Who's there?"}
	{"Name": "Ed", "Text": "Go fmt."}
	{"Name": "Sam", "Text": "Go fmt who?"}
	{"Name": "Ed", "Text": "Go fmt yourself!"}
`
	type Message struct {
		Name, Text string
	}
	dec := json.NewDecoder(strings.NewReader(jsonStream))
	for {
		var m Message
		if err := dec.Decode(&m); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s: %s\n", m.Name, m.Text)
	}

}
