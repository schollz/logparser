package logparser

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const logs = `66.1.2.3 - - [22/May/2018:15:52:28 -0600] "GET /book/belle-de-jour-diary-of-an-unlikely-call-girl-by-belle-de-jour/ HTTP/1.1" 403 38
69.162.2.1 - - [22/May/2018:15:52:32 -0600] "GET / HTTP/1.1" 403 38
91.121.1.3 - - [22/May/2018:15:52:36 -0600] "GET /book/amigurumi-by-lan-anh-bui/ HTTP/1.1" 200 112423
66.1.2.3 - - [22/May/2018:15:52:38 -0600] "GET /book/gefahr-am-strand-mit-audio-cd-by-andrea-maria-wagner/ HTTP/1.1" 403 38
66.1.2.3 - - [22/May/2018:15:52:48 -0600] "GET /book/la-tela-del-ragno-il-delitto-moro-by-sergio-flamigni/ HTTP/1.1" 403 38
91.121.1.5 - - [22/May/2018:15:52:55 -0600] "GET /book/aminta-by-torquato-tasso/ HTTP/1.1" 200 75659
66.1.2.3 - - [22/May/2018:15:52:57 -0600] "GET /book/perfect-together-by-lisa-plumley/ HTTP/1.1" 403 38
66.249.5.2 - - [22/May/2018:15:53:07 -0600] "GET /book/fahrerflucht-h%C3%A3%C2%B6rspiel-ein-liebhaber-des-halbschattens-by-alfred-andersch/ HTTP/1.1" 403 38
66.1.2.3 - - [22/May/2018:15:53:17 -0600] "GET /book/selected-to-live-by-johanna-ruth-dobschiner/ HTTP/1.1" 403 38
91.121.1.5 - - [22/May/2018:15:53:17 -0600] "GET /book/ammonite-by-nicola-griffith/ HTTP/1.1" 200 74852`

const logs1 = `66.1.2.3 - - [22/May/2018:15:52:57 -0600] "GET /book HTTP/1.1" 403 38`

func TestParseLogReader(t *testing.T) {
	lp, err := New(strings.NewReader(logs))
	assert.Nil(t, err)
	err = lp.Parse()
	assert.Nil(t, err)
}

func TestParseBook(t *testing.T) {
	ti := time.Now()
	f, err := os.Open("www.booksuggestions.ninja.log")
	assert.Nil(t, err)
	lp, err := New(f, OptionName("booksuggestions"))
	err = lp.Parse()
	assert.Nil(t, err)
	fmt.Println(len(lp.logData))
	fmt.Println(time.Since(ti))
	bJSON, _ := json.MarshalIndent(lp, "", " ")
	fmt.Println(string(bJSON))
}

func TestParseCommon(t *testing.T) {
	ll, err := parseCommon(logs1)
	assert.Nil(t, err)
	fmt.Println(ll)
}

func BenchmarkParseCommon(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parseCommon(logs1)
	}
}
