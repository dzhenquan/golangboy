package model

import (
	"strconv"
)

// Table Article
type Article struct {
	BaseModel
	Title		string		`json:"title"`
	Desc  		string		`json:"desc"`
	Content 	string		`json:"content"`
	UserID		uint64 		`json:"userId"`
	CategoryID	uint64		`json:"categoryId"`
	User  		User		`gorm:"ForeignKey:ID" json:"user"`
	Category	Category 	`gorm:"ForeignKey:ID" json:"category"`
	IsPublished bool
}

// Post
func (article *Article) Insert() error {
	return DB.Create(article).Error
}

// Delete Article
func (article *Article) Delete() error {
	return DB.Delete(article).Error
}

//Update Article
func (article *Article) Update() error {
	return DB.Model(article).Updates(map[string]interface{}{
		"title":        article.Title,
		"content":         article.Content,
		"category_id":	article.CategoryID,
		"is_published": article.IsPublished,
	}).Error
}

func GetArticleById(id string) (*Article, error) {
	pid, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return nil, err
	}
	var article Article
	err = DB.First(&article, "id = ?", pid).Error
	return &article, err
}


func GetArticleCountByUserId(userID uint64) uint64 {

	var articleCount uint64

	row := DB.Raw("select count(*) from article where user_id = ?", userID).Row()
	row.Scan(&articleCount)

	return articleCount
}