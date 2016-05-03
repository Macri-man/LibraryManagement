package main

import (
    "bytes"
    "encoding/base64"
    "fmt"
  //  "net/url"
    "encoding/json"

    "github.com/julienschmidt/httprouter"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"

    "github.com/SlyMarbo/gmail"

    "net/http"
    //"io/ioutil"
    "log"
    //"github.com/scorredoira/email"
    //"net/smtp"
    //"io"
    //"os"
    "html/template"
    "strings"

    _ "google.golang.org/api/books/v1"
	  _ "google.golang.org/api/googleapi"

)


var templates = template.Must(template.ParseGlob("/Users/Macri-man/goWorkSpace/src/github.com/Macri-man/LibraryManagement/templateHTML/*"))

type User struct {
  userName string
  email string
  firstName string
  lastName string
  password string
}

type Book struct{
  title string
  description string
  isbn string
  cover_image string
  categories string
  available uint16
  quantity uint16
}

var db *sql.DB

/*
type Person struct {
  FirstName string
  LastName string
}
*/
/*
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

func Registering(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {

    user := &User{
        userName: req.FormValue("Username"),
        email: req.FormValue("Email"),
        firstName: req.FormValue("FirstName"),
        lastName: req.FormValue("LastName"),
        password: req.FormValue("Password")}

        fmt.Println(user)

      stmt, err := db.Prepare("INSERT INTO users (username, email, firstname, lastname, password) VALUES (?,?,?,?,?)")
      if err != nil {
    	   log.Fatal(err)
      }

      result,err1 := stmt.Exec(user.userName,user.email,user.firstName,user.lastName,user.password)
      if err1 != nil {
    	  log.Fatal(err)
      }
      fmt.Println(result)
    http.Redirect(res,req,"/Home",http.StatusFound)
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

func Contact(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
    err := templates.ExecuteTemplate(res, "contactPage", nil)
    if err != nil {
        http.Error(res, err.Error(), http.StatusInternalServerError)
        log.Fatalln(err)
    }
}


func ContactMail(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {

  body := req.FormValue("Description")
  subject := req.FormValue("Subject")
//  name := req.FromValue("Names")

    email := gmail.Compose(subject, body)
     email.From = "username@gmail.com"
     email.Password = "password"

     email.ContentType = "text/html; charset=utf-8"

     email.AddRecipient("macriad@clarkson.edu")
     err := email.Send()
         if err != nil {
             log.Fatal(err)
         }
  http.Redirect(res,req,"/Home",http.StatusFound)
}


func Login(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {


  username := req.FormValue("Username")
  password := req.FormValue("Password")


  stmt, err := db.Prepare("select username, password from users where username = ? and password = ?")
  if err != nil {
	   log.Fatal(err)
  }
  var varify1,varify2 string

  err = stmt.QueryRow(username,password).Scan(&varify1,&varify2)
  if err != nil {
    if err == sql.ErrNoRows {
		    http.Redirect(res,req,"/Home",http.StatusFound)
	     } else {
		    log.Fatal(err)
	  }
  }
  fmt.Println(password == varify2)
  fmt.Println(varify2 == "")
  fmt.Println(password == varify2 && username == "Admin")
  if password == varify2 && username != "Admin" {
    http.Redirect(res, req, "/Search", http.StatusFound)
  }else if varify2 == "" || username == ""{
    http.Redirect(res,req,"/Home",http.StatusFound)
  }else if password == varify2 && username == "Admin"{
    http.Redirect(res,req,"/Admin",http.StatusFound)
  }
}

func Logout(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
    http.Redirect(res,req,"/Home",http.StatusFound)
}

/** ADMINISTRATOR*/

