package fileUtils

import (
	uuid "github.com/satori/go.uuid"
	"io"
	"mime/multipart"
	"os"
	"strings"
)

const uploadPath = "static/uploadfile/"
const userPicPath = "userpic/"

func UpdateUserPic(accountId, fileName string, file multipart.File) (filePath string, err error) {
	newName := getNewFileName(fileName)
	tofileDirPath := uploadPath + userPicPath + accountId + "/"
	os.RemoveAll(tofileDirPath)
	return UploadFile(file, newName, tofileDirPath)
}

func getNewFileName(fileName string) string {
	arr := strings.Split(fileName, ".")
	fileExt := arr[len(arr)-1]
	return uuid.NewV4().String() + "." + fileExt
}

func UploadFile(file multipart.File, fileName, tofileDirPath string) (filePath string, err error) {
	os.MkdirAll(tofileDirPath, os.ModePerm)
	filePath = tofileDirPath + fileName
	out, err := os.Create(filePath)
	if err != nil {
		return
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		return
	}
	return
}