/**
 * adapted from https://github.com/peak/s3hash/blob/master/cmd/s3hash/main.go
 */

package main

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"strconv"
)

// Calculate calculates the S3 hash of a given io.ReadSeeker with the given chunk size in bytes.
func calculate(f io.ReadSeeker, chunkSize int64) (string, error) {
	dataSize, err := f.Seek(0, io.SeekEnd)
	if err != nil {
		return "", err
	}

	var (
		sumOfSums []byte
		parts     int
	)
	for i := int64(0); i < dataSize; i += chunkSize {
		length := chunkSize
		if i+chunkSize > dataSize {
			length = dataSize - i
		}
		sum, err := md5sum(f, i, length)
		if err != nil {
			return "", err
		}
		sumOfSums = append(sumOfSums, sum...)
		parts++
	}

	var finalSum []byte

	if parts == 1 {
		finalSum = sumOfSums
	} else {
		h := md5.New()
		_, err := h.Write(sumOfSums)
		if err != nil {
			return "", err
		}
		finalSum = h.Sum(nil)
	}

	sumHex := hex.EncodeToString(finalSum)

	if parts > 1 {
		sumHex += "-" + strconv.Itoa(parts)
	}

	return sumHex, nil
}

func md5sum(r io.ReadSeeker, start, length int64) ([]byte, error) {
	r.Seek(start, io.SeekStart)
	h := md5.New()
	if _, err := io.CopyN(h, r, length); err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}
