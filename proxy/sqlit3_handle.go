package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

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

type Mysqldb struct {
	Sqldb *sql.DB
}

func connectDB() (*Mysqldb, error) {
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &Mysqldb{db}, nil
}

func (p Mysqldb) Close() {
	p.Sqldb.Close()
}

func (p Mysqldb) DbSetup() {
	var err error
	const tableDDL = `
	CREATE TABLE IF NOT EXISTS Agents (
		Agent_name				VARCHAR (128) PRIMARY KEY	
								UNIQUE
								NOT NULL,
		Delay_time_min 	INTEGER		NOT NULL
								DEFAULT (0),
		Delay_time_max 	INTEGER		NOT NULL
								DEFAULT (0),
		Create_time			DATETIME	NOT NULL,
		Last_online_time	DATETIME	NOT NULL
	
	);
	CREATE TABLE IF NOT EXISTS Cmds (
		Id				INTEGER		PRIMARY KEY AUTOINCREMENT
								UNIQUE
								NOT NULL,
		Agent_name			VARCHAR (128)	REFERENCES Agents (Id) ON DELETE NO ACTION
											ON UPDATE NO ACTION
											MATCH [FULL]
								NOT NULL,
		Cmd_type			VARCHAR (8)	NOT NULL
								DEFAULT control,
		Cmd_value			TEXT		NOT NULL,
		Cmd_time			DATETIME	NOT NULL,
		Cmd_exec			BOOLEAN		NOT NULL
								DEFAULT (false),
		Cmd_exec_time			DATETIME	NOT NULL
								DEFAULT [1900-00-00 00:00:00],
		Cmd_exec_result			TEXT		NOT NULL
								DEFAULT "",
		Cmd_cycle_period_second 	INTEGER		NOT NULL
								DEFAULT (0)
	);
	`
	_, err = p.Sqldb.Exec(tableDDL)
	if err != nil {
		log.Println(err)
		return
	}
	//select count(*) from System_config
}
func (p Mysqldb) Insert_agent(agent_info Agent_info) error {
	const insertDML = "INSERT INTO Agents (Agent_name,Create_time,Last_online_time) VALUES (?,?,?)"
	stmt, err := p.Sqldb.Prepare(insertDML)
	if err != nil {
		log.Println(err)
		return err
	}
	defer stmt.Close()
	result, err := stmt.Exec(agent_info.Agent_name,
		agent_info.Create_time,
		agent_info.Last_online_time)
	if err != nil {
		log.Println(err)
		return err
	}
	lastID, err := result.LastInsertId()
	if err != nil {
		log.Println(err)
		return err
	}
	nAffected, err := result.RowsAffected()
	if err != nil {
		log.Println(err)
		return err
	}
	log.Printf("Exec result id=%d,affected=%d\n", lastID, nAffected)
	return nil
}

