package lwc

import (
	"reflect"
	"testing"
	"time"
)

type configTest struct {
	args     []string
	expected Config
}

var configTests = []configTest{
	{
		[]string{},
		Config{
			true, true, false, true, false,
			time.Duration(DEFAULT_INTERVAL) * time.Millisecond,
			false, false,
			[]string{},
			nil,
		},
	},
	{
		[]string{"-w", "--lines"},
		Config{
			true, true, false, false, false,
			time.Duration(DEFAULT_INTERVAL) * time.Millisecond,
			false, false,
			[]string{},
			nil,
		},
	},
	{
		[]string{"foo"},
		Config{
			true, true, false, true, false,
			time.Duration(DEFAULT_INTERVAL) * time.Millisecond,
			false, false,
			[]string{"foo"},
			nil,
		},
	},
	{
		[]string{"--", "/path/to/file"},
		Config{
			true, true, false, true, false,
			time.Duration(DEFAULT_INTERVAL) * time.Millisecond,
			false, false,
			[]string{"/path/to/file"},
			nil,
		},
	},
	{
		[]string{"--max-line-length", "--bytes", "/etc/passwd", "/etc/group"},
		Config{
			false, false, false, true, true,
			time.Duration(DEFAULT_INTERVAL) * time.Millisecond,
			false, false,
			[]string{"/etc/passwd", "/etc/group"},
			nil,
		},
	},
	{
		[]string{"-i", "5000"},
		Config{
			true, true, false, true, false,
			time.Duration(5000) * time.Millisecond,
			false, false,
			[]string{},
			nil,
		},
	},
	{
		[]string{"--interval=2000"},
		Config{
			true, true, false, true, false,
			time.Duration(2000) * time.Millisecond,
			false, false,
			[]string{},
			nil,
		},
	},
	{
		[]string{"--interval", "3000"},
		Config{
			true, true, false, true, false,
			time.Duration(3000) * time.Millisecond,
			false, false,
			[]string{},
			nil,
		},
	},
	{
		[]string{"-i", "0"},
		Config{
			true, true, false, true, false,
			time.Duration(0),
			false, false,
			[]string{},
			nil,
		},
	},
}

func TestBuildConfig(t *testing.T) {
	for i, test := range configTests {
		actual := BuildConfig(append([]string{"lwc"}, test.args...))
		// Clear getopt Set because we don't want to compare it
		actual.g = nil
		if !reflect.DeepEqual(test.expected, actual) {
			t.Errorf("Test #%d failed: expecting %#v, got %#v", i, test.expected, actual)
		}
	}
}
