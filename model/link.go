package model


// Table link
type Link struct {
	BaseModel
	Name 		string 	//名称
	Url  		string 	//地址
	Sort 		int    	`json:"default:'0'"` //排序
	View 		int    	//访问次数
	UserID		uint64 		`json:"userId"`
	User  		User		`gorm:"ForeignKey:ID" json:"user"`
}


func (link *Link) Insert() error {

	return DB.Create(&link).Error
}

//Update Link
func (link *Link) Update() error {
	return DB.Model(link).Updates(map[string]interface{}{
		"name":        link.Name,
		"url":         link.Url,
		"sort": link.Sort,
	}).Error
}

// Delete Link
func (link *Link) Delete() error {
	return DB.Delete(link).Error
}