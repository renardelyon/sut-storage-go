package helper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ContextError(c context.Context) error {
	switch c.Err() {
	case context.Canceled:
		return status.Error(codes.Canceled, "Request Canceled")
	case context.DeadlineExceeded:
		return status.Error(codes.DeadlineExceeded, "Deadline Exceeded")
	default:
		return nil
	}
}
