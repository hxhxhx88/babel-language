package main

import (
	"flag"
	"fmt"
	"os"

	"babel/openapi"
)

func main() {
	output := flag.String("output", "", "output path")
	flag.Parse()

	if *output == "" {
		fmt.Fprintf(os.Stderr, "must specify --output")
		return
	}

	babelapi := openapi.NewBabel()

	data, err := babelapi.MarshalJSON()
	mustOk(err)

	mustOk(os.WriteFile(*output, data, 0644))
}

func mustOk(err error) {
	if err != nil {
		panic(err)
	}
}
