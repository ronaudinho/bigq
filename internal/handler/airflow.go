package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/ronaudinho/bigq/internal/model"

	"github.com/labstack/echo/v4"
)

// RecvAirflow receives Airflow related payload to be processed
func (h *Handler) RecvAirflow(c echo.Context) error {
	body := c.Request().Body
	defer body.Close()
	b, err := ioutil.ReadAll(body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "StatusBadRequest")
	}

	// TODO validate
	var task model.Task
	err = json.Unmarshal(b, &task)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "StatusBadRequest")
	}

	// TODO should probably set airflow as exchange instead
	if task.RoutingKey == "" {
		task.RoutingKey = "airflow"
	}
	err = h.service.RecvAirflow(&task)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, task)
	}

	return c.JSON(http.StatusOK, task)
}

// SendAirflow picks up queued result from RecvAirflow on demand
func (h *Handler) SendAirflow(c echo.Context) error {
	res, err := h.service.SendAirflow("airflow_result")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	if res == nil {
		return c.JSON(http.StatusOK, "nothing to do")
	}
	if len(res.Args) == 0 {
		return c.JSON(http.StatusOK, "nothing to do")
	}

	return c.JSON(http.StatusOK, res.Args[0].Value)
}
