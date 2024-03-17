package handler

import (
	"time"

	"github.com/antiloger/nhostel-go/config"
	"github.com/antiloger/nhostel-go/database"
	"github.com/antiloger/nhostel-go/middlewares"
	"github.com/antiloger/nhostel-go/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func Login(c *fiber.Ctx) error {
	loginreq := new(models.LoginRequest)
	if err := c.BodyParser(loginreq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	user, err := middlewares.CheckLogin(loginreq.Email, loginreq.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": err, "data": nil})
	}

	if !user.Approved {
		return c.Status(200).JSON(fiber.Map{"status": "error", "message": "user not approved", "data": nil})
	}

	day := time.Hour * 24
	claims := jwt.MapClaims{
		"id":       user.ID,
		"email":    user.Email,
		"role":     user.Role,
		"approved": user.Approved,
		"exp":      time.Now().Add(day * 1).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(config.Jwt_Secret))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(models.LoginResponse{
		Token: t,
	})
}

func Insertuser(c *fiber.Ctx) error {
	db := database.DB.Db
	user := new(models.UserInfo)
	err := c.BodyParser(user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Somthing's wrong with your input", "data": err})
	}

	err = db.Create(&user).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "could not created the user", "data": err})
	}

	return c.Status(201).JSON(fiber.Map{"status": "success", "message": "user has created", "data": user})
}

func Hello(c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
}

// Home & Search Handler

func HomeLoad(c *fiber.Ctx) error {
	db := database.DB.Db
	hostels := []models.Hostel{}
	if err := db.Where("available = ? AND nsbm_approved = ?", true, true).Find(&hostels).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "could not find the hostels", "data": err})
	}

	if len(hostels) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "hostels not found", "data": nil})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "hostels has found", "data": hostels})
}

func Search(c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
}

func Hosteldetails(c *fiber.Ctx) error {
	db := database.DB.Db
	id := c.Params("ID")
	hostel := models.Hostel{}

	if err := db.Where("id = ?", id).First(&hostel).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "could not find the hostel", "data": err})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "hostel has found", "data": hostel})
}

// user: student handler

func Studentsignup(c *fiber.Ctx) error {
	db := database.DB.Db
	student_sign := new(models.StudentSingup)
	err := c.BodyParser(student_sign)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Somthing's wrong with your input", "data": err})
	}

	// add hash password

	user := models.UserInfo{
		Email:    student_sign.Email,
		Password: student_sign.Password,
		Role:     "student",
		Approved: false,
	}

	if err := db.Create(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "could not created the user", "data": err})
	}

	student := models.Student{
		StdName: student_sign.StdName,
		BOD:     student_sign.BOD,
		Batch:   student_sign.Batch,
		StdNo:   student_sign.StdNo,
		UserID:  user.ID,
	}

	if err := db.Create(&student).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "could not created the student", "data": err})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "student has created",
	})
}

// user: hostel owner handler

func Hostelownersignup(c *fiber.Ctx) error {
	db := database.DB.Db
	hostelowner_sign := new(models.HostelOwnerSingup)
	err := c.BodyParser(hostelowner_sign)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Somthing's wrong with your input", "data": err})
	}

	user := models.UserInfo{
		Email:    hostelowner_sign.Email,
		Password: hostelowner_sign.Password,
		Role:     "hostelowner",
		Approved: false,
	}

	if err := db.Create(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "could not created the user", "data": err})
	}

	owner := models.HostelOwner{
		OwnerName: hostelowner_sign.OwnerName,
		BOD:       hostelowner_sign.BOD,
		Address:   hostelowner_sign.Address,
		PhoneNo:   hostelowner_sign.PhoneNo,
		NIC:       hostelowner_sign.NIC,
		Image:     hostelowner_sign.Image,
		UserID:    user.ID,
	}

	if err := db.Create(&owner).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "could not created the hostel owner", "data": err})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "hostel owner has created",
	})
}

