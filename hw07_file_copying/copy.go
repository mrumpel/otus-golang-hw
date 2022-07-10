package main

import (
	"errors"
	"io"
	"os"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
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

	switch limit {
	case 0:
		_, err = io.Copy(toFile, fromFile)
	default:
		_, err = io.CopyN(toFile, fromFile, limit)
	}

	if err != nil && err != io.EOF {
		return err
	}

	return nil
}
