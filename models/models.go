package models

import (
	"time"

	"gorm.io/gorm"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
	Role  string `json:"role"`
	ID    uint   `json:"id"`
}

type UserInfo struct {
	gorm.Model
	ID       uint   `gorm:"primaryKey" json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
	Approved bool   `json:"approved" gorm:"default:false"`
}

type Student struct {
	gorm.Model
	ID      uint   `gorm:"primaryKey" json:"id"`
	StdName string `json:"stdname"`
	BOD     string `json:"bod"`
	Batch   string `json:"batch"`
	StdNo   string `json:"stdno"`
	UserID  uint   `json:"userid"`
}

type HostelOwner struct {
	gorm.Model
	ID        uint   `gorm:"primaryKey" json:"id"`
	OwnerName string `json:"ownername"`
	BOD       string `json:"bod"`
	Address   string `json:"address"`
	PhoneNo   string `json:"phoneno"`
	NIC       string `json:"nic"`
	Image     string `json:"image"`
	UserID    uint   `json:"userid"`
}

type Hostel struct {
	gorm.Model
	ID           uint      `gorm:"primaryKey" json:"id"`
	HostelName   string    `json:"hostelname"`
	Address      string    `json:"address"`
	Lat          float64   `json:"lat"`
	Lng          float64   `json:"lng"`
	PhoneNo      string    `json:"phoneno"`
	Image1       string    `json:"image1"`
	Image2       string    `json:"image2"`
	Image3       string    `json:"image3"`
	OwnerID      uint      `json:"ownerid"`
	Rooms        int       `json:"rooms"`
	BathRooms    int       `json:"bathrooms"`
	Price        float64   `json:"price"`
	PriceInfo    string    `json:"priceinfo"`
	Description  string    `json:"description"`
	PostedAt     time.Time `json:"postedat"`
	NsbmApproved bool      `json:"nsbmapproved"`
	Available    bool      `json:"available"`
	Rating       float64   `json:"rating"`
}

type Booking struct {
	gorm.Model
	ID        uint      `gorm:"primaryKey" json:"id"`
	StudentID uint      `json:"studentid"`
	OwnerID   uint      `json:"ownerid"`
	HostelID  uint      `json:"hostelid"`
	CheckIn   time.Time `json:"checkin"`
	CheckOut  time.Time `json:"checkout"`
}

type Admin struct {
	gorm.Model
	ID        uint   `gorm:"primaryKey" json:"id"`
	AdminName string `json:"adminname"`
	Priority  int    `json:"priority"`
	Role      string `json:"role"`
	UserID    uint   `json:"userid"`
}

type Warden struct {
	gorm.Model
	ID         uint   `gorm:"primaryKey" json:"id"`
	WardenName string `json:"wardenname"`
	PhoneNo    string `json:"phoneno"`
	NIC        string `json:"nic"`
	Image      string `json:"image"`
	UserID     uint   `json:"userid"`
}

type Article struct {
	gorm.Model
	ID          uint   `gorm:"primaryKey" json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

// internal struct

type StudentSingup struct {
	StdName  string `json:"stdname"`
	BOD      string `json:"bod"`
	Batch    string `json:"batch"`
	StdNo    string `json:"stdno"`
	Email    string `json:"email"`
	Password string `json:"password"`
	StdPno   string `json:"stdpno"`
}

type HostelOwnerSingup struct {
	OwnerName string `json:"ownername"`
	BOD       string `json:"bod"`
	Address   string `json:"address"`
	PhoneNo   string `json:"phoneno"`
	NIC       string `json:"nic"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Image     string `json:"image"`
}

type HostelReg struct {
	HostelName  string  `json:"hostelname"`
	Address     string  `json:"address"`
	Lat         float64 `json:"lat"`
	Lng         float64 `json:"lng"`
	PhoneNo     string  `json:"phoneno"`
	Image1      string  `json:"image1"`
	Image2      string  `json:"image2"`
	Image3      string  `json:"image3"`
	OwnerID     uint    `json:"ownerid"`
	Rooms       int     `json:"rooms"`
	BathRooms   int     `json:"bathrooms"`
	Price       float64 `json:"price"`
	PriceInfo   string  `json:"priceinfo"`
	Description string  `json:"description"`
}

type OwnerView struct {
	ID        uint     `json:"id"`
	OwnerName string   `json:"ownername"`
	Address   string   `json:"address"`
	PhoneNo   string   `json:"phoneno"`
	Email     string   `json:"email"`
	Approved  bool     `json:"approved"`
	Image     string   `json:"image"`
	Hostels   []Hostel `json:"hostels"`
}

type WardenSingup struct {
	WardenName string `json:"wardenname"`
	PhoneNo    string `json:"phoneno"`
	NIC        string `json:"nic"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	Image      string `json:"image"`
}

type AddArticle struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	PostedBy    string `json:"postedby"`
}
