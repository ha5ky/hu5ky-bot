/**
 * @Author Nil
 * @Description model/file.go
 * @Date 2023/3/28 17:18
 **/

package model

import (
	"errors"
	"github.com/ha5ky/hu5ky-bot/model/base"
	"gorm.io/gorm"
)

func (c *Controller) FileModel(m *File) *File {
	m.controller = c.controller
	return m
}

type File struct {
	controller *gorm.DB
	gorm.Model
	Name string
	Path string
}

func init() {
	c := new(File)
	c.Registry()
}

func (c *File) TableName() string {
	return "file"
}

func (c *File) Registry() {
	base.TableRegister = append(base.TableRegister, &File{})
}

type FileQuery struct {
	PageSize  *int
	PageIndex *int
	PreLoad   bool
	ID        *uint
	IDs       *[]uint
	Owners    []string

	ReviewClauseType *string
	Role             *string
	ReviewClause     *string
	ReviewClauseEn   *string
	SerialNumber     *string
	DRNode           *string
	OrderByNo        bool
}

func (c *File) Condition(q *FileQuery) *gorm.DB {
	return c.controller
}

func (c *File) NoExists(no string) (ok bool) {
	res := new(File)
	if errors.Is(c.controller.Where("no=?", no).Last(&res).Error, gorm.ErrRecordNotFound) {
		return false
	}
	return true
}

func (c *File) Exists(data string) (ok bool, err error) {
	var res File
	if errors.Is(c.controller.Where("JSON_CONTAINS(zh, ?)", data).Last(&res).Error, gorm.ErrRecordNotFound) {
		err = nil
		return
	}
	ok = true
	return
}

func (c *File) Save() error {
	return c.controller.Save(c).Error
}

func (c *File) Get(q *FileQuery) (res File, err error) {
	if errors.Is(c.Condition(q).Last(&res).Error, gorm.ErrRecordNotFound) {
		err = nil
		return
	}
	return
}

func (c *File) List(q *FileQuery) (res []*File, total int64, err error) {
	if err = c.Condition(q).Find(&res).Count(&total).Error; err != nil {
		return
	}
	if q.PageIndex != nil {
		c.controller = c.controller.Offset((*q.PageIndex - 1) * *q.PageSize)
	}
	if q.PageSize != nil {
		c.controller = c.controller.Limit(*q.PageSize)
	}
	if err = c.Condition(q).Find(&res).Error; err != nil {
		return
	}
	return
}

func (c *File) Delete(q *FileQuery) error {
	return c.Condition(q).Delete(c).Error
}
