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
		middlewares.BasicAuth(accounts),
		middlewares.CSRF())

	server.Static("/css", "./views/templates/css")
	server.StaticFile("/logo.svg", "./views/assets/logo-1-white.svg")
	server.LoadHTMLGlob("./views/templates/*.html")

	server.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", gin.H{})
	})

	server.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code": "404", "message": "PAGE NOT FOUND"})
	})

	viewRoutes := server.Group("/view")
	{
		viewRoutes.GET("/vehicles", vehicleController.ShowAllVehicles)
		viewRoutes.GET("/vehicles/create", vehicleController.ShowCreateVehicle)
		viewRoutes.GET("/vehicles/edit/:id", vehicleController.ShowEditVehicle)
		viewRoutes.GET("/carmodels", vehicleController.ShowAllCarModels)
	}

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

		// POST /api/vehicles creates new Vehicle object with nested
		// CarModle in it.
		apiRoutes.POST("/vehicles", func(ctx *gin.Context) {
			err := vehicleController.SaveVehicle(ctx)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else if ctx.ContentType() == "application/json" {
				ctx.JSON(http.StatusOK, gin.H{"message": "New Vehicle added OK"})
			} else {
				ctx.Redirect(http.StatusFound, "/view/vehicles")
			}

		})

		// POST /api/update/vehicle/:id updates existing Vehicle with nested
		// CarModel by vehicle ID. Can update only CarModel.
		apiRoutes.POST("/update/vehicle/:id", func(ctx *gin.Context) {
			err := vehicleController.UpdateVehicle(ctx)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else if ctx.ContentType() == "application/json" {
				ctx.JSON(http.StatusOK, gin.H{"message": "Vehicle updated OK"})
			} else {
				ctx.Redirect(http.StatusFound, "/view/vehicles")
			}

		})

		// POST /api/delete/vehicle/:id soft deletes Vehicle with nested CarModel by vehicle ID.
		apiRoutes.POST("/delete/vehicle/:id", func(ctx *gin.Context) {
			err := vehicleController.DeleteVehicle(ctx)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else if ctx.ContentType() == "application/json" || ctx.ContentType() == "" {
				ctx.JSON(http.StatusOK, gin.H{"message": "Vehicle deleted OK"})
			} else {
				ctx.Redirect(http.StatusFound, "/view/vehicles")
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

	fmt.Println(accounts)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8888"
	}
	// vehicleDB.DestructiveReset()
	server.Run(":" + port)
}
