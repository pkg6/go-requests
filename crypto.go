package requests

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
)

func Md5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func Sha1(str string) string {
	h := sha1.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func Sha256(str string) string {
	h := sha256.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func URLEncode(str string) string {
	return url.QueryEscape(str)
}

func URLDecode(str string) (string, error) {
	return url.QueryUnescape(str)
}

func Base64Encode(str string) string {
	return base64.URLEncoding.EncodeToString([]byte(str))
}

func Base64Decode(str string) (string, error) {
	decodeString, err := base64.URLEncoding.DecodeString(str)
	if err != nil {
		return "", err
	}
	return string(decodeString), nil
}

func Md5File(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", errors.New(fmt.Sprintf(`os.Open failed for name "%s"`, path))
	}
	defer f.Close()
	h := md5.New()
	_, err = io.Copy(h, f)
	if err != nil {
		return "", errors.New(`io.Copy failed`)
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
