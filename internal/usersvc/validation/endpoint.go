package validation

import (
	"github.com/go-kit/kit/endpoint"
	"context"
)

type validationRequest struct {
	Email 	string
}

type validationResponse struct {
	Err error            `json:"error,omitempty"`
}


func makeValidationEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(validationRequest)
		err := s.Validate(req.Email)

		return validationResponse{Err: err}, err
	}
}