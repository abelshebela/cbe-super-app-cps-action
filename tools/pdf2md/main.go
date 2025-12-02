package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	pdf "github.com/ledongthuc/pdf"
)

func main() {
	in := flag.String("in", "", "Input PDF path")
	out := flag.String("out", "", "Output Markdown path")
	flag.Parse()

	if *in == "" {
		log.Fatal("-in is required (path to PDF)")
	}

	// Default output path next to input as .md if not provided
	if *out == "" {
		base := filepath.Base(*in)
		ext := filepath.Ext(base)
		name := strings.TrimSuffix(base, ext)
		*out = filepath.Join(filepath.Dir(*in), name+".md")
	}

	if err := extractToMarkdown(*in, *out); err != nil {
		log.Fatalf("conversion failed: %v", err)
	}

	fmt.Printf("Wrote Markdown to %s\n", *out)
}

func extractToMarkdown(inPath, outPath string) error {
	f, r, err := pdf.Open(inPath)
	if err != nil {
		return fmt.Errorf("open pdf: %w", err)
	}
	defer f.Close()

	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		return fmt.Errorf("ensure output dir: %w", err)
	}

	outFile, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("create output: %w", err)
	}
	defer outFile.Close()

	w := bufio.NewWriter(outFile)
	defer w.Flush()

	// Header
	fmt.Fprintln(w, "# Extracted Document")
	fmt.Fprintf(w, "> Source: %s\n\n", filepath.Base(inPath))

	// Iterate pages and write as Markdown paragraphs
	total := r.NumPage()
	for i := 1; i <= total; i++ {
		p := r.Page(i)
		if p.V.IsNull() {
			continue
		}
		content, err := p.GetPlainText(nil)
		if err != nil {
			return fmt.Errorf("read page %d: %w", i, err)
		}

		// Page separator
		fmt.Fprintf(w, "\n---\n\n## Page %d\n\n", i)

		text := normalizeText(content)
		if strings.TrimSpace(text) == "" {
			fmt.Fprintln(w, "_No extractable text on this page (likely image/diagram)._\n")
			// Insert PlantUML placeholder for diagrams when no text
			fmt.Fprintln(w, "```plantuml")
			fmt.Fprintln(w, "' TODO: Replace this placeholder with the actual diagram in PlantUML")
			fmt.Fprintln(w, "' Example:")
			fmt.Fprintln(w, "' @startuml")
			fmt.Fprintln(w, "' actor User")
			fmt.Fprintln(w, "' User -> System: Action")
			fmt.Fprintln(w, "' @enduml")
			fmt.Fprintln(w, "````\n")
			continue
		}

		fmt.Fprintln(w, text)
	}

	return nil
}

func normalizeText(s string) string {
	// Convert Windows newlines, collapse excessive blank lines, trim trailing spaces per line
	s = strings.ReplaceAll(s, "\r\n", "\n")
	lines := strings.Split(s, "\n")
	out := make([]string, 0, len(lines))
	blank := 0
	for _, ln := range lines {
		ln = strings.TrimRight(ln, " \t")
		if strings.TrimSpace(ln) == "" {
			blank++
			if blank > 1 {
				continue
			}
			out = append(out, "")
			continue
		}
		blank = 0
		out = append(out, ln)
	}
	return strings.Join(out, "\n")
}
