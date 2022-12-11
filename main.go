package main

import (
	"car-park/controllers"
	"car-park/middlewares"
	"car-park/models"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

var (
	vehicleDB         models.VehicleDB         = models.NewVehicleDB()
	vehicleService    models.VehicleService    = models.NewVehicleService(vehicleDB)
	enterpriseService models.EnterpriseService = models.NewEnterpriseSerivce(vehicleDB)

	vehicleController    controllers.VehicleController    = controllers.NewVehicleController(vehicleService)
	enterpriseController controllers.EnterpriseController = controllers.NewEnterpriseController(enterpriseService)
)

func setupLogFile() {
	f, err := os.Create("gin.log")
	if err != nil {
		fmt.Println("[WARNING] Cannot create logfile, will use only stdout.")
	}
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

}

func main() {
	setupLogFile()
	defer vehicleDB.CloseConnection()

	managers := enterpriseController.FindAllManagers()
	accounts := gin.Accounts{}
	for _, m := range managers {
		accounts[m.Login] = m.Password
	}
	accounts["superadmin"] = "superadminqwerty"
	fmt.Println(accounts)

	server := gin.New()
	server.Use(gin.Recovery(),
		middlewares.Logger(),
		middlewares.BasicAuth(accounts))

	// doesnot working  ¯\_(ツ)_/¯
	// vehicleDB.DestructiveReset()
	vehicleDB.AutoMigrate()
	// vehicleDB.LoadFixturesGeotracks(25171)

	server.Static("/css", "./views/css")
	server.StaticFile("/logo.svg", "./views/assets/logo-1-white.svg")
	server.LoadHTMLGlob("./views/templates/*.html")

	// server.GET("/", func(ctx *gin.Context) {
	// 	ctx.HTML(http.StatusOK, "index.html", gin.H{})
	// })

	server.GET("/", enterpriseController.RedirectManager)

	server.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code": "404", "message": "PAGE NOT FOUND"})
	})

	// viewRoutes := server.Group("/view").Use(middlewares.CSRF())
	// {
	// 	viewRoutes.GET("/vehicles", vehicleController.ShowAllVehicles)
	// 	viewRoutes.GET("/vehicles/create", vehicleController.ShowCreateVehicle)
	// 	viewRoutes.GET("/vehicles/edit/:id", vehicleController.ShowEditVehicle)
	// 	viewRoutes.GET("/carmodels", vehicleController.ShowAllCarModels)
	// }

	viewRoutes := server.Group("/view").Use(middlewares.CSRF())
	{
		viewRoutes.GET("/manager/:id/vehicles/", enterpriseController.ManagerShowAllVehicles)
		viewRoutes.GET("/manager/:id/vehicle/create", enterpriseController.ManagerShowCreateVehicle)
		viewRoutes.GET("/manager/:id/vehicles/edit/:vehicle_id", enterpriseController.ManagerShowEditVehicle)
		viewRoutes.GET("/carmodels", vehicleController.ShowAllCarModels)
	}
	// TODO: Move "view" CRUD routes from /api to view group for applying CSRF for all of them.
	// TODO: It will be better to implement JWT in place of basic auth, at least for /api
	// TODO: Damn boah, this mess needs to be really reorganized...
	apiRoutes := server.Group("/api")
	{

		// ID here is Manager id
		apiRoutes.GET("/manager/:id/vehicles/", func(ctx *gin.Context) {
			err := enterpriseController.AuthManager(ctx)
			if err != nil {
				ctx.JSON(401, "Unauthorized")
			} else {
				ctx.JSON(http.StatusOK, enterpriseController.ManagerFindAllVehicles(ctx, false))
			}
		})

		apiRoutes.GET("/manager/:id/drivers/", func(ctx *gin.Context) {
			err := enterpriseController.AuthManager(ctx)
			if err != nil {
				ctx.JSON(401, "Unauthorized")
			} else {
				ctx.JSON(http.StatusOK, enterpriseController.ManagerFindAllDrivers(ctx))
			}
		})

		apiRoutes.POST("/manager/:id/vehicle/create", func(ctx *gin.Context) {
			err := enterpriseController.AuthManager(ctx)
			if err != nil {
				ctx.JSON(401, "Unauthorized")
			} else {
				err := enterpriseController.ManagerSaveVehicle(ctx)
				if err != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{
						"ERROR": err.Error()})

				} else if ctx.ContentType() == "application/json" {
					ctx.JSON(http.StatusOK, gin.H{"message": "Vehicle updated OK"})
				} else {
					ctx.Redirect(http.StatusFound, "/")
				}
			}

		})

		apiRoutes.POST("/manager/:id/vehicle/:vehicle_id/update", func(ctx *gin.Context) {
			err := enterpriseController.AuthManager(ctx)
			if err != nil {
				ctx.JSON(401, "Unauthorized")
			} else {
				err := enterpriseController.ManagerUpdateVehicle(ctx)
				if err != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{
						"ERROR": err.Error()})

				} else if ctx.ContentType() == "application/json" {
					ctx.JSON(http.StatusOK, gin.H{"message": "Vehicle updated OK"})
				} else {
					ctx.Redirect(http.StatusFound, "/")
				}

			}

		})

		apiRoutes.POST("/manager/:id/vehicle/:vehicle_id/delete", func(ctx *gin.Context) {
			err := enterpriseController.AuthManager(ctx)
			if err != nil {
				ctx.JSON(401, "Unauthorized")
			} else {
				err := enterpriseController.ManagerDeleteVehicle(ctx)
				if err != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{
						"ERROR": err.Error()})

				} else if ctx.ContentType() == "application/json" {
					ctx.JSON(http.StatusOK, gin.H{"message": "Vehicle deleted OK"})
				} else {
					fmt.Println("DEBUG", ctx.Request.URL)
					ctx.Redirect(http.StatusFound, "/")
				}
			}

		})

		//GEOPOINTS_METHODS
		apiRoutes.POST("/manager/:id/vehicle/:vehicle_id/checkpoint", func(ctx *gin.Context) {
			err := enterpriseController.AuthManager(ctx)
			if err != nil {
				ctx.JSON(401, "Unauthorized")
			} else {
				err := enterpriseController.ManagerSaveVehicleGeoPoint(ctx)
				if err != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{
						"ERROR": err.Error()})
				} else {
					ctx.JSON(http.StatusOK, "GeoPoint added OK.")
				}
			}

		})

		apiRoutes.GET("/manager/:id/vehicle/:vehicle_id/routes_geojson", func(ctx *gin.Context) {
			err := enterpriseController.AuthManager(ctx)
			if err != nil {
				ctx.JSON(401, "Unauthorized")
			} else {
				geoJSON, err := enterpriseController.ManagerGetVehicleRoutesGeoJSON(ctx)
				if err != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{
						"ERROR": err.Error()})
				} else {
					ctx.JSON(http.StatusOK, geoJSON)
				}
			}

		})

		apiRoutes.GET("/manager/:id/vehicle/:vehicle_id/routes_geopoints", func(ctx *gin.Context) {
			err := enterpriseController.AuthManager(ctx)
			if err != nil {
				ctx.JSON(401, "Unauthorized")
			} else {
				geoPoints, err := enterpriseController.ManagerGetVehicleRoutesGeopoints(ctx)
				if err != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{
						"ERROR": err.Error()})
				} else {
					ctx.JSON(http.StatusOK, geoPoints)
				}
			}

		})

		apiRoutes.GET("/manager/:id/vehicle/:vehicle_id/rides_geojson", func(ctx *gin.Context) {
			err := enterpriseController.AuthManager(ctx)
			if err != nil {
				ctx.JSON(401, "Unauthorized")
			} else {
				rides, err := enterpriseController.ManagerVehicleRides(ctx, true)
				if err != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{
						"ERROR": err.Error()})
				} else {
					ctx.JSON(http.StatusOK, rides)
				}
			}

		})

		apiRoutes.GET("/manager/:id/vehicle/:vehicle_id/rides_geopoints", func(ctx *gin.Context) {
			err := enterpriseController.AuthManager(ctx)
			if err != nil {
				ctx.JSON(401, "Unauthorized")
			} else {
				rides, err := enterpriseController.ManagerVehicleRides(ctx, false)
				if err != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{
						"ERROR": err.Error()})
				} else {
					ctx.JSON(http.StatusOK, rides)
				}
			}

		})

		apiRoutes.GET("/manager/:id/vehicle/:vehicle_id/rides_fold_tracks_geojson", func(ctx *gin.Context) {
			err := enterpriseController.AuthManager(ctx)
			if err != nil {
				ctx.JSON(401, "Unauthorized")
			} else {
				_, geoJson, err := enterpriseController.ManagerVehicleRidesFold(ctx, true)
				if err != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{
						"ERROR": err.Error()})
				} else {
					ctx.JSON(http.StatusOK, geoJson)
				}
			}

		})

		apiRoutes.GET("/manager/:id/vehicle/:vehicle_id/rides_fold_tracks_geopoints", func(ctx *gin.Context) {
			err := enterpriseController.AuthManager(ctx)
			if err != nil {
				ctx.JSON(401, "Unauthorized")
			} else {
				geoPoints, _, err := enterpriseController.ManagerVehicleRidesFold(ctx, false)
				if err != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{
						"ERROR": err.Error()})
				} else {
					ctx.JSON(http.StatusOK, geoPoints)
				}
			}

		})

		apiRoutes.GET("/manager/:id/vehicle/:vehicle_id/rides_human_readable", func(ctx *gin.Context) {
			err := enterpriseController.AuthManager(ctx)
			if err != nil {
				ctx.JSON(401, "Unauthorized")
			} else {
				HumanReadableRides, err := enterpriseController.ManagerHumanReadRides(ctx)
				if err != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{
						"ERROR": err.Error()})

				} else {
					ctx.JSON(http.StatusOK, HumanReadableRides)
				}
			}

		})

		// POST /api/vehicles creates new Vehicle object with nested
		// CarModle in it.

		// TODO: НАУЧИТЬСЯ ПРОГРАММИРОВАТЬСЯ
		// apiRoutes.POST("/vehicles", func(ctx *gin.Context) {
		// 	err := vehicleController.SaveVehicle(ctx)
		// 	if err != nil {
		// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		// 	} else if ctx.ContentType() == "application/json" {
		// 		ctx.JSON(http.StatusOK, gin.H{"message": "New Vehicle added OK"})
		// 	} else {
		// 		ctx.Redirect(http.StatusFound, "/view/vehicles")
		// 	}

		// })

		// // POST /api/update/vehicle/:id updates existing Vehicle with nested
		// // CarModel by vehicle ID. Can update only CarModel.
		// apiRoutes.POST("/update/vehicle/:id", middlewares.CSRF(), func(ctx *gin.Context) {
		// 	err := vehicleController.UpdateVehicle(ctx)
		// 	if err != nil {
		// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		// 	} else if ctx.ContentType() == "application/json" {
		// 		ctx.JSON(http.StatusOK, gin.H{"message": "Vehicle updated OK"})
		// 	} else {
		// 		ctx.Redirect(http.StatusFound, "/view/vehicles")
		// 	}

		// })

		// // POST /api/delete/vehicle/:id soft deletes Vehicle with nested CarModel by vehicle ID.
		// apiRoutes.POST("/delete/vehicle/:id", func(ctx *gin.Context) {
		// 	err := vehicleController.DeleteVehicle(ctx)
		// 	if err != nil {
		// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		// 	} else if ctx.ContentType() == "application/json" || ctx.ContentType() == "" {
		// 		ctx.JSON(http.StatusOK, gin.H{"message": "Vehicle deleted OK"})
		// 	} else {
		// 		ctx.Redirect(http.StatusFound, "/view/vehicles")
		// 	}

		// })

		apiRoutes.POST("/save/ride", func(ctx *gin.Context) {
			err := enterpriseController.SaveRide(ctx)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error:": err.Error()})
			} else {
				ctx.JSON(http.StatusOK, gin.H{"message:": "New ride added OK"})
			}

		})

		apiRoutes.POST("/save/enterprises", func(ctx *gin.Context) {
			err := enterpriseController.SaveEnterprise(ctx)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				ctx.JSON(http.StatusOK, gin.H{"message": "New enterprise added OK"})
			}

		})

		apiRoutes.POST("/save/drivers", func(ctx *gin.Context) {
			err := enterpriseController.SaveDriver(ctx)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				ctx.JSON(http.StatusOK, gin.H{"message": "New driver added OK"})
			}

		})

		apiRoutes.POST("/save/managers", func(ctx *gin.Context) {
			err := enterpriseController.SaveManager(ctx)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				ctx.JSON(http.StatusOK, gin.H{"message": "New manager added OK"})
			}
		})

		apiRoutes.GET("/enterprises", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, enterpriseController.FindAllEnterprises())
		})

		apiRoutes.GET("/drivers", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, enterpriseController.FindAllDrivers())
		})
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8888"
	}
	// vehicleDB.DestructiveReset()
	server.Run(":" + port)

}
