package main

import (
	"crypto/sha256"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	tests := []struct {
		name         string
		in           string
		out          string
		offset       int64
		limit        int64
		error        error
		untypedError bool
	}{
		{
			name:   "out_offset0_limit0",
			in:     "testdata/input.txt",
			out:    "testdata/out_offset0_limit0.txt",
			offset: 0,
			limit:  0,
			error:  nil,
		},
		{
			name:   "out_offset0_limit10",
			in:     "testdata/input.txt",
			out:    "testdata/out_offset0_limit10.txt",
			offset: 0,
			limit:  10,
			error:  nil,
		},
		{
			name:   "out_offset0_limit1000",
			in:     "testdata/input.txt",
			out:    "testdata/out_offset0_limit1000.txt",
			offset: 0,
			limit:  1000,
			error:  nil,
		},
		{
			name:   "out_offset0_limit10000",
			in:     "testdata/input.txt",
			out:    "testdata/out_offset0_limit10000.txt",
			offset: 0,
			limit:  10000,
			error:  nil,
		},
		{
			name:   "out_offset100_limit1000",
			in:     "testdata/input.txt",
			out:    "testdata/out_offset100_limit1000.txt",
			offset: 100,
			limit:  1000,
			error:  nil,
		},
		{
			name:   "out_offset6000_limit1000",
			in:     "testdata/input.txt",
			out:    "testdata/out_offset6000_limit1000.txt",
			offset: 6000,
			limit:  1000,
			error:  nil,
		},
		{
			name: "empty file",
			in:   "testdata/empty.txt",
			out:  "testdata/empty.txt",
		},
		{
			name:   "offset exceed filesize",
			in:     "testdata/input.txt",
			offset: 10000,
			error:  ErrOffsetExceedsFileSize,
		},
		{
			name:  "unsupported (irregular) file",
			in:    "/dev/urandom",
			error: ErrUnsupportedFile,
		},
		{
			name:         "path error",
			in:           "file_not_exist_here",
			untypedError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, _ := os.CreateTemp("", "test_res")
			defer os.Remove(res.Name())
			defer res.Close()

			err := Copy(test.in, res.Name(), test.offset, test.limit)

			if test.untypedError {
				require.Error(t, err)
				return
			}
			if err != nil {
				require.EqualError(t, err, test.error.Error())
			} else {
				checkEqual(t, test.out, res.Name())
			}
		})
	}
}

func checkEqual(t *testing.T, f1, f2 string) {
	t.Helper()

	w1 := sha256.New()
	w2 := sha256.New()

	r1, err := os.Open(f1)
	if err != nil {
		t.Fail()
	}
	defer r1.Close()
	r2, err := os.Open(f2)
	if err != nil {
		t.Fail()
	}
	defer r2.Close()

	_, err = io.Copy(w1, r1)
	if err != nil {
		t.Fail()
	}

	_, err = io.Copy(w2, r2)
	if err != nil {
		t.Fail()
	}

	require.Equal(t, w1.Sum(nil), w2.Sum(nil))
}
