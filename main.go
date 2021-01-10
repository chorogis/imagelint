package main

import (
	"errors"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"strings"

	zglob "github.com/mattn/go-zglob"
	flag "github.com/spf13/pflag"
)

var (
	minWidth  int
	maxWidth  int
	minHeight int
	maxHeight int
)

func init() {
	flag.IntVar(&minWidth, "min-width", 10, "横の最小サイズ")
	flag.IntVar(&maxWidth, "max-width", 800, "横の最大サイズ")
	flag.IntVar(&minHeight, "min-height", 10, "縦の最小サイズ")
	flag.IntVar(&maxHeight, "max-height", 2560, "縦の最大サイズ")
}

func lintImage(filePath string) error {
	var st []string

	file, err := os.Open(filePath)
	defer file.Close()
	if err != nil {
		return err
	}
	x, _, err := image.DecodeConfig(file)
	if err != nil {
		return err
	}

	if x.Width < minWidth {
		st = append(st, fmt.Sprintf("min-width: expected(%d<) actual(%d)", minWidth, x.Width))
	}
	if maxWidth < x.Width {
		st = append(st, fmt.Sprintf("max-width: expected(<%d) actual(%d)", maxWidth, x.Width))
	}
	if x.Height < minHeight {
		st = append(st, fmt.Sprintf("min-height: expected(%d<) actual(%d)", minHeight, x.Height))
	}
	if maxHeight < x.Height {
		st = append(st, fmt.Sprintf("max-height: expected(<%d) actual(%d)", maxHeight, x.Height))
	}
	if len(st) > 0 {
		return errors.New(strings.Join(st, ", "))
	}
	return nil
}

func main() {
	flag.Parse()
	retStatus := 0

	for _, x := range flag.Args() {
		files, _ := zglob.Glob(x)
		for _, file := range files {
			if err := lintImage(file); err != nil {
				fmt.Printf("%s: %s\n", file, err)
				retStatus = 1
			}
		}
	}
	os.Exit(retStatus)
}
