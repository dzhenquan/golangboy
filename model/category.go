package model


// Table Category
type Category struct {
	BaseModel
	Name 	string
	Index   uint
}

// Get Category
func GetCategoryById(cateId uint64) (*Category, error) {

	var newCate Category
	err := DB.First(&newCate, "id = ?", cateId).Error
	return &newCate, err
}


func GetCateCountByUserId(userID uint64) uint64 {

	var cateCount uint64

	row := DB.Raw("select COUNT(DISTINCT category.name) from category inner join article on category.id=article.category_id where article.user_id = ?", userID).Row()
	row.Scan(&cateCount)

	return cateCount
}