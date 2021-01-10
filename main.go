package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"

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

type result struct {
	Path   string
	Reason []string
}

func main() {
	var st []result
	flag.Parse()

	for _, x := range flag.Args() {
		files, err := zglob.Glob(x)
		if err != nil {
			log.Fatal(err)
		}
		for _, file := range files {
			if err := checkImage(file); err != nil {
				st = append(st, result{Path: file, Reason: err})
			}
		}
	}

	if len(st) > 0 {
		fmt.Fprintln(os.Stderr, "画像チェック中にエラーが発生しました。")
		for _, x := range st {
			fmt.Fprintf(os.Stderr, "%s\n", x.Path)
			for _, v := range x.Reason {
				fmt.Fprintf(os.Stderr, " %s\n", v)
			}
		}
		os.Exit(1)
	}

	os.Exit(0)
}

func checkImage(file string) []string {
	var f []string

	w, h, err := getImageDimension(file)
	if err != nil {
		return nil
	}

	if w < minWidth {
		f = append(f, fmt.Sprintf("横:小 最小：%d > 実際：%d", minWidth, w))
	}
	if maxWidth < w {
		f = append(f, fmt.Sprintf("横:大 最大：%d < 実際：%d", maxWidth, w))
	}
	if h < minHeight {
		f = append(f, fmt.Sprintf("縦:小 最小：%d > 実際：%d", minHeight, h))
	}
	if maxHeight < h {
		f = append(f, fmt.Sprintf("縦:大 最大：%d < 実際：%d", maxHeight, h))
	}
	if len(f) > 0 {
		return f
	}
	return nil
}

func getImageDimension(imagePath string) (int, int, error) {
	file, err := os.Open(imagePath)
	defer file.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
	x, _, err := image.DecodeConfig(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", imagePath, err)
	}
	return x.Width, x.Height, err
}
