package util

import (
	"crypto/rand"
)

// StdChars characters for generating the CSRF token
var StdChars = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")

// NewCsrfToken create a new CSRF token
func NewCsrfToken(length int) string {
	if length == 0 {
		return ""
	}
	clen := len(StdChars)
	if clen < 2 || clen > 256 {
		panic("uniuri: wrong charset length for NewLenChars")
	}
	maxrb := 255 - (256 % clen)
	b := make([]byte, length)
	r := make([]byte, length+(length/4)) // storage for random bytes.
	i := 0
	for {
		if _, err := rand.Read(r); err != nil {
			panic("uniuri: error reading random bytes: " + err.Error())
		}
		for _, rb := range r {
			c := int(rb)
			if c > maxrb {
				// Skip this number to avoid modulo bias.
				continue
			}
			b[i] = StdChars[c%clen]
			i++
			if i == length {
				return string(b)
			}
		}
	}
}