func HostelOwnerView(c *fiber.Ctx) error {
	db := database.DB.Db
	id := c.Params("ID")
	owner := models.HostelOwner{}

	if err := db.Where("id = ?", id).First(&owner).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "could not find the hostel owner", "data": err})
	}

	user := models.UserInfo{}
	if err := db.Where("id = ?", owner.UserID).First(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "could not find the user", "data": err})
	}

	hostel := []models.Hostel{}
	if err := db.Where("owner_id = ?", id).Find(&hostel).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "could not find the hostels", "data": err})
	}

	ownerview := models.OwnerView{
		ID:        owner.ID,
		OwnerName: owner.OwnerName,
		Address:   owner.Address,
		PhoneNo:   owner.PhoneNo,
		Email:     user.Email,
		Approved:  user.Approved,
		Hostels:   hostel,
	}

	return c.JSON(fiber.Map{"status": "success", "message": "hostel owner has found", "data": ownerview})
}

// user: hostel handler

func Hostelcreate(c *fiber.Ctx) error {
	db := database.DB.Db
	hostel_reg := new(models.HostelReg)
	err := c.BodyParser(hostel_reg)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Somthing's wrong with your input", "data": err})
	}

	hostel := models.Hostel{
		HostelName:   hostel_reg.HostelName,
		Address:      hostel_reg.Address,
		Location:     hostel_reg.Location,
		PhoneNo:      hostel_reg.PhoneNo,
		Image1:       hostel_reg.Image1,
		Image2:       hostel_reg.Image2,
		Image3:       hostel_reg.Image3,
		OwnerID:      hostel_reg.OwnerID,
		Rooms:        hostel_reg.Rooms,
		BathRooms:    hostel_reg.BathRooms,
		Price:        hostel_reg.Price,
		PriceInfo:    hostel_reg.PriceInfo,
		Description:  hostel_reg.Description,
		PostedAt:     time.Now(),
		NsbmApproved: false,
		Available:    true,
		Rating:       0,
	}

	if err := db.Create(&hostel).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "could not created the hostel", "data": err})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "hostel has created",
	})
}

func Hostelupdate(c *fiber.Ctx) error {
	db := database.DB.Db

	id, err := c.ParamsInt("ID")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "invalid id", "data": err})
	}

	hostelreg := new(models.HostelReg)
	if err := c.BodyParser(hostelreg); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "invalid body", "data": err})
	}

	var hostel models.Hostel
	if err := db.Where("id = ?", id).First(&hostel).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "hostel not found", "data": err})
	}

	hostel.HostelName = hostelreg.HostelName
	hostel.Address = hostelreg.Address
	hostel.Location = hostelreg.Location
	hostel.PhoneNo = hostelreg.PhoneNo
	hostel.Image1 = hostelreg.Image1
	hostel.Image2 = hostelreg.Image2
	hostel.Image3 = hostelreg.Image3
	hostel.Rooms = hostelreg.Rooms
	hostel.BathRooms = hostelreg.BathRooms
	hostel.Price = hostelreg.Price
	hostel.PriceInfo = hostelreg.PriceInfo
	hostel.Description = hostelreg.Description
	hostel.NsbmApproved = false

	if err := db.Save(&hostel).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "could not update the hostel", "data": err})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "hostel has updated",
	})
}

func HostelAvailableUpdate(c *fiber.Ctx) error {
	db := database.DB.Db
	id, err := c.ParamsInt("ID")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "invalid id", "data": err})
	}
	var hostel models.Hostel
	if err := db.Where("id = ?", id).First(&hostel).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "hostel not found", "data": err})
	}
	type Availabletype struct {
		Available bool `json:"available"`
	}
	var availablet Availabletype
	if err := c.BodyParser(&availablet); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "invalid body", "data": err})
	}
	hostel.Available = availablet.Available
	if err := db.Save(&hostel).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "could not update the hostel", "data": err})
	}
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "hostel has updated",
	})
}

