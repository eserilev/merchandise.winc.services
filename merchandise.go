package main

import (
	"flag"

	"github.com/eserilev/merchandise.winc.services/campaigns"
)

func main() {
	filePath := flag.String("file", "", "a file path")
	flag.Parse()
	campaigns.BatchUpload(*filePath)
}
