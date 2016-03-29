package main

import (
    "bytes"
    "encoding/base64"
    "fmt"
    "github.com/julienschmidt/httprouter"
    "net/http"
    "log"
    //"os"
    "html/template"
    "strings"

)

var templates = template.Must(template.ParseGlob("/Users/Macri-man/goWorkSpace/src/github.com/Macri-man/LibraryManagement/templateHTML/home.html"))

/*
type Person struct {
  FirstName string
  LastName string
}

type justFilesFilesystem struct {
    fs http.FileSystem
}

func (fs justFilesFilesystem) Open(name string) (http.File, error) {
    f, err := fs.fs.Open(name)
    if err != nil {
        return nil, err
    }
    return neuteredReaddirFile{f}, nil
}

type neuteredReaddirFile struct {
    http.File
}

func (f neuteredReaddirFile) Readdir(count int) ([]os.FileInfo, error) {
    return nil, nil
}
*/


func BasicAuth(h httprouter.Handle, user, pass []byte) httprouter.Handle {
    return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
        const basicAuthPrefix string = "Basic "

        // Get the Basic Authentication credentials
        auth := r.Header.Get("Authorization")
        if strings.HasPrefix(auth, basicAuthPrefix) {
            // Check credentials
            payload, err := base64.StdEncoding.DecodeString(auth[len(basicAuthPrefix):])
            if err == nil {
                pair := bytes.SplitN(payload, []byte(":"), 2)
                if len(pair) == 2 &&
                    bytes.Equal(pair[0], user) &&
                    bytes.Equal(pair[1], pass) {

                    // Delegate request to the given handle
                    h(w, r, ps)
                    return
                }
            }
        }

        // Request Basic Authentication otherwise
        w.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
        http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
    }
}

func Home(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
    err := templates.ExecuteTemplate(res, "homePage", nil)
    if err != nil {
        http.Error(res, err.Error(), http.StatusInternalServerError)
        log.Fatalln(err)
    }
}

func Register(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
    err := templates.ExecuteTemplate(res, "registerPage", nil)
    if err != nil {
        http.Error(res, err.Error(), http.StatusInternalServerError)
        log.Fatalln(err)
    }
}

func Search(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
    err := templates.ExecuteTemplate(res, "searchPage", nil)
    if err != nil {
        http.Error(res, err.Error(), http.StatusInternalServerError)
        log.Fatalln(err)
    }
}

func Profile(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
    err := templates.ExecuteTemplate(res, "profilePage", nil)
    if err != nil {
        http.Error(res, err.Error(), http.StatusInternalServerError)
        log.Fatalln(err)
    }
}

func Admin(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
    err := templates.ExecuteTemplate(res, "adminPage", nil)
    if err != nil {
        http.Error(res, err.Error(), http.StatusInternalServerError)
        log.Fatalln(err)
    }
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    fmt.Fprint(w, "Not protected!\n")
}

func Protected(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    fmt.Fprint(w, "Protected!\n")
}

func main() {

    fmt.Printf("%#v\n", templates)
    user := []byte("gordon")
    pass := []byte("secret!")

    router := httprouter.New()
    router.GET("/", Index)
    router.GET("/Home", Home)
    router.GET("/Register", Register)
    router.GET("/Search", Search)
    router.GET("/Admin", Admin)
/*
    router.GET("/isbn", GetAllByISBN)
    router.GET("/isbn/:isbn", GetByISBN)
    router.POST("/isbn/:isbn", UpdateISBN)
    router.PUT("/isbn/:isbn", AddISBN)
    router.DELETE("/isbn/:isbn", DeleteISBN)

    router.GET("/author", GetAllByAuthor)
    router.GET("/author/:author", GetByAuthor)
    router.POST("/author/:author", UpdateAuthor)
    router.PUT("/author/:author", AddAuthor)
    router.DELETE("/author/:author", DeleteAuthor)

    router.GET("/title", GetAllByTitle)
    router.GET("/title/:title", GetByTitle)
    router.POST("/title/:title", UpdateTitle)
    router.PUT("/title/:title", AddTitle)
    router.DELETE("/title/:title", DeleteAuthor)

    router.GET("/categories", GetAllByCategory)
    router.GET("/categories/:categories", GetByCategory)
    router.POST("/categories/:categories", UpdateCategory)
    router.PUT("/categories/:categories", AddCategory)
    router.DELETE("/categories/:categories", DeleteCategory)
*/
    router.GET("/protected/", BasicAuth(Protected, user, pass))

    log.Fatal(http.ListenAndServe(":8080", router))
}

/*
func main() {

  http.handle("/",HomeHandler)
  http.handle("/Search/",SearchHandler)
  http.handle("/Profile/",ProfileHandler)
  http.handle("/Admin/",AdminHandler)
  fs := justFilesFilesystem{http.Dir("css/stylesheet")}
  http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir(fs))))

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
*/
