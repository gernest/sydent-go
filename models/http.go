package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// EnsureParams makes sure the query params are in params.
func EnsureParams(req *http.Request, args ...string) (map[string]string, error) {
	m, err := loadArgs(req)
	if err != nil {
		return nil, err
	}
	var missing []string
	for _, v := range args {
		x := m[v]
		if x == "" {
			missing = append(missing, v)
		}
	}
	if len(missing) > 0 {
		return nil, NewError(
			ErrMissingParam,
			fmt.Sprintf("Missing parameters: %s", strings.Join(missing, ",")),
		)
	}
	return m, nil
}

func loadArgs(r *http.Request) (map[string]string, error) {
	args := make(map[string]interface{})
	if r.Header.Get("Content-Type") == "application/json" {
		if r.Body != nil {
			b, err := ioutil.ReadAll(r.Body)
			if err != nil {
				//TODO handle error
			} else {
				err = json.Unmarshal(b, &args)
				if err != nil {
					return nil, NewError(
						ErrBadJSON,
						"Malformed JSON",
					)
				}
			}

		}
	}
	if len(args) == 0 {
		q := r.URL.Query()
		for k := range q {
			args[k] = q.Get(k)
		}
	}
	o := make(map[string]string)
	for k, v := range args {
		if s, ok := v.(string); ok {
			o[k] = s
		} else {
			o[k] = fmt.Sprint(v)
		}
	}
	return o, nil
}
