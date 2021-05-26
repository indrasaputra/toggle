package entity

import (
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	togglev1 "github.com/indrasaputra/toggle/proto/indrasaputra/toggle/v1"
)

// ErrInternal returns codes.Internal explained that unexpected behavior occurred in system.
func ErrInternal(message string) error {
	st := status.New(codes.Internal, message)
	te := &togglev1.ToggleError{
		ErrorCode: togglev1.ToggleErrorCode_INTERNAL,
	}
	res, err := st.WithDetails(te)
	if err != nil {
		return st.Err()
	}
	return res.Err()
}

// ErrEmptyToggle returns codes.InvalidArgument explained that the instance is empty or nil.
func ErrEmptyToggle() error {
	st := status.New(codes.InvalidArgument, "")
	br := createBadRequest(&errdetails.BadRequest_FieldViolation{
		Field:       "toggle instance",
		Description: "empty or nil",
	})

	te := &togglev1.ToggleError{
		ErrorCode: togglev1.ToggleErrorCode_EMPTY_TOGGLE,
	}
	res, err := st.WithDetails(br, te)
	if err != nil {
		return st.Err()
	}
	return res.Err()
}

// ErrAlreadyExists returns codes.AlreadyExists explained that the key already exists.
func ErrAlreadyExists() error {
	st := status.New(codes.AlreadyExists, "")
	te := &togglev1.ToggleError{
		ErrorCode: togglev1.ToggleErrorCode_ALREADY_EXISTS,
	}
	res, err := st.WithDetails(te)
	if err != nil {
		return st.Err()
	}
	return res.Err()
}

// ErrInvalidKey returns codes.InvalidArgument explained that the toggle's key is invalid.
func ErrInvalidKey() error {
	st := status.New(codes.InvalidArgument, "")
	br := createBadRequest(&errdetails.BadRequest_FieldViolation{
		Field:       "key",
		Description: "contain character outside of alphanumeric and dash",
	})

	te := &togglev1.ToggleError{
		ErrorCode: togglev1.ToggleErrorCode_INVALID_KEY,
	}
	res, err := st.WithDetails(br, te)
	if err != nil {
		return st.Err()
	}
	return res.Err()
}

// ErrInvalidValue returns codes.InvalidArgument explained that the toggle's value is invalid.
func ErrInvalidValue() error {
	st := status.New(codes.InvalidArgument, "")
	br := createBadRequest(&errdetails.BadRequest_FieldViolation{
		Field:       "is_enabled",
		Description: "value is not boolean",
	})

	te := &togglev1.ToggleError{
		ErrorCode: togglev1.ToggleErrorCode_INVALID_VALUE,
	}
	res, err := st.WithDetails(br, te)
	if err != nil {
		return st.Err()
	}
	return res.Err()
}

// ErrNotFound returns codes.NotFound explained that the toggle is not found.
func ErrNotFound() error {
	st := status.New(codes.NotFound, "")
	te := &togglev1.ToggleError{
		ErrorCode: togglev1.ToggleErrorCode_NOT_FOUND,
	}
	res, err := st.WithDetails(te)
	if err != nil {
		return st.Err()
	}
	return res.Err()
}

// ErrProhibitedToDelete returns codes.FailedPrecondition explained that the toggle's value is true and can't be deleted.
func ErrProhibitedToDelete() error {
	st := status.New(codes.FailedPrecondition, "toggle's is ENABLED hence it can't be deleted")
	te := &togglev1.ToggleError{
		ErrorCode: togglev1.ToggleErrorCode_PROHIBITED_TO_DELETE,
	}
	res, err := st.WithDetails(te)
	if err != nil {
		return st.Err()
	}
	return res.Err()
}

func createBadRequest(details ...*errdetails.BadRequest_FieldViolation) *errdetails.BadRequest {
	return &errdetails.BadRequest{
		FieldViolations: details,
	}
}
