package main

import (
	"encoding/json"
	"flag"
	//	"strings"
	"time"
	//	"time"

	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"sync"

	"os"
	//	"os/exec"
	//	"os/signal"
	"runtime"
	"strconv"

	"github.com/gorilla/handlers"
	"github.com/kardianos/osext"

	"github.com/BurntSushi/toml"

	"github.com/gorilla/mux"

	//"github.com/gorilla/securecookie"
	"github.com/gorilla/websocket"
)

type testConfig struct {
	Title  string
	Port   int
	Server string
}

type DisCoveryInfo struct {
	Server  string
	Port    int
	Label   string
	Enabled bool
}

type ownerInfo struct {
	Name string
}

type serverInfo struct {
	Server  string
	Port    int
	Enabled bool
}

type testInfo struct {
	Program string
	Log     serverInfo
	Cache   serverInfo
	WSIp    string
	WSPath  string
}

var memprofile = flag.String("memprofile", "", "write memory profile to this file")

var router = mux.NewRouter()
var G_LOCALIP = ""

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func ReadCfg(cfg_paths []string) testConfig {
	cfg := flag.String("cfg", "", "Configure file.")
	flag.Parse()

	if *cfg != "" {
		cfg_paths = append(cfg_paths, *cfg)
	}
	var config testConfig
	for _, cfg_path := range cfg_paths {
		if _, err := toml.DecodeFile(cfg_path, &config); err != nil {
			continue
		}
		break
	}
	return config
}

type CPDFService struct {
	code_s      string
	code_i      int
	status      string
	test_status string
	log_s       string
	LogLocker   *sync.Mutex
}

type LogService struct {
	logPath       string
	logFileHandle map[string]*os.File
}

type TestService struct {
	cfg         testConfig
	logServer   map[string]string
	cacheServer map[string]string
	logfile     map[string]*os.File
	logger      map[string]*log.Logger
	wsconn      *websocket.Conn
}

type LogResult struct {
	Info string
	Err  int
}

func (this *TestService) Init() {
	this.logServer = make(map[string]string)
	this.cacheServer = make(map[string]string)
	this.logfile = make(map[string]*os.File)
	this.logger = make(map[string]*log.Logger)

	//	var err error

	//	execlogpath := "/var/log/fxqa/test.log"
	if runtime.GOOS == "windows" {
		//		execlogpath = "test.log"
	} else {
		ret, _ := exists("/var/log/fxqa")
		if !ret {
			os.MkdirAll("/var/log/fxqa", 0777)
		}
	}

}

func (this *TestService) Log(test_key, logstr string) {

}

func Info(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Host)
	//	fmt.Println(GetLocalIP())
	fmt.Fprintf(w, "%v", `{"val":"alert('test')"}`)
}

func (this *CPDFService) SetCode(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)

	if len(this.code_s) > 0 {
		fmt.Fprintf(w, "%v", `{"ret":-1,"msg":"running"}`)
		return
	}

	this.code_s = r.Form["code"][0]
	fmt.Fprintf(w, "%v", `{"ret":0}`)
}

func (this *CPDFService) GetCode(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)

	if len(this.code_s) == 0 {
		fmt.Fprintf(w, "%v", `{"ret":-1}`)
		return
	}
	fmt.Println("GETCODE:", this.code_s)
	//	fmt.Println(fmt.Sprintf(`{"ret":0,"code":"%s"}`, this.code_s))
	//	res_log := this.code_s
	fmt.Fprintf(w, "%v", fmt.Sprintf(`{"ret":0,"code":"%s"}`, this.code_s))

}

func (this *CPDFService) DelCode(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)

	this.code_s = ""

	fmt.Fprintf(w, "%v", `{"ret":0}`)
}

