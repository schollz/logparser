package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
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
	fmt.Println(string(bJSON))

	fmt.Println("finished in", time.Since(ti))
}
