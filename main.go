package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"mime/multipart"
	"net/http"
)

type ApiResponse struct {
	Code   int    `json:"code"`
	Result string `json:"result"`
	Msg    string `json:"msg"`
}

func main() {
	r := gin.Default()

	// 设置模板
	r.LoadHTMLGlob("templates/*")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	r.POST("/upload", func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			c.String(http.StatusBadRequest, "Error reading file")
			return
		}

		uploadedFile, err := file.Open()
		if err != nil {
			c.String(http.StatusBadRequest, "Error opening file")
			return
		}
		defer func(uploadedFile multipart.File) {
			err := uploadedFile.Close()
			if err != nil {

			}
		}(uploadedFile)

		// 创建multipart的writer（用于编写multipart/form-data编码的请求体）
		var requestBody bytes.Buffer
		multiPartWriter := multipart.NewWriter(&requestBody)

		// 添加文件部分到请求体
		fileWriter, err := multiPartWriter.CreateFormFile("file", file.Filename)
		if err != nil {
			fmt.Println("Error creating form file:", err)
			return
		}
		_, err = io.Copy(fileWriter, uploadedFile)
		if err != nil {
			fmt.Println("Error copying file:", err)
			return
		}

		// 关闭multipart writer以完成写入
		err = multiPartWriter.Close()
		if err != nil {
			fmt.Println("Error closing multiPartWriter:", err)
			return
		}

		// 创建新的请求
		req, err := http.NewRequest("POST", "https://api.oioweb.cn/api/ai/matting", &requestBody)
		if err != nil {
			fmt.Println("Error creating request:", err)
			return
		}

		// 设置请求头部信息，指定内容类型为multipart/form-data
		req.Header.Set("Content-Type", multiPartWriter.FormDataContentType())

		// 发送请求
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println("Error sending request:", err)
			return
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {

			}
		}(resp.Body)

		// 读取和解析返回的内容
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response body:", err)
			return
		}

		var response ApiResponse
		err = json.Unmarshal(body, &response)
		if err != nil {
			fmt.Println("Error parsing JSON:", err)
			return
		}

		fmt.Println("Code:", response.Code)
		fmt.Println("Result:", response.Result)
		fmt.Println("Message:", response.Msg)

		c.HTML(http.StatusOK, "index.html", gin.H{
			"imgURL": response.Result,
		})
	})

	err := r.Run(":8080")
	if err != nil {
		return
	}
}
