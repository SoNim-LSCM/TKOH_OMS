package database

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	db_models "github.com/SoNim-LSCM/TKOH_OMS/database/models"
	"github.com/SoNim-LSCM/TKOH_OMS/errors"

	"io/ioutil"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/ssh"
	sql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

type Dialer struct {
	client *ssh.Client
}

type SSH struct {
	Host     string `json:"host"`
	User     string `json:"user"`
	Port     int    `json:"port"`
	Type     string `json:"type"`
	Password string `json:"password"`
	KeyFile  string `json:"key"`
}

type MySQL struct {
	Host     string `json:"host"`
	User     string `json:"user"`
	Port     int    `json:"port"`
	Password string `json:"password"`
	Database string `json:"database"`
}

func (v *Dialer) Dial(address string) (net.Conn, error) {
	return v.client.Dial("tcp", address)
}

func (s *SSH) DialWithPassword() (*ssh.Client, error) {
	address := fmt.Sprintf("%s:%d", s.Host, s.Port)
	config := &ssh.ClientConfig{
		User: s.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(s.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	return ssh.Dial("tcp", address, config)
}

func (s *SSH) DialWithKeyFile() (*ssh.Client, error) {
	address := fmt.Sprintf("%s:%d", s.Host, s.Port)
	config := &ssh.ClientConfig{
		User:            s.User,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	if k, err := ioutil.ReadFile(s.KeyFile); err != nil {
		return nil, err
	} else {
		signer, err := ssh.ParsePrivateKey(k)
		if err != nil {
			return nil, err
		}
		config.Auth = []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		}
	}
	return ssh.Dial("tcp", address, config)
}

func (m *MySQL) New() (db *gorm.DB, err error) {
	// 填写注册的mysql网络
	dsn := fmt.Sprintf("%s:%s@mysql+ssh(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		m.User, m.Password, m.Host, m.Port, m.Database)
	db, err = gorm.Open(sql.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		return
	}
	return
}

func StartMySql() {

	var (
		dial *ssh.Client
		err  error
	)

	sshPort, err := strconv.Atoi(os.Getenv("SSH_PORT"))
	errors.CheckError(err, "translate string to int in mysql")
	client := SSH{
		Host:     os.Getenv("SSH_HOST"),
		User:     os.Getenv("SSH_USERNAME"),
		Port:     sshPort,
		Password: os.Getenv("SSH_PASSWORD"),
		// KeyFile: "~/.ssh/id_rsa",
		Type: "PASSWORD", // PASSWORD or KEY
	}
	dbPort, err := strconv.Atoi(os.Getenv("MYSQL_DB_PORT"))
	errors.CheckError(err, "translate string to int in mysql")
	my := MySQL{
		Host:     os.Getenv("MYSQL_DB_HOST"),
		User:     os.Getenv("MYSQL_DB_USERNAME"),
		Password: os.Getenv("MYSQL_DB_PASSWORD"),
		Port:     dbPort,
		Database: os.Getenv("MYSQL_DB_NAME"),
	}

	switch client.Type {
	case "KEY":
		dial, err = client.DialWithKeyFile()
	case "PASSWORD":
		dial, err = client.DialWithPassword()
	default:
		panic("unknown ssh type.")
	}
	if err != nil {
		log.Printf("ssh connect error: %s", err.Error())
		return
	}
	// defer dial.Close()

	// 注册ssh代理
	mysql.RegisterDial("mysql+ssh", (&Dialer{client: dial}).Dial)

	db, err := my.New()
	if err != nil {
		log.Printf("mysql connect error: %s", err.Error())
		return
	}

	DB = db

	// val := make(map[string]interface{})
	// if err := DB.Table("users").Where("username = ?", "lscm").Find(&val).Error; err != nil {
	// 	log.Printf("mysql query error: %s", err.Error())
	// 	return
	// }
	// fmt.Println(val)
}

func FindUser(username string, password string, userType string) ([]db_models.Users, error) {
	CheckDatabaseConnection()
	var val []db_models.Users
	err := DB.Table("users").Where("username = ?", username).Where("user_type = ?", userType).Find(&val).Error
	fmt.Println(val[0].Password)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(val[0].Password), []byte(password))
	if err != nil {
		return nil, err
	}
	// fmt.Println(val)
	return val, err
}

func UpdateUser(user db_models.Users, token string, tokenExpire int64) {
	CheckDatabaseConnection()
	timeNow := time.Now().Format("2006-01-02 15:04:05")
	timeExpire := time.Unix(tokenExpire, 0).Format("2006-01-02 15:04:05")
	// map[string]interface{}{"token_expiry_datetime": timeExpire, "lastLogin_datetime": timeNow, "token": token}
	DB.Table("users").Where("user_id = ?", user.UserId).Updates(map[string]interface{}{"token_expiry_datetime": timeExpire, "last_login_datetime": timeNow, "token": token})

	// fmt.Println(val)
}

func FindAllDutyRooms() ([]db_models.Locations, error) {
	CheckDatabaseConnection()
	var val []db_models.Locations
	err := DB.Find(&val).Error
	// fmt.Println(val)
	return val, err
}

func CheckDatabaseConnection() {
	sqlDB, err1 := DB.DB()
	err2 := sqlDB.Ping()
	if err1 != nil || err2 != nil {
		StartMySql()
	}
}
