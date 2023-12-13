package gateway

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/textproto"
	"strings"

	"github.com/golang/protobuf/ptypes"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/pkg/errors"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/status"
)

// MuxWrapper A wrapper around the GRPC Gateway ServeMux.
// It removes the basePath from requests and then forwards them to the ServeMux.
type MuxWrapper struct {
	mux            *runtime.ServeMux
	basePath       string
	specialUrlPath []string
}

// NewMuxWrapper creates a Grpc Gateway ServeMux
func NewMuxWrapper(basePath string, specialUrlPath []string, mux *runtime.ServeMux) *MuxWrapper {
	return &MuxWrapper{
		mux:            mux,
		basePath:       basePath,
		specialUrlPath: specialUrlPath,
	}
}

// ServeHTTP spawns an HTTP server
func (m *MuxWrapper) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	var hasSpecial bool
	if len(m.specialUrlPath) > 0 {
		for i := range m.specialUrlPath {
			if m.specialUrlPath[i] != req.URL.Path {
				continue
			}
			hasSpecial = true
			break
		}
	}

	if !hasSpecial {
		req.URL.Path = "/" + strings.TrimPrefix(req.URL.Path, m.basePath)
	}

	m.mux.ServeHTTP(res, req)
}

// Error format of response
type Error struct {
	Code    int32         `json:"code"`
	Message string        `json:"message"`
	Details []interface{} `json:"details"`
}

// HTTPError gateway error handing middleware
func HTTPError(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
	const fallback = `{"error": "failed to marshal error message"}`

	var customStatus *runtime.HTTPStatusError
	if errors.As(err, &customStatus) {
		err = customStatus.Err
	}

	s := status.Convert(err)
	pb := s.Proto()

	w.Header().Del("Trailer")
	w.Header().Del("Transfer-Encoding")

	contentType := marshaler.ContentType(pb)
	w.Header().Set("Content-Type", contentType)

	if s.Code() == codes.Unauthenticated {
		w.Header().Set("WWW-Authenticate", s.Message())
	}

	body := &Error{
		Message: s.Message(),
		Code:    int32(s.Code()),
	}
	details := make([]interface{}, 0)
	for _, detail := range s.Proto().GetDetails() {
		switch detail.GetTypeUrl() {
		case "type.googleapis.com/google.rpc.ErrorInfo":
			d := new(errdetails.ErrorInfo)
			err = ptypes.UnmarshalAny(detail, d)
			if err == nil {
				details = append(details, struct {
					*errdetails.ErrorInfo
					Type string `json:"@type"`
				}{d, detail.GetTypeUrl()})
			}
		case "type.googleapis.com/google.rpc.BadRequest":
			d := new(errdetails.BadRequest)
			err = ptypes.UnmarshalAny(detail, d)
			if err == nil {
				details = append(details, struct {
					*errdetails.BadRequest
					Type string `json:"@type"`
				}{d, detail.GetTypeUrl()})
			}
		case "type.googleapis.com/google.rpc.RequestInfo":
			d := new(errdetails.RequestInfo)
			err = ptypes.UnmarshalAny(detail, d)
			if err == nil {
				details = append(details, struct {
					*errdetails.RequestInfo
					Type string `json:"@type"`
				}{d, detail.GetTypeUrl()})
			}
		case "type.googleapis.com/google.rpc.QuotaFailure":
			d := new(errdetails.QuotaFailure)
			err = ptypes.UnmarshalAny(detail, d)
			if err == nil {
				details = append(details, struct {
					*errdetails.QuotaFailure
					Type string `json:"@type"`
				}{d, detail.GetTypeUrl()})
			}
		}
	}
	body.Details = details
	buf, merr := marshaler.Marshal(body)
	if merr != nil {
		grpclog.Infof("Failed to marshal error message %q: %v", body, merr)
		w.WriteHeader(http.StatusInternalServerError)
		if _, err := io.WriteString(w, fallback); err != nil {
			grpclog.Infof("Failed to write response: %v", err)
		}
		return
	}

	md, ok := runtime.ServerMetadataFromContext(ctx)
	if !ok {
		grpclog.Infof("Failed to extract ServerMetadata from context")
	}

	handleForwardResponseServerMetadata(w, mux, md)

	// RFC 7230 https://tools.ietf.org/html/rfc7230#section-4.1.2
	// Unless the request includes a TE header field indicating "trailers"
	// is acceptable, as described in Section 4.3, a server SHOULD NOT
	// generate trailer fields that it believes are necessary for the user
	// agent to receive.
	var wantsTrailers bool

	if te := r.Header.Get("TE"); strings.Contains(strings.ToLower(te), "trailers") {
		wantsTrailers = true
		handleForwardResponseTrailerHeader(w, md)
		w.Header().Set("Transfer-Encoding", "chunked")
	}

	st := runtime.HTTPStatusFromCode(s.Code())
	w.WriteHeader(st)
	if _, err := w.Write(buf); err != nil {
		grpclog.Infof("Failed to write response: %v", err)
	}

	if wantsTrailers {
		handleForwardResponseTrailer(w, md)
	}
}

func handleForwardResponseServerMetadata(w http.ResponseWriter, mux *runtime.ServeMux, md runtime.ServerMetadata) {
	for k, vs := range md.HeaderMD {
		if strings.HasPrefix(k, "Grpc-Metadata-") {
			for _, v := range vs {
				w.Header().Add(strings.TrimLeft(k, "Grpc-Metadata-"), v)
			}
		}
	}
}

func handleForwardResponseTrailerHeader(w http.ResponseWriter, md runtime.ServerMetadata) {
	for k := range md.TrailerMD {
		tKey := textproto.CanonicalMIMEHeaderKey(fmt.Sprintf("%s%s", runtime.MetadataTrailerPrefix, k))
		w.Header().Add("Trailer", tKey)
	}
}

func handleForwardResponseTrailer(w http.ResponseWriter, md runtime.ServerMetadata) {
	for k, vs := range md.TrailerMD {
		tKey := fmt.Sprintf("%s%s", runtime.MetadataTrailerPrefix, k)
		for _, v := range vs {
			w.Header().Add(tKey, v)
		}
	}
}
