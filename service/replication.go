package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gernest/sydent-go/core"
	"github.com/gernest/sydent-go/models"
	"github.com/gernest/sydent-go/store"
	"github.com/labstack/echo"
)

func Replicate(coreContext *core.Ctx, m Metric) echo.HandlerFunc {
	count := m.CountError("replication")
	db := coreContext.Store
	return func(ctx echo.Context) error {
		req := ctx.Request()
		tls := req.TLS
		name := tls.PeerCertificates[0].Subject.CommonName
		requestContext := req.Context()
		activePeer, err := db.GetPeerByName(requestContext, name)
		if err != nil {
			count.Inc()
			RequestError(coreContext.Log, req, err)
			return ctx.JSON(http.StatusForbidden, models.NewError(
				models.ErrUnknownPeer,
				"This peer is not known to this server",
			))
		}
		key, err := GetVerifyKeyFromPeer(activePeer)
		if err != nil {
			count.Inc()
			RequestError(coreContext.Log, req, err)
			return ctx.JSON(http.StatusForbidden, models.NewError(
				models.ErrUnknownPeer,
				"This peer is not known to this server",
			))
		}
		if req.Header.Get("Content-Type") != "application/json" {
			m := models.NewError(
				models.ErrNotJSON,
				"his endpoint expects JSON",
			)
			RequestError(coreContext.Log, req, m)
			return ctx.JSON(http.StatusForbidden, m)
		}
		b, err := ioutil.ReadAll(req.Body)
		if err != nil {
			count.Inc()
			RequestError(coreContext.Log, req, err)
			return ctx.JSON(http.StatusForbidden, models.NewError(
				models.ErrNotJSON,
				"Missing json payload",
			))
		}
		var o Payload
		err = json.Unmarshal(b, &o)
		if err != nil {
			count.Inc()
			RequestError(coreContext.Log, req, err)
			return ctx.JSON(http.StatusForbidden, models.NewError(
				models.ErrBadJSON,
				"Malformed JSON",
			))
		}
		if o.SignedAssociations == nil {
			m := models.NewError(
				models.ErrBadJSON,
				`No "sgAssocs" key in JSON`,
			)
			RequestError(coreContext.Log, req, m)
			return ctx.JSON(http.StatusForbidden, m)
		}

		var failed []int64
		for _, a := range o.SignedAssociations {
			if err := VerifySignedAssociation(requestContext, key, activePeer.Name, a.SignedAssociation); err != nil {
				failed = append(failed, a.OriginID)
			}
		}
		if len(failed) > 0 {
			m := models.NewError(
				models.ErrVerificationFailed,
				fmt.Sprintf("Verification failed for one or more associations failed_ids=%v", failed),
			)
			RequestError(coreContext.Log, req, m)
			return ctx.JSON(http.StatusBadRequest, m)
		}
		tx, err := db.DB().BeginTx(requestContext, nil)
		if err != nil {
			count.Inc()
			RequestError(coreContext.Log, req, err)
			return ctx.JSON(http.StatusInternalServerError, models.NewError(
				models.ErrUnknown, "Ooops something went wrong please try again",
			))
		}
		idStore := store.New(tx, db.Metric())
		for _, a := range o.SignedAssociations {
			as, err := AssociationFromMap(a.SignedAssociation)
			if err != nil {
				count.Inc()
				RequestError(coreContext.Log, req, err)
				tx.Rollback()
				return internalError(ctx)
			}
			if as.MatrixID != "" {
				b, _ := json.Marshal(as)
				err = idStore.GlobalAddAssociation(requestContext, as,
					activePeer.Name, a.OriginID, string(b),
				)
				if err != nil {
					count.Inc()
					RequestError(coreContext.Log, req, err)
					tx.Rollback()
					return internalError(ctx)
				}
			} else {
				err = idStore.GlobalRemoveAssociation(requestContext, as.Medium, as.Address)
				if err != nil {
					count.Inc()
					RequestError(coreContext.Log, req, err)
					tx.Rollback()
					return internalError(ctx)
				}
			}
		}
		err = tx.Commit()
		if err != nil {
			count.Inc()
			RequestError(coreContext.Log, req, err)
			return internalError(ctx)
		}
		return ctx.JSON(http.StatusOK, models.Success{
			Success: true,
		})
	}
}

func internalError(ctx echo.Context) error {
	return ctx.JSON(http.StatusInternalServerError, models.NewError(
		models.ErrUnknown,
		"Having issue processing the request",
	))
}
