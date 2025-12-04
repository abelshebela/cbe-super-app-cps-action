package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func main() {
	imagesDir := filepath.FromSlash("docs/images")
	ocrDir := filepath.FromSlash("docs/ocr")
	outPath := filepath.FromSlash("docs/lld for bps_full.md")

	if err := assemble(imagesDir, ocrDir, outPath); err != nil {
		log.Fatalf("assemble failed: %v", err)
	}
	fmt.Printf("Wrote %s\n", outPath)
}

func assemble(imagesDir, ocrDir, outPath string) error {
	entries, err := os.ReadDir(imagesDir)
	if err != nil {
		return fmt.Errorf("read images dir: %w", err)
	}
	var files []string
	for _, e := range entries {
		if e.IsDir() { continue }
		name := e.Name()
		low := strings.ToLower(name)
		if strings.HasSuffix(low, ".png") || strings.HasSuffix(low, ".jpg") || strings.HasSuffix(low, ".jpeg") {
			files = append(files, name)
		}
	}
	if len(files) == 0 {
		return fmt.Errorf("no images found in %s", imagesDir)
	}
	// Stable sort: by name
	sort.Strings(files)

	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		return fmt.Errorf("ensure output dir: %w", err)
	}
	f, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("create output: %w", err)
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	defer w.Flush()

	fmt.Fprintln(w, "# LLD for BPS — Extracted with Images, OCR, and PlantUML Placeholders")
	fmt.Fprintln(w, "> Source: docs/lld for bps.pdf\n")

	for idx, name := range files {
		base := name[:len(name)-len(filepath.Ext(name))]
		imgRel := filepath.ToSlash(filepath.Join("docs/images", name))
		ocrPath := filepath.Join(ocrDir, base+".txt")
		ocrText := readTextIfExists(ocrPath)

		fmt.Fprintf(w, "\n---\n\n## Item %d — %s\n\n", idx+1, name)
		fmt.Fprintf(w, "![%s](%s)\n\n", base, imgRel)
		if strings.TrimSpace(ocrText) == "" {
			fmt.Fprintln(w, "_No OCR text extracted._\n")
		} else {
			fmt.Fprintln(w, "### OCR Text\n")
			fmt.Fprintln(w, "````\n"+ocrText+"\n````\n")
		}
		fmt.Fprintln(w, "### PlantUML Placeholder\n")
		fmt.Fprintln(w, "```plantuml")
		fmt.Fprintln(w, "' TODO: Replace with the actual diagram in PlantUML corresponding to this image")
		fmt.Fprintln(w, "' @startuml")
		fmt.Fprintln(w, "' ' ... your diagram here ...")
		fmt.Fprintln(w, "' @enduml")
		fmt.Fprintln(w, "```")
	}
	return nil
}

func readTextIfExists(p string) string {
	b, err := os.ReadFile(p)
	if err != nil { return "" }
	return strings.ReplaceAll(string(b), "\r\n", "\n")
}
