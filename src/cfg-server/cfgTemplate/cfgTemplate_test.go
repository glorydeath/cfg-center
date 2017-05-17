package cfgTemplate

import (
	"encoding/json"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestCfgTemplate(t *testing.T) {
	//Count: 111
	//Material: wooooood
	//Person:
	//name: Jack
	cfg_t := New()
	cfg_t.LoadVars("./../../../test/GLOBAL_VAR.yaml")

	Convey("load var.yaml", t, func() {
		num, _ := cfg_t.vars["Count"].(json.Number)
		num_int, _ := num.Int64()
		So(num_int, ShouldEqual, int64(111))
		So(cfg_t.vars["Material"], ShouldEqual, "wooooood")
		So(cfg_t.vars["Person"].(map[string]interface{})["name"], ShouldEqual, "Jack")
	})
	Convey("translate normal", t, func() {
		tpl := []byte("{{.Person.name}} got {{.Count}} items are made of {{.Material}}")
		out := cfg_t.Translate(tpl)
		So(string(out), ShouldEqual, "Jack got 111 items are made of wooooood")
	})
	Convey("translate non template", t, func() {
		tpl := []byte("fdsafdsafdsafd{{")
		out := cfg_t.Translate(tpl)
		So(string(out), ShouldEqual, string(tpl))
	})
	Convey("translate key error", t, func() {
		tpl := []byte("{{Person.name}} got {{.Count}} items are made of {{.Material}}")
		out := cfg_t.Translate(tpl)
		So(string(out), ShouldEqual, string(tpl))
	})
}
