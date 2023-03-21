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
	Offset int64
}

func NewAdvancedReader(from string, offset int64, limit int64) (reader AdvancedReader, err error) {
	if from == "" {
		reader.file = os.Stdin
	} else {
		reader.file, err = os.Open(from)
	}
	if limit == math.MaxInt64 {
		reader.reader = &io.LimitedReader{R: reader.file, N: limit}
	} else {
		reader.reader = &io.LimitedReader{R: reader.file, N: limit + offset}
	}
	reader.Offset = offset
	return
}

func (r AdvancedReader) SkipOffset(block []byte) (err error) {
	offsetSkipper := io.LimitReader(r.reader, r.Offset)
	var bytesReadCounter int
	for r.Offset > 0 {
		bytesReadCounter, err = offsetSkipper.Read(block)
		if err != nil && err != io.EOF {
			return
		} else if err == io.EOF && r.Offset > 0 {
			return &BadOffsetError{r.Offset, -1}
		}
		r.Offset -= int64(bytesReadCounter)
	}
	return
}

func (r AdvancedReader) Read(block []byte) (bytesRead int, err error) {
	bytesRead, err = r.reader.Read(block)
	return
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
