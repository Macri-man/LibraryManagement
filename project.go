package main;

import (
  "fmt"
  "net/http"
  "html/template"
  "log"
  "reflect"
)

var templates = template.Must(template.ParseGlob("templateHTML/*"))

type Person struct {
  FirstName string
  LastName string
}

func main() {

  mux := http.NewServeMux()
  mux.handle("/",HomeHandler)
  mux.handle("/Search",SearchHandler)
  mux.handle("/Profile",ProfileHandler)
  mux.handle("/Admin",AdminHandler)
/*
  tpl,err :=template.ParseFiles("form.html")
  if err!=nil {
    log.Fatalln(err)
  }

  http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request){
    fName := req.FormValue("first")
    lName := req.FormValue("last")
    fmt.Println("fName: ",fName)
    fmt.Println("[]byte(fName): ", []byte(fName))
    fmt.Println("typeOf: ", reflect.TypeOf(fName))

    err = tpl.Execute(res, Person{fName,lName})
    if err!=nil{
      http.Error(res,err.Error(), 500)
      log.Fatalln(err)
    }
  })
*/


  log.Fatal(http.ListenAndServe(":8080",mux))
}

func HomeHandler(res http.ResponseWriter, req *http.Request) {
    err := templates.ExecuteTemplate(res, "homePage", nil)
    if err != nil {
        http.Error(res, err.Error(), http.StatusInternalServerError)
        log.Fatalln(err)
    }
}

func SearchHandler(res http.ResponseWriter, req *http.Request) {
    err := templates.ExecuteTemplate(res, "searchPage", nil)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        log.Fatalln(err)
    }
}

func ProfileHandler(res http.ResponseWriter, req *http.Request) {
    err := templates.ExecuteTemplate(res, "profilePage", nil)
    if err != nil {
        http.Error(res, err.Error(), http.StatusInternalServerError)
        log.Fatalln(err)
    }
}

func AdminHandler(res http.ResponseWriter, req *http.Request) {
    err := templates.ExecuteTemplate(res, "adminPage", nil)
    if err != nil {
        http.Error(res, err.Error(), http.StatusInternalServerError)
        log.Fatalln(err)
    }
}
