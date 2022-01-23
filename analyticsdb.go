package main

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"fmt"
	"io/ioutil"
	"net/http"

	// "github.com/aws/aws-sdk-go/aws/session"
	// "github.com/aws/aws-sdk-go/service/rds/rdsutils"
	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

func RegisterRDSMysqlCerts(c *http.Client) error {
	resp, err := c.Get("https://s3.amazonaws.com/rds-downloads/rds-combined-ca-bundle.pem")
	if err != nil {
		return err
	}


	pem, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	rootCertPool := x509.NewCertPool()
	if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
		return err
	}

	err = mysql.RegisterTLSConfig("rds", &tls.Config{RootCAs: rootCertPool, InsecureSkipVerify: true})
	if err != nil {
		return err
	}
	return nil
}

func ConnectToSqlDb() *sql.DB {
	// FILL THESE OUT:
	host := "url-analytics.c9iozcypz2w0.ap-south-1.rds.amazonaws.com"
	user := "admin"
    pass := "NfkD4EomFl5vQZXcmul0"
	dbName := "url_analytics"
	// sess := session.Must(session.NewSessionWithOptions(session.Options{
		// SharedConfigState: session.SharedConfigEnable,
	// }))
    // creds := sess.Config.Credentials

	// region := "ap-south-1"

	host = fmt.Sprintf("%s:%d", host, 3306)
	cfg := &mysql.Config{
		User: user,
        Passwd: pass,
		Addr: host,
		Net:  "tcp",
		Params: map[string]string{
			"tls": "rds",
		},
		DBName: dbName,
        AllowNativePasswords: true,
        AllowOldPasswords: true,
	}
	cfg.AllowCleartextPasswords = true

	var err error

	err = RegisterRDSMysqlCerts(http.DefaultClient)
	if err != nil {
		panic(err)
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("ok")

	return db
}

func SetupSqlDbConnection() {
	db = ConnectToSqlDb()
}