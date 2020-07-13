package main
import(
	"log"
	"strings"
	"strconv"
	"errors"
	registry "golang.org/x/sys/windows/registry"
)
func getSubKeyNamesFromRegistry(key_type registry.Key,regKey string) (keyNames[]string){
	k,err := registry.OpenKey(key_type, regKey, registry.ENUMERATE_SUB_KEYS)
	if err != nil {
		log.Println("Can't open registry key "+regKey,err)
		return
	}
	defer k.Close()
	keyNames,err = k.ReadSubKeyNames(0)
	return
}

func getSettingsFromRegistry(key_type registry.Key,regKey string) (settings map[string]string){
	settings = make(map[string]string)
	k,err := registry.OpenKey(key_type,regKey,registry.QUERY_VALUE)
	if err != nil{
		log.Println("Cant't open registry key"+regKey,err)
		return
	}
	defer k.Close()
	params,err := k.ReadValueNames(0)
	if err != nil {
		log.Printf("Can't ReadSubKeyNames %#v", err)
		return
	}
	for _, param := range(params){
		val,err := getRegistryValueAsString(k,param)
		if err != nil{
			log.Println(err)
			return
		}
		settings[param] = val
	}
	return
}

func getRegistryValueAsString(key registry.Key, subKey string) (string,error){
	valString,_,err:= key.GetStringValue(subKey)
	if err ==nil{
		return valString,nil
	}
	valStrings,_,err := key.GetStringsValue(subKey)
	if err == nil{
		return strings.Join(valStrings,"\n"),nil
	}
	valBinary,_,err := key.GetBinaryValue(subKey)
	if err == nil{
		return string(valBinary),nil
	}
	valInteger,_,err := key.GetIntegerValue(subKey)
	if err == nil{
		return strconv.FormatUint(valInteger,10),nil
	}
	return "error", errors.New("Can't get type for sub key"+subKey)
}

func Get_registry_infos() (registry_infos map[string]map[string]string){
	registry_infos = make(map[string]map[string]string)
	CurrentVersion := `SOFTWARE\Microsoft\Windows NT\CurrentVersion`
	registry_infos[CurrentVersion] = getSettingsFromRegistry(registry.LOCAL_MACHINE,CurrentVersion)
	Client_Default := `Software\Microsoft\Terminal Server Client\Default`
	registry_infos[Client_Default] = getSettingsFromRegistry(registry.CURRENT_USER,Client_Default)
	servers := `Software\Microsoft\Terminal Server Client\Servers`
	subKeys := getSubKeyNamesFromRegistry(registry.CURRENT_USER,servers)
	for _,subKey := range(subKeys){
		registry_infos[servers+`\`+subKey]=getSettingsFromRegistry(registry.CURRENT_USER,servers+`\`+subKey)
	}
	for v,k := range(registry_infos){
		log.Println(v,k)
	}
	return
}
