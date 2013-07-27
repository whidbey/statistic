package main

import "fmt"
import "net/http"
import "database/sql"
import _ "github.com/go-sql-driver/mysql"
import "html/template"
//import "strconv"
//import "math/rand"

type MyHandler struct {
	stmt *sql.Stmt
	count int
}

func (h *MyHandler)ServeHTTP(w http.ResponseWriter, r *http.Request) {

	param := r.URL.Path[6:]
	fmt.Println(param)
	if param != "91" && param != "pp" {
		param = "91"
	}
	
	fmt.Println(param)
	
	fmt.Println("count=%d", h.count)
	h.count++
	
	row, err := h.stmt.Query(param);
	if err != nil {
		panic(err.Error())
	}
	
	
	tdate := []string{}
	tonlineno := []int{}
	tnewno := []int{}
	
	td1 := []float32{}
	td3 := []float32{}
	td7 := []float32{}
	td15 := []float32{}
	td30 := []float32{}
	
	var date string
	var onlineno, newno int
	var d1,d3,d7,d15,d30 float32
	
	for row.Next() {
		err := row.Scan(&date,&onlineno,&newno,&d1,&d3,&d7,&d15,&d30)
		if err != nil {
			panic(err.Error())
		}
		//fmt.Fprintf(w, "%s %d %d %f %f %f %f %f\n", date,onlineno,newno,d1,d3,d7,d15,d30)
		tdate = append(tdate, date)
		tonlineno = append(tonlineno, onlineno)
		tnewno = append(tnewno, newno)
		
		if d1 > 0.1 {
			td1 = append(td1, d1)
		}
		
		if d3 > 0.1 {
			td3 = append(td3, d3)
		}
		
		if d7 > 0.1 {
			td7 = append(td7, d7)
		}
		
		if d15 > 0.1 {
			td15 = append(td15, d15)
		}
		
		if d30 > 0.1 {
			td30 = append(td30, d30)
		}
	}

	t, _ := template.ParseFiles("./templ.html")
    t.Execute(w, map[string]interface{}{"date":tdate, "onlineno":tonlineno, "newno":tnewno, "d1":td1, "d3":td3, "d7":td7, "d15":td15, "d30":td30})
}

func main() {
	db, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/test")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	
	// Prepare statement for inserting data
	stmtIns, err := db.Prepare("INSERT INTO stat(date,servname,onlineno,newno,d1,d3,d7,d15,d30) VALUES( ?,?,?,?,?,?,?,?,? )") // ? = placeholder
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmtIns.Close() // Close the statement when we leave main() / the program terminates

	// Prepare statement for reading data
	stmtOut, err := db.Prepare("SELECT date,onlineno,newno,d1,d3,d7,d15,d30 FROM stat where servname=?")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmtOut.Close()

/*
	// Insert square numbers for 0-24 in the database
	for i := 0; i < 25; i++ {
		_, err = stmtIns.Exec(strconv.Itoa(i+1),"91",rand.Intn(1000),rand.Intn(1000),rand.Float32()*100.0,rand.Float32()*100.0,rand.Float32()*100.0,rand.Float32()*100.0,rand.Float32()*100.0) // Insert tuples (i, i^2)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
	}
*/
	
	handler := &MyHandler{stmt:stmtOut}

	http.Handle("/game/", handler)
	http.ListenAndServe(":8080", nil)
	fmt.Println("ok")
}
