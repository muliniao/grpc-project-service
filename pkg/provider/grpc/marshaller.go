package grpc

import (
	"io"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/protobuf/encoding/protojson"
)

// JsonPbMarshaller jsonPb marshaller wraps the better gogo-gateway marshaller in a way it can be used as golang JsonPb marshaller.
type JsonPbMarshaller struct {
	jsonpb.Marshaler
	runtime.JSONPb
}

// NewJsonPbMarshaller create a JsonPbMarshaller
func NewJsonPbMarshaller() *JsonPbMarshaller {
	conf := NewConfigFromEnv()

	jb := runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			EmitUnpopulated: true,
			Indent:          " ",
			UseProtoNames:   true,
			UseEnumNumbers:  false,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{},
	}

	if conf.DisableEmitDefaults {
		jb.MarshalOptions.EmitUnpopulated = false
	}
	if conf.UseEnumAsInt {
		jb.MarshalOptions.UseEnumNumbers = true
	}

	return &JsonPbMarshaller{
		Marshaler: jsonpb.Marshaler{},
		JSONPb:    jb,
	}
}

// Marshal returns the JSON encoding and write io.writer
func (m *JsonPbMarshaller) Marshal(out io.Writer, pb proto.Message) error {
	bt, err := m.JSONPb.Marshal(pb)
	if err != nil {
		return err
	}
	if _, err := out.Write(bt); err != nil {
		return err
	}
	return nil
}
