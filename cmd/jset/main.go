package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/xyproto/jpath"
	"log"
	"os"
)

func main() {
	flag.Parse()

	if len(flag.Args()) != 3 {
		fmt.Println("Syntax: jset [filename] [JSON path] [value]")
		fmt.Println("Example: jset books.json x[1].author Suzanne")
		os.Exit(1)
	}

	filename := flag.Args()[0]
	JSONpath := flag.Args()[1]
	value := flag.Args()[2]

	// Try to interpret the value as JSON (number, bool, null, object, array)
	var jsonVal any
	if err := json.Unmarshal([]byte(value), &jsonVal); err != nil {
		// Not valid JSON, treat as a plain string
		jsonVal = value
	} else if _, ok := jsonVal.(string); ok {
		// A quoted JSON string like `"hello"` should be stored without quotes
		// but a bare word like `hello` (which failed Unmarshal above) is also a string.
		// This branch handles the quoted case; bare words are handled by the fallback above.
	}

	jf, err := jpath.NewFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	node, parentNode, nodeErr := jf.GetRootNode().GetNodes(JSONpath)
	_ = node
	if nodeErr != nil {
		log.Fatal(nodeErr)
	}
	m, ok := parentNode.CheckMap()
	if !ok {
		log.Fatalf("Parent is not a map: %s", JSONpath)
	}
	m[jpath.LastPart(JSONpath)] = jsonVal

	data, err := jf.JSON()
	if err != nil {
		log.Fatal(err)
	}
	err = jf.Write(data)
	if err != nil {
		log.Fatal(err)
	}
}
