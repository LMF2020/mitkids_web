package fileUtils

import (
	uuid "github.com/satori/go.uuid"
	"io"
	"mime/multipart"
	"os"
	"strings"
)

const uploadPath = "/apistatic/uploadfile/"
const userPicPath = "userpic/"
const localPath = "."

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
	err = os.MkdirAll(localPath+tofileDirPath, os.ModePerm)
	if err != nil {
		return
	}
	filePath = tofileDirPath + fileName
	out, err := os.Create(localPath + filePath)
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

func IsExist(f string) bool {
	_, err := os.Stat(f)
	return err == nil || os.IsExist(err)
}
