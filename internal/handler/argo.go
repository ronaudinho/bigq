package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/ronaudinho/bigq/internal/model"

	"github.com/labstack/echo/v4"
)

// RecvArgo receives Argo related payload to be processed
func (h *Handler) RecvArgo(c echo.Context) error {
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

	// NOTE enforce task.Name since we are splitting the handler function
	task.Name = "argo"
	// TODO should probably set argo as exchange instead
	if task.RoutingKey == "" {
		task.RoutingKey = "argo"
	}
	err = h.service.RecvArgo(&task)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, task)
	}

	return c.JSON(http.StatusOK, task)
}

// SendArgo picks up queued result from RecvArgo on demand
func (h *Handler) SendArgo(c echo.Context) error {
	res, err := h.service.SendArgo("argo_result")
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
