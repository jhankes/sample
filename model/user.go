package model

import (
	"log"
	"regexp"

	"github.com/jhankes/sample/db"

	"github.com/gofiber/fiber"
)

// CheckUserExists : check if a user exists in the db by userid
func CheckUserExists(userid string) bool {
	exists := true
	database := db.DBConn

	//valid userid, check if it already exists
	var existing User
	database.Where("user_id = ?", userid).First(&existing)
	if existing.UserID == "" {
		exists = false
	}

	return exists
}

// GetUser : get a user and it's associated groups
func GetUser(c *fiber.Ctx) {

	// get userid and validate
	userid := c.Params("userid")
	matched, err := regexp.MatchString(`^[A-Za-z0-9]+$`, userid)
	if err != nil || !matched {
		log.Printf("GetUser: Invalid userid: %s", userid)
		c.Status(400).Send("Invalid userid")
		return
	}

	database := db.DBConn
	var user User
	database.Where("user_id = ?", userid).First(&user)
	if user.First == "" {
		log.Printf("GetUser: no user with userid %s", userid)
		c.Status(404).Send("No user with userid " + userid)
		return
	}

	groups := GetAssociationsByUser(userid)
	c.JSON(UserJSON{UserID: userid, First: user.First, Last: user.Last, Groups: groups})
}

// NewUser : create a new user
func NewUser(c *fiber.Ctx) {
	c.Accepts("json", "text")
	c.Accepts("application/json")
	if (c.Get("Content-Type") != "application/json") && (c.Get("Content-Type") != "json") {
		log.Println("NewUser: Content-Type header for json must be present")
		c.Status(400).Send("Missing/incorrect Content-Type header")
		return
	}

	userJSON := new(UserJSON)
	if err := c.BodyParser(userJSON); err != nil {
		//log.Fatal(err)
		c.Status(400).Send(err)
		return
	}

	database := db.DBConn

	//valid user body, check if it already exists
	var existing User
	database.Where("user_id = ?", userJSON.UserID).First(&existing)
	if existing.First != "" {
		log.Printf("NewUser: Duplicate userid: %s", userJSON.UserID)
		c.Status(409).Send("Duplicate userid" + userJSON.UserID)
		return
	}

	// ensure groups exist already for this new user
	groups := userJSON.Groups
	for _, g := range groups {

		//check if group exists, fail out if there is a group that does not
		if !CheckGroupExists(g) {
			log.Printf("NewUser: cannot create user, group does not exist %s", g)
			c.Status(404).Send("Cannot create user, group does not exist " + g)
			return
		}
	}

	// create the user
	user := User{UserID: userJSON.UserID, First: userJSON.First, Last: userJSON.Last}
	database.Create(&user)

	// create user associations
	for _, g := range groups {
		CreateAssociation(userJSON.UserID, g)
	}

	c.Status(201).JSON(userJSON)
}

// UpdateUser : updates an existing group and its associations
func UpdateUser(c *fiber.Ctx) {
	c.Accepts("json", "text")
	c.Accepts("application/json")
	if (c.Get("Content-Type") != "application/json") && (c.Get("Content-Type") != "json") {
		log.Println("UpdateUser: Content-Type header for json must be present")
		c.Status(400).Send("Missing/incorrect Content-Type header")
		return
	}

	// get userid and validate
	userid := c.Params("userid")
	matched, err := regexp.MatchString(`^[A-Za-z0-9]+$`, userid)
	if err != nil || !matched {
		log.Printf("GetUser: Invalid userid: %s", userid)
		c.Status(400).Send("Invalid userid")
		return
	}

	userJSON := new(UserJSON)
	if err := c.BodyParser(userJSON); err != nil {
		//log.Fatal(err)
		c.Status(400).Send(err)
		return
	}

	database := db.DBConn

	//valid user body, check if it already exists
	var existing User
	database.Where("user_id = ?", userJSON.UserID).First(&existing)
	if existing.First == "" {
		log.Printf("UpdateUser: User does not exist: %s", userJSON.UserID)
		c.Status(404).Send("User does not exist: userid" + userJSON.UserID)
		return
	}

	// make sure id matches content
	if userid != userJSON.UserID {
		log.Printf("UpdateUser: User id and content mismatch: %s %s", userid, userJSON.UserID)
		c.Status(404).Send("User does not exist: userid " + " " + userid + userJSON.UserID)
		return
	}

	// ensure groups exist already for this new user
	groups := userJSON.Groups
	newGroupMap := make(map[string]bool)
	for _, g := range groups {
		newGroupMap[g] = true
		//check if group exists, fail out if there is a group that does not
		if !CheckGroupExists(g) {
			log.Printf("UpdateUser: Cannot update user, group does not exist: %s", g)
			c.Status(404).Send("Cannot update user, group does not exist " + g)
			return
		}
	}

	// update the user if needed
	if (userJSON.First != existing.First) || (userJSON.Last != existing.Last) {
		user := User{UserID: userJSON.UserID, First: userJSON.First, Last: userJSON.Last}
		database.Save(&user)
	}

	// update user associations
	existingGroups := GetAssociationsByUser(userJSON.UserID)
	existingGroupMap := make(map[string]bool)
	for _, existingGroup := range existingGroups {
		existingGroupMap[existingGroup] = true
		// group has been removed for user, delete association
		if !newGroupMap[existingGroup] {
			DeleteAssociationByID(userJSON.UserID + existingGroup)
		}
	}
	for _, newGroup := range groups {
		// group has been added to user, add association
		if !existingGroupMap[newGroup] {
			CreateAssociation(userJSON.UserID, newGroup)
		}
	}

	c.JSON(userJSON)

}

// DeleteUser : deletes a user and its group associations
func DeleteUser(c *fiber.Ctx) {

	//get userid and validate
	userid := c.Params("userid")
	matched, err := regexp.MatchString(`^[A-Za-z0-9]+$`, userid)
	if err != nil || !matched {
		log.Printf("DeleteUser: invalid userid: %s", userid)
		c.Status(400).Send("Invalid userid")
		return
	}

	// check database, if user exists delete
	database := db.DBConn
	var existing User
	database.Where("user_id = ?", userid).First(&existing)
	if existing.First == "" {
		log.Printf("DeleteUser: userid does not exist %s", userid)
		c.Status(404).Send("userid does not exist" + userid)
		return
	}

	// delete associations
	DeleteAssociationByUser(userid)

	// delete user
	database.Delete(&existing)

	c.Send("User deleted")
}

// User : user db model
type User struct {
	UserID string `gorm:"primary_key" json:"userid" validator:"nonzero,regexp=^[a-zA-Z0-9]+$"`
	First  string `json:"first_name" validator:"nonzero,regexp=^[a-zA-Z]+$"`
	Last   string `json:"last_name" validator:"nonzero,regexp=^[a-zA-Z]+$"`
}

// UserJSON : user web model
type UserJSON struct {
	UserID string   `json:"userid" validator:"nonzero,regexp=^[a-zA-Z0-9]+$"`
	First  string   `json:"first_name" validator:"nonzero,regexp=^[a-zA-Z]+$"`
	Last   string   `json:"last_name" validator:"nonzero,regexp=^[a-zA-Z]+$"`
	Groups []string `json:"groups" validator:"nonzero,regexp=^[a-zA-Z0-9]+$"`
}
