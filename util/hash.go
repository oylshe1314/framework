package util

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"hash"
	"os"
	"strings"
)

func MD5(src string) string {
	var h = md5.New()
	h.Write([]byte(src))
	return hex.EncodeToString(h.Sum(nil))
}

func SHA1(src string) string {
	var h = sha1.New()
	h.Write([]byte(src))
	return hex.EncodeToString(h.Sum(nil))
}

func SHA256(src string) string {
	var h = sha256.New()
	h.Write([]byte(src))
	return hex.EncodeToString(h.Sum(nil))
}

func SHA512(src string) string {
	var h = sha512.New()
	h.Write([]byte(src))
	return hex.EncodeToString(h.Sum(nil))
}

func HMACSHA256(key, src string) string {
	var h = hmac.New(sha256.New, []byte(key))
	h.Write([]byte(src))
	return hex.EncodeToString(h.Sum(nil))
}

// HashAll 计算给定内容的hash值，包括字符串，目录(目录中的文件)，文件
// merge 是否合并计算成一个hash值
// strs 字符串
// dirs 目录
// files 文件
func HashAll(h hash.Hash, merge bool, strs, dirs, files []string) ([]string, []string, error) {
	var hashs []string
	var infos []string
	if merge {
		h.Reset()
		for _, str := range strs {
			if len(str) == 0 {
				continue
			}
			h.Write([]byte(str))
		}

		for _, dir := range dirs {
			if dir[len(dir)-1] == '/' || dir[len(dir)-1] == '\\' {
				dir = dir[:len(dir)-1]
			}
			subFiles, err := ReadFiles(dir)
			if err != nil {
				return nil, nil, err
			}

			files = append(files, subFiles...)
		}

		for _, file := range files {
			buf, err := os.ReadFile(file)
			if err != nil {
				return nil, nil, err
			}

			h.Write(buf)
		}
		hashs = append(hashs, strings.ToUpper(hex.EncodeToString(h.Sum(nil))))
		infos = append(infos, "")
	} else {
		for _, str := range strs {
			if len(str) == 0 {
				continue
			}
			h.Reset()
			h.Write([]byte(str))
			hashs = append(hashs, strings.ToUpper(hex.EncodeToString(h.Sum(nil))))
			infos = append(infos, str)
		}

		for _, dir := range dirs {
			dir = strings.ReplaceAll(dir, "\\", "/")
			if dir[len(dir)-1] == '/' {
				dir = dir[:len(dir)-1]
			}
			subFiles, err := ReadFiles(dir)
			if err != nil {
				return nil, nil, err
			}

			files = append(files, subFiles...)
		}

		for _, file := range files {
			buf, err := os.ReadFile(file)
			if err != nil {
				return nil, nil, err
			}

			h.Reset()
			h.Write(buf)

			hashs = append(hashs, strings.ToUpper(hex.EncodeToString(h.Sum(nil))))
			infos = append(infos, file)
		}
	}

	return hashs, infos, nil
}
