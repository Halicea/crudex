package crudex

import "gorm.io/gorm"


type IModel interface {
	GetID() uint
	SetID(id uint)
}

type BaseModel struct {
	gorm.Model
}

func (self BaseModel) GetID() uint {
    return self.ID
}

func (self BaseModel) SetID(id uint) {
    self.ID = id
}
	
