package main

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/ant0ine/go-json-rest/rest"
)

type Agent_info struct {
	Agent_name       string
	Delay_time_min   int
	Delay_time_max   int
	Create_time      string
	Last_online_time string
}

type Register_info struct {
	Action     string
	Agent_name string
}

var lock = sync.RWMutex{}
var Mysql *Mysqldb

func init() {
	var err error
	Mysql, err = connectDB()
	if err != nil {
		log.Println(err)
	}
	Mysql.DbSetup()
}

func main() {
	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)

	auth_post_flag := &rest.AuthBasicMiddleware{
		Realm: "test zone",
		Authenticator: func(userId string, password string) bool {
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

	router, err := rest.MakeRouter(
		rest.Post("/register", auth_post_flag.MiddlewareFunc(Register_handle)),
		rest.Post("/create_cmd", auth_post_flag.MiddlewareFunc(Create_new_cmd)),
		rest.Post("/update_cmd_result", auth_post_flag.MiddlewareFunc(Update_cmd_result)),
		rest.Get("/get_cmd_not_exec/:agent_name", auth_user.MiddlewareFunc(Get_cmd_not_exec_by_agent_name)),
		rest.Get("/et_all_cmds_by_agent_name/:agent_name", auth_user.MiddlewareFunc(Get_all_cmds_by_agent_name)),
		rest.Get("/get_all_agent_info", auth_user.MiddlewareFunc(get_all_agent_info)),
		rest.Get("/get_agent_info/:agent_name", auth_user.MiddlewareFunc(get_agent_info)),
	)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)
	log.Fatal(http.ListenAndServe(":8080", api.MakeHandler()))
}

func Register_handle(w rest.ResponseWriter, r *rest.Request) {
	register_info := Register_info{}

	err := r.DecodeJsonPayload(&register_info)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if register_info.Action == "register" {
		agent_infos, err := Mysql.Get_Agents_info(register_info.Agent_name)
		if err == nil {
			if len(agent_infos) > 0 {
				for _, agent_info := range agent_infos {
					agent_info.Last_online_time = time.Now().Format("2006-01-02 15:04:05")
					err := Mysql.Update_agent(agent_info)
					if err == nil {
					}
				}
				w.WriteJson(&register_info)
				return

			} else {
				agent_info := Agent_info{}
				agent_info.Agent_name = register_info.Agent_name
				agent_info.Create_time = time.Now().Format("2006-01-02 15:04:05")
				agent_info.Last_online_time = agent_info.Create_time
				err = Mysql.Insert_agent(agent_info)
				if err == nil {
					w.WriteJson(&register_info)
					return
				}
			}

		} else {
			rest.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
	}
	rest.Error(w, "reg error", http.StatusInternalServerError)
	return

}

func Create_new_cmd(w rest.ResponseWriter, r *rest.Request) {
	cmd_info := Cmd_info{}
	err := r.DecodeJsonPayload(&cmd_info)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = Mysql.Create_cmd(cmd_info)
	if err != nil {
		rest.Error(w, err.Error(), 400)
	}
	w.WriteJson(&cmd_info)
}

func Update_cmd_result(w rest.ResponseWriter, r *rest.Request) {
	cmd_info := Cmd_info{}
	err := r.DecodeJsonPayload(&cmd_info)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = Mysql.Update_cmd_result(cmd_info)
	if err != nil {
		rest.Error(w, err.Error(), 400)
	}
	w.WriteJson(&cmd_info)
}

func Get_cmd_not_exec_by_agent_name(w rest.ResponseWriter, r *rest.Request) {
	agent_name := r.PathParam("agent_name")
	log.Println(agent_name)
	lock.RLock()
	cmd_infos, err := Mysql.Get_cmds_not_exec(agent_name)
	log.Println(err)
	lock.RUnlock()
	w.WriteJson(&cmd_infos)
}

func Get_all_cmds_by_agent_name(w rest.ResponseWriter, r *rest.Request) {
	agent_name := r.PathParam("agent_name")
	log.Println(agent_name)
	lock.RLock()
	cmd_infos, err := Mysql.Get_all_cmds_by_agent_name(agent_name)
	log.Println(err)
	lock.RUnlock()
	w.WriteJson(&cmd_infos)
}

func get_all_agent_info(w rest.ResponseWriter, r *rest.Request) {
	lock.RLock()
	agent_infos, err := Mysql.Get_All_Agents_info()
	if err != nil {
		log.Println(err)
	}
	lock.RUnlock()
	w.WriteJson(&agent_infos)
}

func get_agent_info(w rest.ResponseWriter, r *rest.Request) {
	agent_name := r.PathParam("agent_name")
	log.Println(agent_name)
	lock.RLock()
	agent_infos, err := Mysql.Get_Agents_info(agent_name)
	log.Println(err)
	lock.RUnlock()
	w.WriteJson(&agent_infos)
}
