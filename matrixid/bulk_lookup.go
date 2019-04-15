package matrixid

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gernest/sydent-go/core"
	"github.com/gernest/sydent-go/models"
	"github.com/labstack/echo"
)

var ErrThreepidNotAList = models.NewError(
	models.ErrInvalidParam,
	"threepids must be a list",
)
var ErrMissingThreepid = models.NewError(
	models.ErrMissingParam,
	"missing threepids in request body",
)

func BulkLookup(coreContext *core.Ctx, m Metric) echo.HandlerFunc {
	count := m.CountError("bulk_lookup")
	db := coreContext.Store
	lg := coreContext.Log
	return func(ctx echo.Context) error {
		req := ctx.Request()
		b, err := ioutil.ReadAll(req.Body)
		if err != nil {
			count.Inc()
			RequestError(lg, req, err)
			return ctx.JSON(http.StatusBadRequest, ErrMissingThreepid)
		}
		var o models.BulkLookupRequest
		err = json.Unmarshal(b, &o)
		if err != nil {
			count.Inc()
			RequestError(lg, req, err)
			return ctx.JSON(http.StatusBadRequest, ErrThreepidNotAList)
		}
		if len(o.Threepids) > 0 {
			a, err := db.GlobalGetMxids(ctx.Request().Context(), o.Threepids)
			if err != nil {
				count.Inc()
				RequestError(lg, req, err)
				return ctx.JSON(http.StatusOK, map[string]interface{}{})
			}
			var result models.BulkLookupRequest
			for _, v := range a {
				result.Threepids = append(result.Threepids, []string{
					v.Medium, v.Address, v.MatrixID,
				})
			}
			return ctx.JSON(http.StatusOK, result)
		}
		return ctx.JSON(http.StatusOK, models.BulkLookupRequest{})
	}
}
