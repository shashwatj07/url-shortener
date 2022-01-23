package main

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
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

// Utility function to handle the logic of saving short links
// to the linked Analytics DB Instance along with the expiry time.
func saveUrltoAnalyticsDB(newUrlStruct urlStruct, shortUrl string) {
	var temp urlStruct = newUrlStruct
	temp.ShortURL = shortUrl
	newUrlStruct.ShortURL = HOST_URL + shortUrl
	days := newUrlStruct.ExpDate
	// Expiry date period has to be at least one day
	if days < 1 {
		log.Printf("expiry date out of range")
	} else {
		TTL := time.Now().AddDate(0, 0, days).Format("2006-01-02 15:04:05") 
		SQLStatement := "INSERT INTO urls (ShortUrl,TTL) VALUES ('" + shortUrl + "','" + TTL + "');"
		_,err := db.Query(SQLStatement)
		if err!=nil {
			log.Printf(err.Error())
		}
	}
}

func incrementRedirCount(shortUrl string) {
	curTime := time.Now().Format("2006-01-02")
	SQLStatement := "SELECT COUNT(*) FROM CountPerDay WHERE ShortUrl='"+shortUrl+"' AND Date='"+ curTime +"';"
	resp,err := db.Query(SQLStatement)
	if err!=nil {
		log.Printf(err.Error())
	} else{
		var count int
		for resp.Next(){
			resp.Scan(&count)
			log.Printf("%d",count)
		}
		if count==0 {
			SQLStatement = "INSERT INTO CountPerDay (ShortUrl,Date,Clicks) VALUES ('" + shortUrl + "','" + curTime + "',1);"
			resp,err = db.Query(SQLStatement)
			if err!=nil {
				log.Printf(err.Error())
			} 
		} else {
			SQLStatement = "UPDATE CountPerDay SET Clicks = Clicks + 1 WHERE ShortUrl='"+shortUrl+"' AND Date='"+curTime+"';"
			resp,err = db.Query(SQLStatement)
			if err!=nil {
				log.Printf(err.Error())
			}
		}
	}
}

func SetupSqlDbConnection() {
	db = ConnectToSqlDb()
}