func (p Mysqldb) Update_agent(agent_info Agent_info) error {
	const insertDML = "UPDATE Agents SET Last_online_time=? WHERE Agent_name=?"
	stmt, err := p.Sqldb.Prepare(insertDML)
	if err != nil {
		log.Println(err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(agent_info.Last_online_time, agent_info.Agent_name)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (p Mysqldb) Get_Agents_info(agent_name string) ([]Agent_info, error) {
	agent_infos := []Agent_info{}
	agent_info := Agent_info{}
	const selectDML = "SELECT * FROM Agents WHERE Agent_name=?" //TDO need youhua
	stmt, err := p.Sqldb.Prepare(selectDML)
	if err != nil {
		log.Println(err)
		return agent_infos, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(agent_name)
	if err != nil {
		log.Println(err)
		return agent_infos, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&agent_info.Agent_name,
			&agent_info.Delay_time_min,
			&agent_info.Delay_time_max,
			&agent_info.Create_time,
			&agent_info.Last_online_time)
		if err != nil {
			log.Println(err)
			continue
		}
		agent_infos = append(agent_infos, agent_info)
	}
	if err := rows.Err(); err != nil {
		log.Println(err)
		return agent_infos, err
	}
	return agent_infos, nil
}

func (p Mysqldb) Get_All_Agents_info() ([]Agent_info, error) {
	agent_infos := []Agent_info{}
	agent_info := Agent_info{}
	const selectDML = "SELECT * FROM Agents" //TDO need youhua
	stmt, err := p.Sqldb.Prepare(selectDML)
	if err != nil {
		log.Println(err)
		return agent_infos, err
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		log.Println(err)
		return agent_infos, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&agent_info.Agent_name,
			&agent_info.Delay_time_min,
			&agent_info.Delay_time_max,
			&agent_info.Create_time,
			&agent_info.Last_online_time)
		if err != nil {
			log.Println(err)
			continue
		}
		agent_infos = append(agent_infos, agent_info)
	}
	if err := rows.Err(); err != nil {
		log.Println(err)
		return agent_infos, err
	}
	return agent_infos, nil
}

func (p Mysqldb) Create_cmd(cmd_info Cmd_info) error {
	const insertDML = "INSERT INTO Cmds (Agent_name,Cmd_type,Cmd_value,Cmd_time) VALUES (?,?,?,?)"
	stmt, err := p.Sqldb.Prepare(insertDML)
	if err != nil {
		log.Println(err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(cmd_info.Agent_name, cmd_info.Cmd_type, cmd_info.Cmd_value, cmd_info.Cmd_time)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (p Mysqldb) Update_cmd_result(cmd_info Cmd_info) error {
	const insertDML = "UPDATE Cmds SET Cmd_exec=?,Cmd_exec_time=?,Cmd_exec_result=? WHERE id=?"
	stmt, err := p.Sqldb.Prepare(insertDML)
	if err != nil {
		log.Println(err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(cmd_info.Cmd_exec,
		cmd_info.Cmd_exec_time,
		cmd_info.Cmd_exec_result,
		cmd_info.Id)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (p Mysqldb) Get_cmds_not_exec(agent_name string) ([]Cmd_info, error) {
	cmd_infos := []Cmd_info{}
	cmd_info := Cmd_info{}

	selectDML := "SELECT * FROM Cmds WHERE Cmd_exec=? AND Agent_name=?"
	stmt, err := p.Sqldb.Prepare(selectDML)
	if err != nil {
		log.Println(err)
		return cmd_infos, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(0, agent_name)
	if err != nil {
		log.Println(err)
		return cmd_infos, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&cmd_info.Id,
			&cmd_info.Agent_name,
			&cmd_info.Cmd_type,
			&cmd_info.Cmd_value,
			&cmd_info.Cmd_time,
			&cmd_info.Cmd_exec,
			&cmd_info.Cmd_exec_time,
			&cmd_info.Cmd_exec_result,
			&cmd_info.Cmd_cycle_period_second)
		if err != nil {
			log.Println(err)
			continue
		}
		cmd_infos = append(cmd_infos, cmd_info)
	}
	if err := rows.Err(); err != nil {
		log.Println(err)
		return cmd_infos, err
	}
	return cmd_infos, nil
}
func (p Mysqldb) Get_all_cmds_by_agent_name(agent_name string) ([]Cmd_info, error) {
	cmd_infos := []Cmd_info{}
	cmd_info := Cmd_info{}

	selectDML := "SELECT * FROM Cmds WHERE Agent_name=?"
	stmt, err := p.Sqldb.Prepare(selectDML)
	if err != nil {
		log.Println(err)
		return cmd_infos, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(agent_name)
	if err != nil {
		log.Println(err)
		return cmd_infos, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&cmd_info.Id,
			&cmd_info.Agent_name,
			&cmd_info.Cmd_type,
			&cmd_info.Cmd_value,
			&cmd_info.Cmd_time,
			&cmd_info.Cmd_exec,
			&cmd_info.Cmd_exec_time,
			&cmd_info.Cmd_exec_result,
			&cmd_info.Cmd_cycle_period_second)
		if err != nil {
			log.Println(err)
			continue
		}
		cmd_infos = append(cmd_infos, cmd_info)
	}
	if err := rows.Err(); err != nil {
		log.Println(err)
		return cmd_infos, err
	}
	return cmd_infos, nil
}