func getAllBooks(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {

    stmt, err := db.Prepare("select * from books")
    if err != nil {
	     log.Fatal(err)
    }
    defer stmt.Close()
    rows, err := stmt.Query()
    if err != nil {
	    log.Fatal(err)
    }
    defer rows.Close()
    books := make([]*Book, 0)
    for rows.Next() {
        newbook := new(Book)
	      err := rows.Scan(&newbook.title,&newbook.description,&newbook.isbn,&newbook.cover_image,&newbook.available,&newbook.quantity)
        if err != nil {
          log.Fatal(err)
        }
        books = append(books,newbook)

    }
    if err = rows.Err(); err != nil {
	    log.Fatal(err)
    }
    fmt.Println(books)
    err = templates.ExecuteTemplate(res, "adminPage", books)
    if err != nil {
        http.Error(res, err.Error(), http.StatusInternalServerError)
        log.Fatalln(err)
    }
}

func addBook(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
    isbn := req.FormValue("AddBookISBN")

    fmt.Println(params.ByName("isbn"))
    fmt.Println(isbn)
    bookISBN := getISBNfromAPI(isbn)
    book := new(Book)
    stmt,err := db.Prepare("Select * books where isbn = ?")
    err = stmt.QueryRow().Scan(&book.title,&book.description,&book.isbn,&book.cover_image,&book.available,&book.quantity,&book.categories)
    if err != nil {
      if err == sql.ErrNoRows {
        stmt,err := db.Prepare("INSERT INTO books (title,description,isbn,cover_image,available,quantity,categories) VALUES (?,?,?,?,?,?,?)")
        if err != nil {
          log.Fatal(err)
        }
        result,err1 := stmt.Exec(&bookISBN.title,&bookISBN.description,&bookISBN.isbn,&bookISBN.cover_image,&bookISBN.available,&bookISBN.quantity,&bookISBN.categories)
        if err1 != nil {
      	  log.Fatal(err)
        }
        fmt.Println(result)
      } else {
          log.Fatal(err)
      }
    }else{
      stmt,err := db.Prepare("UPDATE books SET quantity=book.quantity+1 WHERE isbn=book.isbn")
      if err != nil {
        log.Fatal(err)
      }
      result,err1 := stmt.Exec(&book.isbn)
      if err1 != nil {
        log.Fatal(err)
      }
      fmt.Println(result)
    }
    err = templates.ExecuteTemplate(res, "adminPage", nil)
    if err != nil {
        http.Error(res, err.Error(), http.StatusInternalServerError)
        log.Fatalln(err)
    }
}


type BookFromJSON struct {
	Kind string `json:"kind"`
	TotalItems int `json:"totalItems"`
	Items []struct {
		Kind string `json:"kind"`
		ID string `json:"id"`
		Etag string `json:"etag"`
		SelfLink string `json:"selfLink"`
		VolumeInfo struct {
			Title string `json:"title"`
			Authors []string `json:"authors"`
			Publisher string `json:"publisher"`
			PublishedDate string `json:"publishedDate"`
			Description string `json:"description"`
			IndustryIdentifiers []struct {
				Type string `json:"type"`
				Identifier string `json:"identifier"`
			} `json:"industryIdentifiers"`
			ReadingModes struct {
				Text bool `json:"text"`
				Image bool `json:"image"`
			} `json:"readingModes"`
			PageCount int `json:"pageCount"`
			PrintType string `json:"printType"`
			Categories []string `json:"categories"`
			AverageRating float64 `json:"averageRating"`
			RatingsCount int `json:"ratingsCount"`
			MaturityRating string `json:"maturityRating"`
			AllowAnonLogging bool `json:"allowAnonLogging"`
			ContentVersion string `json:"contentVersion"`
			ImageLinks struct {
				SmallThumbnail string `json:"smallThumbnail"`
				Thumbnail string `json:"thumbnail"`
			} `json:"imageLinks"`
			Language string `json:"language"`
			PreviewLink string `json:"previewLink"`
			InfoLink string `json:"infoLink"`
			CanonicalVolumeLink string `json:"canonicalVolumeLink"`
		} `json:"volumeInfo"`
		SaleInfo struct {
			Country string `json:"country"`
			Saleability string `json:"saleability"`
			IsEbook bool `json:"isEbook"`
		} `json:"saleInfo"`
		AccessInfo struct {
			Country string `json:"country"`
			Viewability string `json:"viewability"`
			Embeddable bool `json:"embeddable"`
			PublicDomain bool `json:"publicDomain"`
			TextToSpeechPermission string `json:"textToSpeechPermission"`
			Epub struct {
				IsAvailable bool `json:"isAvailable"`
			} `json:"epub"`
			Pdf struct {
				IsAvailable bool `json:"isAvailable"`
			} `json:"pdf"`
			WebReaderLink string `json:"webReaderLink"`
			AccessViewStatus string `json:"accessViewStatus"`
			QuoteSharingAllowed bool `json:"quoteSharingAllowed"`
		} `json:"accessInfo"`
		SearchInfo struct {
			TextSnippet string `json:"textSnippet"`
		} `json:"searchInfo"`
	} `json:"items"`
}

func getISBNfromAPI(isbn string) *Book{

  //var url *url.URL

//  url, err := url.Parse("https://www.googleapis.com/books/v1/volumes?")

fmt.Println("HELO ISBMN")

//var apikey = "AIzaSyB6ktH7OiCAzz-hGgUbUwEVst_0mWApEA4";
  //url := fmt.Sprintf("https://www.googleapis.com/books/v1/volumes?q=isbn:%s&key%s",isbn,apikey)

  var isbnn = "9781451648546"; // Steve Jobs book

  var url = "https://www.googleapis.com/books/v1/volumes?q=isbn:" + isbnn;
//  parameters := make(url.Values)
//  parameters.Add("q", "isbn"+isbn)
//  parameters.Add("key", "5")
//  parameters.Add("vegetable", "potato")
//  url += parameters.Encode()
//data := map[string]interface{}{}


r, _ := http.Get(url)

defer r.Body.Close()
/*
body, _ := ioutil.ReadAll(r.Body)
json.Unmarshal(body, &data)
//fmt.Println(data)
fmt.Println(data["items"])
*/
  temp := BookFromJSON{}
  dec := json.NewDecoder(r.Body)
  dec.Decode(&temp)
  fmt.Println("Print Temp")
  fmt.Println(temp)
  fmt.Println("Print Volume")
  fmt.Println(temp.Items[0].VolumeInfo)
  return &Book{
        title: temp.Items[0].VolumeInfo.Title,
        description: temp.Items[0].VolumeInfo.Description,
        isbn: temp.Items[0].VolumeInfo.IndustryIdentifiers[1].Identifier,
        cover_image: temp.Items[0].VolumeInfo.ImageLinks.Thumbnail,
        categories: temp.Items[0].VolumeInfo.Categories[0],
        available: 0,
        quantity:1}

}

func deleteBook(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
    //isbn := req.FormValue("AddBookISBN")

    id :=  req.FormValue("AddBook")

    stmt, err := db.Prepare("DELETE from books where isbn = ?")
    if err != nil {
	     log.Fatal(err)
    }
    defer stmt.Close()

    result, err1 := stmt.Exec(id)
    if err1 != nil {
	     log.Fatal(err1)
    }

    affected, err2 := result.RowsAffected()
    if err2 != nil {
	     log.Fatal(err2)
    }
    fmt.Println(affected)

    if(affected == 0){
      err = templates.ExecuteTemplate(res, "adminPage", id)
      if err != nil {
          http.Error(res, err.Error(), http.StatusInternalServerError)
          log.Fatalln(err)
      }
    }else{
      err = templates.ExecuteTemplate(res, "adminPage", nil)
      if err != nil {
          http.Error(res, err.Error(), http.StatusInternalServerError)
          log.Fatalln(err)
      }
    }

}

//"https://www.googleapis.com/books/v1/volumes?q=isbn:{}&key={}"

func getAllStudents(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {

    stmt, err := db.Prepare("select * from users")
    if err != nil {
	     log.Fatal(err)
    }
    defer stmt.Close()
    rows, err := stmt.Query()
    if err != nil {
	    log.Fatal(err)
    }
    defer rows.Close()
    users := make([]*User, 0)
    for rows.Next() {
        newstudent := new(User)
	      err := rows.Scan(&newstudent.userName,&newstudent.email,&newstudent.firstName,&newstudent.lastName,&newstudent.password)
        if err != nil {
          log.Fatal(err)
        }
        users = append(users,newstudent)

    }

    if err = rows.Err(); err != nil {
	    log.Fatal(err)
    }

    err = templates.ExecuteTemplate(res, "adminPage", users)
    if err != nil {
        http.Error(res, err.Error(), http.StatusInternalServerError)
        log.Fatalln(err)
    }
}

func deleteStudent(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
  //  isbn := req.FormValue("AddBookISBN")


}




/**  GET BOOKS ISBN*/

func GetByISBN(res http.ResponseWriter, req *http.Request, param httprouter.Params) {

    fmt.Println(param)

    isbn := req.FormValue("isbn")
    fmt.Println(isbn)

    stmt,err := db.Prepare("Select * books where isbn = ?")
    if err != nil {
  	   log.Fatal(err)
    }
    book := new(Book)
    err = stmt.QueryRow().Scan(&book.title,&book.description,&book.isbn,&book.cover_image,&book.available,&book.quantity)
    if err != nil {
      if err == sql.ErrNoRows {
  		    http.Redirect(res,req,"/Search",http.StatusFound)
  	     } else {
  		    log.Fatal(err)
  	  }
    }
}


func UpdateByISBN(res http.ResponseWriter, req *http.Request, param httprouter.Params) {

}


/**  GET BOOKS AUTHOR*/



/**  GET BOOKS TITLE*/


/**  GET BOOKS CATEGORY*/


func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    fmt.Fprint(w, "Not protected!\n")
}

func Protected(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    fmt.Fprint(w, "Protected!\n")
}

func Hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
}


