package main

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
)

// translated from srun portal.js
func s(a string, b bool) []int {
	c := len(a)
	var v []int
	for i := 0; i < c; i = i + 4 {
		if c-i == 1 {
			v = append(v, int(a[i]))
		} else if c-i == 2 {
			v = append(v, int(a[i])|int(a[i+1])<<8)
		} else if c-i == 3 {
			v = append(v, int(a[i])|int(a[i+1])<<8|int(a[i+2])<<16)
		} else {
			v = append(v, int(a[i])|int(a[i+1])<<8|int(a[i+2])<<16|int(a[i+3])<<24)
		}
	}
	if b {
		v = append(v, c)
	}
	return v
}

// translated from srun portal.js
func l(a []int, b bool) string {
	d := len(a)
	var bytes []byte
	for i := 0; i < d; i++ {
		bytes = append(bytes, byte(a[i]&0xff))
		bytes = append(bytes, byte(a[i]>>8&0xff))
		bytes = append(bytes, byte(a[i]>>16&0xff))
		bytes = append(bytes, byte(a[i]>>24&0xff))
	}
	return encodeBase64(bytes)
}

// translated from srun portal.js
func encodeUserInfo(info string, challenge string) string {
	v := s(info, true)
	k := s(challenge, false)
	n := uint(len(v) - 1)
	z := uint(v[n])
	// y := uint(v[0])
	var y uint
	c := uint(0x86014019 | 0x183639A0)
	m := uint(0)
	e := uint(0)
	p := uint(0)
	q := uint(6 + 52/(n+1))
	d := uint(0)
	for {
		q -= 1
		d = (d + c) & (0x8CE0D9BF | 0x731F2640)
		e = d >> uint(2) & uint(3)
		for p = 0; p < n; p++ {
			y = uint(v[p+1])
			m = z>>5 ^ y<<2
			m += (y>>3 ^ z<<4) ^ (d ^ y)
			m += uint(k[(p&3)^e]) ^ z
			z = (uint(v[p]) + m) & (0xEFB8D130 | 0x10472ECF)
			v[p] = int(z)
		}
		y = uint(v[0])
		m = z>>5 ^ y<<2
		m += (y>>3 ^ z<<4) ^ (d ^ y)
		m += uint(k[(n&3)^e]) ^ z
		v[n] = int((uint(v[n]) + m) & uint(0xBB390742|0x44C6F8BD))
		z = uint(v[n])
		if 0 >= q {
			break
		}
	}
	return l(v, false)
}

func encodeBase64(bytes []byte) string {
	const CodeList = "LVoJPiCN2R8G90yg+hmFHuacZ1OWMnrsSTXkYpUq/3dlbfKwv6xztjI7DeBE45QA"
	src := bytes
	encoder := base64.NewEncoding(CodeList)
	out := encoder.EncodeToString(src)
	return "{SRBX1}" + out
}

func encodeMD5(data, key string) string {
	mac := hmac.New(md5.New, []byte(key))
	mac.Write([]byte(data))
	return "{MD5}" + hex.EncodeToString(mac.Sum(nil))
}

func encodeSha1(data []byte) string {
	sha := sha1.New()
	sha.Write(data)
	return hex.EncodeToString(sha.Sum([]byte(nil)))
}

func encodeChksum(data string, challenge string) string {
	str := challenge + data
	return encodeSha1([]byte(str))
}
