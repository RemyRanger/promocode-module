package server

import (
	"APIs/internal/common/entities/custom_errors"
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	"github.com/pkg/errors"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
)

const render_name = "render"

const fieldErrMsg = "Field validation failed on the '%s' tag"

type ErrorResponse struct {
	ErrorMessage   string         `json:"error"`
	Details        []ErrorDetails `json:"details,omitempty"`
	Detail         string         `json:"detail,omitempty"`
	HTTPStatusCode int            `json:"status"`
}

func (a *ErrorResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type ErrorDetails struct {
	Field string `json:"field"`
	Msg   string `json:"msg"`
}

func RespondError(w http.ResponseWriter, r *http.Request, err error, code int) {
	log.Error().Str("service", render_name).Stack().Err(err).Msg("Error during process")

	if code == http.StatusBadRequest { // Only for errors on input fields validation
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make([]ErrorDetails, len(ve))
			for i, fe := range ve {
				out[i] = ErrorDetails{fe.Field(), fmt.Sprintf(fieldErrMsg, fe.Tag())}
			}

			RespondWithBody(w, r, &ErrorResponse{
				HTTPStatusCode: code,
				ErrorMessage:   http.StatusText(code),
				Details:        out,
			}, code)
		}
	} else {
		var customErr custom_errors.CustomError
		if errors.As(err, &customErr) { // Business custom errors case
			RespondWithBody(w, r, &ErrorResponse{
				HTTPStatusCode: code,
				ErrorMessage:   http.StatusText(code),
				Detail:         err.Error(),
			}, code)
		} else {
			RespondWithBody(w, r, &ErrorResponse{
				HTTPStatusCode: code,
				ErrorMessage:   http.StatusText(code),
			}, code)
		}
	}
}

func RespondWithBody(w http.ResponseWriter, r *http.Request, respData render.Renderer, code int) {
	render.Status(r, code)
	if err := render.Render(w, r, respData); err != nil {
		RespondError(w, r, err, http.StatusInternalServerError)
	}
}
