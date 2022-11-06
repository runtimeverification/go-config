package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultValueInference(t *testing.T) {
	type asdf struct {
		X string
		Y int
		Z bool
	}

	c := asdf{X: "asdfasdf", Y: 54, Z: true}

	assert.NoError(t, Init(&c))

	assert.EqualValues(t, "", c.X, "Could not infer string default value.")
	assert.EqualValues(t, 0, c.Y, "Could not infer int default value.")
	assert.EqualValues(t, false, c.Z, "Could not infer bool default value.")
}

func TestEnvOnly(t *testing.T) {
	type asdf struct {
		X string `env:"X"`
		Y int    `env:"Y"`
		Z bool   `env:"Z"`
	}

	assert.NoError(t, os.Setenv("X", "FOO"))
	assert.NoError(t, os.Setenv("Y", "11"))
	assert.NoError(t, os.Setenv("Z", "True"))

	var c asdf

	assert.NoError(t, Init(&c))

	assert.EqualValues(t, "FOO", c.X)
	assert.EqualValues(t, 11, c.Y)
	assert.EqualValues(t, true, c.Z)

	assert.NoError(t, os.Unsetenv("X"))
	assert.NoError(t, os.Unsetenv("Y"))
	assert.NoError(t, os.Unsetenv("Z"))
}

func TestCliOnly(t *testing.T) {
	type asdf struct {
		X string `cli:"x"`
		Y string `cli:"yyy"`
		Z string `cli:"zzz"`
	}

	var c asdf

	os.Args = append(os.Args, "-x", "BAR")
	os.Args = append(os.Args, "--yyy", "yval")
	os.Args = append(os.Args, "--zzz=zval")

	assert.NoError(t, Init(&c))

	assert.EqualValues(t, "BAR", c.X)
	assert.EqualValues(t, "yval", c.Y)
	assert.EqualValues(t, "zval", c.Z)
}

func TestProvidedDefaultOnly(t *testing.T) {
	type asdf struct {
		X string `default:"XXX"`
		Y int    `default:"10"`
		Z bool   `default:"true"`
	}

	var c asdf

	assert.NoError(t, Init(&c))

	assert.EqualValues(t, "XXX", c.X)
	assert.EqualValues(t, 10, c.Y)
	assert.EqualValues(t, true, c.Z)
}

func TestEnvOverDefault(t *testing.T) {
	type asdf struct {
		X string `default:"XXX" env:"X"`
		Y int    `default:"10" env:"Y"`
		Z bool   `default:"true" env:"Z"`
	}

	var c asdf

	assert.NoError(t, os.Setenv("X", "EnvX"))
	assert.NoError(t, os.Setenv("Y", "11"))
	assert.NoError(t, os.Setenv("Z", "False"))

	assert.NoError(t, Init(&c))

	assert.EqualValues(t, "EnvX", c.X)
	assert.EqualValues(t, 11, c.Y)
	assert.EqualValues(t, false, c.Z)

	assert.NoError(t, os.Unsetenv("X"))
	assert.NoError(t, os.Unsetenv("Y"))
	assert.NoError(t, os.Unsetenv("Z"))
}

func TestCliOverEnv(t *testing.T) {
	type asdf struct {
		X string `env:"X" cli:"xxxx"`
		Y int    `env:"Y" cli:"yyyy"`
		Z bool   `env:"Z" cli:"zzzz"`
	}

	var c asdf

	assert.NoError(t, os.Setenv("X", "EnvX"))
	assert.NoError(t, os.Setenv("Y", "10"))
	assert.NoError(t, os.Setenv("Z", "TRUE"))

	os.Args = append(os.Args, "--xxxx=cli_xxx")
	os.Args = append(os.Args, "--yyyy", "155")
	os.Args = append(os.Args, "--zzzz=false")

	assert.NoError(t, Init(&c))

	assert.EqualValues(t, "cli_xxx", c.X)
	assert.EqualValues(t, 155, c.Y)
	assert.False(t, c.Z)

	assert.NoError(t, os.Unsetenv("X"))
	assert.NoError(t, os.Unsetenv("Y"))
	assert.NoError(t, os.Unsetenv("Z"))
}
