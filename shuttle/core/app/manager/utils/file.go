package utils

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/xid"
	"google.dev/google/common/pkg/logger"
)

func UploadFile(w http.ResponseWriter, r *http.Request) {
	// 解析表单域限制文件大小为 8MB
	r.ParseMultipartForm(8 << 20)

	// 获取上传的文件
	file, handler, err := r.FormFile("file")
	if err != nil {
		logger.Error(err)
		return
	}
	defer file.Close()

	// 生成UUID文件名
	newFileName := xid.New().String() + filepath.Ext(handler.Filename)

	// 打开文件，准备写入
	targetFile, err := os.OpenFile("./static/"+newFileName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return
	}
	defer targetFile.Close()

	// 从上传的文件中读取数据，并写入到打开的文件中
	fileData, err := ioutil.ReadAll(file)
	if err != nil {
		logger.Error(err)
		return
	}
	_, err = io.Copy(targetFile, strings.NewReader(string(fileData)))
	if err != nil {
		logger.Error(err)
		return
	}

	// 构建响应
	responseBody := map[string]string{
		"filename": newFileName,
		"url":      r.Host + "/static/" + newFileName,
	}

	// 返回响应
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(responseBody)
}
