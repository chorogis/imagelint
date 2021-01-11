package main

import (
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

type lintError struct {
	FilePath string
	Results  []string
}

func (e *lintError) Error() string {
	return e.ConvertMDList()
}

func (e *lintError) ConvertMDList() string {
	m := []string{e.FilePath, "\n"}
	for _, x := range e.Results {
		m = append(m, fmt.Sprintf("- %s\n", x))
	}
	return strings.Join(m, "")
}

func lintImage(filePath string) error {
	st := &lintError{FilePath: filePath}

	file, err := os.Open(filePath)
	defer file.Close()
	if err != nil {
		st.Results = append(st.Results, err.Error())
		return st
	}
	x, _, err := image.DecodeConfig(file)
	if err != nil {
		st.Results = append(st.Results, err.Error())
		return st
	}

	if x.Width < minWidth {
		st.Results = append(st.Results,
			fmt.Sprintf("min-width: expected(%d<) actual(%d)", minWidth, x.Width))
	}
	if maxWidth < x.Width {
		st.Results = append(st.Results,
			fmt.Sprintf("max-width: expected(<%d) actual(%d)", maxWidth, x.Width))
	}
	if x.Height < minHeight {
		st.Results = append(st.Results,
			fmt.Sprintf("min-height: expected(%d<) actual(%d)", minHeight, x.Height))
	}
	if maxHeight < x.Height {
		st.Results = append(st.Results,
			fmt.Sprintf("max-height: expected(<%d) actual(%d)", maxHeight, x.Height))
	}
	if len(st.Results) > 0 {
		return st
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
				fmt.Println(err)
				retStatus = 1
			}
		}
	}
	os.Exit(retStatus)
}
