package apperrors

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"
)

// TODO: add sentinel errors
var (
	ErrInternalServer = &sentinelError{Message: "internal server error", StatusCode: http.StatusInternalServerError}
	ErrNotFound       = &sentinelError{Message: "resource not found", StatusCode: http.StatusNotFound}
	ErrUnauthorized   = &sentinelError{Message: "unauthorized error", StatusCode: http.StatusUnauthorized}
	ErrConflict       = &sentinelError{Message: "uncompleted request due to conflict", StatusCode: http.StatusConflict}
)

type IApiError interface {
	ApiError() (statusCode, ErrorCode int, message string)
}

type Trace struct {
	File     string `json:"file"`
	Line     int    `json:"line"`
	Function string `json:"function"`
}

type Traces []Trace

func (ts Traces) PrintTraces() string {
	var sb strings.Builder
	for _, t := range ts {
		sb.WriteString(fmt.Sprintf("%s\n\t%s:%d\n", t.Function, t.File, t.Line))
	}

	return sb.String()
}

const mainMain string = "main.main"

func getStackTrace() Traces {
	pcs := make([]uintptr, 32)

	npcs := runtime.Callers(4, pcs)
	callers := pcs[:npcs]
	cf := runtime.CallersFrames(callers)
	traces := make([]Trace, 0, npcs)
	for {
		frame, more := cf.Next()
		traces = append(traces, Trace{
			File:     frame.File,
			Line:     frame.Line,
			Function: frame.Function,
		})

		if !more || frame.Function == mainMain {
			break
		}
	}

	return traces
}

func hostInfo() HostInfo {
	hn, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	return HostInfo{
		HostName: hn,
		Pid:      os.Getpid(),
		Stack:    getStackTrace(),
	}
}

type HostInfo struct {
	HostName string `json:"host_name"`
	Pid      int    `json:"pid"`
	Stack    Traces `json:"stack"`
}

type WrappedError struct {
	Err         error          `json:"err"`
	HostInfo    HostInfo       `json:"host_info"`
	SentinelErr *sentinelError `json:"error_data"`
}

type sentinelError struct {
	Message    string `json:"message"`
	ErrorCode  int    `json:"error_code"`
	StatusCode int    `json:"status_code"`
}

func (ae sentinelError) ApiError() (int, int, string) {
	return ae.StatusCode, ae.ErrorCode, ae.Message
}

func (we WrappedError) Error() string {
	return we.Err.Error()
}

func (we WrappedError) PrintStack(w io.Writer) {
	io.WriteString(w, we.HostInfo.Stack.PrintTraces())
}

func (we WrappedError) ApiError() (int, int, string) {
	return we.SentinelErr.ApiError()
}

func WrapError(err error, sentinelErr *sentinelError) WrappedError {
	if unWrapped, ok := err.(WrappedError); ok {
		return unWrapped
	}

	return WrappedError{
		Err:         err,
		HostInfo:    hostInfo(),
		SentinelErr: sentinelErr,
		// API error assign with sentinel error
	}
}

// func (e apiError) ApiError() (statusCode, apiErrorCode int, desc string) {
// 	return e.statusCode, e.apiErrorCode, e.description
// }

// type apiError struct {
// 	statusCode   int
// 	apiErrorCode int
// 	description  string
// }

// func (e apiError) ApiError() (statusCode, apiErrorCode int, desc string) {
// 	return e.statusCode, e.apiErrorCode, e.description
// }

// func (e apiError) Error() string {
// 	return e.description
// }

// type apiErrorWrapError struct {
// 	err    error
// 	apiErr *apiError
// }

// func (e apiErrorWrapError) ApiError() (statusCode, apiErrorCode int, desc string) {
// 	return e.apiErr.ApiError()
// }

// func (e apiErrorWrapError) Error() string {
// 	return e.err.Error()
// }

// func (e apiErrorWrapError) Unwrap() error {
// 	return e.err
// }

// func WrapError(err error, apiErr *apiError) error {
// 	return &apiErrorWrapError{err: err, apiErr: apiErr}
// }
