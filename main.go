package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
)

const (
	VERSION      = "1.2.0"
	NAME         = "gotags"
	URL          = "https://github.com/jstemmer/gotags"
	AUTHOR_NAME  = "Joel Stemmer"
	AUTHOR_EMAIL = "stemmertech@gmail.com"
)

var (
	printVersion bool
	inputFile    string
	sortOutput   bool
	silent       bool
	printTree    bool // for debugging
	stdin        bool
)

// Initialize flags.
func init() {
	flag.BoolVar(&printVersion, "v", false, "print version")
	flag.StringVar(&inputFile, "L", "", "source file names are read from the specified file.")
	flag.BoolVar(&sortOutput, "sort", true, "sort tags")
	flag.BoolVar(&silent, "silent", false, "do not produce any output on error")
	flag.BoolVar(&printTree, "tree", false, "print syntax tree (debugging)")
	flag.BoolVar(&stdin, "stdin", false, "read source from stdin")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "gotags version %s\n\n", VERSION)
		fmt.Fprintf(os.Stderr, "Usage: %s [options] file(s)\n\n", os.Args[0])
		flag.PrintDefaults()
	}
}

func getFileNames() ([]string, error) {
	var names []string

	names = append(names, flag.Args()...)

	if len(inputFile) == 0 {
		return names, nil
	}

	var scanner *bufio.Scanner
	if inputFile != "-" {
		in, err := os.Open(inputFile)
		if err != nil {
			return nil, err
		}

		defer in.Close()
		scanner = bufio.NewScanner(in)
	} else {
		scanner = bufio.NewScanner(os.Stdin)
	}

	for scanner.Scan() {
		names = append(names, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return names, nil
}

func main() {
	flag.Parse()

	if printVersion {
		fmt.Printf("gotags version %s\n", VERSION)
		return
	}

	tags := []Tag{}

	if stdin {

		fmt.Println("reading from stdin...")

		ts, err := Parse("-")
		if err != nil {
			if !silent {
				fmt.Fprintf(os.Stderr, "parse error: %s\n\n", err)
			}
		}
		tags = append(tags, ts...)

	} else {

		files, err := getFileNames()
		if err != nil {
			fmt.Fprintf(os.Stderr, "cannot get specified files\n\n")
			flag.Usage()
			os.Exit(1)
		}

		if len(files) == 0 {
			fmt.Fprintf(os.Stderr, "no file specified\n\n")
			flag.Usage()
			os.Exit(1)
		}

		if printTree {
			PrintTree(flag.Arg(0))
			return
		}

		for _, file := range files {
			ts, err := Parse(file)
			if err != nil {
				if !silent {
					fmt.Fprintf(os.Stderr, "parse error: %s\n\n", err)
				}
				continue
			}
			tags = append(tags, ts...)
		}
	}

	output := createMetaTags()
	for _, tag := range tags {
		output = append(output, tag.String())
	}

	if sortOutput {
		sort.Sort(sort.StringSlice(output))
	}

	for _, s := range output {
		fmt.Println(s)
	}
}

// createMetaTags returns a list of meta tags.
func createMetaTags() []string {
	var sorted int
	if sortOutput {
		sorted = 1
	}
	return []string{
		"!_TAG_FILE_FORMAT\t2\t",
		fmt.Sprintf("!_TAG_FILE_SORTED\t%d\t/0=unsorted, 1=sorted/", sorted),
		fmt.Sprintf("!_TAG_PROGRAM_AUTHOR\t%s\t/%s/", AUTHOR_NAME, AUTHOR_EMAIL),
		fmt.Sprintf("!_TAG_PROGRAM_NAME\t%s\t", NAME),
		fmt.Sprintf("!_TAG_PROGRAM_URL\t%s\t", URL),
		fmt.Sprintf("!_TAG_PROGRAM_VERSION\t%s\t", VERSION),
	}
}
