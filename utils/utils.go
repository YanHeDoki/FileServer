package utils

import (
	"io"
	"os"
)

type FileStream struct {
	file     os.File
	size     int
	readSize int
}

func New(file os.File, size int) *FileStream {
	return &FileStream{
		file:     file,
		size:     size,
		readSize: 0,
	}
}

func (fs *FileStream) Stream(w io.Writer) (int, error) {

	// 写入计算
	file_info, err := fs.file.Stat()
	if err != nil {
		return 0, err
	}

	// 如果缓存的容量已经大于了文件的大小就调整为当前文件的大小
	if fs.size > int(file_info.Size()) {
		fs.size = int(file_info.Size())
	}

	// 构造缓存
	buf := make([]byte, fs.size)

	// 不断读取写入文件
	for fs.readSize < int(file_info.Size()) {
		// 读取缓存大小的文件
		_, err := fs.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return fs.readSize, err
		}

		// 写入并且累计写入的字节
		read_size, err := w.Write(buf)
		if err != nil {
			return fs.readSize, err
		}
		fs.readSize += read_size
	}
	return fs.readSize, nil
}

func (fs *FileStream) Read(p []byte) (int, error) {
	return fs.file.Read(p)
}

func (fs *FileStream) Close() error {
	return fs.file.Close()
}

func GetFiles(path string) ([]string, error) {
	// 扫描当前目录
	path_list, err := os.ReadDir(path)

	if err != nil {
		return nil, err
	}

	// 过滤文件
	file_list := make([]string, 0, len(path_list))

	for i := 0; i < len(path_list); i++ {
		if path_list[i].IsDir() {
			file_list = append(file_list, "/"+path_list[i].Name())
			continue
		}
		file_list = append(file_list, path_list[i].Name())
	}
	return file_list, nil
}
