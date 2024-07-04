package main

import "github.com/halicea/crudex"

type Car struct {
	crudex.BaseModel
	Name        string
	License     string
	Description string
	Year        int
}

type Driver struct {
	crudex.BaseModel
	Name  string
	CarID uint
	Car   Car `gorm:"foreignKey:CarID"`
}
