package gapi

import (
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// fieldViolation producing new (*errdetails.BadRequest_FieldViolation)
func fieldViolation(field string, err error) *errdetails.BadRequest_FieldViolation {
	return &errdetails.BadRequest_FieldViolation{
		Field:       field,
		Description: err.Error(),
	}
}

// invalidArgumentError wrapper for list of error to create human read-able error massage using standart google.golang.org/grpc lib
func invalidArgumentError(violations []*errdetails.BadRequest_FieldViolation) error {
	badRequest := &errdetails.BadRequest{FieldViolations: violations}
	//main eror status
	statusInvalid := status.New(codes.InvalidArgument, "invalid parameters")
	//Feeding error detail What field and why it's violated?
	statusDetails, err := statusInvalid.WithDetails(badRequest)
	if err != nil {
		return statusInvalid.Err()
	}

	return statusDetails.Err()
}

/*{
    "code": 3,
    "message": "invalid parameters",
    "details": [
        {
            "@type": "type.googleapis.com/google.rpc.BadRequest",
            "field_violations": [
                {
                    "field": "email",
                    "description": "must between 3 - 200 character"
                }
            ]
        }
    ]
}
*/
/*
 mutiple violation
{
    "code": 3,
    "message": "invalid parameters",
    "details": [
        {
            "@type": "type.googleapis.com/google.rpc.BadRequest",
            "field_violations": [
                {
                    "field": "username",
                    "description": "must contain only lowercase letters, digits, or underscore"
                },
                {
                    "field": "password",
                    "description": "must between 6 - 100 character"
                },
                {
                    "field": "full_name",
                    "description": "must contain only letters or spaces"
                },
                {
                    "field": "email",
                    "description": "must between 3 - 200 character"
                }
            ]
        }
    ]
}
*/

// unauthenticatedError producting error codes.Unauthenticated (16)  : err
func unauthenticatedError(err error) error {
	return status.Errorf(codes.Unauthenticated, "unauthorized: %s", err)
}
