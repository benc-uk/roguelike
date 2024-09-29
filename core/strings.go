package core

import "math/rand"

// Simple ID generator 8 characters long
func RandId(idLen int) string {
	// generate a random string
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	id := make([]byte, idLen)
	for i := range id {
		id[i] = chars[rand.Intn(len(chars))]
	}
	return string(id)
}

// MakeStr creates a string of n length with s repeated
func MakeStr(n int, s string) string {
	str := ""
	for i := 0; i < n; i++ {
		str += s
	}
	return str
}
