package main

import (
	"crypto/tls"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"time"

	"flag"

	_ "github.com/denisenkom/go-mssqldb"
)

// EventData ...
type EventData struct {
	EventID           string                 `json:"EventID"`
	SupplierID        string                 `json:"SupplierID"`
	EventTypeID       string                 `json:"EventTypeID"`
	FollowOnEvent     string                 `json:"FollowOnEvent"`
	Data              string                 `json:"Data"`
	EndDate           time.Time              `json:"EndDate"`
	DueDate           time.Time              `json:"DueDate"`
	UserID            string                 `json:"UserID"`
	AccountID         string                 `json:"AccountID"`
	Status            string                 `json:"Status"`
	ContactID         string                 `json:"ContactID"`
	Latitude          string                 `json:"Latitude"`
	Longitude         string                 `json:"Longitude"`
	Deleted           bool                   `json:"Deleted"`
	Notes             string                 `json:"Notes"`
	CheckDuplicate    bool                   `json:"CheckDuplicate"`
	Notifications     string                 `json:"Notifications"`
	Opportunity       string                 `json:"Opportunity"`
	Label             string                 `json:"Label"`
	Duration          float64                `json:"Duration"`
	FormItem          string                 `json:"FormItem"`
	AccountName       string                 `json:"AccountName"`
	RepID             string                 `json:"RepID"`
	UserName          string                 `json:"UserName"`
	AccountGroup      string                 `json:"AccountGroup"`
	Manager           string                 `json:"Manager"`
	CallCycleWeek     float64                `json:"CallCycleWeek"`
	Form              map[string]interface{} `json:"Form"`
	AccountConditions map[string]interface{} `json:"accountConditions"`
	Images            string                 `json:"Images"`
}

var supplierid string

func main() {
	flag.StringVar(&supplierid, "supplierid", "DEMO", "used in where clause of sql query")
	db, err := sql.Open("mssql", "server=sqli.rapidtrade.biz;user id=rapidphp;password=rapidpwd;port=1433;database=rapidtrade")
	if err != nil {
		log.Fatal("Open connection failed:", err.Error())
	}
	defer db.Close()

	row, err := db.Query("select top 3 * from events(nolock) where duedate between '2017-01-01' and getDate() and supplierid = '" + supplierid + "'")
	defer row.Close()

	// col, _ := row.Columns()
	// cols, _ := row.ColumnTypes()

	// for i, x := range cols {
	// 	fmt.Printf("Column @ %d(%v): %v; %v\n", i, col[i], x, x.ScanType())
	// }

	for row.Next() {
		var data EventData
		err := row.Scan(&data.EventID, &data.SupplierID, &data.EventTypeID, &data.FollowOnEvent, &data.Data, &data.EndDate,
			&data.DueDate, &data.UserID, &data.AccountID, &data.Status, &data.ContactID, &data.Deleted, &data.Latitude,
			&data.Longitude, &data.Notes, &data.Notifications, &data.Label, &data.Opportunity, &data.Images, &data.Duration)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
		fmt.Printf("Row data: %v\n", data)

		url := "https://api.rapidtrade.biz/rest2/Post?method=usp_event_modify3"
		jsonByte, err := json.Marshal(&data)
		jsonStr := string(jsonByte)

		req, err := http.NewRequest("POST", url, strings.NewReader(jsonStr))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
		req.Header.Set("Authorization", "Basic REVNTzpERU1P") //REVNTzpERU1P//ZGVtbzpERU1P

		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{Transport: tr}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Error: %v", err)
			return
		}
		defer resp.Body.Close()

		fmt.Println("response Status:", resp.Status)
		fmt.Println("response Headers:", resp.Header)

		//body, _ := ioutil.ReadAll(resp.Body)
		//fmt.Println("response Body:", string(body))
	}

	//fmt.Printf("Columns: %v", col)
}
