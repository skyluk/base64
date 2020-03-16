/*
Base64 encoder/decoder.
Copyright (C) 2020  Ross Fischer

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/
package main

import (
	"fmt"
	"os"
	"strings"
)

var translation [64]string
var rtranslation map[byte]byte
var args []string

//
// usage:
//
// base64 [encode|decode] [input]
//
func main() {

	args = os.Args[1:]

	if len(args) < 2 {
		fmt.Println("too few arguments")
		os.Exit(1)
	}

	switch strings.ToLower(args[0]) {
	case "encode":

		if len(args[0]) == 0 {
			fmt.Println("invalid input")
			os.Exit(1)
		}

		out, err := encode(args[1])
		if err != nil {
			panic(err)
		}

		fmt.Print(out)

	case "-d":
		fallthrough
	case "decode":

		if len(args[0]) == 0 {
			fmt.Println("invalid input")
			os.Exit(1)
		}

		out, err := decode(args[1])
		if err != nil {
			panic(err)
		}

		fmt.Print(out + "\n")

	default:
		fmt.Printf("invalid command")
		os.Exit(1)
	}

}

func encode(in string) (string, error) {

	var out strings.Builder
	var t, b, m byte

	n := uint8(2)

	for i := 0; i < len(in); i++ {
		t = in[i]

		b = t >> n
		b = b | m

		out.WriteString(e_lookup(b))

		m = get_emask(n, t)

		if n == 6 {
			// reached the end of a group of characters
			b = t & m
			out.WriteString(e_lookup(b))

			// reset the mask
			m = 0x00
		}

		// n can have the values 2, 4, or 6
		n = (n % 6) + 2
	}

	if n != 2 {
		out.WriteString(e_lookup(m))
	}

	out.WriteString(get_pad(in))
	return out.String(), nil
}

func decode(in string) (string, error) {

	var out strings.Builder
	o := make([]byte, 0, len(in))
	var t, b, m byte

	n := uint8(2)
	l := len(in)

	b = d_lookup(in[0]) << 2

	for i := 1; i < l; i++ {
		if '=' == in[i] {
			break // done!
		}

		t = d_lookup(in[i]) << 2

		m = get_dmask(n, t)
		b = b | m

		o = append(o, b)

		b = t << n

		if n == 6 {
			b = d_lookup(in[i+1]) << 2
			i++
		}

		// n can have the values 2, 4, or 6
		n = (n % 6) + 2
	}

	_, err := out.Write(o)
	if err != nil {
		return "", err
	}

	return out.String(), nil
}

//
func get_emask(n uint8, t byte) byte {

	var m byte

	switch n {
	case 2:
		m = t & 0x03
		m = m << 4
	case 4:
		m = t & 0x0f
		m = m << 2
	case 6:
		m = t & 0x3f
	}

	return m
}

//
func get_dmask(n uint8, t byte) byte {

	var m byte

	switch n {
	case 2:
		m = t & 0xd0
		m = m >> 6
	case 4:
		m = t & 0xf0
		m = m >> 4
	case 6:
		m = t & 0xfd
		m = m >> 2
	}

	return m
}

//
func get_pad(in string) string {

	n := len(in)

	if n == 0 || n%3 == 0 {
		return ""
	}

	return strings.Repeat("=", 3-(n%3))
}

func e_lookup(in byte) string {
	return translation[in]
}

func d_lookup(in byte) byte {
	return rtranslation[in]
}

func init() {

	translation = [...]string{
		"A", "B", "C", "D", "E", "F", "G", "H",
		"I", "J", "K", "L", "M", "N", "O", "P",
		"Q", "R", "S", "T", "U", "V", "W", "X",
		"Y", "Z", "a", "b", "c", "d", "e", "f",
		"g", "h", "i", "j", "k", "l", "m", "n",
		"o", "p", "q", "r", "s", "t", "u", "v",
		"w", "x", "y", "z", "0", "1", "2", "3",
		"4", "5", "6", "7", "8", "9", "+", "/",
	}

	rtranslation = map[byte]byte{
		'A': 0, 'B': 1, 'C': 2, 'D': 3, 'E': 4, 'F': 5, 'G': 6, 'H': 7,
		'I': 8, 'J': 9, 'K': 10, 'L': 11, 'M': 12, 'N': 13, 'O': 14, 'P': 15,
		'Q': 16, 'R': 17, 'S': 18, 'T': 19, 'U': 20, 'V': 21, 'W': 22, 'X': 23,
		'Y': 24, 'Z': 25, 'a': 26, 'b': 27, 'c': 28, 'd': 29, 'e': 30, 'f': 31,
		'g': 32, 'h': 33, 'i': 34, 'j': 35, 'k': 36, 'l': 37, 'm': 38, 'n': 39,
		'o': 40, 'p': 41, 'q': 42, 'r': 43, 's': 44, 't': 45, 'u': 46, 'v': 47,
		'w': 48, 'x': 49, 'y': 50, 'z': 51, '0': 52, '1': 53, '2': 54, '3': 55,
		'4': 56, '5': 57, '6': 58, '7': 59, '8': 60, '9': 61, '+': 62, '/': 63,
	}
}
