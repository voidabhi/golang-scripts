package gocard

import (
	"fmt"
	"io"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
)

type CardNumber string

var DELTA = []int{0, 1, 2, 3, 4, -4, -3, -2, -1, 0}

func (cc *CardNumber) IsCard() bool {
	checksum := 0
	bOdd := false
	card := []byte(*cc)
	for i := len(card) - 1; i > -1; i-- {
		cn := int(card[i]) - 48
		checksum += cn

		if bOdd {
			checksum += DELTA[cn]
		}
		bOdd = !bOdd
	}
	if checksum%10 == 0 {
		return true
	}
	return false
}

func (cc *CardNumber) Last4() string {
	return *cc[len(t)-4 : len(t)]
}

func (cc *CardNumber) MD5() []byte {
	h := md5.New()
	io.WriteString(h, *cc)
	return h.Sum(nil)
}

func (cc *CardNumber) SHA1() []byte {
	h := sha1.New()
	io.WriteString(h, *cc)
	return h.Sum(nil)
}

func (cc *CardNumber) SHA256() []byte {
	h := sha256.New()
	io.WriteString(h, *cc)
	return h.Sum(nil)
}

func (cc *CardNumber) SHA512() []byte {
	h := sha512.New()
	io.WriteString(h, *cc)
	return h.Sum(nil)
}
