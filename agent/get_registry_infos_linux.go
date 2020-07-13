package main
import(
)
func Get_registry_infos()(registry_infos map[string]map[string]string){
	registry_infos = make(map[string]map[string]string)
	errors := make(map[string]string)
	errors["error"]="Only windows support registry!"
	registry_infos["error"]=errors
	return
}
