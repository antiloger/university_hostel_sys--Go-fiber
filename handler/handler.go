package handler

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/antiloger/nhostel-go/config"
	"github.com/antiloger/nhostel-go/database"
	"github.com/antiloger/nhostel-go/middlewares"
	"github.com/antiloger/nhostel-go/models"
	"github.com/antiloger/nhostel-go/util"
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
		Role:  user.Role,
		ID:    user.ID,
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

func SearchHostel(c *fiber.Ctx) error {
	db := database.DB.Db
	hostels := []models.Hostel{}
	search := c.Query("search")
	if err := db.Where("hostel_name LIKE ?", "%"+search+"%").Or("address Like ?", "%"+search+"%").Find(&hostels).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "could not find the hostels", "data": err})
	}

	if len(hostels) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "hostels not found", "data": nil})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "hostels has found", "data": hostels})
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
	fmt.Println(student_sign)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Somthing's wrong with your input", "data": err})
	}

	// add hash password

	user := models.UserInfo{
		Email:    student_sign.Email,
		Password: student_sign.Password,
		Role:     "student",
		Approved: true,
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

	image, err := c.FormFile("image")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "could not get the image", "data": err})
	}

	if image.Size > 0 {
		imagename := "./uploads/hostelowner" + hostelowner_sign.OwnerName + image.Filename

		if err := c.SaveFile(image, imagename); err != nil {
			return c.Status(500).JSON(fiber.Map{"status": "error", "message": "could not save the image", "data": err})
		}

		hostelowner_sign.Image = imagename
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
	user_id := c.Params("ID") // this id is from user_infos table
	owner := models.HostelOwner{}

	if err := db.Where("user_id = ?", user_id).First(&owner).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "could not find the hostel owner", "data": err})
	}

	user := models.UserInfo{}
	if err := db.Where("id = ?", user_id).First(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "could not find the user", "data": err})
	}

	hostel := []models.Hostel{}
	if err := db.Where("owner_id = ?", owner.ID).Find(&hostel).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "could not find the hostels", "data": err})
	}

	ownerview := models.OwnerView{
		ID:        owner.ID,
		OwnerName: owner.OwnerName,
		Address:   owner.Address,
		PhoneNo:   owner.PhoneNo,
		Email:     user.Email,
		Approved:  user.Approved,
		Image:     owner.Image,
		Hostels:   hostel,
	}

	return c.JSON(fiber.Map{"status": "success", "message": "hostel owner has found", "data": ownerview})
}

// user: warden handler

func Wardensignup(c *fiber.Ctx) error {
	db := database.DB.Db
	warden_sign := new(models.WardenSingup)
	err := c.BodyParser(warden_sign)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Somthing's wrong with your input", "data": err})
	}

	image, err := c.FormFile("image")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "could not get the image", "data": err})
	}

	if image.Size > 0 {
		imagename := "./uploads/warden" + warden_sign.WardenName + image.Filename

		if err := c.SaveFile(image, imagename); err != nil {
			return c.Status(500).JSON(fiber.Map{"status": "error", "message": "could not save the image", "data": err})
		}

		warden_sign.Image = imagename
	}

	user := models.UserInfo{
		Email:    warden_sign.Email,
		Password: warden_sign.Password,
		Role:     "warden",
		Approved: true,
	}

	if err := db.Create(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "could not created the user", "data": err})
	}

	warden := models.Warden{
		WardenName: warden_sign.WardenName,
		PhoneNo:    warden_sign.PhoneNo,
		NIC:        warden_sign.NIC,
		Image:      warden_sign.Image,
		UserID:     user.ID,
	}

	if err := db.Create(&warden).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "could not created the warden", "data": err})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "warden has created",
	})
}

// user: hostel handler

func Hostelcreate(c *fiber.Ctx) error {
	db := database.DB.Db
	hostel_reg := new(models.HostelReg)
	err := c.BodyParser(hostel_reg)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Somthing's wrong with your input", "data": err})
	}

	img1, err := c.FormFile("image1")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "could not get the image1", "data": err})
	}
	img2, err := c.FormFile("image2")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "could not get the image2", "data": err})
	}
	img3, err := c.FormFile("image3")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "could not get the image3", "data": err})
	}

	hostel := models.Hostel{
		HostelName:   hostel_reg.HostelName,
		Address:      hostel_reg.Address,
		Lat:          hostel_reg.Lat,
		Lng:          hostel_reg.Lng,
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

	if img1.Size > 0 || img2.Size > 0 || img3.Size > 0 {
		imagefolder := fmt.Sprintf("./uploads/hostel/%s/", hostel_reg.HostelName)
		err := os.Mkdir(imagefolder, 0755)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"status": "error", "message": "could not create the image folder", "data": err})
		}
		img1path := imagefolder + img1.Filename
		img2path := imagefolder + img2.Filename
		img3path := imagefolder + img3.Filename

		if err := c.SaveFile(img1, img1path); err != nil {
			return c.Status(500).JSON(fiber.Map{"status": "error", "message": "could not save the image1", "data": err})
		}
		if err := c.SaveFile(img2, img2path); err != nil {
			return c.Status(500).JSON(fiber.Map{"status": "error", "message": "could not save the image2", "data": err})
		}

		if err := c.SaveFile(img3, img3path); err != nil {
			return c.Status(500).JSON(fiber.Map{"status": "error", "message": "could not save the image3", "data": err})
		}

		hostel.Image1 = img1path
		hostel.Image2 = img2path
		hostel.Image3 = img3path
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
	hostel.Lat = hostelreg.Lat
	hostel.Lng = hostelreg.Lng
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

