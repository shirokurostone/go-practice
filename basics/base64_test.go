package main

import "testing"

func TestEncode(t *testing.T) {

	cases := []struct{ raw, encoded byte }{
		{0, 'A'},
		{25, 'Z'},
		{26, 'a'},
		{51, 'z'},
		{52, '0'},
		{61, '9'},
		{62, '+'},
		{63, '/'},
	}

	for _, data := range cases {
		result, err := Encode(data.raw)
		if err != nil {
			t.Errorf("encode (err=%s)", err)
		}
		if result != data.encoded {
			t.Errorf("encode (expected=%d, actual=%d)", data.encoded, result)
		}
	}
	for _, data := range cases {
		result, err := Decode(data.encoded)
		if err != nil {
			t.Errorf("decode (err=%s)", err)
		}
		if result != data.raw {
			t.Errorf("decode (expected=%d, actual=%d)", data.raw, result)
		}
	}
}

func TestEncodeError(t *testing.T) {
	if _, err := Encode(64); err == nil {
		t.Errorf("encode (err == nil)")
	}
	if _, err := Decode('?'); err == nil {
		t.Errorf("decode (err == nil)")
	}
}

func TestBase64Encode(t *testing.T) {

	testCases := []struct{ raw, encoded string }{
		{"A", "QQ=="},
		{"AB", "QUI="},
		{"ABC", "QUJD"},
	}

	for _, data := range testCases {
		result, err := Base64encode(data.raw)

		if err != nil {
			t.Errorf("base64encode (err=%s)", err)
		}
		if result != data.encoded {
			t.Errorf("base64encode (expected=%s, actual=%s)", data.encoded, result)
		}

	}
	for _, data := range testCases {
		result, err := Base64decode(data.encoded)
		if err != nil {
			t.Errorf("base64dencode (err=%s)", err)
		}
		if result != data.raw {
			t.Errorf("base64decode (expected=%s, actual=%s)", data.raw, result)
		}
	}
}

func TestBase64DecodeError(t *testing.T) {
	testCases := []string{
		"123",
		"@@@@",
	}

	for _, data := range testCases {
		_, err := Base64decode(data)
		if err == nil {
			t.Errorf("base64decode (err=nil)")
		}
	}
}
