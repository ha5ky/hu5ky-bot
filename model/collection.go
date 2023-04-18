/**
 * @Author Nil
 * @Description model/collection.go
 * @Date 2023/4/2 18:08
 **/

package model

import (
	"errors"
	"github.com/ha5ky/hu5ky-bot/model/base"
	"gorm.io/gorm"
)

func (c *Controller) CollectionModel(m *Collection) *Collection {
	m.controller = c.controller
	return m
}

type Collection struct {
	controller *gorm.DB
	gorm.Model
	Name        string
	FileId      uint
	Description string
}

func init() {
	c := new(Collection)
	c.Registry()
}

func (c *Collection) TableName() string {
	return "collection"
}

func (c *Collection) Registry() {
	base.TableRegister = append(base.TableRegister, &Collection{})
}

type CollectionQuery struct {
	PageSize  *int
	PageIndex *int
	PreLoad   bool
	ID        *uint
	IDs       *[]uint
}

func (c *Collection) Condition(q *CollectionQuery) *gorm.DB {
	if q.ID != nil {
		c.controller = c.controller.Where("id = ?", *q.ID)
	}
	return c.controller
}

func (c *Collection) NoExists(no string) (ok bool) {
	res := new(Collection)
	if errors.Is(c.controller.Where("no=?", no).Last(&res).Error, gorm.ErrRecordNotFound) {
		return false
	}
	return true
}

func (c *Collection) Exists(data string) (ok bool, err error) {
	var res Collection
	if errors.Is(c.controller.Where("JSON_CONTAINS(zh, ?)", data).Last(&res).Error, gorm.ErrRecordNotFound) {
		err = nil
		return
	}
	ok = true
	return
}

func (c *Collection) Save() error {
	return c.controller.Save(c).Error
}

func (c *Collection) Get(q *CollectionQuery) (res Collection, err error) {
	if errors.Is(c.Condition(q).Last(&res).Error, gorm.ErrRecordNotFound) {
		err = nil
		return
	}
	return
}

func (c *Collection) List(q *CollectionQuery) (res []*Collection, total int64, err error) {
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

func (c *Collection) Delete(q *CollectionQuery) error {
	return c.Condition(q).Delete(c).Error
}
