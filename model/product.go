package model

import (
	"github.com/jinzhu/gorm"
	"github.com/myrachanto/allmicro/gormmicro/categorymicroservice/httperors"
)

type Product struct {
	Name string `gorm:"not null"`
	Title string `gorm:"not null"`
	Description string `gorm:"not null"`
	SubCategory SubCategory `gorm:"foreignKey:UserID; not null"`
	SubCategoryID uint `gorm:"not null"`
	Picture string 
	gorm.Model
}
type MajorCategory struct {
	Name string `gorm:"not null"`
	Title string `gorm:"not null"`
	Description string `gorm:"not null"`
	Category []Category `gorm:"foreignKey:UserID; not null"`
	gorm.Model
}
type Category struct {
	Name string `gorm:"not null"`
	Title string `gorm:"not null"`
	Description string `gorm:"not null"`
	MajorCategory MajorCategory `gorm:"foreignKey:UserID; not null"`
	MajorCategoryID uint `json:"userid"`
	SubCategory []SubCategory
	gorm.Model
}
type SubCategory struct {
	Name string `gorm:"not null"`
	Title string `gorm:"not null"`
	Description string `gorm:"not null"`
	Category Category `gorm:"foreignKey:UserID; not null"`
	CategoryID uint `json:"userid"`
	Product []Product
	gorm.Model
}

// func (product Product) Validate() *httperors.HttpError{
// 	if product.Name == "" && len(product.Name) < 3 {
// 		return httperors.NewNotFoundError("Invalid name")
// 	}
// 	if product.Title == "" && len(product.Title) < 5 {
// 		return httperors.NewNotFoundError("Invalid title")
// 	}
// 	if product.Description == "" && len(product.Description) < 10 {
// 		return httperors.NewNotFoundError("Invalid description")
// 	}
// 	if product.Picture == "" {
// 		return httperors.NewNotFoundError("Invalid picture")
// 	}
// 	return nil
// }
// func (majorCategory MajorCategory) Validate() *httperors.HttpError{
// 	if majorCategory.Name == "" && len(majorCategory.Name) < 3 {
// 		return httperors.NewNotFoundError("Invalid name")
// 	}
// 	if majorCategory.Title == "" && len(majorCategory.Title) < 5 {
// 		return httperors.NewNotFoundError("Invalid title")
// 	}
// 	if majorCategory.Description == "" && len(majorCategory.Description) < 10 {
// 		return httperors.NewNotFoundError("Invalid description")
// 	}
// 	return nil
// }
func (category Category) Validate() *httperors.HttpError{
	if category.Name == "" && len(category.Name) < 3 {
		return httperors.NewNotFoundError("Invalid name")
	}
	if category.Title == "" && len(category.Title) < 5 {
		return httperors.NewNotFoundError("Invalid title")
	}
	if category.Description == "" && len(category.Description) < 10 {
		return httperors.NewNotFoundError("Invalid description")
	}
	return nil
}
// func (subcategory SubCategory) Validate() *httperors.HttpError{
// 	if subcategory.Name == "" && len(subcategory.Name) < 3 {
// 		return httperors.NewNotFoundError("Invalid name")
// 	}
// 	if subcategory.Title == "" && len(subcategory.Title) < 5 {
// 		return httperors.NewNotFoundError("Invalid title")
// 	}
// 	if subcategory.Description == "" && len(subcategory.Description) < 10 {
// 		return httperors.NewNotFoundError("Invalid description")
// 	}
// 	return nil
// }