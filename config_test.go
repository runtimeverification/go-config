package config

import (
	"os"
	"testing"
)

func testInit(c interface{}, t *testing.T) {
	err := Init(c)

	if err != nil {
		t.Fatalf("Init() failed -- %v", err)
	}
}

func TestDefaultValueInference(t *testing.T) {
	type asdf struct {
		X string
		Y int
		Z bool
	}

	var c asdf

	testInit(&c, t)

	if c.X != "" {
		t.Fatalf("Could not infer default value of empty string")
	}

	if c.Y != 0 {
		t.Fatalf("Could not infer default value of 0 for int field")
	}

	if c.Z != false {
		t.Fatalf("Could not infer default value of FALSE for bool field")
	}
}

func TestEnvOnly(t *testing.T) {
	type asdf struct {
		X string `env:"X"`
	}

	os.Setenv("X", "FOO")

	var c asdf

	testInit(&c, t)

	if c.X != "FOO" {
		t.Fatalf("Expected \"%s\", found \"%s\"", "FOO", c.X)
	}

	os.Unsetenv("X")
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

	testInit(&c, t)

	if c.X != "BAR" {
		t.Fatalf("Expected \"%s\", found \"%s\"", "BAR", c.X)
	}

	if c.Y != "yval" {
		t.Fatalf("Expected \"%s\", found \"%s\"", "yval", c.X)
	}

	if c.Z != "zval" {
		t.Fatalf("Expected \"%s\", found \"%s\"", "zval", c.X)
	}
}

func TestProvidedDefaultOnly(t *testing.T) {
	type asdf struct {
		X string `default:"XXX"`
		Y int    `default:"10"`
		Z bool   `default:"true"`
	}

	var c asdf

	testInit(&c, t)

	if c.X != "XXX" {
		t.Fatalf("Failed provided string default \"XXX\", found %s", c.X)
	}

	if c.Y != 10 {
		t.Fatalf("Failed provided int default 10, found %d", c.Y)
	}

	if c.Z != true {
		t.Fatalf("Failed provided bool default TRUE, found %v", c.Z)
	}
}
