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
	must(err)

	must(os.WriteFile(*output, data, 0644))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
