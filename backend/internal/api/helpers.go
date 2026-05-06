package api

import (
	"fmt"
	"strconv"

	"devhelper/internal/utils"

	"github.com/gin-gonic/gin"
)

func parseUintParam(c *gin.Context, name string, out *uint) (uint, error) {
	s := c.Param(name)
	n, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		utils.BadRequest(c, fmt.Sprintf("invalid %s", name))
		return 0, err
	}
	*out = uint(n)
	return uint(n), nil
}
