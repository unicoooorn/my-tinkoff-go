package main

import (
	"errors"
	"io"
	"math"
	"os"
	"strings"
)

type AlreadyExistsError struct {
	to string
}

func (err AlreadyExistsError) Error() string {
	return err.to + " already exists"
}

type AdvancedReader struct {
	file   *os.File
	reader *io.LimitedReader
	Off    int64
}

func NewAdvancedReader(from string, offset int64, lim int64) (reader AdvancedReader, err error) {
	if from == "" {
		reader.file = os.Stdin
	} else {
		reader.file, err = os.Open(from)
	}
	if lim == math.MaxInt64 {
		reader.reader = &io.LimitedReader{R: reader.file, N: lim}
	} else {
		reader.reader = &io.LimitedReader{R: reader.file, N: lim + offset}
	}
	reader.Off = offset
	return
}

func (r AdvancedReader) SkipOffset(block []byte) (err error) {
	offsetSkipper := io.LimitReader(r.reader, r.Off)
	var skipCount int
	for r.Off > 0 {
		skipCount, err = offsetSkipper.Read(block)
		if err != nil && err != io.EOF {
			return
		} else if err == io.EOF && r.Off > 0 {
			return &BadOffsetError{r.Off, -1}
		}
		r.Off -= int64(skipCount)
	}
	return
}

func (r AdvancedReader) Read(block []byte) (int, error) {
	return r.reader.Read(block)
}

func (r AdvancedReader) Close() error {
	return r.file.Close()
}

type AdvancedWriter struct {
	file *os.File
	conv convFlag
}

func NewAdvancedWriter(to string, conv convFlag) (reader AdvancedWriter, err error) {
	if to == "" {
		reader.file = os.Stdout
	} else {
		if _, err = os.Stat(to); !errors.Is(err, os.ErrNotExist) {
			err = &AlreadyExistsError{}
			return
		}
		reader.file, err = os.OpenFile(to, os.O_RDWR|os.O_CREATE, 0755)
	}
	reader.conv = conv
	return
}

func (r AdvancedWriter) performConv(p []byte) []byte {
	buf := string(p)
	if r.conv.TrimSpaces {
		buf = strings.TrimSpace(buf)
	}
	if r.conv.Uppercase {
		buf = strings.ToTitle(buf)
	}
	if r.conv.Lowercase {
		buf = strings.ToLower(buf)
	}
	return []byte(buf)
}

func (r AdvancedWriter) Write(p []byte) (n int, err error) {
	p = r.performConv(p)
	n, err = r.file.Write(p)
	return
}

func (r AdvancedWriter) Close() error {
	return r.file.Close()
}
