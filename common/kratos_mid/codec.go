package mid

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-kratos/kratos/v2/encoding"
	"github.com/go-kratos/kratos/v2/errors"
	http2 "github.com/go-kratos/kratos/v2/transport/http"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"io"
	"net/http"
	"reflect"
)

// Name is the name registered for the json codec.
const Name = "json"

var (
	// MarshalOptions is a configurable JSON format marshaller.
	MarshalOptions = protojson.MarshalOptions{
		EmitUnpopulated: true,
		UseEnumNumbers:  true,
	}
	// UnmarshalOptions is a configurable JSON format parser.
	UnmarshalOptions = protojson.UnmarshalOptions{
		DiscardUnknown: true,
	}
)

// codec is a Codec implementation with json.
type codec struct{}

func NewCodec() encoding.Codec {
	return codec{}
}

func (codec) Marshal(v interface{}) ([]byte, error) {
	switch m := v.(type) {
	case json.Marshaler:
		return m.MarshalJSON()
	case proto.Message:
		return MarshalOptions.Marshal(m)
	default:
		return json.Marshal(m)
	}
}

func (codec) Unmarshal(data []byte, v interface{}) error {
	switch m := v.(type) {
	case json.Unmarshaler:
		return m.UnmarshalJSON(data)
	case proto.Message:
		return UnmarshalOptions.Unmarshal(data, m)
	default:
		rv := reflect.ValueOf(v)
		for rv := rv; rv.Kind() == reflect.Ptr; {
			if rv.IsNil() {
				rv.Set(reflect.New(rv.Type().Elem()))
			}
			rv = rv.Elem()
		}
		if m, ok := reflect.Indirect(rv).Interface().(proto.Message); ok {
			return UnmarshalOptions.Unmarshal(data, m)
		}
		return json.Unmarshal(data, m)
	}
}

func (codec) Name() string {
	return Name
}

type ResponseEncoder struct {
	Status int         `json:"status" form:"status" `
	Msg    string      `json:"msg" form:"msg" `
	Data   interface{} `json:"data" form:"data" `
}

func NewResponseEncoder() *ResponseEncoder {
	return &ResponseEncoder{
		Status: 0,
		Msg:    "",
	}
}

func (e *ResponseEncoder) Encoder(w http.ResponseWriter, r *http.Request, v interface{}) error {
	data, err := json.Marshal(e)
	if err != nil {
		return err
	}
	var buff []byte
	if v != nil {
		switch m := v.(type) {
		//case proto.Message:
		//	buff, err = MarshalOptions.Marshal(m)
		default:
			buff, err = json.Marshal(m)
		}
		if err != nil {
			return err
		}
		data = append(data[:len(data)-5], buff...)
		data = append(data, []byte("}")...)
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	return err
}

// HTTPError is a http error.
type HTTPError struct {
	Status int    `json:"status" form:"status" `
	Msg    string `json:"msg" form:"msg" `
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("status: %d, msg: %s", e.Status, e.Msg)
}

// FromError converts error to HTTPError.
func FromError(err error) *HTTPError {
	if err == nil {
		return nil
	}
	if e, ok := err.(*HTTPError); ok {
		return e
	}
	return &HTTPError{Status: 99, Msg: err.Error()}
}

func ErrorEncoder(w http.ResponseWriter, r *http.Request, err error) {
	se := FromError(err)
	codec, _ := http2.CodecForRequest(r, "Accept")
	body, err := codec.Marshal(se)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/"+codec.Name())
	_, _ = w.Write(body)
}

type RequestDecoder struct {
}

func NewRequestDecoder() *RequestDecoder {
	return &RequestDecoder{}
}

func (d *RequestDecoder) Decoder(r *http.Request, v interface{}) error {
	codec, _ := http2.CodecForRequest(r, "Content-Type")
	data, err := io.ReadAll(r.Body)
	// reset body.
	r.Body = io.NopCloser(bytes.NewBuffer(data))

	if err != nil {
		return errors.BadRequest("CODEC", err.Error())
	}
	if len(data) == 0 {
		return nil
	}
	if err = codec.Unmarshal(data, v); err != nil {
		return errors.BadRequest("CODEC", fmt.Sprintf("body unmarshal %s", err.Error()))
	}
	return nil
}
