package httpApi

import (
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/4paradigm/cfg-center/src/cfg-server/cfgLoader"
	log "github.com/auxten/logrus"
)

var hello_page = `
<h1>Welconme to cfg-server</h1>
`

func hello(w http.ResponseWriter, r *http.Request, h *myHandler) {
	io.WriteString(w, hello_page)
}

func conf(w http.ResponseWriter, r *http.Request, h *myHandler) {
	clean_url := filepath.Clean(r.URL.Path)
	clean_url_list := strings.Split(clean_url, "/")

	log.Debug(clean_url, clean_url_list)

	req_args := r.URL.Query()
	var resp_str string
	if req_args.Get("flat") == "1" {
		resp_str = string(h.cfgm.GetCfg_json(clean_url_list[2:], cfgLoader.Flat))
	} else {
		resp_str = string(h.cfgm.GetCfg_json(clean_url_list[2:], cfgLoader.Raw))
	}
	log.Debug(resp_str)
	io.WriteString(w, resp_str)
}

func reloadGit(w http.ResponseWriter, r *http.Request, h *myHandler) {
	resp_str := h.cfgm.ReloadGitCfg()
	log.Debug(resp_str)
	io.WriteString(w, resp_str)
}

var mux map[string]func(http.ResponseWriter, *http.Request, *myHandler)

func HTTPServerStart(listenPort int, cfgm *cfgLoader.CfgManager) error {
	strListenPort := fmt.Sprintf(":%d", listenPort)
	server := http.Server{
		Addr: strListenPort,
		Handler: &myHandler{
			cfgm: cfgm,
		},
	}

	mux = make(map[string]func(http.ResponseWriter, *http.Request, *myHandler))
	mux[""] = hello
	mux["conf"] = conf
	mux["reload-git"] = reloadGit

	return server.ListenAndServe()
}

type myHandler struct {
	cfgm *cfgLoader.CfgManager
}

/*
 - Implementation of
	type Handler interface {
		ServeHTTP(ResponseWriter, *Request)
	}
 - See https://gowalker.org/net/http#Handler
*/
func (http_handler *myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	/*
	  request URL like "http://127.0.0.1:2120" makes r.URL.String == "/"
	*/

	clean_url := filepath.Clean(r.URL.Path)
	clean_url_list := strings.Split(clean_url, "/")

	url_index := clean_url_list[1]
	//
	//log.Debug(clean_url_list)
	//log.Debug(url_index)

	if h, ok := mux[url_index]; ok {
		h(w, r, http_handler)
		return
	}

	w.WriteHeader(http.StatusNotFound)
	io.WriteString(w, clean_url+" conf not found")
}
