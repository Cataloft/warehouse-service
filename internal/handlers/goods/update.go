package goods

import (
	"errors"
	"fmt"
	errs "lamoda_task/internal/constants/errors"
	"lamoda_task/internal/tools/requests/goods"
	"lamoda_task/internal/tools/validation"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type Response struct {
	Data []goods.Data `json:"data"`
}

type Meta struct {
	UniqueCode string `json:"uniqueCode"`
	ErrTitle   string `json:"errTitle"`
}

type Errors struct {
	Status int    `json:"status"`
	Meta   []Meta `json:"meta"`
}

type ErrResponse struct {
	Errors []Errors `json:"errors"`
}

type StorageManager interface {
	UpdateGood(uniqueCode string, amount int32) (int32, error)
	GetAvailableGood(uniqueCode string, amount int32) error
}

func Update(manager StorageManager, log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		log = log.With(slog.String("op", "handlers.goods.Update"))

		var (
			req     goods.Request
			resp    Response
			errResp ErrResponse
		)

		errResp.Errors = make([]Errors, 1)

		err := c.ShouldBindJSON(&req)
		if err != nil {
			errResp.Errors[0].Status = http.StatusBadRequest
			errResp.Errors[0].Meta = append(errResp.Errors[0].Meta, Meta{
				ErrTitle: fmt.Sprintf("%e", err),
			})
			c.JSON(http.StatusBadRequest, errResp)
			log.Error("Error requested data", "error", err)

			return
		}

		if code, err := validation.ValidateRequestStruct(req); err != nil {
			errResp.Errors[0].Status = http.StatusBadRequest
			errResp.Errors[0].Meta = append(errResp.Errors[0].Meta, Meta{
				UniqueCode: code,
				ErrTitle:   err.Error(),
			})
		}

		for _, dataVal := range req.Data {
			if err = manager.GetAvailableGood(dataVal.ID, dataVal.Attributes.Amount); err != nil {
				switch {
				case errors.Is(err, pgx.ErrNoRows):
					errResp.Errors[0].Meta = append(errResp.Errors[0].Meta, Meta{
						UniqueCode: dataVal.ID,
						ErrTitle:   errs.WarehouseAvailabilityError,
					})
				case errors.Is(err, fmt.Errorf("%s", errs.GoodError)):
					errResp.Errors[0].Meta = append(errResp.Errors[0].Meta, Meta{
						UniqueCode: dataVal.ID,
						ErrTitle:   err.Error(),
					})
				default:
					errResp.Errors[0].Status = http.StatusInternalServerError
					errResp.Errors[0].Meta = append(errResp.Errors[0].Meta, Meta{
						ErrTitle: err.Error(),
					})

					c.JSON(http.StatusInternalServerError, errResp)
					log.Error("Error getting available goods", "error", err)

					return
				}
			}

			updatedAmount, err := manager.UpdateGood(dataVal.ID, dataVal.Attributes.Amount)
			if err != nil {
				errResp.Errors[0].Status = http.StatusInternalServerError
				errResp.Errors[0].Meta = append(errResp.Errors[0].Meta, Meta{
					ErrTitle: err.Error(),
				})
				c.JSON(http.StatusInternalServerError, errResp)
				log.Error("Error getting available goods", "error", err)

				return
			}

			resp.Data = append(resp.Data, goods.Data{
				Type:       dataVal.Type,
				ID:         dataVal.ID,
				Attributes: goods.Attributes{Amount: updatedAmount},
			})
		}

		if errResp.Errors[0].Status != 0 {
			c.JSON(http.StatusMultiStatus, errResp)
			log.Info("Goods partially updated")

			return
		}

		c.JSON(http.StatusOK, resp)
		log.Info("Goods successfully updated")
	}
}
