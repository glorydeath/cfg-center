package mapFlattener

import (
	"encoding/json"
	. "github.com/smartystreets/goconvey/convey"
	"reflect"
	"testing"
)

func TestCfgTemplate(t *testing.T) {
	//Count: 111
	//Material: wooooood
	//Person:
	//name: Jack
	input := `
	{
		"outterJSON":{
			"innerJSON1":{
				"value1":10,
				"value2":22,
				"InnerInnerJSONArray": [ {"fld1" : "val1"} , {"fld2" : "val2"} ],
				"InnerInnerArray": [ "test1" , "test2", [ "test3"]]
			},
			"InnerJSON2":"NoneValue"
		}
	}
	`
	output := `
	{
		"outterJSON.InnerJSON2":"NoneValue",
		"outterJSON.innerJSON1.InnerInnerArray":["test1","test2",["test3"]],
		"outterJSON.innerJSON1.InnerInnerJSONArray":[{"fld1":"val1"},{"fld2":"val2"}],
		"outterJSON.innerJSON1.value1":10,"outterJSON.innerJSON1.value2":22
	}
	`
	//Creating the maps for JSON
	in := map[string]interface{}{}
	out := map[string]interface{}{}

	//Parsing/Unmarshalling JSON encoding/json
	json.Unmarshal([]byte(input), &in)
	json.Unmarshal([]byte(output), &out)
	//fmt.Println(Flatten(in))
	jin, _ := json.Marshal(Flatten(in))
	jout, _ := json.Marshal(out)

	Convey("iter map", t, func() {
		So(string(jin), ShouldEqual, string(jout))
		So(reflect.DeepEqual(Flatten(in), out), ShouldBeTrue)
	})
}
