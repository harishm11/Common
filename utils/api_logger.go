package utils

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/harishm11/Common/logger"
)

func RequestLogger(c *fiber.Ctx) error {
	requestMethod := c.Method()
	requestPath := c.Path()

	// Log request headers
	headers := make([]string, 0)
	c.Request().Header.VisitAll(func(key, value []byte) {
		headers = append(headers, string(key)+": "+string(value))
	})

	// Access the request body
	requestBody := string(c.Request().Body())

	// var requestBody interface{}
	// if len(c.Request().Body()) > 0 {
	// 	if err := json.Unmarshal(c.Request().Body(), &requestBody); err != nil {
	// 		return c.Status(fiber.StatusBadRequest).SendString("Request body unmarshal error")
	// 	}
	// }

	// Log request details
	logger.GetLogger().LogAnything("Request:" + requestMethod + requestPath)
	logger.GetLogger().LogAnything("Request headers:" + strings.Join(headers, ", "))
	logger.GetLogger().LogAnything("Request body:")
	logger.GetLogger().LogAnything(requestBody)

	// Continue to the next middleware or handler
	return c.Next()
}

func ResponseLogger(c *fiber.Ctx) error {

	// Call the next handler
	if err := c.Next(); err != nil {
		return err // Return if there's an error
	}

	// Log the response status code
	logger.GetLogger().LogAnything("Response status:" + fmt.Sprint(c.Response().StatusCode()))

	// Log the response body
	responseBody := string(c.Response().Body())
	// responseJsonBody, err := jsonStringToJSON(responseBody)
	// if err != nil {
	//     return c.Status(fiber.StatusBadRequest).SendString("Response body unmarshal error")
	// }
	// var responseBody interface{}
	// if len(c.Response().Body()) > 0 {
	// 	if err := json.Unmarshal(c.Response().Body(), &responseBody); err != nil {
	// 		return c.Status(fiber.StatusBadRequest).SendString("Response body unmarshal error")
	// 	}
	// }

	logger.GetLogger().LogAnything("Response body:")
	logger.GetLogger().LogAnything(responseBody)
	return nil

}

func jsonStringToJSON(jsonStr string) (interface{}, error) {
	var jsonObj interface{}
	err := json.Unmarshal([]byte(jsonStr), &jsonObj)
	if err != nil {
		return nil, err
	}
	return jsonObj, nil
}
