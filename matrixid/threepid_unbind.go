package matrixid

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"

	"github.com/gernest/sydent-go/config"
	"github.com/gernest/sydent-go/core"
	"github.com/gernest/sydent-go/models"
	"github.com/gernest/signedjson"
	"github.com/labstack/echo"
)

type UnbindOption struct {
	Threepid *models.Association `json:"threepid,omitempty"`
	MatrixID string              `json:"mxid"`
}

func Unbind(coreContext *core.Ctx, fedClient config.HTTPClient) echo.HandlerFunc {
	cache := make(map[string]*ServerKeys)
	var mu sync.RWMutex
	getKey := func(server string) *ServerKeys {
		mu.RLock()
		k := cache[server]
		mu.RUnlock()
		return k
	}

	setKey := func(server string, key *ServerKeys) {
		mu.Lock()
		cache[server] = key
		mu.Unlock()
	}

	name := coreContext.Config.Server.Name
	unbind := RemoveBinding(coreContext)
	lg := coreContext.Log
	return func(ctx echo.Context) error {
		req := ctx.Request()
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			RequestError(lg, req, err)
			return ctx.JSON(http.StatusBadRequest, models.NewError(
				models.ErrBadJSON,
				"Malformed JSON",
			))
		}
		var opts UnbindOption
		err = json.Unmarshal(body, &opts)
		if err != nil {
			RequestError(lg, req, err)
			return ctx.JSON(http.StatusBadRequest, models.NewError(
				models.ErrBadJSON,
				"Malformed JSON",
			))
		}
		var missing []string
		if opts.Threepid == nil {
			missing = append(missing, "threepid")
		}
		if opts.MatrixID == "" {
			missing = append(missing, "mxid")
		}
		if len(missing) > 0 {
			return ctx.JSON(http.StatusBadRequest, models.NewError(
				models.ErrMissingParam,
				fmt.Sprintf("Missing parameters: %s", strings.Join(missing, ",")),
			))
		}
		if opts.Threepid.Medium == "" {
			missing = append(missing, "medium")
		}
		if opts.Threepid.Address == "" {
			missing = append(missing, "address")
		}
		if len(missing) > 0 {
			return ctx.JSON(http.StatusBadRequest, models.NewError(
				models.ErrMissingParam,
				fmt.Sprintf("Threepid lack: %s", strings.Join(missing, ",")),
			))
		}

		authorization := req.Header.Get("Authorization")
		if authorization == "" {
			return ctx.JSON(http.StatusUnauthorized, models.NewError(
				models.ErrForbidden,
				"Missing Authorization headers",
			))
		}
		if !strings.HasPrefix(authorization, "X-Matrix") {
			return ctx.JSON(http.StatusUnauthorized, models.NewError(
				models.ErrForbidden,
				"Missing X-Matrix Authorization header",
			))
		}
		parts := strings.Split(authorization, " ")
		var origin, key, sig string
		if len(parts) == 2 {
			for _, v := range strings.Split(parts[1], ",") {
				kv := strings.Split(v, "=")
				if len(kv) == 2 {
					switch kv[0] {
					case "origin":
						origin = stripQuote(kv[1])
					case "key":
						key = stripQuote(kv[1])
					case "sig":
						sig = stripQuote(kv[1])
					}
				}
			}
		}
		if origin == "" {
			missing = append(missing, "origin")
		}
		if key == "" {
			missing = append(missing, "key")
		}
		if sig == "" {
			missing = append(missing, "sig")
		}
		if len(missing) > 0 {
			return ctx.JSON(http.StatusUnauthorized, models.NewError(
				models.ErrForbidden,
				"Bad X-Matrix Authorization header, missing "+strings.Join(missing, ","),
			))
		}
		request := map[string]interface{}{
			"method":         req.Method,
			"uri":            req.URL.String(),
			"destination_is": name,
			"content":        string(body),
			"signatures": map[string]interface{}{
				origin: map[string]string{
					key: sig,
				},
			},
		}
		keys := getKey(origin)
		if keys == nil {
			uri := fmt.Sprintf("matrix://%s/_matrix/key/v2/server/", origin)
			freq, _ := http.NewRequest(http.MethodGet, uri, nil)
			res, err := fedClient.Do(freq)
			if err != nil {
				RequestError(lg, req, err)
				return ctx.JSON(http.StatusUnauthorized, models.NewError(
					models.ErrForbidden,
					"Failed to retrieve verification keys ",
				))
			}
			defer res.Body.Close()
			rb, err := ioutil.ReadAll(res.Body)
			if err != nil {
				RequestError(lg, req, err)
				return ctx.JSON(http.StatusUnauthorized, models.NewError(
					models.ErrForbidden,
					"Failed to retrieve verification keys ",
				))
			}
			var rs ServerKeys
			err = json.Unmarshal(rb, &rs)
			if err != nil {
				RequestError(lg, req, err)
				return ctx.JSON(http.StatusUnauthorized, models.NewError(
					models.ErrForbidden,
					"Failed to retrieve verification keys ",
				))
			}
			setKey(origin, &rs)
			keys = &rs
		}
		if k, ok := keys.VerifyKeys[key]; ok {
			byt, err := signedjson.DecodeBase64(k.Key)
			if err != nil {
				RequestError(lg, req, err)
				return InternalError(ctx)
			}
			vk, err := signedjson.DecodeVerifyKeyBytes(key, byt)
			if err != nil {
				RequestError(lg, req, err)
				return InternalError(ctx)
			}
			err = vk.Verify(request, origin)
			if err != nil {
				RequestError(lg, req, err)
				return InternalError(ctx)
			}
		} else {
			return ctx.JSON(http.StatusUnauthorized, models.NewError(
				models.ErrForbidden,
				"No matching signature found",
			))
		}

		if !strings.HasSuffix(opts.MatrixID, ":"+origin) {
			return ctx.JSON(http.StatusForbidden, models.NewError(
				models.ErrForbidden,
				"Origin server name does not match mxid",
			))
		}
		err = unbind(req.Context(), opts.Threepid)
		if err != nil {
			RequestError(lg, req, err)
			return InternalError(ctx)
		}
		return ctx.JSON(http.StatusOK, map[string]interface{}{})
	}
}

type Key struct {
	Key       string `json:"key"`
	ExpiredTS int64  `json:"expired_ts,omitempty"`
}

// ServerKeys is an object returned when asking for verification keys on a
// matrix server
type ServerKeys struct {
	Name          string                       `json:"server_name"`
	VerifyKeys    map[string]Key               `json:"verify_keys"`
	OldVerifyKeys map[string]Key               `json:"old_verify_keys"`
	Signatures    map[string]map[string]string `json:"signatures"`
	ValidUntil    int64                        `json:"valid_until_ts"`
}

func stripQuote(s string) string {
	return strings.TrimFunc(s, func(r rune) bool {
		return r == '"'
	})
}
