package validation

import (
	"net/http"
	"context"
	"encoding/json"

	"github.com/gorilla/mux"
	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/kit/endpoint"
)

func MakeHandler(service Service, logger kitlog.Logger) http.Handler {

	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encodeError),
	}


	var validationEndpoint endpoint.Endpoint
	validationEndpoint = makeValidationEndpoint(service)
	validationEndpoint = loggingMiddlware(logger)(validationEndpoint)
	validationHandler := kithttp.NewServer(
		validationEndpoint,
		decodeValidationRequest,
		encodeResponse,
		opts...,
	)

	router := mux.NewRouter()
	router.Handle("/api/users/validate", validationHandler).Methods("POST")

	return router
}

func decodeValidationRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var body struct {
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, err
	}

	return validationRequest{
		Email: body.Email,
	}, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {

	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
		return nil
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

type errorer interface {
	error() error
}

// encode errors from business-logic
func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err {
	case ErrEmptyEmail:
		w.WriteHeader(http.StatusBadRequest)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
