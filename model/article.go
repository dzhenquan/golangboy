package model

import (
	"strconv"
	"strings"
	"time"
	"encoding/json"
	"fmt"
)

// Table Article
type Article struct {
	BaseModel
	Title			string		`json:"title"`
	Desc  			string		`json:"desc"`
	Content 		string		`json:"content"`
	UserID			uint64 		`json:"userId"`
	ViewCount		uint64 		`json:"viewCount"`
	CommentCount	uint64 		`json:"commentCount"`
	CategoryID		uint64		`json:"categoryId"`
	User  			User		`gorm:"ForeignKey:ID" json:"user"`
	Category		Category 	`gorm:"ForeignKey:ID" json:"category"`
	IsPublished 	bool
}

// ArticleData
type ArticleData struct {
	Title 			string 	`json:"title"`
	ReadCount		uint64 	`json:"readCount"`
	ArticleID		uint64 	`json:"articleID"`
	CommentCount	uint64 	`json:"commentCount"`
	Desc			string 	`json:"desc"`
	Category		string 	`json:"category"`
	CategoryID 		uint64	`json:"categoryID"`
	CreateTime 		string 	`json:"createTime"`
}

// Article Category Info
type ArticleCateInfo struct {
	ID        	uint64 		`gorm:"primary_key"`
	Title		string		`json:"title"`
	Name		string		`json:"name"`
	IsPublished bool
	Count 		uint64 		`json:"count"`
	CreatedAt	time.Time
	UpdatedAt	time.Time
}


// Article Category
type ArticleCate struct {
	CateID 		uint64 	`json:"cateId"`
	CateCount 	uint64	`json:"cateCount"`
	CateName	string 	`json:"cateName"`
}

//Article Archive
type ArticleArch struct {
	Year		int			`json:"year"`
	Month 		string		`json:"month"`
	YearMonth	string 		`json:"yearMonth"`
	Articles	[]*Article	`json:"articles"`
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
		"title"			: article.Title,
		"desc"			: article.Desc,
		"content"		: article.Content,
		"category_id"	: article.CategoryID,
		"is_published"	: article.IsPublished,
	}).Error
}

func (article *Article) UpdateView() error {

	return DB.Model(article).Updates(map[string]interface{}{
		"view_count" : article.ViewCount,
	}).Error
}

func GetArticleById(id string) (*Article, error) {

	var article Article

	pid, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return nil, err
	}

	err = DB.First(&article, "id = ?", pid).Error
	return &article, err
}


func GetArticleJsonData(perPageInt int, offset int) []interface {} {
	var data []interface{}

	rows, err := DB.Raw("select * from article ORDER BY updated_at desc LIMIT ? offset ?",
		perPageInt, offset).Rows()
	if err != nil {
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var article Article
		DB.ScanRows(rows, &article)

		if article.IsPublished {

			user, _ := GetUserById(article.UserID)
			cate, _ := GetCategoryById(article.CategoryID)

			article.User = *user
			article.Category = *cate

			dataJson := GetArticleJson(&article, cate)

			var result map[string]interface{}

			decoder := json.NewDecoder(strings.NewReader(string(dataJson)))
			decoder.Decode(&result)

			data = append(data, result)
		}
	}
	return data
}


func GetArticleJson(article *Article, cate *Category) string {

	articleTime := article.CreatedAt.Format(("2006-01-02 15:04"))

	artData := &ArticleData{
		Title		: article.Title,
		ReadCount	: article.ViewCount,
		ArticleID	: article.ID,
		CommentCount: article.CommentCount,
		Desc		: article.Desc,
		Category	: cate.Name,
		CategoryID	: cate.ID,
		CreateTime	: articleTime,
	}

	jsonData, err := json.Marshal(artData)
	if err != nil {
		return ""
	}

	return string(jsonData)
}


func GetArtileCount() uint64 {
	var articleTotal uint64

	row := DB.Raw("select count(*) from article where is_published = 1;").Row()

	err := row.Scan(&articleTotal)
	if err != nil {
		articleTotal = 0
	}
	return articleTotal
}

func GetArticleCountByCateId(cateId uint64) uint64 {

	var articleCount uint64

	row := DB.Raw("select count(*) from article where category_id = ?", cateId).Row()
	err := row.Scan(&articleCount)
	if err != nil {
		articleCount = 0
	}

	return articleCount
}

func GetArticleCountByUserId(userID uint64) uint64 {

	var articleCount uint64

	row := DB.Raw("select count(*) from article where user_id = ?", userID).Row()
	err := row.Scan(&articleCount)
	if err != nil {
		articleCount = 0
	}

	return articleCount
}


