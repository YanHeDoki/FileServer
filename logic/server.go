package logic

import (
	"fmt"
	"html/template"
	"io"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"

	"FileServerFiber/utils"
)

// GetFileList 遍历当前目录返回
func GetFileList(c *fiber.Ctx) error {

	// 检查参数
	in_path := c.Params("*")
	get_path := "./"
	// 有期望进入的目录
	if in_path != "" {
		get_path += in_path
	}

	file_list, err := utils.GetFiles(get_path)
	if err != nil {
		fmt.Println(err)
		return c.SendStatus(fiber.StatusNotFound)
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
	fmt.Println("file read size:", read_size)
	return nil
}

func Tmpl(c *fiber.Ctx) error {

	// 定义startsWith函数
	funcMap := template.FuncMap{
		"starts_with": func(s, prefix string) bool {
			return strings.HasPrefix(s, prefix)
		},
		"add_path": func(p string) string {
			original_url := c.OriginalURL()
			if original_url == "/" {
				original_url = ""
			}
			return c.BaseURL() + original_url + p
		},
		"download_path": func(p string) string {
			original_url := c.OriginalURL()
			if original_url == "/" {
				original_url = ""
			}
			return c.BaseURL() + "/download" + original_url + "/" + p
		},
	}

	t1, err := template.New("index.html").Funcs(funcMap).ParseFiles("./tmpl/index.html")
	if err != nil {
		panic(err)
	}

	// 检查参数
	in_path := c.Params("*")
	get_path := "./"

	// 有期望进入的目录
	if in_path != "" {
		get_path += in_path
	}

	// 获取当前请求的路径
	uri := string(c.Request().URI().Path())
	parent_patch := getParentPath(uri)

	file_list, err := utils.GetFiles(get_path)
	if err != nil {
		fmt.Println(err)
		return c.SendStatus(fiber.StatusNotFound)
	}

	dateMap := map[string]any{"file_list": file_list, "current_path": uri, "parent_patch": parent_patch}
	c.Set("Content-Type", "text/html")
	return t1.Execute(c.Response().BodyWriter(), dateMap)
}

// getParentPath 返回给定路径的父级路径
func getParentPath(p string) string {
	// 找到最后一个斜杠的位置
	lastSlashIdx := strings.LastIndexByte(p, '/')
	// 如果没有斜杠，返回根路径
	if lastSlashIdx == -1 {
		return "/"
	} else if lastSlashIdx == 0 { // 如果最后一个 / 是在头部也去掉
		return "/"
	}
	// 返回从开始到最后一个斜杠的部分
	return p[:lastSlashIdx]
}
