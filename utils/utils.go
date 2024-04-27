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
