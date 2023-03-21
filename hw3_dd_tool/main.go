package main

import (
	"flag"
	"fmt"
	"io"
	"math"
	"os"
)

type Options struct {
	From      string
	To        string
	Offset    int64
	Limit     int64
	BlockSize int64
	Conv      convFlag
}

func checkFatalError(err error) {
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func ParseFlags() (*Options, error) {
	var opts Options
	flag.StringVar(&opts.From, "from", "", "file to read. by default - stdin")
	flag.StringVar(&opts.To, "to", "", "file to write. by default - stdout")
	flag.Int64Var(&opts.Offset, "offset", 0, "number of bytes to skip when copy")
	flag.Int64Var(&opts.Limit, "limit", math.MaxInt64, "number of bytes to copy")
	flag.Int64Var(&opts.BlockSize, "block-size", 1000, "size of copy block")
	flag.Var(&opts.Conv, "conv", "conversion to be performed when copy")
	flag.Parse()
	if err := opts.ValidateOffset(); err != nil {
		return nil, err
	}
	return &opts, nil
}

func main() {
	opts, err := ParseFlags()
	checkFatalError(err)
	reader, err := NewAdvancedReader(opts.From, opts.Offset, opts.Limit)
	checkFatalError(err)
	defer reader.Close()
	writer, err := NewAdvancedWriter(opts.To, opts.Conv)
	checkFatalError(err)
	defer writer.Close()

	result := make([]byte, 0)
	fixedBuf := make([]byte, opts.BlockSize)
	err = reader.SkipOffset(fixedBuf)
	checkFatalError(err)
	for {
		bytesCount, err := reader.Read(fixedBuf)
		buf := make([]byte, bytesCount)
		copy(buf, fixedBuf)
		result = append(result, buf...)
		if err == io.EOF {
			break
		} else if err != nil {
			checkFatalError(err)
		} else if bytesCount == 0 {
			break
		}
	}
	_, err = writer.Write(result)
	checkFatalError(err)
}
