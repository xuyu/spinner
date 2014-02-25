package spinner

import (
	"crypto/md5"
	"encoding/hex"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func MustInt64(s string) int64 {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		panic(err)
	}
	return i
}

func MustInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return i
}

func ReadVersion(filename string) string {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}

func FileMd5(filename string) (string, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	h := md5.New()
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil)), nil
}

func FilePerm(filename string, perm os.FileMode) (os.FileMode, error) {
	info, err := os.Stat(filename)
	if err != nil {
		return perm, err
	}
	return info.Mode(), nil
}
