package main

import
(
	"github.com/gorilla/mux"
	"log"
	"net/http"
	_"encoding/json"
	_"fmt"
	"database/sql"
	_"github.com/go-sql-driver/mysql"
	"html/template"
	"encoding/json"
	"fmt"
)

var db *sql.DB
var err error
var tpl *template.Template


type UserDetail struct{
	ID string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Position string `json:"position,omitempty"`
	Address string `json:"address,omitempty"`
}

func init(){
 tpl = template.Must(template.ParseGlob("*.html"))
}


func OpenTemplate(w http.ResponseWriter, req *http.Request){
	err:=tpl.ExecuteTemplate(w,"index.html",nil)
	checkError(err)
}



func InsertData(w http.ResponseWriter, req *http.Request){

	id:=req.FormValue("ID")
	name:=req.FormValue("Name")
	position:=req.FormValue("Position")
	address:=req.FormValue("Address")

	_,err:=db.Exec(`INSERT INTO emptable VALUES(?,?,?,?)`,id,name,position,address)
	checkError(err)
	error:=tpl.ExecuteTemplate(w,"index.html",nil)
	checkError(error)
}



func DeleteData(w http.ResponseWriter, req *http.Request){
	id:=req.FormValue("ID")
	name:=req.FormValue("Name")

   _,err:= db.Exec(`DELETE FROM emptable WHERE ID = ? AND NAME = ? `,id,name)
   checkError(err)
	error:=tpl.ExecuteTemplate(w,"index.html",nil)
	checkError(error)
}


func FetchData(w http.ResponseWriter, req *http.Request){
	id:=req.FormValue("ID")
	name:=req.FormValue("Name")
	fmt.Println(id, name)
	res,err:= db.Query(`SELECT * FROM emptable WHERE ID = ? AND NAME = ? `,id,name)
	checkError(err)

     var userinfo []UserDetail
     var userinfo2 UserDetail
	for res.Next(){

		err:= res.Scan(&userinfo2.ID,&userinfo2.Name,&userinfo2.Position,&userinfo2.Address)
		userinfo = append(userinfo,userinfo2)
		checkError(err)
	}
	   for _,v:=range userinfo {
		   json.NewEncoder(w).Encode(v.ID)
		   json.NewEncoder(w).Encode(v.Name)
		   json.NewEncoder(w).Encode(v.Position)
		   json.NewEncoder(w).Encode(v.Address)
	   }

}



func UpdateData(w http.ResponseWriter, req *http.Request){
	id:=req.FormValue("ID")
	name:=req.FormValue("Name")

	_,err:= db.Exec(`UPDATE emptable SET NAME = ? WHERE ID = ?`,name,id)
	checkError(err)
	error:=tpl.ExecuteTemplate(w,"index.html",nil)
	checkError(error)
}


func checkError(e error){
	if e!=nil{
		log.Fatalln(e)
	}
}


func main(){
 //Connecting Database
	db,err = sql.Open("mysql","root:password@tcp(127.0.0.1:3306)/sample_db")
	checkError(err)
	defer db.Close()
	err = db.Ping()
	checkError(err)

	//Routing URLs
	router := mux.NewRouter()
	router.HandleFunc("/",OpenTemplate).Methods("GET")
	router.HandleFunc("/insert",InsertData).Methods("POST")
	router.HandleFunc("/query",FetchData).Methods("POST")
	router.HandleFunc("/update",UpdateData).Methods("POST")
	router.HandleFunc("/delete",DeleteData).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080",router))

}
