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