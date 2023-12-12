package database

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	"io/ioutil"

	errorHandler "github.com/SoNim-LSCM/TKOH_OMS/errors"
	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/ssh"
	sql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB
var is_init = false

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

func StartMySqlSSH() {

	var (
		dial *ssh.Client
		err  error
	)

	sshPort, err := strconv.Atoi(os.Getenv("SSH_PORT"))
	errorHandler.CheckError(err, "translate string to int in mysql")
	client := SSH{
		Host:     os.Getenv("SSH_HOST"),
		User:     os.Getenv("SSH_USERNAME"),
		Port:     sshPort,
		Password: os.Getenv("SSH_PASSWORD"),
		// KeyFile: "~/.ssh/id_rsa",
		Type: "PASSWORD", // PASSWORD or KEY
	}
	dbPort, err := strconv.Atoi(os.Getenv("MYSQL_DB_PORT"))
	errorHandler.CheckError(err, "translate string to int in mysql")
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
}

func StartMySql(db_connected chan<- bool) {
	host := os.Getenv("MYSQL_DB_HOST")
	user := os.Getenv("MYSQL_DB_USERNAME")
	password := os.Getenv("MYSQL_DB_PASSWORD")
	port := os.Getenv("MYSQL_DB_PORT")
	database := os.Getenv("MYSQL_DB_NAME")
	dsn := user + ":" + password + "@tcp(" + host + ":" + port + ")/" + database + "?charset=utf8mb4&parseTime=True&loc=Local&tls=true"

	db, err := gorm.Open(sql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("mysql connect error: %s", err.Error())
		return
	}

	DB = db
	is_init = true
	db_connected <- true
}

// path to cert-files hard coded
// Most of this is copy pasted from the internet
// and used without much reflection
func createTLSConf() tls.Config {

	rootCertPool := x509.NewCertPool()
	pem, err := ioutil.ReadFile("tbmsuvmdev01_key.pem")
	if err != nil {
		log.Fatal(err)
	}
	if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
		log.Fatal("Failed to append PEM!")
	}
	clientCert := make([]tls.Certificate, 0, 1)

	certs, err := tls.LoadX509KeyPair("cert/client-cert.pem", "cert/client-key.pem")
	if err != nil {
		log.Fatal(err)
	}

	clientCert = append(clientCert, certs)

	return tls.Config{
		RootCAs:            rootCertPool,
		Certificates:       clientCert,
		InsecureSkipVerify: true, // needed for self signed certs
	}
}

func CheckDatabaseConnection() bool {
	if is_init {
		sqlDB, err1 := DB.DB()
		err2 := sqlDB.Ping()
		if err1 != nil || err2 != nil {
			StartMySqlSSH()
		}
	}
	return is_init
}
