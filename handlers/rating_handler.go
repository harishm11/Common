package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/harishm11/Common/builder"
	"github.com/harishm11/Common/clients"
	"github.com/harishm11/Common/logger"
	models "github.com/harishm11/Common/models/PolicyModels"
)

func RateRequestHandler(bundle *models.Bundle, ctx *fiber.Ctx) (interface{}, error) {
	logger.GetLogger().Info("Executing Rating")

	rateRequest, err := builder.BuildRateRequest(bundle)
	if err != nil {
		logger.GetLogger().Error(err, "Failed to prepare rate request")
		return nil, ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to prepare rate request"})
	}

	rateResponse, err := clients.CallRatingService(rateRequest)
	if err != nil {
		logger.GetLogger().Error(err, "Failed to get rate response")
		return nil, ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get rate response"})
	}

	logger.GetLogger().Info("Rating response: ", rateResponse)

	var response interface{} = rateResponse
	return response, nil

}