func Hosteldelete(c *fiber.Ctx) error {
	db := database.DB.Db
	id, err := c.ParamsInt("ID")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "invalid id", "data": err})
	}
	var hostel models.Hostel

	if err := db.Where("id = ?", id).First(&hostel).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "hostel not found", "data": err})
	}

	if err := db.Delete(&hostel).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "could not delete the hostel", "data": err})
	}
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "hostel has deleted",
	})
}

// user: admin handler

func Adminregister(c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
}

func HostelApproveTable(c *fiber.Ctx) error {
	db := database.DB.Db
	hostels := []models.Hostel{}
	if err := db.Where("nsbm_approved = ?", false).Find(&hostels).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "could not find the hostels", "data": err})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "hostels has found", "data": hostels})
}

func HostelOwnerApproveTable(c *fiber.Ctx) error {
	db := database.DB.Db

	owners := []models.HostelOwner{}
	if err := db.Table("hostel_owners").Joins("INNER JOIN user_infos ON hostel_owners.user_id = user_infos.id").Where("user_infos.approved = ?", false).Find(&owners).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "could not find the hostel owners", "data": err})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "hostel owners has found", "data": owners})
}

func StudentApproveTable(c *fiber.Ctx) error {
	db := database.DB.Db

	students := []models.Student{}
	if err := db.Table("students").Joins("INNER JOIN user_infos ON students.user_id = user_infos.id").Where("user_infos.approved = ?", false).Find(&students).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "could not find the students", "data": err})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "students has found", "data": students})
}

func StudentApprove(c *fiber.Ctx) error {
	// this is for admin to approve the student
	db := database.DB.Db
	id, err := c.ParamsInt("ID")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "invalid id", "data": err})
	}
	var student models.Student
	if err := db.Where("id = ?", id).First(&student).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "student not found", "data": err})
	}
	type Approvetype struct {
		Approved bool `json:"approved"`
	}
	var approvet Approvetype
	if err := c.BodyParser(&approvet); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "invalid body", "data": err})
	}
	user := models.UserInfo{}
	if err := db.Where("id = ?", student.UserID).First(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "could not find the user", "data": err})
	}
	user.Approved = approvet.Approved
	if err := db.Save(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "could not approve the student", "data": err})
	}
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "student has approved",
	})
}

func OwnerApprove(c *fiber.Ctx) error {
	// this func is for approve the hostel owner
	db := database.DB.Db
	id, err := c.ParamsInt("ID")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "invalid id", "data": err})
	}
	var owner models.HostelOwner
	if err := db.Where("id = ?", id).First(&owner).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "hostel owner not found", "data": err})
	}
	type Approvetype struct {
		Approved bool `json:"approved"`
	}
	var approvet Approvetype
	if err := c.BodyParser(&approvet); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "invalid body", "data": err})
	}
	user := models.UserInfo{}
	if err := db.Where("id = ?", owner.UserID).First(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "could not find the user", "data": err})
	}
	user.Approved = approvet.Approved
	if err := db.Save(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "could not approve the hostel owner", "data": err})
	}
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "hostel owner has approved",
	})
}

func HostelApprove(c *fiber.Ctx) error {
	db := database.DB.Db
	id, err := c.ParamsInt("ID")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "invalid id", "data": err})
	}

	var hostel models.Hostel
	if err := db.Where("id = ?", id).First(&hostel).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "hostel not found", "data": err})
	}

	type Approvetype struct {
		NsbmApproved bool `json:"nsbmapproved"`
	}

	var approvet Approvetype
	if err := c.BodyParser(&approvet); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "invalid body", "data": err})
	}

	hostel.NsbmApproved = approvet.NsbmApproved

	if err := db.Save(&hostel).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "could not approve the hostel", "data": err})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "hostel has approved",
	})
}
