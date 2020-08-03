package model

import (
	"github.com/jhankes/sample/db"
)

//CreateAssociation : by userid and group
func CreateAssociation(userid string, group string) {
	// create association
	database := db.DBConn
	association := Association{ID: userid + group, UserID: userid, GroupName: group}
	database.Create(&association)
}

//DeleteAssociationByID
func DeleteAssociationByID(id string) {
	database := db.DBConn
	database.Where("id = ?", id).Delete(Association{})
}

//DeleteAssociationByUser
func DeleteAssociationByUser(userid string) {
	database := db.DBConn
	database.Where("user_id = ?", userid).Delete(Association{})
}

//GetAssociationsByGroupName
func GetAssociationsByGroupName(name string) UserIDs {
	uids := []string{}

	var associations []Association

	database := db.DBConn
	database.Where("group_name = ?", name).Find(&associations)

	for _, association := range associations {
		uids = append(uids, association.UserID)
	}

	return UserIDs{UserID: uids}
}

//GetAssociationsByUser
func GetAssociationsByUser(userid string) []string {
	uids := []string{}

	var associations []Association

	database := db.DBConn
	database.Where("user_id = ?", userid).Find(&associations)

	for _, association := range associations {
		uids = append(uids, association.GroupName)
	}

	return uids
}

//DeleteAssociationsByGroupName
func DeleteAssociationsByGroupName(name string) {
	database := db.DBConn
	database.Where("group_name = ?", name).Delete(Association{})
}

//Association
type Association struct {
	ID        string `gorm:"primary_key" json:"userid" validator:"nonzero"`
	UserID    string `gorm:"index:userid_idx" json:"userid" validator:"nonzero,regexp=^[a-zA-Z0-9]+$"`
	GroupName string `gorm:"index:group_idx" json:"group" validator:"nonzero,regexp=^[a-zA-Z0-9]+$"`
}

//UserIDs for the group association
type UserIDs struct {
	UserID []string `json:"userids"`
}
