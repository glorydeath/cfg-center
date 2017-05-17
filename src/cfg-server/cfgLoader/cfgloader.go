package cfgLoader

import (
	"sync"
	//"gopkg.in/yaml.v2"
	log "github.com/auxten/logrus"
	ym "github.com/ghodss/yaml"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	//"errors"
	"encoding/json"
	"fmt"
	"github.com/4paradigm/cfg-center/src/cfg-server/cfgTemplate"
	"github.com/4paradigm/cfg-center/src/cfg-server/mapFlattener"
	"github.com/bitly/go-simplejson"
	"github.com/coreos/etcd/pkg/fileutil"
	"os/exec"
)

type cfgData map[string]interface{}

const (
	CONF_EXT   = ".yaml"
	GLOBAL_VAR = "GLOBAL_VAR.yaml"
)

type CfgManager struct {
	cfg_file_path string
	cfg_data      *simplejson.Json
	cfg_data_lock sync.RWMutex
}

func (cfgm *CfgManager) ReloadGitCfg() string {
	cmd := fmt.Sprintf("cd %s && git pull --recurse-submodules=yes origin master && git rev-parse HEAD", cfgm.cfg_file_path)
	//cmd := fmt.Sprintf("sh -c \"pwd\"")
	output, err := exec.Command("sh", "-c", cmd).CombinedOutput()
	outstr := string(output)
	if err != nil {
		log.Error(err)
		outstr = fmt.Sprintf("%s %s", outstr, err.Error())
	}
	cfgm.LoadCfg()
	return outstr
}

func New(cfgDataPath string) *CfgManager {
	return &CfgManager{
		cfg_file_path: cfgDataPath,
	}
}

type Fmt int

const (
	Raw Fmt = 1 + iota
	Flat
)

const NIL_JSON string = "null"

/*
	For URL "http://cfg-center:2120/conf/deploy/prophet/task-manager/"
	'conf_path_list' should be [deploy, prophet, task-manager]
*/
func (cfgm *CfgManager) GetCfg_json(conf_path_list []string, format Fmt) (js_ret []byte) {
	log.Debug("GetCfg_json ", conf_path_list)

	cfgm.cfg_data_lock.RLock()
	defer cfgm.cfg_data_lock.RUnlock()

	jin := cfgm.cfg_data
	for _, p := range conf_path_list {
		jin = jin.Get(p)
	}

	if format == Flat {
		m, err := jin.Map()
		if err != nil {
			log.Error(err)
			return []byte(NIL_JSON)
		}
		flat_m := mapFlattener.Flatten(m)
		js_ret, _ = json.Marshal(flat_m)
	} else if format == Raw {
		js_ret, _ = jin.MarshalJSON()
	}
	return js_ret
}

