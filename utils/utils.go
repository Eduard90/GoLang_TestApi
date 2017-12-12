package utils

import (
	"fmt"
	"api_test_2/db"
	"api_test_2/models"
	"time"
	"golang.org/x/crypto/bcrypt"
)

func AddCategory(name string, parentId uint) {
	//-- UPDATE categories SET lft = lft + 2, rgt = rgt + 2 WHERE lft > 5
	//
	//-- UPDATE categories SET rgt = rgt + 2 WHERE rgt >= 5 AND lft < 5
	//
	//INSERT INTO categories SET lft = 5, rgt = 5 + 1, level = 0 + 1, name = "Processors" , parent_id = 1, created_at="2017-12-05 15:40:53", updated_at = "2017-12-05 15:40:53";


	databaseConnect := db.GetDatabaseConnectInstance()

	dbConnect := databaseConnect.GetDbConnect()

	var parentCategory models.Category
	dbConnect.First(&parentCategory, parentId)
	if parentCategory.ID != parentId {  // WTF?
		return  // Need return error
	}

	dbConnect.Exec("UPDATE categories SET lft = lft + 2, rgt = rgt + 2 WHERE lft > ?", parentCategory.Rgt)
	dbConnect.Exec("UPDATE categories SET rgt = rgt + 2 WHERE rgt >= ? AND lft < ?", parentCategory.Rgt, parentCategory.Rgt)
	dbConnect.Exec("INSERT INTO categories SET lft = ?, rgt = ? + 1, level = ? + 1, name = ? , parent_id = ?, created_at=?, updated_at = ?",
		parentCategory.Rgt, parentCategory.Rgt, parentCategory.Level, name, parentId, time.Now(), time.Now())

	fmt.Println("Rebuild!")
}

func DeleteCategory(categoryID int) {
	databaseConnect := db.GetDatabaseConnectInstance()
	dbConnect := databaseConnect.GetDbConnect()

	var category models.Category
	dbConnect.First(&category, categoryID)

	dbConnect.Exec("DELETE FROM categories WHERE lft >= ? AND rgt <= ?", category.Lft, category.Rgt)
	dbConnect.Exec("UPDATE categories SET rgt = rgt - (? - ? + 1) WHERE rgt > ? AND lft < ?", category.Rgt, category.Lft,
		category.Rgt, category.Lft)
	dbConnect.Exec("UPDATE categories SET lft = lft - (? - ? + 1), rgt = rgt - (? - ? + 1) WHERE lft > ?", category.Rgt,
		category.Lft, category.Rgt, category.Lft, category.Rgt)
}

func UserIsAdmin(userID int) bool {
	// Function for check admin rights by user id
	databaseConnect := db.GetDatabaseConnectInstance()
	dbConnect := databaseConnect.GetDbConnect()
	var user models.User

	userResult := dbConnect.First(&user, userID)
	if userResult.Error != nil {
		fmt.Println(userResult.Error)
		return false
	}

	return user.IsAdmin
}

func UpdateProduct(productID int, name string, price float32) {
	// Function for update product name and price
	databaseConnect := db.GetDatabaseConnectInstance()
	dbConnect := databaseConnect.GetDbConnect()

	var product models.Product
	userResult := dbConnect.First(&product, productID)
	if userResult.Error != nil {
		fmt.Println(userResult.Error)
		return
	}

	if name != "" {
		product.Name = name
	}
	if price != 0 {
		product.Price = price
	}

	dbConnect.Save(&product)
}


func GeneratePasswordHash(password string) string {
	// Function for generate hash from custom string (ex. password)
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
	}

	return string(hash)
}