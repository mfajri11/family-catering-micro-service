package app

// var (
// 	allowedPath string = "/api/v1/login"
// 	httpOK      string = "200"
// )

// func writeSIDToCookieIFSuccess(ctx context.Context, w http.ResponseWriter, _ protoreflect.ProtoMessage) error {
// 	var sid string

// 	md, ok := runtime.ServerMetadataFromContext(ctx)
// 	if !ok {
// 		return errors.New("metadata not found")
// 	}
// 	// X-status-Code & X-Path were set at rpc call in handler
// 	if !(len(md.HeaderMD.Get("X-Status-Code")) > 0) &&
// 		!(md.HeaderMD.Get("X-Status-Code")[0] == httpOK) &&
// 		!(len(md.HeaderMD.Get("X-Path")) > 0) &&
// 		!(md.HeaderMD.Get("X-Path")[0] == allowedPath) {

// 		if len(md.HeaderMD.Get("X-sid")) == 0 {
// 			return errors.New("sid headers is empty")
// 		}

// 		sid = md.HeaderMD.Get("X-Sid")[0]
// 		// delete these headers if exists
// 		delete(w.Header(), "Grpc-Metadata-X-Path")
// 		delete(w.Header(), "Grpc-Metadata-X-Status-Code")
// 		delete(w.Header(), "Grpc-Metadata-X-Sid")
// 		delete(w.Header(), "Grpc-Metadata-Content-Type")
// 		delete(w.Header(), "X-Path")
// 		delete(w.Header(), "X-Status-Code")

// 		http.SetCookie(w, &http.Cookie{
// 			Name:     string(consts.SidCtxKey),
// 			Value:    sid,
// 			HttpOnly: true,
// 			Secure:   true,
// 		})

// 	}

// 	return nil
// }
