package warehouses

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	errs "github.com/Cataloft/warehouse-service/internal/constants/errors"
	"github.com/gin-gonic/gin"
)

type Attributes struct {
	Name         string `json:"name"`
	Availability bool   `json:"availability"`
}

type Meta struct {
	TotalAmount int32 `json:"totalAmount"`
}

type Data struct {
	Type       string     `json:"type"`
	ID         int64      `json:"id"`
	Attributes Attributes `json:"attributes"`
	Meta       Meta       `json:"meta"`
}

type Response struct {
	Data []Data `json:"data"`
}

type Errors struct {
	Status int    `json:"status"`
	Title  string `json:"title"`
}

type ErrResponse struct {
	Errors []Errors `json:"errors"`
}

type TotalAmountGetter interface {
	GetWarehouse(warehouseID int64) (Attributes, int32, error)
}

func GetOne(getter TotalAmountGetter, log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		log = log.With(slog.String("op", "handlers.warehouses.GetOne"))

		var (
			resp    Response
			errResp ErrResponse
		)

		param := c.Param("id")

		warehouseID, err := strconv.ParseInt(param, 10, 64)
		if err != nil {
			errResp.Errors = append(errResp.Errors, Errors{
				Status: http.StatusBadRequest,
				Title:  err.Error(),
			})

			c.JSON(http.StatusBadRequest, errResp)
			log.Error("Error requested data", "error", err)

			return
		}

		attrs, total, err := getter.GetWarehouse(warehouseID)
		if err != nil {
			if errors.Is(err, fmt.Errorf("%s", errs.WarehouseError)) {
				errResp.Errors = append(errResp.Errors, Errors{
					Status: http.StatusNotFound,
					Title:  err.Error(),
				})

				c.JSON(http.StatusNotFound, errResp)
				log.Error("Error db", "error", err)

				return
			}

			errResp.Errors = append(errResp.Errors, Errors{
				Status: http.StatusInternalServerError,
				Title:  err.Error(),
			})

			c.JSON(http.StatusInternalServerError, errResp)
			log.Error("Error db", "error", err)

			return
		}

		resp.Data = append(resp.Data, Data{
			Type:       "warehouses",
			ID:         warehouseID,
			Attributes: attrs,
			Meta:       Meta{TotalAmount: total},
		})

		c.JSON(http.StatusOK, resp)
		log.Info("Warehouses successfully got")
	}
}
