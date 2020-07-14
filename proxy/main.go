
package main
import (
	"github.com/ant0ine/go-json-rest/rest"
	"encoding/base64"
	"log"
	"net/http"
	"sync"
)

var lock = sync.RWMutex{}
var Mysql *Mysqldb
func init() {
	var err error
	Mysql,err = connectDB()
	if err != nil{
		log.Println(err)
	}
	Mysql.DbSetup()
}

func send_flag_task(){
	//TODO
}

func main(){
	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)

	auth_post_flag := &rest.AuthBasicMiddleware{
		Realm: "test zone",
		Authenticator: func(userId string,password string)bool {
			if userId == "admin" && password == "post" {
				return true
			}
			return false
		},
	}

	auth_user := &rest.AuthBasicMiddleware{
		Realm: "test zone",
		Authenticator: func(userId string, password string) bool {
			if userId == "admin" && password == "get" {
				return true
			}
			return false
		},
	}

	router,err := rest.MakeRouter(
		rest.Post("/create_target", auth_post_flag.MiddlewareFunc(CreateTarget)),
		rest.Get("/get_system_config",auth_user.MiddlewareFunc(Get_system_config)),
		rest.Get("/get_target_info",auth_user.MiddlewareFunc(Get_target_info)),
		rest.Get("/get_target_info/:ip",auth_user.MiddlewareFunc(Get_target_info_by_ip)),
		rest.Get("/get_cmd_not_exec/:ip",auth_user.MiddlewareFunc(Get_cmd_not_exec_by_ip)),
		rest.Post("/update_cmd_result",auth_post_flag.MiddlewareFunc(Update_cmd_result)),
		rest.Post("/create_cmd",auth_post_flag.MiddlewareFunc(Create_new_cmd)),
	)
	if err != nil{
		log.Fatal(err)
	}
	api.SetApp(router)
	go send_flag_task()
	log.Fatal(http.ListenAndServe(":8080",api.MakeHandler()))
}

func Create_new_cmd(w rest.ResponseWriter, r *rest.Request) {
	cmd_info := Cmd_info{}
	err := r.DecodeJsonPayload(&cmd_info)
	if err != nil {
		rest.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}
	err = Mysql.Create_cmd(cmd_info)
	if err != nil {
		rest.Error(w, err.Error(),400)
	}
	w.WriteJson(&cmd_info)
}

func Update_cmd_result(w rest.ResponseWriter, r *rest.Request) {
	cmd_info := Cmd_info{}
	err := r.DecodeJsonPayload(&cmd_info)
	if err != nil {
		rest.Error(w, err.Error() , http.StatusInternalServerError)
		return
	}
	err = Mysql.Update_cmd_result(cmd_info)
	if err  != nil{
		rest.Error(w, err.Error(), 400)
	}
	w.WriteJson(&cmd_info)
}

func Get_cmd_not_exec_by_ip(w rest.ResponseWriter, r *rest.Request) {
	sip := r.PathParam("ip")
	ip_byte,_:= base64.StdEncoding.DecodeString(sip)
	ip := string(ip_byte)
	log.Println(ip)
	lock.RLock()
	cmd_info,err := Mysql.Get_cmds_not_exec(ip)
	log.Println(err)
	lock.RUnlock()
	w.WriteJson(&cmd_info)
}

func Get_target_info_by_ip(w rest.ResponseWriter, r *rest.Request) {
	sip := r.PathParam("ip")
	ip_byte,_ := base64.StdEncoding.DecodeString(sip)
	ip := string(ip_byte)
	log.Println(ip)
	lock.RLock()
	target_info, err := Mysql.Get_targets_info(ip)
	log.Println(err)
	lock.RUnlock()
	w.WriteJson(&target_info)
}

func Get_target_info( w rest.ResponseWriter, r*rest.Request) {
	lock.RLock()
	target_info,err := Mysql.Get_targets_info("all")
	if err != nil{
		log.Println(err)
	}
	lock.RUnlock()
	w.WriteJson(&target_info)
}

func Get_system_config(w rest.ResponseWriter, r *rest.Request) {
	lock.RLock()
	system_config_info, err := Mysql.Get_System_Config_Info()
	if err != nil {
		log.Println(err)
	}
	lock.RUnlock()
	w.WriteJson(&system_config_info)
}

func CreateTarget(w rest.ResponseWriter, r *rest.Request) {
	target_info := Target_info{}
	err := r.DecodeJsonPayload(&target_info.Target_base_info)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//net.ParseIP(ipv4)
	if target_info.Target_base_info.Ip ==""{
		rest.Error(w, "flag Ip required", 400)
		return
	}
	lock.RLock()
	target_info.System_config_info,_=Mysql.Get_System_Config_Info()
	Mysql.Insert_targets(target_info)
	lock.RUnlock()
	w.WriteJson(&target_info)
}
