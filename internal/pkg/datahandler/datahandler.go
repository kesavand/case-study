package datahandler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/kesavand/case-study/internal/pkg/utils"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	dbName = "user_info"
)

type Datahandler struct {
	db *gorm.DB
}

type DataHandlerInterface interface {
	Write(userInfo *utils.UserInfo) error
	GetAllusers(w http.ResponseWriter, users *[]utils.UserInfo)
	GetUsersByAge(w http.ResponseWriter, users *[]utils.UserInfo, age int)
	GetUsersByName(w http.ResponseWriter, users *[]utils.UserInfo, name string)
	GetUsersByNameAndAge(w http.ResponseWriter, users *[]utils.UserInfo, name string, age int)
}

var once sync.Once
var db *gorm.DB

func NewDBConnection() *gorm.DB {
	once.Do(func() { // <-- atomic, does not allow repeating

		dsn := "root@tcp(mariadb:3306)/?parseTime=true"

		sqlDB, err := sql.Open("mysql", dsn)

		if err != nil {
			panic(err)
		}

		_, err = sqlDB.Exec("CREATE DATABASE IF NOT EXISTS " + dbName)
		if err != nil {
			panic(err)
		}

		_, err = sqlDB.Exec("USE " + dbName)
		if err != nil {
			panic(err)
		}

		if err != nil {
			fmt.Println(err.Error())
			panic("failed to connect database")
		}

		db, err = gorm.Open(mysql.New(mysql.Config{
			Conn: sqlDB,
		}), &gorm.Config{})

		if err != nil {
			fmt.Println(err.Error())
			panic("failed to connect database")
		}

		user := &utils.UserInfo{}

		db.AutoMigrate(user)

	})

	return db
}

func NewDatahandler() DataHandlerInterface {

	dh := &Datahandler{
		db: NewDBConnection(),
	}

	return dh
}

func (dh *Datahandler) Write(userInfo *utils.UserInfo) error {
	if dh != nil {
		dh.db.Create(userInfo)
	} else {
		return fmt.Errorf("the api handler is nil")
	}
	fmt.Printf("Successfully written data %v\n", userInfo)
	return nil
}

func (dh *Datahandler) Getusers(w http.ResponseWriter, users *utils.UserInfo) {
	dh.db.Find(&users)

	json.NewEncoder(w).Encode(users)
}

func (dh *Datahandler) Getuser(w http.ResponseWriter) {
	var users []utils.UserInfo
	dh.db.Find(users)
	fmt.Println("{}", users)

	json.NewEncoder(w).Encode(users)
}

func (dh *Datahandler) GetAllusers(w http.ResponseWriter, users *[]utils.UserInfo) {
	dh.db.Find(users)
	fmt.Println("{}", users)

	json.NewEncoder(w).Encode(users)
}

func (dh *Datahandler) GetUsersByName(w http.ResponseWriter, users *[]utils.UserInfo, name string) {
	dh.db.Where("name = ?", name).Find(users)
	fmt.Println("{}", users)

	json.NewEncoder(w).Encode(users)
}

func (dh *Datahandler) GetUsersByAge(w http.ResponseWriter, users *[]utils.UserInfo, age int) {
	dh.db.Where("age = ?", age).Find(users)
	fmt.Println("{}", users)

	json.NewEncoder(w).Encode(users)
}

func (dh *Datahandler) GetUsersByNameAndAge(w http.ResponseWriter, users *[]utils.UserInfo, name string, age int) {
	dh.db.Where("age = ?", age).Where("name = ?", name).Find(users)
	fmt.Println("{}", users)

	json.NewEncoder(w).Encode(users)
}
