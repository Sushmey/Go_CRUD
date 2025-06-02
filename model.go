package main

type User struct{
	ID uint `JSON:"id" gorm:"primaryKey"`
	Name string `JSON:"name" gorm:"not null"`
	Email string `JSON:"email" gorm:"not null;unique"`
	Age int `JSON:"age"`
}
 
//Serves as a Data Transfer Object (DTO) for input handling in web APIs
type UserParam struct{
	Name string `JSON:"name" binding:"required"`
	Email string `JSON:"email" binding:"required"`
	Age int `JSON:"age"`
}