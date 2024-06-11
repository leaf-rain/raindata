package encoding

import (
	"github.com/leaf-rain/raindata/common/config/encoding/dotenv"
	"github.com/leaf-rain/raindata/common/config/encoding/hcl"
	"github.com/leaf-rain/raindata/common/config/encoding/ini"
	"github.com/leaf-rain/raindata/common/config/encoding/javaproperties"
	"github.com/leaf-rain/raindata/common/config/encoding/json"
	"github.com/leaf-rain/raindata/common/config/encoding/toml"
	"github.com/leaf-rain/raindata/common/config/encoding/yaml"
)

type Encoding struct {
	keyDelim        string
	iniLoadOptions  ini.LoadOptions
	encoderRegistry *EncoderRegistry
	decoderRegistry *DecoderRegistry
}

var AllConfigType = []string{
	"yaml",
	"yml",
	"json",
	"toml",
	"hcl",
	"tfvars",
	"ini",
	"properties",
	"props",
	"prop",
	"dotenv",
	"env",
}

func NewEncoding() *Encoding {
	var encoding = &Encoding{
		keyDelim: ".",
	}
	encoderRegistry := NewEncoderRegistry()
	decoderRegistry := NewDecoderRegistry()

	{
		codec := yaml.Codec{}

		encoderRegistry.RegisterEncoder("yaml", codec)
		decoderRegistry.RegisterDecoder("yaml", codec)

		encoderRegistry.RegisterEncoder("yml", codec)
		decoderRegistry.RegisterDecoder("yml", codec)
	}

	{
		codec := json.Codec{}

		encoderRegistry.RegisterEncoder("json", codec)
		decoderRegistry.RegisterDecoder("json", codec)
	}

	{
		codec := toml.Codec{}

		encoderRegistry.RegisterEncoder("toml", codec)
		decoderRegistry.RegisterDecoder("toml", codec)
	}

	{
		codec := hcl.Codec{}

		encoderRegistry.RegisterEncoder("hcl", codec)
		decoderRegistry.RegisterDecoder("hcl", codec)

		encoderRegistry.RegisterEncoder("tfvars", codec)
		decoderRegistry.RegisterDecoder("tfvars", codec)
	}

	{
		codec := ini.Codec{
			KeyDelimiter: encoding.keyDelim,
			LoadOptions:  encoding.iniLoadOptions,
		}

		encoderRegistry.RegisterEncoder("ini", codec)
		decoderRegistry.RegisterDecoder("ini", codec)
	}

	{
		codec := &javaproperties.Codec{
			KeyDelimiter: encoding.keyDelim,
		}

		encoderRegistry.RegisterEncoder("properties", codec)
		decoderRegistry.RegisterDecoder("properties", codec)

		encoderRegistry.RegisterEncoder("props", codec)
		decoderRegistry.RegisterDecoder("props", codec)

		encoderRegistry.RegisterEncoder("prop", codec)
		decoderRegistry.RegisterDecoder("prop", codec)
	}

	{
		codec := &dotenv.Codec{}

		encoderRegistry.RegisterEncoder("dotenv", codec)
		decoderRegistry.RegisterDecoder("dotenv", codec)

		encoderRegistry.RegisterEncoder("env", codec)
		decoderRegistry.RegisterDecoder("env", codec)
	}
	encoding.encoderRegistry = encoderRegistry
	encoding.decoderRegistry = decoderRegistry
	return encoding
}

func (e *Encoding) Encode(format string, v map[string]any) ([]byte, error) {
	return e.encoderRegistry.Encode(format, v)
}

func (e *Encoding) Decode(format string, b []byte, v map[string]any) error {
	return e.decoderRegistry.Decode(format, b, v)
}

// KeyDelimiter sets the delimiter used for determining key parts.
// By default it's value is ".".
func (e *Encoding) KeyDelimiter(d string) {
	e.keyDelim = d
}

func (e *Encoding) IniLoadOptions(in ini.LoadOptions) {
	e.iniLoadOptions = in
}
