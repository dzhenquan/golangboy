package model

// Table link
type Link struct {
	BaseModel
	Name 		string 		`json:"name"`			//名称
	Url  		string 		`json:"url"`			//地址
	Sort 		int    		`json:"default:'0'"` 	//排序
	View 		int    		`json:"view"`			//访问次数
	UserID		uint64 		`json:"userId"`
	User  		User		`gorm:"ForeignKey:ID" json:"user"`
}

// Insert Link
func (link *Link) Insert() error {
	return DB.Create(&link).Error
}

//Update Link
func (link *Link) Update() error {
	return DB.Model(link).Updates(map[string]interface{}{
		"name"	:	link.Name,
		"url"	:	link.Url,
		"sort"	:	link.Sort,
	}).Error
}

// Delete Link
func (link *Link) Delete() error {
	return DB.Delete(link).Error
}


// Get Link Querys By UserID
func GetLinkQuerysByUserID(userID uint64) ([]*Link, error) {

	var links []*Link

	rows, err := DB.Raw("select * from link where user_id = ? ORDER BY updated_at desc", userID).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var link Link

		err := DB.ScanRows(rows, &link)
		if err != nil {
			return nil, err
		}

		links = append(links, &link)
	}
	return links, nil
}