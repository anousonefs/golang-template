package helper

import (
	"errors"
	"net/http"

	hspb "github.com/anousoneFS/clean-architecture/proto/http"
	"google.golang.org/genproto/googleapis/rpc/code"
	edpb "google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	// ErrUnauthorized is returned when the user is not authorized.
	ErrUnauthorized = errors.New("unauthorized")
	// ErrPermissionDenied is returned when the user
	// does not have permission to perform the action.
	ErrPermissionDenied = errors.New("permission denied")
	// ErrNoInfo is returned when no info is found.
	ErrNoInfo = errors.New("no info")
	// ErrUnProcessAbleEntity is returned when it unprocessable entity
	ErrUnProcessAbleEntity = errors.New("unprocessable entity")
	// ErrInternalServerError is returned when it internal server error
	ErrInternalServerError = errors.New("internal server error")
	// ErrAlreadyExist is returned when resource was existed
	ErrAlreadyExist = errors.New("already exist")
)

var StatusInvalidENUM = func() *status.Status {
	s, _ := status.New(codes.InvalidArgument, "Invalid ENUM").
		WithDetails(&edpb.ErrorInfo{
			Reason: "INVALID_ENUM",
			Domain: "laopost",
		})
	return s
}()

var StatusInvalidUUID = func() *status.Status {
	s, _ := status.New(codes.InvalidArgument, "Invalid UUID").
		WithDetails(&edpb.ErrorInfo{
			Reason: "INVALID_UUID",
			Domain: "laopost",
		})
	return s
}()

var StatusBindingFailure = func() *status.Status {
	s, _ := status.New(codes.InvalidArgument, "Binding JSON body failure. Please pass a valid JSON body").
		WithDetails(&edpb.ErrorInfo{
			Reason: "BINDING_FAILURE",
			Domain: "laopost",
		})
	return s
}()

var StatusInvalidPassword = func() *status.Status {
	s, _ := status.New(codes.InvalidArgument, "password length should greater than 6").
		WithDetails(&edpb.ErrorInfo{
			Reason: "INVALID_PASSWORD",
			Domain: "laopost",
		})
	return s
}()

var StatusUnauthenticated = func() *status.Status {
	s, _ := status.New(codes.Unauthenticated, "ID token not valid. Please pass a valid ID token.").
		WithDetails(
			&edpb.ErrorInfo{
				Reason: "TOKEN_INVALID",
				Domain: "laopost.com.la",
				Metadata: map[string]string{
					"service": "laopost",
				},
			})
	return s
}()

var StatusSessionExpired = func() *status.Status {
	s, _ := status.New(codes.Unauthenticated, "Session has been expired. Please make a new session and try again.").
		WithDetails(
			&edpb.ErrorInfo{
				Reason: "SESSION_EXPIRED",
				Domain: "laopost",
			})
	return s
}()

var StatusPermissionDenied = func() *status.Status {
	s, _ := status.New(codes.PermissionDenied, "You does't have sufficient permission to perform action.").
		WithDetails(
			&edpb.ErrorInfo{
				Reason: "INSUFFICIENT_PERMISSION",
				Domain: "laopost",
			})
	return s
}()

var StatusNoInfo = func() *status.Status {
	s, _ := status.New(codes.NotFound, "Info not found").
		WithDetails(
			&edpb.ErrorInfo{
				Reason: "NOT_FOUND",
				Domain: "laopost",
			})
	return s
}()

var StatusOutOfRange = func() *status.Status {
	s, _ := status.New(codes.OutOfRange, "Unprocessable Entity").
		WithDetails(
			&edpb.ErrorInfo{
				Reason: "UNPROCESSABLE_ENTITY",
				Domain: "laopost",
			})
	return s
}()

var StatusInternalServerError = func() *status.Status {
	s, _ := status.New(codes.Internal, "Internal Server Error").
		WithDetails(
			&edpb.ErrorInfo{
				Reason: "INTERNAL_SERVER_ERROR",
				Domain: "laopost",
			})
	return s
}()

var StatusAlreadyExist = func() *status.Status {
	s, _ := status.New(codes.AlreadyExists, "already exists").
		WithDetails(
			&edpb.ErrorInfo{
				Reason: "ALREADY_EXISTS",
				Domain: "laopost",
			})
	return s
}()

func GRPCStatusFromErr(err error) *status.Status {
	switch {
	case err == nil:
		return status.New(codes.OK, "OK")
	case errors.Is(err, ErrUnauthorized):
		return StatusUnauthenticated
	case errors.Is(err, ErrPermissionDenied):
		return StatusPermissionDenied
	case errors.Is(err, ErrNoInfo):
		return StatusNoInfo
	case errors.Is(err, ErrUnProcessAbleEntity):
		return StatusOutOfRange
	case errors.Is(err, ErrInternalServerError):
		return StatusInternalServerError
	case errors.Is(err, ErrAlreadyExist):
		return StatusAlreadyExist
	}

	return StatusInternalServerError
}

func HttpStatusPbFromRPC(s *status.Status) *hspb.Error {
	return &hspb.Error{
		Error: &hspb.Error_Status{
			Code:    int32(httpStatusFromCode(s.Code())),
			Status:  code.Code(s.Code()),
			Message: s.Message(),
			Details: s.Proto().Details,
		},
	}
}

// httpStatusFromCode converts a gRPC error code into the corresponding HTTP response status.
// See: https://github.com/googleapis/googleapis/blob/master/google/rpc/code.proto
func httpStatusFromCode(code codes.Code) int {
	switch code {
	case codes.OK:
		return http.StatusOK
	case codes.Canceled:
		return http.StatusRequestTimeout
	case codes.Unknown:
		return http.StatusInternalServerError
	case codes.InvalidArgument:
		return http.StatusBadRequest
	case codes.DeadlineExceeded:
		return http.StatusGatewayTimeout
	case codes.NotFound:
		return http.StatusNotFound
	case codes.AlreadyExists:
		return http.StatusConflict
	case codes.PermissionDenied:
		return http.StatusForbidden
	case codes.Unauthenticated:
		return http.StatusUnauthorized
	case codes.ResourceExhausted:
		return http.StatusTooManyRequests
	case codes.FailedPrecondition:
		// Note, this deliberately doesn't translate to the similarly named '412 Precondition Failed' HTTP response status.
		return http.StatusBadRequest
	case codes.Aborted:
		return http.StatusConflict
	case codes.OutOfRange:
		return http.StatusBadRequest
	case codes.Unimplemented:
		return http.StatusNotImplemented
	case codes.Internal:
		return http.StatusInternalServerError
	case codes.Unavailable:
		return http.StatusServiceUnavailable
	case codes.DataLoss:
		return http.StatusInternalServerError
	}

	return http.StatusInternalServerError
}
