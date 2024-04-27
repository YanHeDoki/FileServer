package logic

import (
	"FileServerFiber/utils"
	"fmt"
	"io"
	"os"

	"github.com/gofiber/fiber/v2"
)

// GetFileList 遍历当前目录返回
func GetFileList(c *fiber.Ctx) error {

	// 扫描当前目录
	path_list, err := os.ReadDir(".")

	if err != nil {
		return err
	}

	// 过滤文件
	file_list := make([]string, 0, len(path_list))

	for i := 0; i < len(path_list); i++ {
		if path_list[i].IsDir() {
			continue
		}
		file_list = append(file_list, path_list[i].Name())
	}

	// 返回目录文件列表
	c.JSON(file_list)

	return nil
}

// 下载文件
func DownloadFile(c *fiber.Ctx) error {

	// 确认文件是否存在
	file_name := c.Params("+")

	_, err := os.Stat(file_name)
	if os.IsNotExist(err) {
		return c.SendStatus(fiber.StatusNotFound)
	}

	// 文件存在,流式返回文件
	file, err := os.OpenFile(file_name, os.O_RDONLY, 0666)

	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	defer file.Close()
	// 设置HTTP响应头

	c.Set("Content-Type", "application/octet-stream")
	_, err = io.CopyBuffer(c, file, make([]byte, 1024))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return nil
}

// 自定义文件流下载文件
func MyDownloadFile(c *fiber.Ctx) error {

	// 确认文件是否存在
	file_name := c.Params("+")

	_, err := os.Stat(file_name)
	if os.IsNotExist(err) {
		return c.SendStatus(fiber.StatusNotFound)
	}

	// 文件存在,流式返回文件
	file, err := os.OpenFile(file_name, os.O_RDONLY, 0666)

	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	fs := utils.New(*file, 1024)

	defer fs.Close()
	// 设置HTTP响应头

	c.Set("Content-Type", "application/octet-stream")
	read_size, err := fs.Stream(c)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	fmt.Println(read_size)
	return nil
}
