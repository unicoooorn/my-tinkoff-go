package main

import (
	"os"
	"strconv"
	"strings"
)

type BadOffsetError struct {
	Offset   int64
	filesize int64
}

func (e *BadOffsetError) Error() string {
	return "Offset " + strconv.FormatInt(e.Offset, 10) + " is out of file bounds " + strconv.FormatInt(e.filesize, 10)
}

type convFlag struct {
	Uppercase  bool
	Lowercase  bool
	TrimSpaces bool
}

func (opts *Options) ValidateOffset() (err error) {
	var fs os.FileInfo
	if opts.Offset < 0 {
		return &BadOffsetError{opts.Offset, 0}
	} else if opts.From == "" {
		return nil
	} else {
		if fs, err = os.Stat(opts.From); err != nil {
			return err
		}
		if fs.Size() < opts.Offset {
			return &BadOffsetError{opts.Offset, fs.Size()}
		}
	}
	return nil
}

type BadConvError struct {
	Conv string
}

func (e *BadConvError) Error() string {
	return "Cannot perform " + e.Conv + " conversion"
}

func (i *convFlag) String() string {
	var res string
	if i.TrimSpaces {
		res += " trimmed"
	}
	if i.Uppercase {
		res += " uppercase"
	}
	if i.Lowercase {
		res += " lowercase"
	}
	if res != "" {
		return "Text must be" + res
	} else {
		return "No conversions required"
	}
}

func (i *convFlag) Set(value string) error {
	flags := strings.Split(value, ",")
	for _, conv := range flags {
		switch strings.TrimSpace(conv) {
		case "upper_case":
			i.Uppercase = true
		case "lower_case":
			i.Lowercase = true
		case "trim_spaces":
			i.TrimSpaces = true
		default:
			return &BadConvError{value}
		}
	}
	if i.Uppercase && i.Lowercase {
		return &BadConvError{"uppercase and lowercase"}
	}
	return nil
}
