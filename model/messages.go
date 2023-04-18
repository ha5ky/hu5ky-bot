/**
 * @Author Nil
 * @Description model/messages.go
 * @Date 2023/4/11 21:28
 **/

package model

import (
	"errors"
	"github.com/ha5ky/hu5ky-bot/model/base"
	"gorm.io/gorm"
)

func (c *Controller) MessageModel(m *Message) *Message {
	m.controller = c.controller
	return m
}

type Message struct {
	controller *gorm.DB
	gorm.Model
	Prompt     string
	Completion string
}

func init() {
	c := new(Message)
	c.Registry()
}

func (c *Message) TableName() string {
	return "message"
}

func (c *Message) Registry() {
	base.TableRegister = append(base.TableRegister, &Message{})
}

type MessageQuery struct {
	PageSize  *int
	PageIndex *int
	PreLoad   bool
	ID        *uint
	IDs       *[]uint

	Desc bool
}

func (c *Message) Condition(q *MessageQuery) *gorm.DB {
	if q.ID != nil {
		c.controller = c.controller.Where("id = ?", *q.ID)
	}
	if q.Desc {
		c.controller = c.controller.Order("created_at desc")
	}
	return c.controller
}

func (c *Message) NoExists(no string) (ok bool) {
	res := new(Message)
	if errors.Is(c.controller.Where("no=?", no).Last(&res).Error, gorm.ErrRecordNotFound) {
		return false
	}
	return true
}

func (c *Message) Exists(data string) (ok bool, err error) {
	var res Message
	if errors.Is(c.controller.Where("JSON_CONTAINS(zh, ?)", data).Last(&res).Error, gorm.ErrRecordNotFound) {
		err = nil
		return
	}
	ok = true
	return
}

func (c *Message) Save() error {
	return c.controller.Save(c).Error
}

func (c *Message) Get(q *MessageQuery) (res Message, err error) {
	if errors.Is(c.Condition(q).Last(&res).Error, gorm.ErrRecordNotFound) {
		err = nil
		return
	}
	return
}

func (c *Message) List(q *MessageQuery) (res []*Message, total int64, err error) {
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

func (c *Message) Delete(q *MessageQuery) error {
	return c.Condition(q).Delete(c).Error
}
