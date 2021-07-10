package main

import (
	"flag"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"strings"

	zglob "github.com/mattn/go-zglob"
)

type LintOption func(*LintImage) error
type OutputFunc func(*LintImage) string

type LintImage struct {
	*image.Config
	FilePath   string
	errors     []error
	formatFunc OutputFunc
}

func (i *LintImage) Error() string {
	return i.formatFunc(i)
}

func OutputText(i *LintImage) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("filename: %s\n", i.FilePath))
	for _, x := range i.errors {
		sb.WriteString(fmt.Sprintf(" %s\n", x))
	}
	return sb.String()
}

func OutputMarkDown(i *LintImage) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("* %s\n", i.FilePath))
	for _, x := range i.errors {
		sb.WriteString(fmt.Sprintf("  * %s\n", x))
	}
	return sb.String()
}

func WithMinWidth(minWidth int) LintOption {
	return func(i *LintImage) error {
		if i.Width < minWidth {
			return fmt.Errorf("min-width: expected(%d<) actual(%d)", minWidth, i.Width)
		}
		return nil
	}
}

func WithMaxWidth(maxWidth int) LintOption {
	return func(i *LintImage) error {
		if maxWidth < i.Width {
			return fmt.Errorf("max-width: expected(<%d) actual(%d)", maxWidth, i.Width)
		}
		return nil
	}
}

func WithMinHeight(minHeight int) LintOption {
	return func(i *LintImage) error {
		if i.Height < minHeight {
			return fmt.Errorf("min-height: expected(%d<) actual(%d)", minHeight, i.Height)
		}
		return nil
	}
}

func WithMaxHeight(maxHeight int) LintOption {
	return func(i *LintImage) error {
		if maxHeight < i.Height {
			return fmt.Errorf("max-width: expected(<%d) actual(%d)", maxHeight, i.Height)
		}
		return nil
	}
}

func NewLintImage(path string, format OutputFunc, opts ...LintOption) error {
	lint := &LintImage{
		FilePath:   path,
		formatFunc: format,
	}
	f, err := os.Open(path)
	if err != nil {
		lint.errors = append(lint.errors, err)
		return lint
	}
	defer f.Close()
	img, _, err := image.DecodeConfig(f)
	if err != nil {
		lint.errors = append(lint.errors, err)
		return lint
	}
	lint.Config = &img
	for _, o := range opts {
		if err := o(lint); err != nil {
			lint.errors = append(lint.errors, err)
		}
	}
	if len(lint.errors) > 0 {
		return lint
	}
	return nil
}

func main() {
	var (
		minWidth, maxWidth, minHeight, maxHeight int
		markdown                                 bool
		outputFormat                             = OutputMarkDown
	)
	flag.IntVar(&minWidth, "min-width", 1, "min width")
	flag.IntVar(&maxWidth, "max-width", 1920, "max width")
	flag.IntVar(&minHeight, "min-height", 1, "min height")
	flag.IntVar(&maxHeight, "max-height", 1080, "max height")
	flag.BoolVar(&markdown, "markdown", true, "show errors in markdown list format")
	flag.Usage = func() {
		fmt.Printf("Usage: %s\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Printf("\n  globs(**) ...\n")
	}
	flag.Parse()
	if len(flag.Args()) < 1 {
		flag.Usage()
		os.Exit(1)
	}
	rc := 0
	if !markdown {
		outputFormat = OutputText
	}
	for _, x := range flag.Args() {
		files, _ := zglob.Glob(x)
		for _, file := range files {
			err := NewLintImage(file, outputFormat,
				WithMinWidth(minWidth),
				WithMaxWidth(maxWidth),
				WithMinHeight(minHeight),
				WithMaxHeight(maxHeight))
			if err != nil {
				fmt.Println(err)
				rc = 1
			}
		}
	}
	os.Exit(rc)
}
