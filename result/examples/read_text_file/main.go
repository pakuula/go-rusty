package main

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/pakuula/go-rusty/result"
	"github.com/pakuula/go-rusty/result/io_r"
)

func ReadTextFile(fname string) (string, error) {
	f, err := os.Open(fname)
	if err != nil {
		return "", err
	}
	defer f.Close()

	buf := bytes.Buffer{}
	_, err = io.Copy(&buf, f)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func ReadTextFileR(fname string) (res result.Result[string]) {
	defer result.Catch(&res)
	f := result.Wrap(os.Open(fname)).Mustf("failed to open file: %s", fname)
	defer f.Close()

	buf := bytes.Buffer{}
	io_r.Copy(&buf, f).Mustf("failed to read into the buffer")

	return result.Val(buf.String())
}

func main() {
	if len(os.Args) == 1 {
		fmt.Fprintln(os.Stderr,
			"Expected at least one file name",
		)
		os.Exit(1)
	}
	for _, fname := range os.Args[1:] {
		text := ReadTextFileR(fname)
		if text.IsError() {
			fmt.Fprintf(os.Stderr, "%s: %s\n", fname, text.Err().Error())
			continue
		}
		fmt.Print(text.Unwrap())
	}
}
