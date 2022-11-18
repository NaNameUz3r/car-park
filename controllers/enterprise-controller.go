package controllers

import (
	"car-park/models"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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
	// authheader := ctx.Request.Header["Authorization"]
	// encoded := strings.Split(authheader[0], " ")[1]
	// fmt.Println(base64.StdEncoding.DecodeString(encoded))
	accessibleEnterprises := manager.AccessibleEnterprises
	return c.service.ManagerFindAllVehicles(accessibleEnterprises, preload)
}

func (c *enterpriseController) ManagerFindAllDrivers(ctx *gin.Context) []models.Driver {
	managerId, _ := strconv.ParseUint(ctx.Param("id"), 0, 0)
	manager := c.service.ManagerByID(uint(managerId))
	// authheader := ctx.Request.Header["Authorization"]
	// encoded := strings.Split(authheader[0], " ")[1]
	// fmt.Println(base64.StdEncoding.DecodeString(encoded))
	accessibleEnterprises := manager.AccessibleEnterprises
	return c.service.ManagerFindAllDrivers(accessibleEnterprises)
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

	fmt.Println("THIS IS FUCKING MANAGER ID", manager.ID)
	fmt.Println("THIS IS FUCKING IN FROM URL", managerIdFromURL)
	if len(credsPair) != 2 || manager.ID != uint(managerIdFromURL) {
		manager = models.Manager{}
		return fmt.Errorf("Unauthorized")
	}
	return nil
}

// func (c *enterpriseController) authenticateManager(username, password string) bool {

// 	manager := c.service.ManagerByCreds(username, password)
// 	fmt.Println(manager)
// 	return true
// }
