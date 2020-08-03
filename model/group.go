package model

import (
	"log"
	"regexp"
	"strings"

	"github.com/jhankes/sample/db"

	"github.com/gofiber/fiber"
)

// CheckGroupExists : Check if a group exists in the database by name
func CheckGroupExists(name string) bool {
	exists := true
	database := db.DBConn

	//valid group name, check if it already exists
	var existing Group
	database.Where("name = ?", name).First(&existing)
	if existing.Name == "" {
		exists = false
	}

	return exists
}

// GetGroup :  get a group's userids
func GetGroup(c *fiber.Ctx) {
	// get group name and validate
	name := c.Params("name")
	matched, err := regexp.MatchString(`^[A-Za-z0-9]+$`, name)
	if err != nil || !matched {
		log.Printf("GetGroup: Invalid group name: %s", name)
		c.Status(400).Send("Invalid group name")
		return
	}

	database := db.DBConn

	// valid group name, check if it already exists
	var existing Group
	database.Where("name = ?", name).First(&existing)
	if existing.Name == "" {
		log.Printf("GetGroup: Group does not exist: %s", name)
		c.Status(404).Send("Group does not exist" + name)
		return
	}

	// group exists, get all user association userids and return as JSON list
	userids := GetAssociationsByGroupName(name)

	c.JSON(userids)

}

// NewGroup : make a new group with a name only
func NewGroup(c *fiber.Ctx) {

	body := c.Body()
	matched, err := regexp.MatchString(`^name=[A-Za-z0-9]+$`, body)
	if err != nil || !matched {
		log.Println("NewGroup: invalid group name")
		c.Status(400).Send("Invalid group name")
		return
	}

	// parse group name
	group := new(Group)
	group.Name = strings.Replace(body, "name=", "", -1)

	database := db.DBConn

	//valid group name, check if it already exists
	var existing Group
	database.Where("name = ?", group.Name).First(&existing)
	if existing.Name != "" {
		log.Printf("NewGroup: Duplicate group: %s", group.Name)
		c.Status(409).Send("Duplicate group name" + group.Name)
		return
	}

	//otherwise create the group
	database.Create(&group)
	userids := GetAssociationsByGroupName(group.Name)

	c.JSON(userids)
}

// UpdateGroup : modify a group's users
func UpdateGroup(c *fiber.Ctx) {
	c.Accepts("json", "text")
	c.Accepts("application/json")
	if (c.Get("Content-Type") != "application/json") && (c.Get("Content-Type") != "json") {
		log.Println("UpdateGroup: Content-Type header for json must be present")
		c.Status(400).Send("Missing/incorrect Content-Type header")
		return
	}

	// get group name and validate
	name := c.Params("name")
	matched, err := regexp.MatchString(`^[A-Za-z0-9]+$`, name)
	if err != nil || !matched {
		log.Printf("UpdateGroup: invalid group name: %s", name)
		c.Status(400).Send("Invalid group name")
		return
	}

	database := db.DBConn

	//valid group name, check if it already exists
	var existing Group
	database.Where("name = ?", name).First(&existing)
	if existing.Name == "" {
		log.Printf("UpdateGroup: Group does not exist: %s", name)
		c.Status(404).Send("Group does not exist" + name)
		return
	}

	// validate userids and check users exist
	userIDs := new(UserIDs)
	if err := c.BodyParser(userIDs); err != nil {
		log.Printf("UpdateGroup: Group body invalid: %s", name)
		c.Status(400).Send(err)
		return
	}

	newUseridMap := make(map[string]bool)
	// generate map for new userids and check users exist
	for _, userid := range userIDs.UserID {
		newUseridMap[userid] = true
		log.Printf("adding %s userid to new userid map", userid)
		//check if user exists, fail out if it does not
		if !CheckUserExists(userid) {
			log.Printf("UpdateGroup: User does not exist: %s", userid)
			c.Status(404).Send("Cannot update group, user does not exist " + userid)
			return
		}
	}

	// check existing ids with new map, create existing map
	existingUseridMap := make(map[string]bool)
	existingUserids := GetAssociationsByGroupName(name).UserID
	for _, userid := range existingUserids {
		existingUseridMap[userid] = true
		// membership has changed, delete existing association
		if !newUseridMap[userid] {
			log.Printf("deleting %s userid association, not in new map", userid)
			DeleteAssociationByID(userid + name)
		}
	}

	// check new userids and if they are not on the existing map add them
	for _, userid := range userIDs.UserID {
		if !existingUseridMap[userid] {
			log.Printf("creating %s and %s association not in existing map", userid, name)
			CreateAssociation(userid, name)
		}
	}

	c.JSON(userIDs)
}

// DeleteGroup : delete a group and any associations
func DeleteGroup(c *fiber.Ctx) {
	// get group name and validate
	name := c.Params("name")
	matched, err := regexp.MatchString(`^[A-Za-z0-9]+$`, name)
	if err != nil || !matched {
		log.Printf("DeleteGroup: Invalid group name: %s", name)
		c.Status(400).Send("Invalid group name")
		return
	}

	database := db.DBConn

	//valid group name, check if it already exists
	var existing Group
	database.Where("name = ?", name).First(&existing)
	if existing.Name == "" {
		log.Printf("DeleteGroup: group does not exist: %s", name)
		c.Status(404).Send("Group does not exist: " + name)
		return
	}

	// delete group
	database.Delete(&existing)

	// delete group user associations
	DeleteAssociationsByGroupName(name)

	// update to send back delete msg
	c.Send("Group deleted")
}

// Group : a named group
type Group struct {
	Name string `gorm:"primary_key" validator:"nonzero,regexp=^[a-zA-Z0-9]+$"`
}
