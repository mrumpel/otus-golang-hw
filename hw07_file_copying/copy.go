package main

import (
	"errors"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	// Files + checks
	fromFile, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer fromFile.Close()

	fromFileStat, err := fromFile.Stat()
	if err != nil {
		return err
	}
	if offset > fromFileStat.Size() {
		return ErrOffsetExceedsFileSize
	}
	if !fromFileStat.Mode().IsRegular() {
		return ErrUnsupportedFile
	}

	_, err = fromFile.Seek(offset, io.SeekStart)
	if err != nil {
		return err
	}

	toFile, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer toFile.Close()

	// Progress bar init
	expectedLen := fromFileStat.Size() - offset
	if expectedLen > limit {
		expectedLen = limit
	}

	bar := pb.Simple.Start64(expectedLen)
	defer bar.Finish()

	barWriter := bar.NewProxyWriter(toFile)

	// Main copy logic
	switch limit {
	case 0:
		_, err = io.Copy(barWriter, fromFile)
	default:
		_, err = io.CopyN(barWriter, fromFile, limit)
	}

	if err != nil && err != io.EOF {
		return err
	}

	return nil
}