func (this *CPDFService) SetLog(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)

	defer fmt.Fprintf(w, "%v", `{"ret":0}`)
	log_s := LogResult{}
	err_i, _ := strconv.Atoi(r.FormValue("err"))
	log_s.Err = err_i
	log_s.Info = r.FormValue("str")

	_data, err := json.Marshal(log_s)
	if err != nil {
		this.log_s = `{"Info":"internal error.","Err":-2}`
		return
	}
	this.LogLocker.Lock()
	this.log_s = string(_data)
	this.LogLocker.Unlock()
	fmt.Println("SET LOG:", this.log_s)
}

func (this *CPDFService) GetLog(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)

	//	fmt.Println(this.log_s)
	this.LogLocker.Lock()
	res_log := this.log_s
	this.log_s = ""
	this.LogLocker.Unlock()
	fmt.Println("GET LOG:", res_log)
	fmt.Println("GLOG:", this.log_s)
	fmt.Fprintf(w, "%v", res_log)
	//fmt.Fprintf(w, "%v", fmt.Sprintf(`{"ret":0,"log":"%s"}`, strings.Replace(this.log_s, `"`, `\"`, -1)))
}

func (this *CPDFService) Status(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)

	action_s := r.Form["action"][0]
	if action_s == "start" {
		this.status = "start"
		this.test_status = ""
		this.code_s = ""
		this.log_s = ""
		fmt.Fprintf(w, "%v", `{"ret":0}`)
	} else if action_s == "stop" {
		this.status = ""
		fmt.Fprintf(w, "%v", `{"ret":0}`)
	}
}

func (this *CPDFService) SetTestStatus(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)

	action_s := r.Form["action"][0]
	if action_s == "start" {
		this.test_status = "start"
	} else if action_s == "end" {
		this.test_status = "end"
	}
	fmt.Fprintf(w, "%v", `{"ret":0}`)
}

func (this *CPDFService) GetTestStatus(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)

	fmt.Fprintf(w, "%v", fmt.Sprintf(`{"ret":0,"status":"%s"}`, this.test_status))
}

func (this *CPDFService) Wait(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)

	for len(this.code_s) > 0 {
		time.Sleep(1e7)
	}
	fmt.Fprintf(w, "%v", fmt.Sprintf(`{"ret":0,"code":"%s"}`, this.code_s))

}

func main() {
	//	cfg := ReadCfg([]string{"/etc/fxqa-test.conf", "fxqa-test.conf"})
	portPtr := flag.Int("port", 9092, "port")
	flag.Parse()

	cpdf_serv := new(CPDFService)
	cpdf_serv.LogLocker = new(sync.Mutex)

	//	test_serv := TestService{cfg: cfg}
	//	test_serv.Init()

	//	interrupt := make(chan os.Signal, 1)
	//signal.Notify(interrupt, os.Interrupt)

	excfodler, _ := osext.ExecutableFolder()
	router.PathPrefix("/workspace/").Handler(
		http.StripPrefix("/workspace", http.FileServer(http.Dir(excfodler))))

	router.HandleFunc("/info", Info).Methods("GET")
	router.HandleFunc("/code", cpdf_serv.SetCode).Methods("POST")
	router.HandleFunc("/code", cpdf_serv.GetCode).Methods("GET")
	router.HandleFunc("/code", cpdf_serv.DelCode).Methods("DELETE")
	router.HandleFunc("/clearcode", cpdf_serv.DelCode).Methods("GET")

	router.HandleFunc("/wait", cpdf_serv.Wait).Methods("GET")

	router.HandleFunc("/log", cpdf_serv.SetLog).Methods("POST")
	router.HandleFunc("/log", cpdf_serv.GetLog).Methods("GET")

	router.HandleFunc("/status", cpdf_serv.Status).Methods("POST")

	router.HandleFunc("/test-status", cpdf_serv.SetTestStatus).Methods("POST")
	router.HandleFunc("/test-status", cpdf_serv.GetTestStatus).Methods("GET")

	http.Handle("/", router)
	fmt.Println("START")
	http.ListenAndServe(":"+strconv.Itoa(*portPtr), handlers.CORS()(router))
}
