package test

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/antiloger/nhostel-go/models"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestHostelcreate(t *testing.T) {
	app := fiber.New()
	// Initialize your app and set up the necessary dependencies

	// Test case 1: Successful creation of a hostel
	t.Run("Create Hostel - Success", func(t *testing.T) {
		// Prepare the request body and form files
		hostelReg := models.HostelReg{
			HostelName:  "Test Hostel",
			Address:     "Test Address",
			Lat:         1.2345,
			Lng:         2.3456,
			PhoneNo:     "1234567890",
			Image1:      "image1.jpg",
			Image2:      "image2.jpg",
			Image3:      "image3.jpg",
			OwnerID:     1,
			Rooms:       10,
			BathRooms:   5,
			Price:       1000,
			PriceInfo:   "Test Price Info",
			Description: "Test Description",
		}
		body, _ := json.Marshal(hostelReg)
		req := httptest.NewRequest(http.MethodPost, "/hostelcreate", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		form := multipart.NewWriter(req.Body)
		img1, _ := form.CreateFormFile("image1", "image1.jpg")
		img1.Write([]byte("Test Image 1"))
		img2, _ := form.CreateFormFile("image2", "image2.jpg")
		img2.Write([]byte("Test Image 2"))
		img3, _ := form.CreateFormFile("image3", "image3.jpg")
		img3.Write([]byte("Test Image 3"))
		form.Close()

		req.Header.Set("Content-Type", form.FormDataContentType())

		// Make the request to the endpoint
		resp, _ := app.Test(req)

		// Assert the response status code and body
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Assert the response JSON
		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)
		assert.Equal(t, "success", result["status"])
		assert.Equal(t, "hostel has created", result["message"])
	})

	// Test case 2: Invalid input
	t.Run("Create Hostel - Invalid Input", func(t *testing.T) {
		// Prepare the request body with invalid input
		hostelReg := models.HostelReg{
			HostelName:  "",
			Address:     "Test Address",
			Lat:         1.2345,
			Lng:         2.3456,
			PhoneNo:     "1234567890",
			Image1:      "image1.jpg",
			Image2:      "image2.jpg",
			Image3:      "image3.jpg",
			OwnerID:     1,
			Rooms:       10,
			BathRooms:   5,
			Price:       1000,
			PriceInfo:   "Test Price Info",
			Description: "Test Description",
		}
		body, _ := json.Marshal(hostelReg)
		req := httptest.NewRequest(http.MethodPost, "/hostelcreate", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		form := multipart.NewWriter(req)
		img1, _ := form.CreateFormFile("image1", "image1.jpg")
		img1.Write([]byte("Test Image 1"))
		img2, _ := form.CreateFormFile("image2", "image2.jpg")
		img2.Write([]byte("Test Image 2"))
		img3, _ := form.CreateFormFile("image3", "image3.jpg")
		img3.Write([]byte("Test Image 3"))
		form.Close()

		req.Header.Set("Content-Type", form.FormDataContentType())

		// Make the request to the endpoint
		resp, _ := app.Test(req)

		// Assert the response status code and body
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		// Assert the response JSON
		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)
		assert.Equal(t, "error", result["status"])
		assert.Equal(t, "Somthing's wrong with your input", result["message"])
	})

	// Add more test cases as needed
}
