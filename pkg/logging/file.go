package logging

import (
	"fmt"
	"gin_demo/pkg/file"
	"gin_demo/pkg/setting"
	"github.com/pkg/errors"
	"os"
	"time"
)

func getLogFilePath() string {
	return fmt.Sprintf("%s%s", setting.AppSetting.RuntimeRootPath, setting.AppSetting.LogSavePath)
}

func getLogFileName() string {
	return fmt.Sprintf("%s%s.%s",
		setting.AppSetting.LogSaveName,
		time.Now().Format(setting.AppSetting.TimeFormat),
		setting.AppSetting.LogFileExt,
	)
}

func openLogFile(fileName, filePath string) (*os.File, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, errors.Errorf("os.GetWD err : %s", err)
	}

	src := dir + "/" + filePath
	if perm := file.CheckPermission(src); perm == true {
		return nil, fmt.Errorf("file.CheckPermission Permission denied src: %s", src)
	}

	err = file.IsNotExistMkDir(src)
	if err != nil {
		return nil, fmt.Errorf("file.IsNotExistMkDir src: %s, err: %v", src, err)
	}

	f, err := file.Open(src+fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)

	if err != nil {
		return nil, fmt.Errorf("Fail to OpenFile :%v", err)
	}
	return f, nil
}

//var (
//	LogSavePath = "runtime/logs/"
//	LogFileName = "log"
//	LogFileExt  = "log"
//	TimeFormat  = "20060102"
//)

//
//func getLogFileFullPath() string {
//	prefixPath := getLogFilePath()
//	suffixPath := fmt.Sprintf("%s%s.%s",
//		setting.AppSetting.LogSaveName,
//		time.Now().Format(setting.AppSetting.TimeFormat),
//		setting.AppSetting.LogFileExt,
//	)
//
//	return fmt.Sprintf("%s%s", prefixPath, suffixPath)
//}

//func mkDir() {
//	dir, _ := os.Getwd()
//	err := os.MkdirAll(dir+"/"+getLogFilePath(), os.ModePerm)
//	if err != nil {
//		panic(err)
//	}
//}
