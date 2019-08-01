package router

import (
	"CarServiceCenter/src/controller"
	"CarServiceCenter/src/db"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rs/cors"
)

// Initialize chi mux router
func Initialize() *chi.Mux {
	mongoStruct := &db.MongoStruct{}
	appointmentsController := controller.AppointmentsController{DB: mongoStruct}
	muxRouter := chi.NewRouter()

	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET"},
		AllowCredentials: true,
	})

	muxRouter.Use(cors.Handler)
	muxRouter.Use(middleware.RequestID)
	muxRouter.Use(middleware.RealIP)
	muxRouter.Use(middleware.Logger)
	muxRouter.Use(middleware.Recoverer)
	muxRouter.Use(middleware.Timeout(200 * time.Second))

	muxRouter.Get("/appointment/{id}", appointmentsController.GetAppointment)
	muxRouter.Post("/appointment/", appointmentsController.CreateAppointment)
	muxRouter.Patch("/appointment/{id}", appointmentsController.UpdateAppointmentStatus)
	muxRouter.Delete("/appointment/{id}", appointmentsController.DeleteAppointment)
	muxRouter.Get("/appointments/range/", appointmentsController.GetAppointmentsWithinDateRange)

	return muxRouter
}
