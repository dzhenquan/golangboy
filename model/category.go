package model

import (
	"errors"
	"fmt"
)

// Table Category
type Category struct {
	BaseModel
	Name 	string		`json:"name"`
	Count   uint64		`json:"count"`
}


// Insert Category
func (cate *Category) Insert() error {
	return DB.Create(cate).Error
}

// Get Category By Id
func GetCategoryById(cateId uint64) (*Category, error) {

	var newCate Category

	err := DB.First(&newCate, "id = ?", cateId).Error

	return &newCate, err
}

// Get Category by name
func GetCategoryByName(cateName string) (*Category, error) {

	if len(cateName) == 0 {
		err := errors.New("分类名不为空!")
		return nil, err
	}

	var findCate Category

	err := DB.Where("name = ?", cateName).First(&findCate).Error
	if err == nil {
		return &findCate, nil
	}

	return nil, err
}

func GetCateCountByUserId(userID uint64) uint64 {

	var cateCount uint64

	row := DB.Raw("select COUNT(DISTINCT category.name) from category inner join " +
		"article on category.id=article.category_id where article.user_id = ?", userID).Row()
	err := row.Scan(&cateCount)
	if err != nil {
		cateCount = 0
	}

	return cateCount
}


func GetCategoryQuerys() ([]*Category, error) {

	var categorys []*Category

	sql := "select DISTINCT category.id, category.created_at, category.updated_at, " +
		"category.name, category.count from article INNER JOIN category " +
			"WHERE article.category_id=category.id;"
	rows, err := DB.Raw(sql).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {

		var cate Category

		err := DB.ScanRows(rows, &cate)
		if err != nil {
			return nil, err
		}

		categorys = append(categorys, &cate)
	}
	return categorys, nil
}

func GetCategoryQuerysByUserId(userId uint64) ([]*Category, error) {

	var categorys []*Category

	sql := fmt.Sprintf("select DISTINCT category.id, category.created_at, category.updated_at, " +
		"category.name, category.count from article INNER JOIN category " +
		"WHERE article.category_id=category.id and article.user_id = %d ORDER BY category.updated_at desc;", userId)
	rows, err := DB.Raw(sql).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {

		var cate Category

		err := DB.ScanRows(rows, &cate)
		if err != nil {
			return nil, err
		}

		cate.Count = GetArticleCountByUserCateId(userId, cate.ID)

		categorys = append(categorys, &cate)
	}
	return categorys, nil
}


func AdminGetArticleByCategory(userId uint64, cateId uint64) ([]*ArticleCateInfo, error) {

	var articles []*ArticleCateInfo

	sql := fmt.Sprintf("select DISTINCT article.id,article.title,category.name,article.is_published, " +
		"article.created_at, article.updated_at from category INNER JOIN article WHERE article.category_id=category.id and " +
			"article.category_id=%d and article.user_id=%d ORDER BY article.updated_at desc;", cateId, userId)
	rows, err := DB.Raw(sql).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var article ArticleCateInfo

		err := DB.ScanRows(rows, &article)
		if err != nil {
			return nil, err
		}

		articles = append(articles, &article)
	}
	return articles, nil
}