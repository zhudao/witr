//go:build linux || darwin || freebsd

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/pranshuparmar/witr/internal/app"
	"github.com/spf13/cobra/doc"
)

func main() {
	out := flag.String("out", "./docs/cli", "output directory")
	format := flag.String("format", "markdown", "markdown|man|rest")
	front := flag.Bool("frontmatter", false, "prepend simple YAML front matter to markdown")
	flag.Parse()

	if err := os.MkdirAll(*out, 0o755); err != nil {
		log.Fatal(err)
	}

	root := app.Root()
	root.DisableAutoGenTag = true // stable, reproducible files (no timestamp footer)
	root.InitDefaultVersionFlag() // cobra only adds -v/--version when the command runs

	switch *format {
	case "markdown":
		if *front {
			prep := func(filename string) string {
				base := filepath.Base(filename)
				name := strings.TrimSuffix(base, filepath.Ext(base))
				title := strings.ReplaceAll(name, "_", " ")
				return fmt.Sprintf("---\ntitle: %q\nslug: %q\ndescription: \"CLI reference for %s\"\n---\n\n", title, name, title)
			}
			link := func(name string) string { return strings.ToLower(name) }
			if err := doc.GenMarkdownTreeCustom(root, *out, prep, link); err != nil {
				log.Fatal(err)
			}
		} else {
			if err := doc.GenMarkdownTree(root, *out); err != nil {
				log.Fatal(err)
			}
		}
		var md strings.Builder
		md.WriteString("\n### Exit codes\n\n| Code | Meaning |\n|------|---------|\n")
		for _, c := range exitCodes {
			fmt.Fprintf(&md, "| %s | %s |\n", c[0], c[1])
		}
		md.WriteString("\n" + severityNote + "\n")
		insertSection(filepath.Join(*out, root.Name()+".md"), "", md.String())
	case "man":
		hdr := &doc.GenManHeader{Title: strings.ToUpper(root.Name()), Section: "1"}
		if err := doc.GenManTree(root, hdr, *out); err != nil {
			log.Fatal(err)
		}
		var man strings.Builder
		man.WriteString(".SH EXIT STATUS\n")
		for _, c := range exitCodes {
			fmt.Fprintf(&man, ".TP\n\\fB%s\\fP\n%s\n", c[0], c[1])
		}
		man.WriteString(".PP\n" + severityNote + "\n\n\n")
		insertSection(filepath.Join(*out, root.Name()+".1"), ".SH EXAMPLE", man.String())
	case "rest":
		if err := doc.GenReSTTree(root, *out); err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatalf("unknown format: %s", *format)
	}
}

// exitCodes mirrors the exit codes in internal/app and the README's Exit Codes
// table; cobra has no exit-status section of its own.
var exitCodes = [][2]string{
	{"0", "Clean: process found, no warnings"},
	{"1", "Warnings: process found but has one or more warnings"},
	{"2", "Not found: no matching process or service"},
	{"3", "Permission denied: insufficient privileges (retry with sudo)"},
	{"4", "Invalid input: bad arguments or ambiguous match"},
	{"5", "Internal error: an unexpected failure occurred"},
	{"6", "Cause unknown: process found, but what started it can't be traced (not used on Windows)"},
}

const severityNote = "With several targets, the most severe result wins. Cause unknown ranks above warnings but below every failure (2 to 5)."

// insertSection adds text to a generated page just before marker, or at the
// end when marker is empty or missing.
func insertSection(path, marker, text string) {
	b, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	s := string(b)
	if i := strings.Index(s, marker); marker != "" && i >= 0 {
		s = s[:i] + text + s[i:]
	} else {
		s += text
	}
	if err := os.WriteFile(path, []byte(s), 0o644); err != nil {
		log.Fatal(err)
	}
}