func GetArticleCountByUserCateId(userId uint64, cateId uint64) uint64 {

	var articleCount uint64

	row := DB.Raw("select count(*) from article where user_id = ? and category_id = ?", userId, cateId).Row()
	err := row.Scan(&articleCount)
	if err != nil {
		articleCount = 0
	}

	return articleCount
}

func GetArticleQuerysByCateId(cateId uint64) ([]*ArticleData, error){
	var articles []*ArticleData

	sql := fmt.Sprintf("select * from article where category_id = %d ORDER BY " +
		"created_at desc;", cateId)

	rows, err := DB.Raw(sql).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var article Article

		err := DB.ScanRows(rows, &article)
		if err != nil {
			return nil, err
		}

		createTime := article.CreatedAt.Format(("2006-01-02 15:04"))
		cate, err := GetCategoryById(article.CategoryID)
		if err != nil {
			return nil, err
		}
		articleData := &ArticleData{
			Title			:	article.Title,
			ReadCount		:	article.ViewCount,
			ArticleID		:	article.ID,
			CommentCount	:	article.CommentCount,
			Desc			:	article.Desc,
			Category		:	cate.Name,
			CreateTime		:	createTime,
		}

		articles = append(articles, articleData)
	}
	return articles, nil
}

func GetRecentArticleQuerys() ([]*Article, error) {

	var articles []*Article

	rows, err := DB.Raw("select * from article where is_published = 1 ORDER BY created_at desc LIMIT 6").Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var article Article

		err := DB.ScanRows(rows, &article)
		if err != nil {
			return nil, err
		}

		articles = append(articles, &article)
	}
	return articles, nil
}


func GetArticleCategoryQuerys() ([]*ArticleCate, error) {

	var articleCates []*ArticleCate

	categorys, err := GetCategoryQuerys()
	if err != nil {
		return nil, err
	}

	for _, cate := range categorys {

		articleCate := &ArticleCate{
			CateID	: cate.ID,
			CateName: cate.Name,
		}
		articleCate.CateCount = GetArticleCountByCateId(cate.ID)

		articleCates = append(articleCates, articleCate)
	}
	return articleCates, nil
}


func GetArticleArchiveQuerys() ([]*ArticleArch, error) {

	var articleArchs []*ArticleArch

	rows, err := DB.Raw("select * from article ORDER BY created_at desc").Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {

		var article Article
		err := DB.ScanRows(rows, &article)
		if err != nil {
			return nil, err
		}

		year  := article.CreatedAt.Year()
		month := article.CreatedAt.Format("01")
		yearMonth	 := article.CreatedAt.Format("2006-01")

		articleArch := &ArticleArch{
			Year		: year,
			Month		: month,
			YearMonth	: yearMonth,
		}

		if len(articleArchs) == 0 {
			articleArch.Articles = append(articleArch.Articles, &article)
			articleArchs = append(articleArchs, articleArch)
			continue
		}

		flag := false

		for _, arch := range articleArchs {
			if strings.Compare(arch.YearMonth, yearMonth) == 0 {
				arch.Articles = append(arch.Articles, &article)
				flag = true
				break
			}
		}

		if flag == true {
			continue
		}

		articleArch.Articles = append(articleArch.Articles, &article)
		articleArchs = append(articleArchs, articleArch)
	}

	return articleArchs, nil
}

func GetArticleArchiveQuerysByTime(archTime string) ([]*ArticleData, error) {

	var articleDatas []*ArticleData

	rows, err := DB.Raw("select * from article ORDER BY created_at desc").Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {

		var article Article
		err := DB.ScanRows(rows, &article)
		if err != nil {
			return nil, err
		}

		yearMonth	 := article.CreatedAt.Format("2006-01")
		if strings.Compare(yearMonth, archTime) == 0 {

			createTime := article.CreatedAt.Format(("2006-01-02 15:04"))
			cate, err := GetCategoryById(article.CategoryID)
			if err != nil {
				return nil, err
			}


			articleData := &ArticleData{
				Title		: article.Title,
				ReadCount	: article.ViewCount,
				ArticleID	: article.ID,
				CommentCount: article.CommentCount,
				Desc		: article.Desc,
				Category	: cate.Name,
				CreateTime	: createTime,
			}

			articleDatas = append(articleDatas, articleData)
		}
	}

	return articleDatas, nil
}


func GetArticleQuerysByUserId(userID uint64) ([]*Article, error) {

	var articles []*Article

	sql := fmt.Sprintf("select * from article where user_id = %d ORDER BY updated_at desc;",
		userID)

	rows, err := DB.Raw(sql ).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var article Article

		err := DB.ScanRows(rows, &article)
		if err != nil {
			return nil, err
		}

		category, _ := GetCategoryById(article.CategoryID)
		article.Category = *category

		articles = append(articles, &article)
	}
	return articles, nil
}