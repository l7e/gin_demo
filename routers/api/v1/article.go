package v1

import (
	"gin_demo/models"
	"gin_demo/pkg/e"
	"gin_demo/pkg/setting"
	"gin_demo/pkg/util"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

func GetArticles(c *gin.Context) {
	data := make(map[string]interface{})
	maps := make(map[string]interface{})

	valid := validation.Validation{}
	state := -1

	if arg := c.Query("state"); arg != "" {
		state, _ = strconv.Atoi(arg)
		valid.Range(state, 0, 1, "state").Message("状态只能为0或1")
		maps["state"] = state
	}

	if arg := c.Query("tag_id"); arg != "" {
		tagId, _ := strconv.Atoi(arg)
		valid.Min(tagId, 1, "tag_id").Message("tag_id不能小于1")
		maps["tag_id"] = tagId
	}

	code := e.INVALID_PARAMS

	if !valid.HasErrors() {
		data["lists"] = models.GetArticles(util.GetPage(c), setting.AppSetting.PageSize, maps)
		data["total"] = models.GetArticleTotal(maps)
		code = e.SUCCESS
	} else {
		for _, err := range valid.Errors {
			log.Printf("%s", err)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})
}

func GetArticle(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	valid := validation.Validation{}

	valid.Required(id, "id").Message("id不能为空")
	valid.Min(id, 1, "id").Message("id必须大于0")

	code := e.INVALID_PARAMS

	var data interface{}

	if !valid.HasErrors() {
		if models.ExistArticleByID(id) {
			data = models.GetArticle(id)
			code = e.SUCCESS
		} else {
			code = e.ERROR_NOT_EXIST_ARTICLE
		}
	} else {
		for _, err := range valid.Errors {
			log.Printf("%s", err)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})
}

func AddArticle(c *gin.Context) {
	tagId, _ := strconv.Atoi(c.Query("tag_id"))
	title := c.Query("title")
	desc := c.Query("desc")
	content := c.Query("content")
	createdBy := c.Query("created_by")
	state, _ := strconv.Atoi(c.DefaultQuery("state", "0"))

	valid := validation.Validation{}
	valid.Min(tagId, 1, "tag_id").Message("tag_id不能小于1")
	valid.Required(title, "title").Message("文章标题不能为空")
	valid.Required(desc, "desc").Message("简述不能为空")
	valid.Required(content, "content").Message("内容不能为空")
	valid.Required(createdBy, "created_by").Message("创建人不能为空")
	valid.Required(state, "state").Message("状态不能为空")

	code := e.INVALID_PARAMS

	if !valid.HasErrors() {
		if models.ExistTagByID(tagId) {
			data := make(map[string]interface{})
			data["tag_id"] = tagId
			data["title"] = title
			data["desc"] = desc
			data["content"] = content
			data["created_by"] = createdBy
			data["state"] = state

			models.AddArticle(data)
			code = e.SUCCESS
		} else {
			code = e.ERROR_NOT_EXIST_TAG
		}
	} else {
		for _, err := range valid.Errors {
			log.Printf("err.key : %s ,err.message : %s\n", err.Key, err.Message)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]interface{}),
	})
}

func UpdateArticle(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	tagId, _ := strconv.Atoi(c.Query("tag_id"))
	title := c.Query("title")
	desc := c.Query("desc")
	content := c.Query("content")
	modifiedBy := c.Query("modified_by")

	valid := validation.Validation{}

	state := -1
	if arg := c.Query("state"); arg != "" {
		state, _ := strconv.Atoi(arg)
		valid.Range(state, 0, 1, "state").Message("状态只允许0或1")
	}

	valid.Min(id, 0, "id").Message("id不能小于1")
	valid.Min(tagId, 0, "id").Message("tag_id不能小于1")
	valid.MaxSize(title, 100, "title").Message("标题最长100字符")
	valid.MaxSize(desc, 100, "desc").Message("简述最长100字符")
	valid.MaxSize(content, 100, "content").Message("内容最长100字符")
	valid.Required(modifiedBy, "created_by").Message("修改人不能为空")
	valid.MaxSize(modifiedBy, 100, "created_by").Message("修改人最长100字符")

	code := e.INVALID_PARAMS

	if !valid.HasErrors() {
		if models.ExistArticleByID(id) {
			if models.ExistTagByID(tagId) {
				code = e.SUCCESS

				data := make(map[string]interface{})

				if tagId > 0 {
					data["tag_id"] = tagId
				}

				if title != "" {
					data["title"] = title
				}

				if desc != "" {
					data["desc"] = desc
				}

				if content != "" {
					data["content"] = content
				}

				if state != -1 {
					data["state"] = state
				}

				data["modified_by"] = modifiedBy

				models.UpdateArticle(id, data)
			} else {
				code = e.ERROR_NOT_EXIST_TAG
			}
		} else {
			code = e.ERROR_NOT_EXIST_ARTICLE
		}
	} else {
		for _, err := range valid.Errors {
			log.Printf("err.key : %s ,err.message : %s\n", err.Key, err.Message)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]interface{}),
	})
}

func DeleteArticle(c *gin.Context) {
	id, _ := strconv.Atoi(c.Query("id"))

	valid := validation.Validation{}
	valid.Min(id, 0, "id").Message("id不能小于0")

	code := e.INVALID_PARAMS
	if !valid.HasErrors() {
		if models.ExistArticleByID(id) {
			models.DeleteArticle(id)
			code = e.SUCCESS
		} else {
			code = e.ERROR_NOT_EXIST_ARTICLE
		}
	} else {
		for _, err := range valid.Errors {
			log.Printf("key : %s,value : %s", err.Key, err.Message)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]interface{}),
	})
}
