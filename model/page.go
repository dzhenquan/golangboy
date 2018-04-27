package model

import "strconv"

// Table Page
type Page struct {
	BaseModel
	Title 		string		`json:"title"`
	Content		string		`json:"content"`
	View		int			`json:"view"`
	UserID 		uint64		`json:"userId"`
	User		User 		`gorm:"ForeignKey:ID" json:"user"`
	IsPublished bool
}


// Insert Page
func (page *Page) Insert() error {
	return DB.Create(&page).Error
}

func (page *Page) Update() error {
	return DB.Model(page).Updates(map[string]interface{}{
		"title":		page.Title,
		"content":		page.Content,
		"is_published":	page.IsPublished,
	}).Error
}

func (page *Page) Delete() error {
	return DB.Delete(page).Error
}

func GetPageById(id string) (*Page, error) {
	pid, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return nil, err
	}
	var page Page
	err = DB.First(&page, "id = ?", pid).Error
	return &page, err
}