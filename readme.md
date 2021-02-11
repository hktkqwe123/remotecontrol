先占坑

curl -u admin:get http://127.0.0.1:8080/get_all_agent_info


curl -u admin:post -H "Content-Type:application/json" -X POST --data '{"Agent_name": "test","Cmd_type":"cmd","Cmd_value":"cmd.exe /c dir"}' http://localhost:8080/create_cmd


curl -u admin:get http://127.0.0.1:8080/et_all_cmds_by_agent_name/test