package envloader

import (
	"fmt"
	"math"
	"os"
	"testing"
)

type Config struct {
	String  string     `env:"STRING"`
	Int     int        `env:"INT"`
	Int8    int8       `env:"INT8"`
	Int16   int16      `env:"INT16"`
	Int32   int32      `env:"INT32"`
	Int64   int64      `env:"INT64"`
	Uint    uint       `env:"UINT"`
	Uint8   uint8      `env:"UINT8"`
	Uint16  uint16     `env:"UINT16"`
	Uint32  uint32     `env:"UINT32"`
	Uint64  uint64     `env:"UINT64"`
	Float32 float32    `env:"FLOAT32"`
	Float64 float64    `env:"FLOAT64"`
	Bool    bool       `env:"BOOL,required"`
	Complex complex128 `env:"COMPLEX"`
}

func TestBindStruct(t *testing.T) {

	testCases := []struct {
		name    string
		env     map[string]string
		wantErr bool
		want    Config
	}{
		{
			name: "AllTypes",
			env: map[string]string{
				"STRING":  "test",
				"INT":     "1",
				"INT8":    "2",
				"INT16":   "3",
				"INT32":   "4",
				"INT64":   "5",
				"UINT":    "6",
				"UINT8":   "7",
				"UINT16":  "8",
				"UINT32":  "9",
				"UINT64":  "10",
				"FLOAT32": "11.1",
				"FLOAT64": "12.2",
				"BOOL":    "true",
				"COMPLEX": "1+2i",
			},
			wantErr: false,
			want: Config{
				String:  "test",
				Int:     1,
				Int8:    2,
				Int16:   3,
				Int32:   4,
				Int64:   5,
				Uint:    6,
				Uint8:   7,
				Uint16:  8,
				Uint32:  9,
				Uint64:  10,
				Float32: 11.1,
				Float64: 12.2,
				Bool:    true,
				Complex: complex(1, 2),
			},
		},
		{
			name: "BoolTrueValues",
			env: map[string]string{
				"BOOL": "Y",
			},
			wantErr: false,
			want: Config{
				Bool: true,
			},
		},
		{
			name: "BoolFalseValues",
			env: map[string]string{
				"BOOL": "n",
			},
			wantErr: false,
			want: Config{
				Bool: false,
			},
		},
		{
			name:    "TestRequired",
			env:     map[string]string{},
			wantErr: true,
			want:    Config{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Set environment variables using data provider mock
			for k, v := range tc.env {
				os.Setenv(k, v)
			}
			defer func() {
				for k := range tc.env {
					os.Unsetenv(k)
				}
			}()

			var cfg Config
			err := BindStruct(&cfg, func(key string) (string, error) {
				value, ok := tc.env[key]
				if !ok {
					return "", fmt.Errorf("environment variable %s not found", key)
				}
				return value, nil
			})

			if tc.wantErr && err == nil {
				t.Errorf("Expected an error but got nil")
			} else if !tc.wantErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			} else if !tc.wantErr {
				assertConfig(t, cfg, tc.want)
			}
		})
	}
}

func assertConfig(t *testing.T, got, want Config) {
	if got.String != want.String {
		t.Errorf("String: got %v, want %v", got.String, want.String)
	}
	if got.Int != want.Int {
		t.Errorf("Int: got %v, want %v", got.Int, want.Int)
	}
	if got.Int8 != want.Int8 {
		t.Errorf("Int8: got %v, want %v", got.Int8, want.Int8)
	}
	if got.Int16 != want.Int16 {
		t.Errorf("Int16: got %v, want %v", got.Int16, want.Int16)
	}
	if got.Int32 != want.Int32 {
		t.Errorf("Int32: got %v, want %v", got.Int32, want.Int32)
	}
	if got.Int64 != want.Int64 {
		t.Errorf("Int64: got %v, want %v", got.Int64, want.Int64)
	}
	if got.Uint != want.Uint {
		t.Errorf("Uint: got %v, want %v", got.Uint, want.Uint)
	}
	if got.Uint8 != want.Uint8 {
		t.Errorf("Uint8: got %v, want %v", got.Uint8, want.Uint8)
	}
	if got.Uint16 != want.Uint16 {
		t.Errorf("Uint16: got %v, want %v", got.Uint16, want.Uint16)
	}
	if got.Uint32 != want.Uint32 {
		t.Errorf("Uint32: got %v, want %v", got.Uint32, want.Uint32)
	}
	if got.Uint64 != want.Uint64 {
		t.Errorf("Uint64: got %v, want %v", got.Uint64, want.Uint64)
	}
	if !equalFloat(float64(got.Float32), float64(want.Float32)) {
		t.Errorf("Float32: got %v, want %v", got.Float32, want.Float32)
	}
	if !equalFloat(got.Float64, want.Float64) {
		t.Errorf("Float64: got %v, want %v", got.Float64, want.Float64)
	}
	if got.Bool != want.Bool {
		t.Errorf("Bool: got %v, want %v", got.Bool, want.Bool)
	}
	if !equalComplex(got.Complex, want.Complex) {
		t.Errorf("Complex: got %v, want %v", got.Complex, want.Complex)
	}
}

const epsilon = 1e-6

func equalFloat(a, b float64) bool {
	return math.Abs(a-b) < epsilon
}

func equalComplex(a, b complex128) bool {
	return equalFloat(real(a), real(b)) && equalFloat(imag(a), imag(b))
}
