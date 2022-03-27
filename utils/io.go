package utils

import (
	"fmt"
	"io"
	"os"
)

const DefaultStream = "-"

func WriteToFile(file string, write func(writer io.Writer) error) error {
	var out io.Writer
	if file == DefaultStream {
		out = os.Stdout
	} else {
		outFile, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			return fmt.Errorf("open output file error: %w", err)
		}

		defer outFile.Close()
		out = outFile
	}

	return write(out)
}

func ReadFromFile(file string, read func(reader io.Reader) error) error {
	var in io.Reader
	if file == DefaultStream {
		in = os.Stdin
	} else {
		inFile, err := os.OpenFile(file, os.O_RDONLY, 0644)
		if err != nil {
			return fmt.Errorf("open input file error: %w", err)
		}
		defer inFile.Close()
		in = inFile
	}

	return read(in)
}
