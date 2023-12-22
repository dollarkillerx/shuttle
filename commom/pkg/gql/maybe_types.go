package gql

import (
	"encoding/json"
	"fmt"
	"io"
	"math"

	"github.com/99designs/gqlgen/graphql"
	proto "github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/pkg/errors"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// MarshalMaybeString ...
func MarshalMaybeString(s wrappers.StringValue) graphql.Marshaler {
	return graphql.MarshalString(s.Value)
}

// UnmarshalMaybeString ...
func UnmarshalMaybeString(v interface{}) (wrappers.StringValue, error) {
	switch v := v.(type) {
	case string:
		return wrappers.StringValue{Value: v}, nil
	default:
		return wrappers.StringValue{Value: ""}, errors.Errorf("%T is not a StringValue", v)
	}
}

// MarshalMaybeBool ...
func MarshalMaybeBool(s wrappers.BoolValue) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		io.WriteString(w, fmt.Sprintf("%v", s.Value))
	})
}

// UnmarshalMaybeBool ...
func UnmarshalMaybeBool(v interface{}) (wrappers.BoolValue, error) {
	switch v.(type) {
	case bool:
		return wrappers.BoolValue{Value: v.(bool)}, nil
	default:
		return wrappers.BoolValue{Value: false}, errors.Errorf("%T is not a BoolValue", v)
	}
}

// MarshalMaybeBytes ...
func MarshalMaybeBytes(s wrappers.BytesValue) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		// TODO: implement the marshal functionality
	})
}

// UnmarshalMaybeBytes ...
func UnmarshalMaybeBytes(v interface{}) (wrappers.BytesValue, error) {
	switch v.(type) {
	// TODO: implement the unmarshal functionality
	default:
		return wrappers.BytesValue{Value: []byte{}}, errors.Errorf("%T is not a BytesValue", v)
	}
}

// MarshalMaybeFloat64 ...
func MarshalMaybeFloat64(d wrappers.DoubleValue) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		io.WriteString(w, fmt.Sprintf("%v", d.Value))
	})
}

// UnmarshalMaybeFloat64 ...
func UnmarshalMaybeFloat64(v interface{}) (wrappers.DoubleValue, error) {
	switch v := v.(type) {
	case float32:
		return wrappers.DoubleValue{Value: float64(v)}, nil
	case float64:
		return wrappers.DoubleValue{Value: v}, nil
	case int:
		return wrappers.DoubleValue{Value: float64(v)}, nil
	case int32:
		return wrappers.DoubleValue{Value: float64(v)}, nil
	case int64:
		return wrappers.DoubleValue{Value: float64(v)}, nil
	case json.Number:
		i, err := v.Float64()
		return wrappers.DoubleValue{Value: i}, err
	default:
		return wrappers.DoubleValue{Value: 0}, errors.Errorf("%T is not a DoubleValue", v)
	}
}

// MarshalMaybeFloat32 ...
func MarshalMaybeFloat32(d wrappers.FloatValue) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		io.WriteString(w, fmt.Sprintf("%v", d.Value))
	})
}

// UnmarshalMaybeFloat32 ...
func UnmarshalMaybeFloat32(v interface{}) (wrappers.FloatValue, error) {
	switch v := v.(type) {
	case float32:
		return wrappers.FloatValue{Value: v}, nil
	case float64:
		if v > math.MaxFloat32 {
			return wrappers.FloatValue{Value: 0}, errors.Errorf("%T is out of FloatValue range", v)
		}
		return wrappers.FloatValue{Value: float32(v)}, nil
	case int:
		return wrappers.FloatValue{Value: float32(v)}, nil
	case int32:
		return wrappers.FloatValue{Value: float32(v)}, nil
	case int64:
		return wrappers.FloatValue{Value: float32(v)}, nil
	case json.Number:
		i, err := v.Float64()
		if i > math.MaxFloat32 {
			return wrappers.FloatValue{Value: 0}, errors.Errorf("%T is out of FloatValue range", v)
		}
		return wrappers.FloatValue{Value: float32(i)}, err
	default:
		return wrappers.FloatValue{Value: 0}, errors.Errorf("%T is not a FloatValue", v)
	}
}

// MarshalMaybeInt64 ...
func MarshalMaybeInt64(d wrappers.Int64Value) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		io.WriteString(w, fmt.Sprintf("%v", d.Value))
	})
}

// UnmarshalMaybeInt64 ...
func UnmarshalMaybeInt64(v interface{}) (wrappers.Int64Value, error) {
	switch v := v.(type) {
	case int64:
		return wrappers.Int64Value{Value: v}, nil
	case int:
		return wrappers.Int64Value{Value: int64(v)}, nil
	case int32:
		return wrappers.Int64Value{Value: int64(v)}, nil
	case json.Number:
		i, err := v.Float64()
		if i > math.MaxInt64 {
			return wrappers.Int64Value{Value: 0}, errors.Errorf("%T is out of Int64Value range", v)
		}
		return wrappers.Int64Value{Value: int64(i)}, errors.WithStack(err)
	default:
		return wrappers.Int64Value{Value: 0}, errors.Errorf("%T is not a Int64Value", v)
	}
}

// MarshalMaybeUInt64 ...
func MarshalMaybeUInt64(d wrappers.UInt64Value) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		// TODO: implement the marshal functionality
	})
}

// UnmarshalMaybeUInt64 ...
func UnmarshalMaybeUInt64(v interface{}) (wrappers.UInt64Value, error) {
	switch v.(type) {
	// TODO: implement the unmarshal functionality
	default:
		return wrappers.UInt64Value{Value: 0}, errors.Errorf("%T is not a UInt64Value", v)
	}
}

// MarshalMaybeInt32 ...
func MarshalMaybeInt32(d wrappers.Int32Value) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		io.WriteString(w, fmt.Sprintf("%v", d.Value))
	})
}

// UnmarshalMaybeInt32 ...
func UnmarshalMaybeInt32(v interface{}) (wrappers.Int32Value, error) {
	switch v := v.(type) {
	case int:
		return wrappers.Int32Value{Value: int32(v)}, nil
	case int32:
		return wrappers.Int32Value{Value: v}, nil
	case int64:
		if v > math.MaxInt32 || v < math.MinInt32 {
			return wrappers.Int32Value{Value: 0}, errors.Errorf("%T is out of Int32Value range", v)
		}
		return wrappers.Int32Value{Value: int32(v)}, nil
	case json.Number:
		i, err := v.Int64()
		if i > math.MaxInt32 || i < math.MinInt32 {
			return wrappers.Int32Value{Value: 0}, errors.Errorf("%T is out of Int32Value range", v)
		}
		return wrappers.Int32Value{Value: int32(i)}, err
	default:
		return wrappers.Int32Value{Value: 0}, errors.Errorf("%T is not a Int32Value", v)
	}
}

// MarshalMaybeUInt32 ...
func MarshalMaybeUInt32(d wrappers.UInt32Value) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		// TODO: implement the marshal functionality
	})
}

// UnmarshalMaybeUInt32 ...
func UnmarshalMaybeUInt32(v interface{}) (wrappers.UInt32Value, error) {
	switch v.(type) {
	// TODO: implement the unmarshal functionality
	default:
		return wrappers.UInt32Value{Value: 0}, errors.Errorf("%T is not a UInt32Value", v)
	}
}
