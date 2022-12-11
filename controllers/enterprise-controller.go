package controllers

import (
	"car-park/models"
	"car-park/utils"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/csrf"
	"github.com/imdario/mergo"
)

type EnterpriseController interface {
	SaveEnterprise(ctx *gin.Context) error
	SaveDriver(ctx *gin.Context) error
	SaveManager(ctx *gin.Context) error
	FindAllEnterprises() []models.Enterprise
	FindAllDrivers() []models.Driver
	FindAllManagers() []models.Manager

	RedirectManager(ctx *gin.Context)
	ManagerShowAllVehicles(ctx *gin.Context)
	ManagerShowCreateVehicle(ctx *gin.Context)
	ManagerShowEditVehicle(ctx *gin.Context)

	ManagerFindAllVehicles(ctx *gin.Context, preload bool) []models.Vehicle
	ManagerFindAllDrivers(ctx *gin.Context) []models.Driver
	ManagerSaveVehicle(ctx *gin.Context) (err error)
	ManagerUpdateVehicle(ctx *gin.Context) error
	ManagerDeleteVehicle(ctx *gin.Context) error

	ManagerGetVehicleRoutesGeoJSON(ctx *gin.Context) (string, error)
	ManagerGetVehicleRoutesGeopoints(ctx *gin.Context) ([]models.GeoPoint, error)
	ManagerSaveVehicleGeoPoint(ctx *gin.Context) error

	SaveRide(ctx *gin.Context) error
	ManagerVehicleRides(ctx *gin.Context, inGeoJson bool) ([]models.Ride, error)
	ManagerVehicleRidesFold(ctx *gin.Context, inGeoJsons bool) ([]models.GeoPoint, string, error)

	ManagerHumanReadRides(ctx *gin.Context) ([]models.HumanReadRide, error)

	AuthManager(ctx *gin.Context) error
}

type enterpriseController struct {
	service models.EnterpriseService
}

var timeLayout string = "2006-01-02 15:04:05"

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

