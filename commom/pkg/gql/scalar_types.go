package gql

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"strconv"

	"github.com/99designs/gqlgen/graphql"
	proto "github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/pkg/errors"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// MarshalBytes ...
func MarshalBytes(b []byte) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		_, _ = fmt.Fprintf(w, "%q", string(b))
	})
}

// UnmarshalBytes ...
func UnmarshalBytes(v interface{}) ([]byte, error) {
	switch v := v.(type) {
	case string:
		return []byte(v), nil
	case *string:
		return []byte(*v), nil
	case []byte:
		return v, nil
	case json.RawMessage:
		return []byte(v), nil
	default:
		return nil, errors.Errorf("%T is not []byte", v)
	}
}

// MarshalInt32 ...
func MarshalInt32(num int32) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		_, _ = w.Write([]byte(strconv.Itoa(int(num))))
	})
}

// UnmarshalInt32 ...
func UnmarshalInt32(v interface{}) (int32, error) {
	switch v := v.(type) {
	case int:
		return int32(v), nil
	case int32:
		return v, nil
	case int64:
		return int32(v), nil
	case json.Number:
		i, err := v.Int64()
		return int32(i), err
	default:
		return 0, errors.Errorf("%T is not int32", v)
	}
}

// MarshalInt64 ...
func MarshalInt64(num int64) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		_, _ = w.Write([]byte(strconv.Itoa(int(num))))
	})
}

// UnmarshalInt64 ...
func UnmarshalInt64(v interface{}) (int64, error) {
	switch v := v.(type) {
	case int:
		return int64(v), nil
	case int32:
		return int64(v), nil
	case int64:
		return v, nil
	case json.Number:
		i, err := v.Int64()
		return i, err
	default:
		return 0, errors.Errorf("%T is not int64", v)
	}
}

// MarshalUint32 ...
func MarshalUint32(any uint32) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		_, _ = w.Write([]byte(strconv.Itoa(int(any))))
	})
}

// UnmarshalUint32 ...
func UnmarshalUint32(v interface{}) (uint32, error) {
	switch v := v.(type) {
	case int:
		return uint32(v), nil
	case uint32:
		return v, nil
	case json.Number:
		i, err := v.Int64()
		return uint32(i), err
	default:
		return 0, errors.Errorf("%T is not int32", v)
	}
}

// MarshalUint64 ...
func MarshalUint64(any uint64) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		_, _ = w.Write([]byte(strconv.Itoa(int(any))))
	})
}

// UnmarshalUint64 ...
func UnmarshalUint64(v interface{}) (uint64, error) {
	switch v := v.(type) {
	case int:
		return uint64(v), nil
	case uint64:
		return v, nil //TODO add an unmarshal mechanism
	case json.Number:
		i, err := v.Int64()
		return uint64(i), err
	default:
		return 0, errors.Errorf("%T is not uint64", v)
	}
}

// MarshalFloat32 ...
func MarshalFloat32(any float32) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		_, _ = w.Write([]byte(strconv.FormatFloat(float64(any), 'f', 2, 32)))
	})
}

// UnmarshalFloat32 ...
func UnmarshalFloat32(v interface{}) (float32, error) {
	switch v := v.(type) {
	case int:
		return float32(v), nil
	case float32:
		return v, nil
	case json.Number:
		f, err := v.Float64()
		return float32(f), err
	default:
		return 0, errors.Errorf("%T is not float32", v)
	}
}

// MarshalFloat64 ...
func MarshalFloat64(any float64) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		_, _ = w.Write([]byte(strconv.FormatFloat(any, 'f', 2, 64)))
	})
}

// UnmarshalFloat64 ...
func UnmarshalFloat64(v interface{}) (float64, error) {
	switch v := v.(type) {
	case int:
		return float64(v), nil
	case float64:
		return v, nil
	case json.Number:
		f, err := v.Float64()
		return float64(f), err
	default:
		return 0, errors.Errorf("%T is not float64", v)
	}
}

// MarshalTimestamp ...
func MarshalTimestamp(t timestamp.Timestamp) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		_, _ = w.Write([]byte(strconv.Itoa(int(t.Seconds))))
	})
}

// UnmarshalTimestamp ...
func UnmarshalTimestamp(v interface{}) (timestamp.Timestamp, error) {
	if seconds, ok := v.(int64); ok {
		return timestamp.Timestamp{Seconds: seconds}, nil
	}
	return timestamp.Timestamp{}, errors.New("time should be a unix timestamp")
}

// MarshalStringMap ...
func MarshalStringMap(m map[string]string) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		d, _ := json.Marshal(m)
		_, _ = w.Write(d)
	})
}

// UnmarshalStringMap ...
func UnmarshalStringMap(v interface{}) (map[string]string, error) {
	if m, ok := v.(map[string]string); ok {
		return m, nil
	}
	return nil, errors.New("value should be a map[string]string")
}