func (cfgm *CfgManager) LoadCfg() (js *simplejson.Json, err error) {
	log.Info("Reloading ", cfgm.cfg_file_path)
	js = simplejson.New()
	defer func() {
		if r := recover(); r != nil {
			//todo if cfgm.cfg_file_path is a file
			log.Panic("Config path error, should be dir")
		}
	}()

	// 载入环境公共GLOBAL_VAR.yaml
	g_var_path := cfgm.cfg_file_path + "/" + GLOBAL_VAR
	cfg_t := cfgTemplate.New()
	if fileutil.Exist(g_var_path) {
		err = cfg_t.LoadVars(g_var_path)
		if err != nil {
			log.Info("No global var yaml loaded")
		}
	}

	var scan_func = func(path_one string, info os.FileInfo, e error) (err error) {
		defer func() {
			if r := recover(); r != nil {
				log.Errorf("%s load failed, continue.", path_one)
			}
		}()

		if e == nil && info.Mode().IsRegular() &&
			strings.ToLower(filepath.Ext(info.Name())) == CONF_EXT {
			var (
				cfg_file_content, jsonBytes []byte
				rel_path                    string
			)
			if cfg_file_content, err = ioutil.ReadFile(path_one); err == nil {
				if rel_path, err = filepath.Rel(cfgm.cfg_file_path, path_one); err == nil {
					rel_path_without_ext := rel_path[:len(rel_path)-len(CONF_EXT)]
					path_list := strings.Split(rel_path_without_ext, string(os.PathSeparator))

					// Init template
					inited_content := cfg_t.Translate(cfg_file_content)
					log.Debug(string(inited_content))

					// YAML objects are not completely compatible with JSON objects (e.g. you
					// can have non-string keys in YAML). So, convert the YAML-compatible object
					// to a JSON-compatible object, failing with an error if irrecoverable
					// incompatibilties happen along the way.
					if jsonBytes, err = ym.YAMLToJSON(inited_content); err == nil {
						jsFromCfg := simplejson.New()

						if err = jsFromCfg.UnmarshalJSON(jsonBytes); err == nil {
							js.SetPath(path_list, jsFromCfg.Interface())

							jsout, err := js.EncodePretty()
							log.Debug(string(jsout), err)

							return nil
						}
					}
				}
			}
			log.Error(fmt.Sprintf("Error parsing %s, Skipped", path_one))
			log.Error(err.Error())
		}
		return nil // 永远返回nil,这样就不会因为解析一个文件出错导致不正常运行
	}

	//todo if cfgm.cfg_file_path is a file
	err = filepath.Walk(cfgm.cfg_file_path, scan_func)
	if err != nil {
		log.Error(err.Error())
		return js, err
	}

	jsout, err := js.EncodePretty()
	log.Debug(string(jsout))

	//todo load cfg_data
	cfgm.cfg_data_lock.Lock()
	cfgm.cfg_data = js
	cfgm.cfg_data_lock.Unlock()

	return js, err
}

/*
func (cfgm *CfgManager) LoadCfg() (cfgData, error) {
	theCfg := make(map[string]interface{})
	var scan_func = func(path_one string, info os.FileInfo, err error) error {
		if err == nil && info.Mode().IsRegular() &&
		strings.ToLower(filepath.Ext(info.Name())) == CONF_EXT {
			cfg_file_content, err := ioutil.ReadFile(path_one)
			if err != nil {
				log.Error(err.Error())
				return err
			}
			rel_path, err := filepath.Rel(cfgm.cfg_file_path, path_one)
			if err != nil {
				log.Error(err.Error())
				return err
			}
			rel_path_without_ext := rel_path[:len(rel_path)-len(CONF_EXT)]
			path_list := strings.Split(rel_path_without_ext, string(os.PathSeparator))

			// construct the map dynamically
			var cfg_p = theCfg
			for k, p := range path_list {
				tmp_cfg_p, ok := cfg_p[p]
				if !ok {
					if k == len(path_list) - 1 {
						// *.yaml file
						tmpCfg := make(map[string]map[string]interface{})
						newCfg := make(map[string]interface{})
						if err := yaml.Unmarshal(cfg_file_content, &newCfg); err != nil {
							log.Error(err.Error())
							return err
						}
						tmpCfg[p] = newCfg
						cfg_p[p] = tmpCfg
						//log.Debug(rel_path, " ", "#", path_list, "#", cfg_p)
					} else {
						cfg_p[p] = make(map[string]interface{})
						var success bool
						cfg_p, success = cfg_p[p].(map[string]interface{})
						if !success {
							err := errors.New("construct cfg map")
							log.Error(err.Error())
							return err
						}
					}
				} else {
					var success bool
					cfg_p, success = tmp_cfg_p.(map[string]interface{})
					if !success {
						err := errors.New("construct cfg map")
						log.Error(err.Error())
						return err
					}
				}
			}
		}
		return nil
	}

	err := filepath.Walk(cfgm.cfg_file_path, scan_func)

	log.Debug(theCfg)

	cfgm.cfg_data_lock.Lock()
	cfgm.cfg_data = theCfg
	cfgm.cfg_data_lock.Unlock()


	return theCfg, err
}
*/
