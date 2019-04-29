package main

import (
	"errors"
	"strings"
)

func Encode(b byte) (byte, error) {
	switch {
	case 0 <= b && b <= 25:
		return 'A' + b, nil
	case 26 <= b && b <= 51:
		return 'a' + b - 26, nil
	case 52 <= b && b <= 61:
		return '0' + b - 52, nil
	case b == 62:
		return '+', nil
	case b == 63:
		return '/', nil
	default:
		return 0, errors.New("illegal argument")
	}
}

func Decode(b byte) (byte, error) {
	switch {
	case 'A' <= b && b <= 'Z':
		return b - 'A', nil
	case 'a' <= b && b <= 'z':
		return b - 'a' + 26, nil
	case '0' <= b && b <= '9':
		return b - '0' + 52, nil
	case b == '+':
		return 62, nil
	case b == '/':
		return 63, nil
	default:
		return 0, errors.New("illegal argument")
	}
}

func Base64encode(input string) (string, error) {
	data := []byte(input)
	var sb strings.Builder

	for i := 0; i < len(data)/3*3; i += 3 {
		b1 := data[i] >> 2
		b2 := data[i]&0x03<<4 | data[i+1]&0xf0>>4
		b3 := data[i+1]&0x0f<<2 | data[i+2]&0xc0>>6
		b4 := data[i+2] & 0x3f

		for _, b := range []byte{b1, b2, b3, b4} {
			c, err := Encode(b)
			if err != nil {
				return "", err
			}
			sb.WriteByte(c)
		}
	}

	if len(data)%3 == 1 {
		i := len(data) / 3 * 3
		b1 := data[i] >> 2
		b2 := data[i] & 0x03 << 4
		for _, b := range []byte{b1, b2} {
			c, err := Encode(b)
			if err != nil {
				return "", err
			}
			sb.WriteByte(c)
		}
		sb.WriteString("==")
	} else if len(data)%3 == 2 {
		i := len(data) / 3 * 3
		b1 := data[i] >> 2
		b2 := data[i]&0x03<<4 | data[i+1]&0xf0>>4
		b3 := data[i+1] & 0x0f << 2

		for _, b := range []byte{b1, b2, b3} {
			c, err := Encode(b)
			if err != nil {
				return "", err
			}
			sb.WriteByte(c)
		}
		sb.WriteString("=")
	}

	return sb.String(), nil
}

func Base64decode(input string) (string, error) {
	byteData := []byte(input)
	var sb strings.Builder

	if len(byteData)%4 != 0 {
		return "", errors.New("illegal argument")
	}

	for i := 0; i < len(byteData); i += 4 {

		c := [4]byte{}
		for j := 0; j < 4; j++ {
			v, err := Decode(byteData[i+j])
			if byteData[i+j] != '=' && err != nil {
				return "", err
			}
			c[j] = v
		}

		b1 := c[0]<<2 | c[1]&0x30>>4
		b2 := c[1]&0x0f<<4 | c[2]&0x3c>>2
		b3 := c[2]&0x03<<6 | c[3]

		sb.WriteByte(b1)
		if byteData[i+2] != '=' {
			sb.WriteByte(b2)
		}
		if byteData[i+3] != '=' {
			sb.WriteByte(b3)
		}
	}
	return sb.String(), nil
}
