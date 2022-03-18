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
		file, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("open output file error: %w", err)
		}

		defer file.Close()
		out = file
	}

	return write(out)
}

func ReadFromFile(file string, read func(reader io.Reader) error) error {
	var in io.Reader
	if file == DefaultStream {
		in = os.Stdin
	} else {
		file, err := os.OpenFile(file, os.O_RDONLY, 0644)
		if err != nil {
			return fmt.Errorf("open input file error: %w", err)
		}
		defer file.Close()
		in = file
	}

	return read(in)
}