func Hostelownertable(c *fiber.Ctx) error {
	db := database.DB.Db
	owners := []models.HostelOwner{}
	if err := db.Find(&owners).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "could not find the hostel owners", "data": err})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "hostel owners has found", "data": owners})
}

func StudentTable(c *fiber.Ctx) error {
	db := database.DB.Db
	students := []models.Student{}
	if err := db.Find(&students).Error; err != nil {
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

func CreateArticle(c *fiber.Ctx) error {
	db := database.DB.Db
	article := new(models.Article)
	err := c.BodyParser(article)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Somthing's wrong with your input", "data": err})
	}

	if err := db.Create(&article).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "could not created the article", "data": err})
	}

	return c.Status(201).JSON(fiber.Map{"status": "success", "message": "article has created", "data": article})
}

func GetArticles(c *fiber.Ctx) error {
	db := database.DB.Db
	articles := []models.Article{}
	if err := db.Find(&articles).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "could not find the articles", "data": err})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "articles has found", "data": articles})
}

func GetArticle(c *fiber.Ctx) error {
	db := database.DB.Db
	id, err := c.ParamsInt("ID")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "invalid id", "data": err})
	}
	article := models.Article{}
	if err := db.Where("id = ?", id).First(&article).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "article not found", "data": err})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "article has found", "data": article})
}

func DeleteArticle(c *fiber.Ctx) error {
	db := database.DB.Db
	id, err := c.ParamsInt("ID")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "invalid id", "data": err})
	}
	var article models.Article

	if err := db.Where("id = ?", id).First(&article).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "article not found", "data": err})
	}

	if err := db.Delete(&article).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "could not delete the article", "data": err})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "article has deleted"})
}

func GetMyProfile(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	fmt.Print(tokenString)
	userID, role, err := util.ValidateTokenAndExtractClaims(tokenString)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid or expired JWT", "data": err})
	}

	if role == "hostelowner" {
		db := database.DB.Db
		owner := models.HostelOwner{}
		if err := db.Table("hostel_owners").Joins("INNER JOIN user_infos ON hostel_owners.user_id = user_infos.id").Where("user_infos.id = ?", userID).First(&owner).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "could not find the hostel owner", "data": err})
		}

		return c.JSON(fiber.Map{"status": "success", "message": "hostel owner has found", "data": owner, "role": role})
	} else if role == "warden" {
		db := database.DB.Db
		warden := models.Warden{}
		if err := db.Table("wardens").Joins("INNER JOIN user_infos ON wardens.user_id = user_infos.id").Where("user_infos.id = ?", userID).First(&warden).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "could not find the warden", "data": err})
		}

		return c.JSON(fiber.Map{"status": "success", "message": "warden has found", "data": warden, "role": role})
	} else if role == "student" {
		return c.JSON(fiber.Map{"status": "success", "message": "student has found", "data": nil, "role": role})
	} else {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid role"})
	}
}

func Getwardendetails(c *fiber.Ctx) error {
	db := database.DB.Db
	warden := models.WardenSingup{}
	war_user := models.UserInfo{}
	war_info := models.Warden{}
	user_id := c.Params("ID")

	if err := db.Table("user_infos").Where("id = ?", user_id).First(&war_user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "could not find the warden user", "data": err})
	}

	if err := db.Table("wardens").Where("user_id = ?", user_id).First(&war_info).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "could not find the warden info", "data": err})
	}

	warden.Email = war_user.Email
	warden.WardenName = war_info.WardenName
	warden.PhoneNo = war_info.PhoneNo
	warden.NIC = war_info.NIC
	warden.Image = war_info.Image

	return c.JSON(fiber.Map{"status": "success", "message": "warden has found", "data": warden})
}

func WardenTable(c *fiber.Ctx) error {
	db := database.DB.Db
	wardens := []models.Warden{}

	if err := db.Find(&wardens).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "could not find the wardens", "data": err})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "wardens has found", "data": wardens})

}
