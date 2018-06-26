package sqlSvc

import (
	"fmt"
	"database/sql"
	"log"
	"io"
	"context"
	"time"
	"gopkg.in/goracle.v2"
	"net/http"
	"encoding/json"
)

type DBConfig struct {
	driver   string
	dns      string
	db       string
	user     string
	password string
}

func InitRestAPI() {
	http.HandleFunc("/dumpTable", func(w http.ResponseWriter, r *http.Request) {
		tableNames, ok := r.URL.Query()["tableName"]
		if (ok) && len(tableNames) > 0 {
			tableName := tableNames[0]
			excecuteQueryAndDump(tableName, w)
		} else
		{
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, "Table is needed as a parameter !")
		}
	})
	log.Println("InitRestAPI for sqlSvc registered")
}

func buildConnectString(conf *DBConfig) string {
	return conf.user + "/" + conf.password + "@" + conf.db
}

func connectToDb(conf DBConfig) (*sql.DB) {
	dbConnectionString := buildConnectString(&conf)
	fmt.Println("Connect to DB using " + dbConnectionString)

	db, err := sql.Open(conf.driver, dbConnectionString)
	if err != nil {
		log.Printf("sql.Open(%s, %s)\n\t%s\n",
			conf.driver, dbConnectionString, err.Error())
		fmt.Println("Cannot onnect to DB " + err.Error())
		return nil
	}

	context.WithTimeout(context.Background(), 30*time.Second)
	ctxt, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	db.Conn(ctxt)

	return db
}

func dumpTable(rows *sql.Rows, out io.Writer) error {
	colNames, err := rows.Columns()
	if err != nil {
		return err
	}
	jsonEncoder := json.NewEncoder(out)

	//writer := csv.NewWriter(out)
	//writer.Comma = ','

	readCols := make([]interface{}, len(colNames))

	colStringValues := make([]sql.NullString, len(colNames))
	for i := range colStringValues {
		readCols[i] = &colStringValues[i]
	}

	//writer.Write(colNames)
	jsonEncoder.Encode(colNames)
	for rows.Next() {
		err := rows.Scan(readCols...)
		if err != nil {
			return err
		}
		//writer.Write(convertToString(colStringValues))
		jsonEncoder.Encode(convertToString(colStringValues))
	}
	if err = rows.Err(); err != nil {
		panic(err)
	}
	//writer.Flush()

	defer rows.Close()
	return nil
}

func excecuteQueryAndDump(tableName string, out io.Writer) {
	if error := dumpTable(executeSelect(tableName, out), out); error != nil {
		fmt.Fprintln(out, "Error :"+error.Error())
	}
}

func executeSelect(tableName string, out io.Writer) (*sql.Rows) {
	// Execution Context
	db := connect()

	if (tableName == "*") {
		log.Println("Displaying rows from table " + tableName)
		rows, error := db.QueryContext(context.Background(), "SELECT TABLE_NAME FROM USER_TABLES")
		log.Println("Displaying all tables in schema")
		if (error != nil) {
			log.Println("Cannot get all tables " + error.Error())
			return nil
		}
		return rows
	} else {
		log.Println("Displaying rows from table " + tableName)
		rows, error := db.QueryContext(context.Background(), "SELECT * FROM "+tableName, goracle.ClobAsString())
		if (error != nil) {
			log.Println("Cannot retrieve rows from " + tableName + ". Error -> " + error.Error())
			return nil
		}
		return rows
	}
	db.Close()
	return nil
}

func connect() (*sql.DB) {
	conf := DBConfig{"goracle", "host", "instance", "SCHEMA_OWNER", "PASSWORD"}
	return connectToDb(conf)
}
