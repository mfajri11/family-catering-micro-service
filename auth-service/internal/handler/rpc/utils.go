package rpc

import (
	"context"
	"errors"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/mfajri11/family-catering-micro-service/auth-service/pkg/consts"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var (
	allowedPath string = "/api/v1/login"
	httpOK      string = "200"
)

type headerOptions map[string]string

func setHeader(ctx context.Context, opt headerOptions) error {
	var (
		md metadata.MD
		ok bool
	)
	md, ok = metadata.FromOutgoingContext(ctx)
	if !ok {
		md = metadata.New(opt)
	}

	for k, v := range opt {
		md.Append(k, v)
	}

	return grpc.SetHeader(ctx, md)

}

func WriteSIDToCookieIFSuccess(ctx context.Context, w http.ResponseWriter, p protoreflect.ProtoMessage) error {
	return writeSIDToCookieIFSuccess(ctx, w, p)
}

func writeSIDToCookieIFSuccess(ctx context.Context, w http.ResponseWriter, _ protoreflect.ProtoMessage) error {
	var sid string

	md, ok := runtime.ServerMetadataFromContext(ctx)
	if !ok {
		return errors.New("metadata not found")
	}
	// X-status-Code & X-Path were set at rpc call in handler
	if !(len(md.HeaderMD.Get("X-Status-Code")) > 0) &&
		!(md.HeaderMD.Get("X-Status-Code")[0] == httpOK) &&
		!(len(md.HeaderMD.Get("X-Path")) > 0) &&
		!(md.HeaderMD.Get("X-Path")[0] == allowedPath) {

		if len(md.HeaderMD.Get("X-sid")) == 0 {
			return errors.New("sid headers is empty")
		}

		sid = md.HeaderMD.Get("X-Sid")[0]
		// delete these headers if exists
		delete(w.Header(), "Grpc-Metadata-X-Path")
		delete(w.Header(), "Grpc-Metadata-X-Status-Code")
		delete(w.Header(), "Grpc-Metadata-X-Sid")
		delete(w.Header(), "Grpc-Metadata-Content-Type")
		delete(w.Header(), "X-Path")
		delete(w.Header(), "X-Status-Code")

		http.SetCookie(w, &http.Cookie{
			Name:     string(consts.SidCtxKey),
			Value:    sid,
			HttpOnly: true,
			Secure:   true,
		})

	}

	return nil
}
