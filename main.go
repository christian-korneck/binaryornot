package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func Usage() {
	exe := filepath.Base(os.Args[0])
	fmt.Printf("Usage: %s [-s] [file] \n", exe)
	flag.CommandLine.SetOutput(os.Stdout)
	flag.PrintDefaults()
	flag.CommandLine.SetOutput(os.Stderr)
	fmt.Printf("\nDetects if an input from file or stdin is text or binary. Exit codes: 0: text; 1: error; 5: binary.\n\nExamples:\n%s file.txt\necho hello | %s -s \n%s -s < test2.dump\n\n", exe, exe, exe)
}

func main() {

	var path string
	var input io.Reader

	flag.Usage = Usage
	FlagStdin := flag.Bool("s", false, "use stdin (instead of file)")
	flag.Parse()
	UseStdin := *FlagStdin

	switch UseStdin {
	case true:
		input = os.Stdin
		path = "<STDIN>"

	case false:
		if flag.NArg() != 1 {
			flag.Usage()
			os.Exit(1)
		}
		path = flag.Args()[0]
		file, err := os.Open(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "can't open file %s\n", path)
			os.Exit(1)
		}
		defer file.Close()
		input = file
	}

	// read up to 512 bytes
	sample, err := io.ReadAll(io.LimitReader(input, 512))
	if err != nil {
		fmt.Fprintf(os.Stderr, "can't read from file %s\n", path)
		os.Exit(1)
	}

	// treat empty file as text, but no data via stdin as error
	if UseStdin && len(sample) == 0 {
		fmt.Fprintf(os.Stderr, "no data received from %s\n", path)
		os.Exit(1)
	}

	contentType := http.DetectContentType(sample)

	contentTypeSplit := strings.Split(contentType, "/")
	if len(contentTypeSplit) < 2 {
		fmt.Fprintf(os.Stderr, "could not determine content-type of file %s\n", path)
		os.Exit(1)
	}

	switch strings.ToLower(contentTypeSplit[0]) {
	case "text":
		fmt.Println("text")
		os.Exit(0)
	default:
		fmt.Println("binary")
		os.Exit(5)
	}

	os.Exit(1)

}
