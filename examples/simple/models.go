package main

import "github.com/halicea/crudex"

type Car struct {
	crudex.BaseModel
	Name        string `crud-input:"text" crud-placeholder:"Enter name"`
	License     string `crud-input:"html" crud-placeholder:"Enter the license plate"`
	Description string `crud-input:"wysiwyg" crud-placeholder:"Describe it"`
	Year        int    `crud-input:"number" crud-placeholder:"Model year of the car"`
}

type Driver struct {
	crudex.BaseModel
	Name  string
	CarID uint
	Car   Car `gorm:"foreignKey:CarID"`
}
