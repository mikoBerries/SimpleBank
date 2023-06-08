package gapi

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GRPCLogger this is option logger for gRPC only catching unary request
func GRPCLogger(ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {

	startTime := time.Now()
	//pass original control flow request
	result, err := handler(ctx, req)
	//done requesting and get callback
	duration := time.Since(startTime)

	statusCode := codes.Unknown
	if st, ok := status.FromError(err); ok {
		statusCode = st.Code()
	}

	//init logger as log-info type
	logger := log.Info()
	if err != nil { //change logger to log-err type and holding err keymap value
		logger = log.Error().Err(err)
	}
	code := fmt.Sprintf(`%v %s`, int(statusCode), statusCode.String())
	//adding standart logging info we desire
	logger.
		Str("protocol", "gRPC-unary").
		Str("method", info.FullMethod).
		Str("status_code", code).
		Dur("duration", duration).
		Msg("received a gRPC request")

	return result, err
}

// overtake data needed from response writer to use
type ResponseRecorder struct {
	http.ResponseWriter
	StatusCode int
	Body       []byte
}

func (rec *ResponseRecorder) WriteHeader(statusCode int) {
	rec.StatusCode = statusCode
	rec.ResponseWriter.WriteHeader(statusCode)
}

func (rec *ResponseRecorder) Write(body []byte) (int, error) {
	rec.Body = body
	return rec.ResponseWriter.Write(body)
}
