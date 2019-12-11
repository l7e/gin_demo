package upload

import (
	"fmt"
	"gin_demo/pkg/file"
	"gin_demo/pkg/setting"
	"gin_demo/pkg/util"
	"github.com/pkg/errors"
	"log"
	"mime/multipart"
	"os"
	"path"
	"strings"
	"time"
)

//获取图片完整的URL路径
func GetImageFullUrl(name string) string {
	return setting.AppSetting.ImagePrefixUrl + "/" + GetImageSavePath() + name
}

//获取图片保存的路径
func GetImageSavePath() string {
	return setting.AppSetting.ImageSavePath
}

//获取MD5编码后的图片名
func GetImageName(name string) string {
	ext := path.Ext(name)
	fileName := strings.TrimSuffix(name, ext)
	fileName = util.EncodeMD5(fileName)

	return time.Now().Format("20060102150405") + fileName + ext
}

//获取图片本地完整路径
func GetImageFullPath() string {
	return setting.AppSetting.RuntimeRootPath + GetImageSavePath()
}

//检查图片后缀
func CheckImageExt(filename string) bool {
	ext := file.GetExt(filename)
	for _, allowExt := range setting.AppSetting.ImageAllowExts {
		if strings.ToLower(ext) == strings.ToLower(allowExt) {
			return true
		}
	}
	return false
}

func CheckImageSize(f multipart.File) bool {
	size, err := file.GetSize(f)
	if err != nil {
		log.Println(err)
		return false
	}
	return size <= setting.AppSetting.ImageMaxSize
}


//检查路径是否存在,是否有权限
func CheckImage(src string) error {
	dir, err := os.Getwd()
	if err != nil {
		return errors.Errorf("os.Getwd err: %v", err)
	}

	err = file.IsNotExistMkDir(dir + "/" + src)
	if err != nil {
		return errors.Errorf("file.IsNotExistMkDir err: %v", err)
	}

	perm := file.CheckPermission(src)
	if perm == true {
		return fmt.Errorf("file.CheckPermission Permission denied src: %s", src)
	}

	return nil
}
