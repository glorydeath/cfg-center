package cfgLoader

import (
	//"encoding/json"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestCfgManager_GetCfg_json(t *testing.T) {
	cfg_mgr := New("./../../../test")
	cfg_mgr.LoadCfg()
	cfg_mgr.ReloadGitCfg()
	Convey("test yaml load", t, func() {
		jsout, _ := cfg_mgr.cfg_data.Encode()
		So(string(jsout), ShouldContainSubstring, "jdbc:mariadb://localhost:3306/TM")
		So(string(cfg_mgr.GetCfg_json([]string{""}, Raw)), ShouldEqual, "null")
		Convey("test var and template both yaml", func(){
			So(string(cfg_mgr.GetCfg_json([]string{"test", "dummyWorker", "slotStatement"}, Raw)),
				ShouldEqual, "\"select * from A\"")
		})
	})
}

func TestCfgManager_GetCfg_json_flat(t *testing.T) {
	cfg_mgr := New("./../../../test")
	cfg_mgr.LoadCfg()
	cfg_mgr.ReloadGitCfg()
	Convey("test yaml load", t, func() {
		jsout, _ := cfg_mgr.cfg_data.Encode()
		So(string(jsout), ShouldContainSubstring, "jdbc:mariadb://localhost:3306/TM")
		So(string(cfg_mgr.GetCfg_json([]string{""}, Flat)), ShouldEqual, "null")
		Convey("test var and template both yaml flat", func(){
			So(string(cfg_mgr.GetCfg_json([]string{"test"}, Flat)),
				ShouldContainSubstring, "dummyWorker.slotStatement")
		})
	})
}

func TestCfgManager_GetCfg_json_path_error(t *testing.T) {
	cfg_mgr := New("./../../../testfdsaf")
	cfg_mgr.LoadCfg()
	cfg_mgr.ReloadGitCfg()
	Convey("test yaml load not exist", t, func() {
		jsout, _ := cfg_mgr.cfg_data.Encode()
		So(string(jsout), ShouldNotContainSubstring, "jdbc:mariadb://localhost:3306/TM")
		So(string(cfg_mgr.GetCfg_json([]string{""}, Raw)), ShouldEqual, "null")
		So(string(cfg_mgr.GetCfg_json([]string{"test", "dummyWorker", "slotStatement"}, Raw)),
			ShouldEqual, "null")
	})
}

func TestCfgManager_GetCfg_json_path_error2(t *testing.T) {

	cfg_mgr := New("./../../../test/aws.yaml")
	Convey("test yaml load error", t, func() {
		defer func() {
			if r := recover(); r != nil {
				t.Log("Config path error, should be dir")
			}
		}()
		cfg_mgr.LoadCfg()
		jsout, _ := cfg_mgr.cfg_data.Encode()
		So(string(jsout), ShouldEqual, "{}")
		So(string(cfg_mgr.GetCfg_json([]string{""}, Raw)), ShouldEqual, "null")
		So(string(cfg_mgr.GetCfg_json([]string{"test", "dummyWorker", "slotStatement"}, Raw)),
			ShouldEqual, "null")
	})
}
