package utils

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func UrlQTimeStampsToUTCStrings(ctx *gin.Context) (string, string, error) {

	notBefore := ""
	notAfter := ""

	urlQuery := ctx.Request.URL.Query()
	for key, value := range urlQuery {
		queryValue := value[len(value)-1]
		switch key {
		case "notBefore":
			notBefore = queryValue
			break
		case "notAfter":
			notAfter = queryValue
			break
		}
	}
	notBeforeRFC, err := time.Parse("2006-01-02", notBefore)
	notAfterRFC, err := time.Parse("2006-01-02", notAfter)
	if err == nil {
		notBefore = notBeforeRFC.Format("2006-01-02T15:04:05Z")
		notAfter = notAfterRFC.Format("2006-01-02T15:04:05Z")
	}
	fmt.Println("XXXXXXXXXXXXXXXXXXXXXXXXXXXAAAAAAAAAAAAAA")
	fmt.Println(notBeforeRFC, notAfterRFC, err)
	// RFC3339local := "2006-01-02T15:04:05Z"
	utcLoc, _ := time.LoadLocation("UTC")
	if notBefore != "" {
		timeNotBefore, err1 := time.ParseInLocation(time.RFC3339, notBefore, utcLoc)
		if err1 != nil {
			return "", "", fmt.Errorf("notBefore invalid timestamp. Please refer to RFC3339")
		}
		timeNotBefore = timeNotBefore.In(utcLoc)
		notBefore = timeNotBefore.Format("2006-01-02 15:04:05")
	}

	if notAfter != "" {
		timeNotAfter, err2 := time.ParseInLocation(time.RFC3339, notAfter, utcLoc)
		if err2 != nil {
			return "", "", fmt.Errorf("notAfter invalid timestamp. Please refer to RFC3339")
		}
		timeNotAfter = timeNotAfter.In(utcLoc)
		notAfter = timeNotAfter.Format("2006-01-02 15:04:05")
	}

	return notBefore, notAfter, nil

}

// func UrlQTimeStampsToUTCStringsUI(ctx *gin.Context) (string, string, error) {

// 	notBefore := ""
// 	notAfter := ""

// 	urlQuery := ctx.Request.URL.Query()
// 	for key, value := range urlQuery {
// 		queryValue := value[len(value)-1]
// 		switch key {
// 		case "notBefore":
// 			notBefore = queryValue
// 			break
// 		case "notAfter":
// 			notAfter = queryValue
// 			break
// 		}
// 	}

// 	if notBefore == "" || notAfter == "" {
// 		return "", "", fmt.Errorf("Empty date range in rides request")
// 	}

// 	notBeforeRFC, _ := time.Parse("2006-01-02", notBefore)
// 	notAfterRFC, _ := time.Parse("2006-01-02", notAfter)
// 	notBefore = notBeforeRFC.Format("2006-01-02T15:04:05Z")
// 	notAfter = notAfterRFC.Format("2006-01-02T15:04:05Z")
// 	// RFC3339local := "2006-01-02T15:04:05Z"
// 	// utcLoc, _ := time.LoadLocation("UTC")
// 	// if notBefore != "" {
// 	// 	timeNotBefore, err1 := time.ParseInLocation(time.RFC3339, notBefore, utcLoc)
// 	// 	if err1 != nil {
// 	// 		return "", "", fmt.Errorf("notBefore invalid timestamp. Please refer to RFC3339")
// 	// 	}
// 	// 	timeNotBefore = timeNotBefore.In(utcLoc)
// 	// 	notBefore = timeNotBefore.Format("2006-01-02 15:04:05")
// 	// }

// 	// if notAfter != "" {
// 	// 	timeNotAfter, err2 := time.ParseInLocation(time.RFC3339, notAfter, utcLoc)
// 	// 	if err2 != nil {
// 	// 		return "", "", fmt.Errorf("notAfter invalid timestamp. Please refer to RFC3339")
// 	// 	}
// 	// 	timeNotAfter = timeNotAfter.In(utcLoc)
// 	// 	notAfter = timeNotAfter.Format("2006-01-02 15:04:05")
// 	// }

// 	return notBefore, notAfter, nil

// }

func TimeStampsToUTCStrings(notBefore time.Time, notAfter time.Time) (string, string, error) {

	// RFC3339local := "2006-01-02T15:04:05Z"
	utcLoc, _ := time.LoadLocation("UTC")
	notBefore = notBefore.In(utcLoc)
	notBeforeStr := notBefore.Format("2006-01-02 15:04:05")

	notAfter = notAfter.In(utcLoc)
	notAfterStr := notAfter.Format("2006-01-02 15:04:05")

	return notBeforeStr, notAfterStr, nil

}
