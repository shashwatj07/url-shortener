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

// Database object
var db *sql.DB

// Generate certificates for authentication in RDS SQL database
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

// Connect to database using the given config
func ConnectToSqlDb() *sql.DB {
	host := "url-analytics.c9iozcypz2w0.ap-south-1.rds.amazonaws.com"
	user := "admin"
    pass := "NfkD4EomFl5vQZXcmul0"
	dbName := "url_analytics"

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
		log.Println("expiry date out of range")
	} else {
		TTL := time.Now().AddDate(0, 0, days).Format("2006-01-02 15:04:05") 
		SQLStatement := "INSERT INTO urls (ShortUrl,TTL) VALUES ('" + shortUrl + "','" + TTL + "');"
		_,err := db.Query(SQLStatement)
		if err!=nil {
			log.Println(err.Error())
		}
	}
}

// Increment the count for the date when a url is visited
func incrementRedirCount(shortUrl string) {
	curTime := time.Now().Format("2006-01-02")
	SQLStatement := "SELECT COUNT(*) FROM CountPerDay WHERE ShortUrl='"+shortUrl+"' AND Date='"+ curTime +"';"
	resp,err := db.Query(SQLStatement)
	if err!=nil {
		log.Println(err.Error())
	} else{
		var count int
		for resp.Next(){
			resp.Scan(&count)
			log.Printf("%d",count)
		}
		if count==0 {
			SQLStatement = "INSERT INTO CountPerDay (ShortUrl,Date,Clicks) VALUES ('" + shortUrl + "','" + curTime + "',1);"
			_,err = db.Query(SQLStatement)
			if err!=nil {
				log.Println(err.Error())
			} 
		} else {
			SQLStatement = "UPDATE CountPerDay SET Clicks = Clicks + 1 WHERE ShortUrl='"+shortUrl+"' AND Date='"+curTime+"';"
			_,err = db.Query(SQLStatement)
			if err!=nil {
				log.Println(err.Error())
			}
		}
	}
}

// Call ConnectToSqlDb() and store the db instance returned by it
func SetupSqlDbConnection() {
	db = ConnectToSqlDb()
}

// Get analytics from the sql db
func GetAnalyticsFromDb(ShortUrl string) ([]DateCountStruct, error) {
	var dateCountList []DateCountStruct
	SQLStatement := "SELECT Date, Clicks FROM CountPerDay WHERE ShortUrl='"+ShortUrl+"';"
	rows, err := db.Query(SQLStatement)
	if err!=nil {
		log.Println(err.Error())
		return nil, err
	} else {
		var count int
		var date string
		
		for rows.Next(){
			rows.Scan(&date,&count)
			curDateCount := DateCountStruct{date,count}
			dateCountList = append(dateCountList, curDateCount)
		}
	}
	return dateCountList, err
}