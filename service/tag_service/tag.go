package tag_service

import (
	"encoding/json"
	"gin_demo/models"
	"gin_demo/pkg/excel"
	"gin_demo/pkg/gredis"
	"gin_demo/pkg/logging"
	"gin_demo/service/cache_service"
	"github.com/tealeg/xlsx"
	"strconv"
	"time"
)

type Tag struct {
	ID         int
	Name       string
	CreatedBy  string
	ModifiedBy string
	State      int

	PageNum  int
	PageSize int
}

func (t *Tag) Export() (string, error) {
	tags, err := t.GetAll()
	if err != nil {
		return "", err
	}

	file := xlsx.NewFile()
	sheet, err := file.AddSheet("标签信息")
	if err != nil {
		return "", err
	}

	titles := []string{"ID", "名称", "创建人", "创建时间", "修改人", "修改时间"}
	row := sheet.AddRow()

	for _, title := range titles {
		cell := row.AddCell()
		cell.Value = title
	}

	for _, tag := range tags {
		values := []string{
			strconv.Itoa(tag.ID),
			tag.Name,
			tag.CreatedBy,
			strconv.Itoa(tag.CreatedOn),
			tag.ModifiedBy,
			strconv.Itoa(tag.ModifiedOn),
		}

		row := sheet.AddRow()
		for _, v := range values {
			cell := row.AddCell()
			cell.Value = v
		}
	}

	times := strconv.Itoa(int(time.Now().Unix()))
	filename := "tag-" + times + ".xlsx"
	fullPath := excel.GetExcelFullPath() + filename
	err = file.Save(fullPath)
	logging.Info(err,11111)
	if err != nil {
		return "", err
	}
	return filename, nil
}

func (t *Tag) GetAll() ([]models.Tag, error) {
	var (
		tags, cacheTags []models.Tag
	)

	cache := cache_service.Tag{
		State:    t.State,
		PageNum:  t.PageNum,
		PageSize: t.PageSize,
	}

	keys := cache.GetTagsKey()

	if gredis.Exists(keys) {
		data, err := gredis.Get(keys)
		if err != nil {
			logging.Info(err)
		} else {
			json.Unmarshal(data, &cacheTags)
			return cacheTags, nil
		}
	}

	tags, err := models.GetTags(t.PageNum, t.PageSize, t.getMaps())
	if err != nil {
		return nil, err
	}

	gredis.Set(keys, tags, 3600)
	return tags, nil
}

func (t *Tag) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})
	maps["deleted_on"] = 0

	if t.Name != "" {
		maps["name"] = t.Name
	}
	if t.State >= 0 {
		maps["state"] = t.State
	}

	return maps
}
