package model

import (
	"github.com/jankes/sample/db"
)

func CreateAssociation(userid string, group string) {
	// create association
	database := db.DBConn
	association := Association{ID: userid + group, UserID: userid, GroupName: group}
	database.Create(&association)
}

func DeleteAssociationByID(id string) {
	database := db.DBConn
	database.Where("id = ?", id).Delete(Association{})
}

func DeleteAssociationByUser(userid string) {
	database := db.DBConn
	database.Where("user_id = ?", userid).Delete(Association{})
}

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

func DeleteAssociationsByGroupName(name string) {
	database := db.DBConn
	database.Where("group_name = ?", name).Delete(Association{})
}

type Association struct {
	ID        string `gorm:"primary_key" json:"userid" validator:"nonzero"`
	UserID    string `gorm:"index:userid_idx" json:"userid" validator:"nonzero,regexp=^[a-zA-Z0-9]+$"`
	GroupName string `gorm:"index:group_idx" json:"group" validator:"nonzero,regexp=^[a-zA-Z0-9]+$"`
}

type UserIDs struct {
	UserID []string `json:"userids"`
}
