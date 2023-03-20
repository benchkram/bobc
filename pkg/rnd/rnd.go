package rnd

import (
	"crypto/sha1"
	"encoding/hex"
	"math/rand"
	"time"

	"github.com/pkg/errors"
)

var src = rand.NewSource(time.Now().UnixNano())

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

// RandStringBytesMaskImprSrc
// http://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-golang
func RandStringBytesMaskImprSrc(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

func RandSHA1(n int) string {
	hash := sha1.New()
	hash.Write([]byte(RandStringBytesMaskImprSrc(n)))
	return hex.EncodeToString(hash.Sum(nil))
}

func GenerateRandomBytes(n uint) ([]byte, error) {
	b := make([]byte, n)

	_, err := rand.Read(b)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read random bytes")
	}

	return b, nil
}
