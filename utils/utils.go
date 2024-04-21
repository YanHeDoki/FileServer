package utils

import (
	"bufio"
	"os"
)

type FileSteam struct {
	file     os.File
	reader   bufio.Reader
	size     int
	readSize int
}

func New(file os.File, size int) *FileSteam {
	return &FileSteam{
		file:     file,
		size:     size,
		readSize: 0,
	}
}

func (fs *FileSteam) Read(p []byte) (int, error) {
	return fs.file.Read(p)
}

func (fs *FileSteam) Close() error {
	return fs.file.Close()
}
