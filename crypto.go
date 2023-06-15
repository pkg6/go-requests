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

func Base64File(path string) (string, error) {
	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		return "", errors.New(fmt.Sprintf(`os.Open failed for name "%s"`, path))
	}
	return Base64Reader(f)
}

func Base64Reader(reader io.Reader) (string, error) {
	fd, err := io.ReadAll(reader)
	if err != nil {
		return "", errors.New(fmt.Sprintf(`io.ReadAll failed  "%v"`, err))
	}
	return base64.StdEncoding.EncodeToString(fd), nil
}

func Md5File(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", errors.New(fmt.Sprintf(`os.Open failed for name "%s"`, path))
	}
	defer f.Close()
	return Md5Reader(f)
}

func Md5Reader(reader io.Reader) (string, error) {
	h := md5.New()
	if _, err := io.Copy(h, reader); err != nil {
		return "", errors.New(`io.Copy failed`)
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
