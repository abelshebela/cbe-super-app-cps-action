package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/pdfcpu/pdfcpu/pkg/api"
)

func main() {
	in := flag.String("in", "", "Input PDF path")
	outDir := flag.String("out", "", "Output directory for extracted images")
	flag.Parse()

	if *in == "" {
		log.Fatal("-in is required (path to PDF)")
	}
	if *outDir == "" {
		*outDir = filepath.Join(filepath.Dir(*in), "images")
	}
	if err := os.MkdirAll(*outDir, 0o755); err != nil {
		log.Fatalf("create out dir: %v", err)
	}

	// Extract images from all pages using default configuration
	// selectedPages=nil => all pages; conf=nil => default conf
	err := api.ExtractImagesFile(*in, *outDir, nil, nil)
	if err != nil {
		log.Fatalf("extract images: %v", err)
	}

	fmt.Printf("Images extracted to %s\n", *outDir)
}
