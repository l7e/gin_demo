package file

import (
	"io/ioutil"
	"mime/multipart"
	"os"
	"path"
)

func GetSize(f multipart.File) (int, error) {
	bytes, e := ioutil.ReadAll(f)
	return len(bytes), e
}

func GetExt(filename string) string {
	return path.Ext(filename)
}

func CheckExist(src string) bool {
	_, err := os.Stat(src)
	return os.IsNotExist(err)
}

func CheckPermission(src string) bool {
	_, err := os.Stat(src)
	return os.IsPermission(err)
}

func MkDir(src string) error {
	err := os.MkdirAll(src, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func IsNotExistMkDir(src string) error {
	if notExist := CheckExist(src); notExist == true {
		if err := MkDir(src); err != nil {
			return err
		}
	}
	return nil
}

func Open(f string, flag int, perm os.FileMode) (*os.File, error) {
	file, err := os.OpenFile(f, flag, perm)
	if err != nil {
		return nil, err
	}
	return file, nil
}
