package main

import (
	"flag"
	"fmt"
	"os"

	"gopkg.in/yaml.v2"

	"babel/openapi"
)

func main() {
	output := flag.String("output", "", "output path")
	flag.Parse()

	if *output == "" {
		fmt.Fprintf(os.Stderr, "must specify --output")
		return
	}

	spec := openapi.New()

	data, err := yaml.Marshal(spec)
	mustOk(err)

	mustOk(os.WriteFile(*output, data, 0644))
}

func mustOk(err error) {
	if err != nil {
		panic(err)
	}
}
