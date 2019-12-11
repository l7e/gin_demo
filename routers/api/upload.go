package api

import (
	"gin_demo/pkg/e"
	"gin_demo/pkg/upload"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UploadImage(c *gin.Context) {
	code := e.SUCCESS
	data := make(map[string]interface{})

	file, header, err := c.Request.FormFile("image")
	if err != nil {
		code = e.ERROR
	}

	if header == nil {
		code = e.INVALID_PARAMS
	} else {
		imageName := upload.GetImageName(header.Filename)
		fullPath := upload.GetImageFullPath()
		savePath := upload.GetImageSavePath()

		src := fullPath + imageName
		if !upload.CheckImageExt(imageName) || !upload.CheckImageSize(file) {
			code = e.ERROR_UPLOAD_CHECK_IMAGE_FAIL
		} else {
			err := upload.CheckImage(fullPath)
			if err != nil {
				code = e.ERROR_UPLOAD_CHECK_IMAGE_FAIL
			} else if err := c.SaveUploadedFile(header, src); err != nil {
				code = e.ERROR_UPLOAD_SAVE_IMAGE_FAIL
			} else {
				data["image_url"] = upload.GetImageFullUrl(imageName)
				data["image_save_url"] = savePath + imageName
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})
}