func (c *enterpriseController) SaveRide(ctx *gin.Context) error {
	var ride models.Ride
	err := ctx.ShouldBindJSON(&ride)
	if err != nil {
		return err
	}

	err = validate.Struct(ride)
	if err != nil {
		return err
	}

	if !ride.RideStart.Before(ride.RideFinish) {
		return fmt.Errorf("RideFinish is earlier that RideStart")
	}

	err = c.service.SaveRide(ride)
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

func (c *enterpriseController) RedirectManager(ctx *gin.Context) {
	auth := strings.SplitN(ctx.Request.Header.Get("Authorization"), " ", 2)
	if len(auth) != 2 || auth[0] != "Basic" {
		ctx.JSON(401, "Unauthorized")
		return
	}

	payload, _ := base64.StdEncoding.DecodeString(auth[1])
	credsPair := strings.SplitN(string(payload), ":", 2)

	manager := c.service.ManagerByCreds(credsPair[0], credsPair[1])
	redirectURL := "/view/manager/" + strconv.FormatUint(uint64(manager.ID), 10) + "/vehicles"

	cookiePage, err := ctx.Cookie("last_page")
	if err != nil {
		log.Println("Нету куки, вася!", err)
	} else {
		redirectURL = redirectURL + "/?page=" + cookiePage
	}
	ctx.Redirect(http.StatusFound, redirectURL)
}

func (c *enterpriseController) ManagerFindAllVehicles(ctx *gin.Context, preload bool) []models.Vehicle {
	managerId, _ := strconv.ParseUint(ctx.Param("id"), 0, 0)
	manager := c.service.ManagerByID(uint(managerId))
	accessibleEnterprises := manager.AccessibleEnterprises

	pagination := models.GenPaginationFromRequest(ctx)
	vehicles := c.service.ManagerFindAllVehicles(accessibleEnterprises, pagination, preload)

	return vehicles
}

func (c *enterpriseController) ManagerFindAllDrivers(ctx *gin.Context) []models.Driver {
	managerId, _ := strconv.ParseUint(ctx.Param("id"), 0, 0)
	manager := c.service.ManagerByID(uint(managerId))
	accessibleEnterprises := manager.AccessibleEnterprises
	return c.service.ManagerFindAllDrivers(accessibleEnterprises)
}

func (c *enterpriseController) ManagerSaveVehicle(ctx *gin.Context) (err error) {
	var vehicle models.Vehicle
	err = ctx.Bind(&vehicle)
	vehicleBrand := vehicle.CarModel.Brand
	if vehicleBrand == "Choose car Brand" || vehicleBrand == "" {
		vehicle.CarModel.Brand = "No Brand"
	}

	err = validate.Struct(vehicle)
	if err != nil {
		return err
	}

	vehicle.CommissioningDate = time.Now().UTC()
	fmt.Println("THIS IS THIS TIIIIMEEE", vehicle.CommissioningDate)

	if c.authManagerEntUpdates(ctx, int(vehicle.EnterpriseID)) == true {
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
	if c.authManagerEntUpdates(ctx, int(newVehicle.EnterpriseID)) != true {
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

	if c.authManagerVehicleUpdates(ctx, vehicle.ID) != true {
		err = fmt.Errorf("Wrong Vehicle ID")
		return err
	} else {
		c.service.DeleteVehicle(vehicle)
		return nil
	}

}

func (c *enterpriseController) ManagerShowAllVehicles(ctx *gin.Context) {
	if len(ctx.Request.URL.Query()["page"]) == 0 || ctx.Request.URL.Query()["page"][0] == "0" {
		ctx.Redirect(http.StatusMovedPermanently, ctx.Request.URL.String()+"?page=1")
		return
	}

	managerId, _ := strconv.ParseUint(ctx.Param("id"), 0, 0)
	manager := c.service.ManagerByID(uint(managerId))
	accessibleEnterprises := manager.AccessibleEnterprises

	pagination := models.GenPaginationFromRequest(ctx)

	currentPage := ctx.Request.URL.Query()["page"][0]
	currentPageInt, _ := strconv.Atoi(currentPage)
	nextPage := ""
	prevPage := ""
	fmt.Println(currentPage)

	fmt.Println(currentPageInt, reflect.TypeOf(currentPageInt))

	vehicles := c.service.ManagerFindAllVehicles(accessibleEnterprises, pagination, false)
	if len(vehicles) < models.DEFAULT_PAGINATION_LIMIT {
		nextPage = "nonext"
	} else {
		nextPage = strings.Replace(ctx.Request.URL.String(), "?page="+currentPage, "?page="+strconv.Itoa(currentPageInt+1), 1)
	}

	if currentPageInt == 1 {
		prevPage = "noprev"
	} else {
		prevPage = strings.Replace(ctx.Request.URL.String(), "?page="+currentPage, "?page="+strconv.Itoa(currentPageInt-1), 1)
	}
	data := gin.H{
		"title":          "Vihecles Stock",
		"vehicles":       vehicles,
		csrf.TemplateTag: csrf.TemplateField(ctx.Request),
		"nextpage":       nextPage,
		"prevpage":       prevPage,
		"managerID":      managerId,
	}

	fmt.Println(len(vehicles))

	fmt.Println(ctx.Request.URL.Query())

	ctx.SetCookie("last_page", currentPage, 15, "/", "localhost", false, true)
	ctx.HTML(http.StatusOK, "vehicles.html", data)
}

func (c *enterpriseController) ManagerShowCreateVehicle(ctx *gin.Context) {
	managerId, _ := strconv.ParseUint(ctx.Param("id"), 0, 0)
	manager := c.service.ManagerByID(uint(managerId))
	accessibleEnterprises := manager.AccessibleEnterprises
	ctx.HTML(http.StatusOK, "create-vehicle.html", gin.H{

		csrf.TemplateTag:        csrf.TemplateField(ctx.Request),
		"managerID":             managerId,
		"accessibleEnterprises": accessibleEnterprises,
	})
}

func (c *enterpriseController) ManagerShowEditVehicle(ctx *gin.Context) {
	var vehicle models.Vehicle
	ctx.ShouldBind(&vehicle)

	managerId, _ := strconv.ParseUint(ctx.Param("id"), 0, 0)
	manager := c.service.ManagerByID(uint(managerId))
	accessibleEnterprises := manager.AccessibleEnterprises

	vehicleId, _ := strconv.ParseUint(ctx.Param("vehicle_id"), 0, 0)

	vehicle = c.service.VehicleByID(uint(vehicleId))
	validate.Struct(vehicle)
	fmt.Println(vehicle)
	data := gin.H{
		"title":                 "Edit car in stock",
		"vehicle":               vehicle,
		csrf.TemplateTag:        csrf.TemplateField(ctx.Request),
		"managerID":             managerId,
		"accessibleEnterprises": accessibleEnterprises,
		"currentEID":            vehicle.EnterpriseID,
	}

	ctx.HTML(http.StatusOK, "edit-vehicle.html", data)
}

func (c *enterpriseController) ManagerGetVehicleRoutesGeoJSON(ctx *gin.Context) (string, error) {
	vehicleID, _ := strconv.ParseUint(ctx.Param("vehicle_id"), 0, 0)

	notBefore, notAfter, err := utils.UrlQTimeStampsToUTCStrings(ctx)
	if err != nil {
		return "", err
	}

	if !c.authManagerVehicleUpdates(ctx, uint(vehicleID)) {
		return "", fmt.Errorf("Wrong vehicle ID")
	} else {
		geoJSON, err := c.service.ManagerGetVehicleRoutesGeoJSON(uint(vehicleID), notBefore, notAfter)
		if err != nil {
			return "", err
		} else {
			return geoJSON, nil
		}
	}

}

func (c *enterpriseController) ManagerGetVehicleRoutesGeopoints(ctx *gin.Context) ([]models.GeoPoint, error) {
	var geoPoints []models.GeoPoint
	var err error
	vehicleID, _ := strconv.ParseUint(ctx.Param("vehicle_id"), 0, 0)
	notBefore, notAfter, err := utils.UrlQTimeStampsToUTCStrings(ctx)
	if err != nil {
		return geoPoints, err
	}

	if !c.authManagerVehicleUpdates(ctx, uint(vehicleID)) {
		return geoPoints, fmt.Errorf("Wrong vehicle ID")
	} else {
		geoPoints, err = c.service.ManagerGetVehicleRoutesGeopoints(uint(vehicleID), notBefore, notAfter)
		if err != nil {
			return geoPoints, err
		} else {
			return geoPoints, nil
		}
	}
}

func (c *enterpriseController) ManagerSaveVehicleGeoPoint(ctx *gin.Context) error {
	var geoPoint models.GeoPoint
	err := ctx.ShouldBind(&geoPoint)
	if err != nil {
		if strings.Contains(err.Error(), "parsing") {
			log.Println(err)
			return fmt.Errorf("Invalid track_time format. Please refer to RFC3339")

		} else {
			return err
		}
	}

	err = validate.Struct(geoPoint)
	if err != nil {
		return err
	}

	vehicleID, _ := strconv.ParseUint(ctx.Param("vehicle_id"), 0, 0)
	geoPoint.VehicleID = uint(vehicleID)
	err = c.service.ManagerSaveGeoPoint(geoPoint)
	return err
}

func (c *enterpriseController) ManagerVehicleRides(ctx *gin.Context, inGeoJsons bool) ([]models.Ride, error) {
	var rides []models.Ride
	vehicleID, _ := strconv.ParseUint(ctx.Param("vehicle_id"), 0, 0)
	notBefore, notAfter, err := utils.UrlQTimeStampsToUTCStrings(ctx)
	if err != nil {
		return rides, err
	}

	if !c.authManagerVehicleUpdates(ctx, uint(vehicleID)) {
		return rides, fmt.Errorf("Wrong vehicle ID")
	} else {
		rides, err = c.service.ManagerGetVehicleRides(uint(vehicleID), notBefore, notAfter, inGeoJsons)
		if err != nil {
			return rides, err
		} else {
			return rides, nil
		}
	}

}

func (c *enterpriseController) ManagerVehicleRidesFold(ctx *gin.Context, inGeoJsons bool) (GeoPoints []models.GeoPoint, GeoJson string, err error) {
	var geoPoints []models.GeoPoint
	var rides []models.Ride
	var geoJson string
	vehicleID, _ := strconv.ParseUint(ctx.Param("vehicle_id"), 0, 0)
	notBefore, notAfter, err := utils.UrlQTimeStampsToUTCStrings(ctx)
	if err != nil {
		return geoPoints, "", err
	}

	if !c.authManagerVehicleUpdates(ctx, uint(vehicleID)) {
		return geoPoints, "", fmt.Errorf("Wrong vehicle ID")
	} else {
		rides, err = c.service.ManagerGetVehicleRides(uint(vehicleID), notBefore, notAfter, inGeoJsons)
		if err != nil {
			return geoPoints, "", err
		}

		nb := rides[0].RideStart
		na := rides[len(rides)-1].RideFinish

		notBeforeFold, notAfterFold, err := utils.TimeStampsToUTCStrings(nb, na)
		if err != nil {
			return geoPoints, "", err
		}
		if inGeoJsons {
			geoJson, err = c.service.ManagerGetVehicleRoutesGeoJSON(uint(vehicleID), notBeforeFold, notAfterFold)
			if err != nil {

				return geoPoints, "", err
			} else {
				return geoPoints, geoJson, nil
			}
		} else {
			geoPoints, err = c.service.ManagerGetVehicleRoutesGeopoints(uint(vehicleID), notBeforeFold, notAfterFold)
			if err != nil {
				return geoPoints, "", err
			} else {
				return geoPoints, "", nil
			}
		}
	}

}

func (c *enterpriseController) ManagerHumanReadRides(ctx *gin.Context) ([]models.HumanReadRide, error) {
	var rides []models.Ride
	var humanReadableRides []models.HumanReadRide
	var err error

	rides, err = c.ManagerVehicleRides(ctx, false)
	if err != nil {
		return humanReadableRides, err
	}

	for _, ride := range rides {
		geoPointsLen := len(ride.GeoPoints)
		humanReadableRides = append(humanReadableRides, models.HumanReadRide{
			VehicleID:    ride.VehicleID,
			RideStart:    ride.RideStart,
			StartAddress: utils.GeocodeToAddress(ride.GeoPoints[0].GeoY, float64(ride.GeoPoints[0].GeoX)),

			RideFinish:    ride.RideFinish,
			FinishAddress: utils.GeocodeToAddress(ride.GeoPoints[geoPointsLen-1].GeoY, ride.GeoPoints[geoPointsLen-1].GeoX),
			RideDuration:  time.Time{}.Add(ride.RideFinish.Sub(ride.RideStart)).Format("15:04:05"),
		})
	}
	return humanReadableRides, err
}

func (c *enterpriseController) AuthManager(ctx *gin.Context) error {
	managerIdFromURL, _ := strconv.ParseUint(ctx.Param("id"), 0, 0)

	// TODO: this is total disaster, please make normal authorization
	if managerIdFromURL == 666 {
		return nil
	}
	auth := strings.SplitN(ctx.Request.Header.Get("Authorization"), " ", 2)
	if len(auth) != 2 || auth[0] != "Basic" {
		ctx.JSON(401, "Unauthorized")
		return fmt.Errorf("Unauthorized")
	}

	payload, _ := base64.StdEncoding.DecodeString(auth[1])
	credsPair := strings.SplitN(string(payload), ":", 2)

	manager := c.service.ManagerByCreds(credsPair[0], credsPair[1])
	fmt.Println(manager)

	if len(credsPair) != 2 {
		manager = models.Manager{}
		fmt.Println("TROUBLE IN PAIR")
		return fmt.Errorf("Unauthorized")
	}

	if manager.ID != uint(managerIdFromURL) {
		manager = models.Manager{}
		fmt.Println("TROUBLE IN ID")
		return fmt.Errorf("Unauthorized")
	}
	return nil
}

func (c *enterpriseController) authManagerEntUpdates(ctx *gin.Context, enterpriseID int) bool {
	access := false

	err := c.AuthManager(ctx)
	if err != nil {
		return access
	}

	managerIdFromURL, _ := strconv.ParseUint(ctx.Param("id"), 0, 0)
	// TODO: I hate myself so much
	if managerIdFromURL == 666 {
		return true
	}
	manager := c.service.ManagerByID(uint(managerIdFromURL))
	array := make([]int64, len(manager.AccessibleEnterprises))
	copy(array, manager.AccessibleEnterprises)
	access = contains(array, int64(enterpriseID))
	return access

}

func (c *enterpriseController) authManagerVehicleUpdates(ctx *gin.Context, vehicleID uint) bool {
	access := false

	err := c.AuthManager(ctx)
	if err != nil {
		return access
	}

	managerIdFromURL, _ := strconv.ParseUint(ctx.Param("id"), 0, 0)
	manager := c.service.ManagerByID(uint(managerIdFromURL))
	array := make([]int64, len(manager.AccessibleEnterprises))
	copy(array, manager.AccessibleEnterprises)

	vehicle := c.service.VehicleByID(uint(vehicleID))

	access = contains(array, int64(vehicle.EnterpriseID))
	return access
}
