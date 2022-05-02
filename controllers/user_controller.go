package controllers

import (
	"Inshorts/configs"
	"Inshorts/models"
	"Inshorts/responses"
	"encoding/json"
	"fmt"
	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/bson"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
)

var stateCollection *mongo.Collection = configs.GetCollection(configs.DB, "states")
var filename = ".\\resource\\static\\data.min.json"
var states = map[string]string{
	"Rajasthan":                   "RJ",
	"Andaman and Nicobar Islands": "AN",
	"Andhra Pradesh":              "AP",
	"Arunachal Pradesh":           "AR",
	"Assam":                       "AS",
	"Bihar":                       "BR",
	"Chandigarh":                  "CH",
	"Chhattisgarh":                "CT",
	"Delhi":                       "DL",
	"Dadra and Nagar Haveli":      "DN",
	"Goa":                         "GA",
	"Gujarat":                     "GJ",
	"Himachal Pradesh":            "HP",
	"Haryana":                     "HR",
	"Jharkhand":                   "JH",
	"Jammu and Kashmir":           "JK",
	"Karnataka":                   "KA",
	"Kerala":                      "KL",
	"Ladakh":                      "LA",
	"Lakshadweep":                 "LD",
	"Maharashtra":                 "MH",
	"Manipur":                     "MN",
	"Meghalaya":                   "ML",
	"Madhya Pradesh":              "MP",
	"Mizoram":                     "MZ",
	"Nagaland":                    "NL",
	"Odisha":                      "OR",
	"Punjab":                      "PB",
	"Puducherry":                  "PY",
	"Sikkim":                      "SK",
	"Telangana":                   "TG",
	"Tamil Nadu":                  "TN",
	"Tripura":                     "TR",
	"Uttar Pradesh":               "UP",
	"Uttarakhand":                 "UT",
	"West Bengal":                 "WB",
}

func GetStateActiveCases(c echo.Context) error {
	fullState := c.Param("state")
	state := states[fullState]
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

	var covidCase models.CovidCases
	defer cancel()

	err := stateCollection.FindOne(ctx, bson.M{"state": state}).Decode(&covidCase)
	if err != nil {

		file, _ := ioutil.ReadFile(filename)
		var result map[string]interface{}
		_ = json.Unmarshal([]byte(file), &result)
		stateMap, _ := result[state].(map[string]interface{})
		if stateMap == nil {
			return c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "Please enter correct state", Data: &echo.Map{"data": ""}})
		}

		totalCases := stateMap["total"].(map[string]interface{})

		confirmedCases := totalCases["confirmed"].(float64)

		newCovidCase := models.CovidCases{
			State:     state,
			TotalCase: confirmedCases,
			Timestamp: time.Now().UTC().String(),
		}

		newResult, newError := stateCollection.InsertOne(ctx, newCovidCase)
		if newError != nil {
			return c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": newError.Error()}})
		}
		fmt.Println(newResult)
		return c.JSON(http.StatusCreated, responses.UserResponse{Status: http.StatusCreated, Message: "success", Data: &echo.Map{"data": newCovidCase}})
	}
	return c.JSON(http.StatusOK, responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"data": covidCase}})
}

func GetStateActiveCasesUsingPosition(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	tr := &http.Transport{}
	client := &http.Client{Transport: tr}
	x := c.Param("x")
	y := c.Param("y")
	state := ""
	var covidCase models.CovidCases

	response, err := client.Get(
		"http://api.positionstack.com/v1/reverse?access_key=d6c6d6181c6441a54d2bf2b751aecbdb&query=" + x + "," + y,
	)
	defer response.Body.Close()
	defer cancel()
	fmt.Println(err)
	if response.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}

		bodyString := string(bodyBytes)
		var result map[string]interface{}
		_ = json.Unmarshal([]byte(bodyString), &result)
		tempData, _ := result["data"]
		for key, value := range tempData.([]interface{})[0].(map[string]interface{}) {

			if key == "region_code" {
				state = value.(string)
				break
			}

		}
		err = stateCollection.FindOne(ctx, bson.M{"state": state}).Decode(&covidCase)
		if err != nil {

			file, _ := ioutil.ReadFile(filename)
			var result map[string]interface{}
			_ = json.Unmarshal([]byte(file), &result)
			stateMap, _ := result[state].(map[string]interface{})
			if stateMap == nil {
				return c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "Please enter correct state", Data: &echo.Map{"data": ""}})
			}

			totalCases := stateMap["total"].(map[string]interface{})

			confirmedCases := totalCases["confirmed"].(float64)

			newCovidCase := models.CovidCases{
				State:     state,
				TotalCase: confirmedCases,
				Timestamp: time.Now().UTC().String(),
			}

			newResult, newError := stateCollection.InsertOne(ctx, newCovidCase)
			if newError != nil {
				return c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": newError.Error()}})
			}
			fmt.Println(newResult)
			return c.JSON(http.StatusCreated, responses.UserResponse{Status: http.StatusCreated, Message: "success", Data: &echo.Map{"data": newCovidCase}})
		}
		return c.JSON(http.StatusOK, responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"data": covidCase}})
	} else {
		return c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "Please enter coordinates", Data: &echo.Map{"data": ""}})
	}

}
