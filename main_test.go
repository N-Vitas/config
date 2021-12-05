package config

import (
	"os"
	"testing"
)

func TestUseFile(t *testing.T) {
	os.Setenv("go111module", "on")
	con := NewWithOptions(
		SetConfigFile("./main_test.json"),
		UseConfigFile(),
	)

	if !con.GetBool("auth.ismart", false) {
		t.Fatal("ismart not is true")
	}
	if con.GetInt64("db.test.port", 3036) != 1443 {
		t.Fatal("port not is 1443")
	}
	if len(con.GetSlice("auth.maps", []string{})) == 0 {
		t.Fatal("maps is empty")
	}
	if con.GetString("go111module", "off") != "on" {
		t.Fatal("wronk read env")
	}
}

func TestUseSystem(t *testing.T) {
	os.Setenv("go111module", "on")
	os.Setenv("ismart", "true")
	os.Setenv("port", "1443")
	os.Setenv("maps", "guest, user, moderator, admin")

	con := NewWithOptions(UseEventSystem())

	if !con.GetBool("ismart", false) {
		t.Fatal("ismart not is true")
	}
	if con.GetInt64("port", 3036) != 1443 {
		t.Fatal("port not is 1443")
	}
	if len(con.GetSlice("maps", []string{})) == 0 {
		t.Fatal("maps is empty")
	}
	if con.GetString("go111module", "off") != "on" {
		t.Fatal("wronk read env")
	}
}
