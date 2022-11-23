package controllers

import (
	"car-park/models"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/imdario/mergo"
)

type EnterpriseController interface {
	SaveEnterprise(ctx *gin.Context) error
	SaveDriver(ctx *gin.Context) error
	SaveManager(ctx *gin.Context) error
	FindAllEnterprises() []models.Enterprise
	FindAllDrivers() []models.Driver
	FindAllManagers() []models.Manager

	ManagerFindAllVehicles(ctx *gin.Context, preload bool) []models.Vehicle
	ManagerFindAllDrivers(ctx *gin.Context) []models.Driver
	ManagerSaveVehicle(ctx *gin.Context) error
	ManagerUpdateVehicle(ctx *gin.Context) error
	ManagerDeleteVehicle(ctx *gin.Context) error

	AuthManager(ctx *gin.Context) error
}

type enterpriseController struct {
	service models.EnterpriseService
}

func NewEnterpriseController(svc models.EnterpriseService) EnterpriseController {
	validate = validator.New()
	return &enterpriseController{
		service: svc,
	}
}

func (c *enterpriseController) SaveEnterprise(ctx *gin.Context) error {
	var enterprise models.Enterprise
	err := ctx.ShouldBindJSON(&enterprise)
	if err != nil {
		return err
	}

	err = validate.Struct(enterprise)
	if err != nil {
		return err
	}
	err = c.service.SaveEnterprise(enterprise)
	return err
}

func (c *enterpriseController) SaveDriver(ctx *gin.Context) error {
	var driver models.Driver
	err := ctx.ShouldBindJSON(&driver)
	if err != nil {
		return err
	}

	err = validate.Struct(driver)
	if err != nil {
		return err
	}
	err = c.service.SaveDriver(driver)
	return err
}

func (c *enterpriseController) SaveManager(ctx *gin.Context) error {
	var manager models.Manager
	err := ctx.ShouldBindJSON(&manager)
	if err != nil {
		return err
	}

	err = validate.Struct(manager)
	if err != nil {
		return err
	}
	err = c.service.SaveManager(manager)
	return err
}

func (c *enterpriseController) FindAllEnterprises() []models.Enterprise {
	return c.service.FindAllEnterprises()
}

func (c *enterpriseController) FindAllDrivers() []models.Driver {
	return c.service.FindAllDrivers()
}
func (c *enterpriseController) FindAllManagers() []models.Manager {
	return c.service.FindAllManagers()
}

func (c *enterpriseController) ManagerFindAllVehicles(ctx *gin.Context, preload bool) []models.Vehicle {
	managerId, _ := strconv.ParseUint(ctx.Param("id"), 0, 0)
	manager := c.service.ManagerByID(uint(managerId))
	accessibleEnterprises := manager.AccessibleEnterprises
	return c.service.ManagerFindAllVehicles(accessibleEnterprises, preload)
}

func (c *enterpriseController) ManagerFindAllDrivers(ctx *gin.Context) []models.Driver {
	managerId, _ := strconv.ParseUint(ctx.Param("id"), 0, 0)
	manager := c.service.ManagerByID(uint(managerId))
	accessibleEnterprises := manager.AccessibleEnterprises
	return c.service.ManagerFindAllDrivers(accessibleEnterprises)
}

func (c *enterpriseController) ManagerSaveVehicle(ctx *gin.Context) error {
	var vehicle models.Vehicle
	err := ctx.Bind(&vehicle)
	vehicleBrand := vehicle.CarModel.Brand
	if vehicleBrand == "Choose car Brand" || vehicleBrand == "" {
		vehicle.CarModel.Brand = "No Brand"
	}

	err = validate.Struct(vehicle)
	if err != nil {
		return err
	}

	if c.authManagerUpdates(ctx, int(vehicle.EnterpriseID)) == true {
		err = c.service.SaveVehicle(vehicle)
	} else {
		err = fmt.Errorf("Wrong Enterprise ID")
	}
	return err

}

func (c *enterpriseController) ManagerUpdateVehicle(ctx *gin.Context) error {
	var newVehicle models.Vehicle
	err := ctx.Bind(&newVehicle)
	err = validate.Struct(newVehicle)
	if c.authManagerUpdates(ctx, int(newVehicle.EnterpriseID)) != true {
		err = fmt.Errorf("Wrong Enterprise ID")
		return err
	} else {

		var vehicleOrig models.Vehicle
		vehcleID, _ := strconv.ParseUint(ctx.Param("vehicle_id"), 0, 0)
		vehicleOrig = c.service.ManagerVehicleByID(uint(vehcleID))
		err = validate.Struct(vehicleOrig)
		if err != nil {
			return err
		}

		mergo.Merge(&newVehicle, vehicleOrig)
		err = validate.Struct(newVehicle)
		if err != nil {
			return err
		}

		return c.service.UpdateVehicle(newVehicle)

	}
}

func (c *enterpriseController) ManagerDeleteVehicle(ctx *gin.Context) error {
	var vehicle models.Vehicle
	id, err := strconv.ParseUint(ctx.Param("vehicle_id"), 0, 0)
	if err != nil {
		return err
	}

	vehicle.ID = uint(id)

	if c.authManagerUpdates(ctx, int(vehicle.ID)) != true {
		err = fmt.Errorf("Wrong Vehicle ID")
		return err
	} else {
		c.service.DeleteVehicle(vehicle)
		return nil
	}

}

func (c *enterpriseController) AuthManager(ctx *gin.Context) error {
	auth := strings.SplitN(ctx.Request.Header.Get("Authorization"), " ", 2)
	if len(auth) != 2 || auth[0] != "Basic" {
		ctx.JSON(401, "Unauthorized")
		return fmt.Errorf("Unauthorized")
	}

	payload, _ := base64.StdEncoding.DecodeString(auth[1])
	credsPair := strings.SplitN(string(payload), ":", 2)

	manager := c.service.ManagerByCreds(credsPair[0], credsPair[1])
	managerIdFromURL, _ := strconv.ParseUint(ctx.Param("id"), 0, 0)

	if len(credsPair) != 2 || manager.ID != uint(managerIdFromURL) {
		manager = models.Manager{}
		return fmt.Errorf("Unauthorized")
	}
	return nil
}

func (c *enterpriseController) authManagerUpdates(ctx *gin.Context, enterpriseID int) bool {
	access := false

	err := c.AuthManager(ctx)
	if err != nil {
		return access
	}

	managerIdFromURL, _ := strconv.ParseUint(ctx.Param("id"), 0, 0)
	manager := c.service.ManagerByID(uint(managerIdFromURL))
	array := make([]int64, len(manager.AccessibleEnterprises))
	copy(array, manager.AccessibleEnterprises)
	access = contains(array, int64(enterpriseID))
	return access

}
