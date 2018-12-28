package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/schollz/logparser"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("no log files")
		os.Exit(1)
	}

	ti := time.Now()
	f, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}

	_, fname := filepath.Split(os.Args[1])
	lp, err := logparser.New(f, logparser.OptionName(fname))
	if err != nil {
		panic(err)
	}

	err = lp.Parse()
	if err != nil {
		panic(err)
	}

	bJSON, _ := json.MarshalIndent(lp, "", " ")
	fmt.Print("\n\n")
	fmt.Println(string(bJSON))

	fmt.Println("finished in", time.Since(ti))

	type kv struct {
		Key   string
		Value int
	}

	var ss []kv
	for k, v := range lp.TotalHitsPerRoute {
		ss = append(ss, kv{k, v})
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value > ss[j].Value
	})

	for i, kv := range ss {
		if strings.Contains(kv.Key, "/img/") {
			continue
		}
		if strings.Contains(kv.Key, "png") {
			continue
		}
		if strings.Contains(kv.Key, "jpg") {
			continue
		}
		if strings.Contains(kv.Key, "css") {
			continue
		}
		if strings.Contains(kv.Key, ".js") {
			continue
		}
		fmt.Printf("%s, %d\n", kv.Key, kv.Value)
		if i > 100 {
			break
		}
	}
}
