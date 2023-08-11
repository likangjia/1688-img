package file

import (
	"fmt"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"os"
)

func Download(url string, fileDir string, fileName string) error {
	// 发起http get请求
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 检查http响应的状态码，如果状态码不是200，说明请求没有成功
	if resp.StatusCode != http.StatusOK {
		return errors.New("failed to download file: status code is not 200")
	}
	err = os.MkdirAll(fileDir, os.ModePerm)
	if err != nil {
		return err
	}

	file, err := os.Create(generateFilePath(fileDir, fileName))
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err = io.Copy(file, resp.Body); err != nil {
		return err
	}

	return nil
}

func generateFilePath(dir, fileName string) string {
	return fmt.Sprintf("%s/%s", dir, fileName)
}
