// Interface with small LCD displays
package display

import (
	actor "github.com/wasmcloud/actor-tinygo"     //nolint
	cbor "github.com/wasmcloud/tinygo-cbor"       //nolint
	msgpack "github.com/wasmcloud/tinygo-msgpack" //nolint
)

type Line struct {
	LineNumber uint8  `json:"lineNumber"`
	Text       string `json:"text"`
}

// MEncode serializes a Line using msgpack
func (o *Line) MEncode(encoder msgpack.Writer) error {
	encoder.WriteMapSize(2)
	encoder.WriteString("lineNumber")
	encoder.WriteUint8(o.LineNumber)
	encoder.WriteString("text")
	encoder.WriteString(o.Text)

	return encoder.CheckError()
}

// MDecodeLine deserializes a Line using msgpack
func MDecodeLine(d *msgpack.Decoder) (Line, error) {
	var val Line
	isNil, err := d.IsNextNil()
	if err != nil || isNil {
		return val, err
	}
	size, err := d.ReadMapSize()
	if err != nil {
		return val, err
	}
	for i := uint32(0); i < size; i++ {
		field, err := d.ReadString()
		if err != nil {
			return val, err
		}
		switch field {
		case "lineNumber":
			val.LineNumber, err = d.ReadUint8()
		case "text":
			val.Text, err = d.ReadString()
		default:
			err = d.Skip()
		}
		if err != nil {
			return val, err
		}
	}
	return val, nil
}

// CEncode serializes a Line using cbor
func (o *Line) CEncode(encoder cbor.Writer) error {
	encoder.WriteMapSize(2)
	encoder.WriteString("lineNumber")
	encoder.WriteUint8(o.LineNumber)
	encoder.WriteString("text")
	encoder.WriteString(o.Text)

	return encoder.CheckError()
}

// CDecodeLine deserializes a Line using cbor
func CDecodeLine(d *cbor.Decoder) (Line, error) {
	var val Line
	isNil, err := d.IsNextNil()
	if err != nil || isNil {
		return val, err
	}
	size, indef, err := d.ReadMapSize()
	if err != nil && indef {
		err = cbor.NewReadError("indefinite maps not supported")
	}
	if err != nil {
		return val, err
	}
	for i := uint32(0); i < size; i++ {
		field, err := d.ReadString()
		if err != nil {
			return val, err
		}
		switch field {
		case "lineNumber":
			val.LineNumber, err = d.ReadUint8()
		case "text":
			val.Text, err = d.ReadString()
		default:
			err = d.Skip()
		}
		if err != nil {
			return val, err
		}
	}
	return val, nil
}

type Display interface {
	DisplayLine(ctx *actor.Context, arg Line) (bool, error)
	Clear(ctx *actor.Context) (bool, error)
}

// DisplayHandler is called by an actor during `main` to generate a dispatch handler
// The output of this call should be passed into `actor.RegisterHandlers`
func DisplayHandler(actor_ Display) actor.Handler {
	return actor.NewHandler("Display", &DisplayReceiver{}, actor_)
}

// DisplayContractId returns the capability contract id for this interface
func DisplayContractId() string { return "wasmcloud:display" }

// DisplayReceiver receives messages defined in the Display service interface
type DisplayReceiver struct{}

func (r *DisplayReceiver) Dispatch(ctx *actor.Context, svc interface{}, message *actor.Message) (*actor.Message, error) {
	svc_, _ := svc.(Display)
	switch message.Method {

	case "DisplayLine":
		{

			d := msgpack.NewDecoder(message.Arg)
			value, err_ := MDecodeLine(&d)
			if err_ != nil {
				return nil, err_
			}

			resp, err := svc_.DisplayLine(ctx, value)
			if err != nil {
				return nil, err
			}

			var sizer msgpack.Sizer
			size_enc := &sizer
			size_enc.WriteBool(resp)
			buf := make([]byte, sizer.Len())
			encoder := msgpack.NewEncoder(buf)
			enc := &encoder
			enc.WriteBool(resp)
			return &actor.Message{Method: "Display.DisplayLine", Arg: buf}, nil
		}
	case "Clear":
		{
			resp, err := svc_.Clear(ctx)
			if err != nil {
				return nil, err
			}

			var sizer msgpack.Sizer
			size_enc := &sizer
			size_enc.WriteBool(resp)
			buf := make([]byte, sizer.Len())
			encoder := msgpack.NewEncoder(buf)
			enc := &encoder
			enc.WriteBool(resp)
			return &actor.Message{Method: "Display.Clear", Arg: buf}, nil
		}
	default:
		return nil, actor.NewRpcError("MethodNotHandled", "Display."+message.Method)
	}
}

// DisplaySender sends messages to a Display service
type DisplaySender struct{ transport actor.Transport }

// NewProvider constructs a client for sending to a Display provider
// implementing the 'wasmcloud:display' capability contract, with the "default" link
func NewProviderDisplay() *DisplaySender {
	transport := actor.ToProvider("wasmcloud:display", "default")
	return &DisplaySender{transport: transport}
}

// NewProviderDisplayLink constructs a client for sending to a Display provider
// implementing the 'wasmcloud:display' capability contract, with the specified link name
func NewProviderDisplayLink(linkName string) *DisplaySender {
	transport := actor.ToProvider("wasmcloud:display", linkName)
	return &DisplaySender{transport: transport}
}

func (s *DisplaySender) DisplayLine(ctx *actor.Context, arg Line) (bool, error) {

	var sizer msgpack.Sizer
	size_enc := &sizer
	arg.MEncode(size_enc)
	buf := make([]byte, sizer.Len())

	var encoder = msgpack.NewEncoder(buf)
	enc := &encoder
	arg.MEncode(enc)

	out_buf, _ := s.transport.Send(ctx, actor.Message{Method: "Display.DisplayLine", Arg: buf})
	d := msgpack.NewDecoder(out_buf)
	resp, err_ := d.ReadBool()
	if err_ != nil {
		return false, err_
	}
	return resp, nil
}
func (s *DisplaySender) Clear(ctx *actor.Context) (bool, error) {
	buf := make([]byte, 0)
	out_buf, _ := s.transport.Send(ctx, actor.Message{Method: "Display.Clear", Arg: buf})
	d := msgpack.NewDecoder(out_buf)
	resp, err_ := d.ReadBool()
	if err_ != nil {
		return false, err_
	}
	return resp, nil
}

// This file is generated automatically using wasmcloud/weld-codegen 0.6.0
