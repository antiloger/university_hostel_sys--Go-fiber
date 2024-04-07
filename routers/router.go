package routers

import (
	"github.com/antiloger/nhostel-go/handler"
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(App *fiber.App) {
	student_route := App.Group("/student")
	owner_route := App.Group("/owner")
	admin_route := App.Group("/admin")
	warden_route := App.Group("/warden")

	App.Get("/", handler.Hello)
	App.Post("/login", handler.Login)
	home := App.Group("/home")

	home.Get("/", handler.HomeLoad)
	home.Get("/search/:search", handler.SearchHostel)
	hostel_route := App.Group("/hostel")

	owner_route.Post("/signup", handler.Hostelownersignup)
	student_route.Post("/signup", handler.Studentsignup)
	hostel_route.Post("/signup", handler.Hostelcreate)
	warden_route.Post("/signup", handler.Wardensignup)

	home.Get("/getmyprofile", handler.GetMyProfile)
	warden_route.Get("/getwardenprofile/:ID", handler.Getwardendetails)
	warden_route.Get("/hostelapprovaltable", handler.HostelApproveTable)
	warden_route.Put("/hostelapproval/:ID", handler.HostelApprove)
	// jwt middleware

	// App.Use(jwtware.New(jwtware.Config{
	// 	SigningKey: jwtware.SigningKey{Key: []byte(config.Jwt_Secret)},
	// }))

	// home routes

	// owner routes
	owner_route.Get("/:ID", handler.HostelOwnerView)

	warden_route.Get("/wardentable", handler.WardenTable)

	// hostel routes
	hostel_route.Get("/:ID", handler.Hosteldetails)
	hostel_route.Put("/:ID", handler.Hostelupdate)
	hostel_route.Delete("/:ID", handler.Hosteldelete)
	hostel_route.Put("/available/:ID", handler.HostelAvailableUpdate)

	// wardens routes

	// admin routes
	admin_route.Get("/studentapprovaltable", handler.StudentTable)
	admin_route.Put("/studentapproval/:ID", handler.StudentApprove)

	admin_route.Get("/ownertable", handler.Hostelownertable)
	admin_route.Put("/ownerapproval/:ID", handler.OwnerApprove)

	admin_route.Post("/addarticle", handler.CreateArticle)
	admin_route.Get("/article", handler.GetArticles)
	admin_route.Get("/article/:ID", handler.GetArticle)
	admin_route.Delete("/article/:ID", handler.DeleteArticle)

	App.Post("/users", handler.Insertuser)
}
