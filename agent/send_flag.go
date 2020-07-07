package main
import (
	"log"
	"net/http"
	"io/ioutil"
	"strings"
	"encoding/base64"
	"encoding/json"
	"errors"
)
func Send_info(server_addr string,data string)(string){
	reqest,err := http.NewRequest("POST", server_addr, strings.NewReader(data))
	reqest.Header.Set("Content-Type","application/json; charset=utf-8")
	reqest.Header.Set("Authorization","Basic "+base64.StdEncoding.EncodeToString([]byte("admin:post")))
	client := &http.Client{}
	resp, err:= client.Do(reqest)
	if err != nil{
		return "error"
	}
	defer resp.Body.Close()
	//body,_:=ioutil.ReadAll(resp.Body)
	log.Println("status",resp.Status)
	return resp.Status
}
func Get_info(server_addr string, ip string)(Target_info, error){
	target_info := []Target_info{}
	server_addr_with_ip := server_addr+base64.StdEncoding.EncodeToString([]byte(ip))
	reqest, err := http.NewRequest("GET", server_addr_with_ip,nil)
	if err != nil{
		log.Println("Get_info",err.Error())
		return Target_info{},err
	}
	reqest.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("admin:get")))
	client := &http.Client{}
	resp,err := client.Do(reqest)
	if err != nil{
		return Target_info{},err
	}
	defer resp.Body.Close()
	body,_ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &target_info)
	if err != nil{
		return Target_info{},err
	}
	if len(target_info) >0{
		return target_info[0],nil
	}else{
		return Target_info{},errors.New("no target_info")
	}
}

func Get_Cmds_not_exec_info(server_addr string, ip string) ([]Cmd_info, error){
	cmd_infos := []Cmd_info{}
	server_addr_with_ip := server_addr+base64.StdEncoding.EncodeToString([]byte(ip))
	reqest,err := http.NewRequest("GET", server_addr_with_ip, nil)
	if err != nil{
		log.Println("Get_info",err.Error())
		return cmd_infos,err
	}
	reqest.Header.Set("Authorization", "Basic "+ base64.StdEncoding.EncodeToString([]byte("admin:get")))
	client := &http.Client{}
	resp, err := client.Do(reqest)
	if err != nil{
		return cmd_infos,err
	}
	defer resp.Body.Close()
	body,_ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &cmd_infos)
	if err != nil{
		return cmd_infos,err
	}
	return cmd_infos,nil
}

