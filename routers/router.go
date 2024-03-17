package routers

import (
	"github.com/antiloger/nhostel-go/config"
	"github.com/antiloger/nhostel-go/handler"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(App *fiber.App) {
	student_route := App.Group("/student")
	owner_route := App.Group("/owner")
	admin_route := App.Group("/admin")

	App.Get("/", handler.Hello)
	App.Post("/login", handler.Login)

	owner_route.Post("/signup", handler.Hostelownersignup)
	student_route.Post("/signup", handler.Studentsignup)

	// jwt middleware

	App.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(config.Jwt_Secret)},
	}))

	home := App.Group("/home")
	hostel_route := App.Group("/hostel")

	// home routes
	home.Get("/", handler.HomeLoad)

	// owner routes
	owner_route.Get("/:ID", handler.HostelOwnerView)

	// hostel routes
	hostel_route.Post("/signup", handler.Hostelcreate)
	hostel_route.Get("/:ID", handler.Hosteldetails)
	hostel_route.Put("/:ID", handler.Hostelupdate)
	hostel_route.Delete("/:ID", handler.Hosteldelete)
	hostel_route.Put("/available/:ID", handler.HostelAvailableUpdate)

	// admin routes
	admin_route.Get("/hostelapprovaltable", handler.HostelApproveTable)
	admin_route.Put("/hostelapproval/:ID", handler.HostelApprove)

	admin_route.Get("/studentapprovaltable", handler.StudentApproveTable)
	admin_route.Put("/studentapproval/:ID", handler.StudentApprove)

	admin_route.Get("/ownerapprovaltable", handler.HostelOwnerApproveTable)
	admin_route.Put("/ownerapproval/:ID", handler.OwnerApprove)

	App.Post("/users", handler.Insertuser)
}