func main() {

    var err error
    db, err = sql.Open("mysql","Macri-man@tcp(127.0.0.1:3306)/HCI")
    if err != nil {
      log.Fatal(err)
    }
    defer db.Close()

    //fmt.Printf("%#v\n", templates)
    user := []byte("gordon")
    pass := []byte("secret!")

    router := httprouter.New()
    //router.GET("/", Index)
    router.GET("/Home", Home)

    router.POST("/Login", Login)
    router.GET("/Logout", Logout)

    router.GET("/Register", Register)
    router.POST("/Registering", Registering)

    router.GET("/Search", Search)
    router.GET("/Profile", Profile)
    router.GET("/Admin", Admin)
    router.GET("/Contact", Contact)
    router.GET("/Contact/sendmail", ContactMail)

    router.GET("/Hello/:name", Hello)

    //router.POST("/checkout",checkout)
    //router.POST("/checkin",checkout)

    router.POST("/AdminBooks/isbn", getAllBooks)
    router.GET("/AdminBooks/isbn", addBook)
    router.DELETE("/AdminBooks/isbn", deleteBook)

    router.POST("/AdminStudents/:student", getAllStudents)
    router.DELETE("/AdminStudents/:student", deleteStudent)

    //router.GET("/isbn", GetAllByISBN)
    router.GET("/Search/Books", GetByISBN)
    router.POST("/Search/Books", UpdateByISBN)
    //router.PUT("/isbn/:isbn", AddISBN)
    //router.DELETE("/isbn/:isbn", DeleteISBN)
/*
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

    //router.NotFound = http.FileServer(http.Dir("public"))
    //router.NotFound = http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("/*filepath"))))

    //router.ServeFiles("/css", http.Dir("/var/www"))

    //fs := justFilesFilesystem(http.Dir("css/"))
    //http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir(fs))))
    http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css/"))))

    fmt.Printf("Serving on LocalHost Port 8080\n")
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
    username := req.FormValue("first")
    password := req.FormValue("last")
    fmt.Println("username: ",username)
    fmt.Println("[]byte(username): ", []byte(username))
    fmt.Println("typeOf: ", reflect.TypeOf(username))

    err = tpl.Execute(res, Person{username,password})
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
