package main

import (
	"flag"
	"fmt"
	"github.com/xyproto/jpath"
	"log"
	"os"
)

func main() {
	flag.Parse()

	if len(flag.Args()) != 2 {
		fmt.Println("Syntax: jget [filename] [JSON path]")
		fmt.Println("Example: jget books.json x[1].author")
		os.Exit(1)
	}

	filename := flag.Args()[0]
	JSONpath := flag.Args()[1]

	jf, err := jpath.NewFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	node, err := jf.GetNode(JSONpath)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(node.StringValue())
}
