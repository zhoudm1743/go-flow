package util

import (
	"io"
	"os"
	"path/filepath"
)

var FileUtil = fileUtil{}

// FileUtil 文件工具类
type fileUtil struct{}

// GetFileExt 获取文件扩展名
func (fu fileUtil) GetFileExt(name string) string {
	return filepath.Ext(name)
}

// GetFileName 获取文件名
func (fu fileUtil) GetFileName(name string) string {
	return filepath.Base(name)
}

// RemoveFile 删除文件
func (fu fileUtil) RemoveFile(name string) error {
	err := os.Remove(name)
	if err != nil {
		return err
	}
	return nil
}

// RemoveFiles 删除多个文件
func (fu fileUtil) RemoveFiles(names []string) error {
	for _, name := range names {
		if err := fu.RemoveFile(name); err != nil {
			return err
		}
	}
	return nil
}

// MoveFile 移动文件
func (fu fileUtil) MoveFile(src string, dst string) error {
	err := os.Rename(src, dst)
	if err != nil {
		return err
	}
	return nil
}

// CopyFile 复制文件
func (fu fileUtil) CopyFile(src string, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}
	return nil
}

// IsExist 判断文件是否存在
func (fu fileUtil) IsExist(name string) bool {
	_, err := os.Stat(name)
	return err == nil || os.IsExist(err)
}

// IsDir 判断是否是目录
func (fu fileUtil) IsDir(name string) bool {
	err := fu.IsExist(name)
	if err {
		fi, _ := os.Stat(name)
		return fi.IsDir()
	}
	return false
}

// IsFile 判断是否是文件
func (fu fileUtil) IsFile(name string) bool {
	err := fu.IsExist(name)
	if err {
		fi, _ := os.Stat(name)
		return !fi.IsDir()
	}
	return false
}

// CreateDir 创建目录
func (fu fileUtil) CreateDir(name string) error {
	err := os.MkdirAll(name, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

// CreateFile 创建文件
func (fu fileUtil) CreateFile(name string) error {
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	defer f.Close()
	return nil
}

// ReadFile 读取文件
func (fu fileUtil) ReadFile(name string) ([]byte, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var data []byte
	_, err = f.Read(data)
	return data, err
}

// WriteFile 写入文件
func (fu fileUtil) WriteFile(name string, data []byte) error {
	f, err := os.OpenFile(name, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(data)
	return err
}
