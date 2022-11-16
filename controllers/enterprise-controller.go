package controllers

import (
	"car-park/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type EnterpriseController interface {
	SaveEnterprise(ctx *gin.Context) error
	SaveDriver(ctx *gin.Context) error
	FindAllEnterprises() []models.Enterprise
	FindAllDrivers() []models.Driver
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

func (c *enterpriseController) FindAllEnterprises() []models.Enterprise {
	return c.service.FindAllEnterprises()
}

func (c *enterpriseController) FindAllDrivers() []models.Driver {
	return c.service.FindAllDrivers()
}
