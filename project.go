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
    //"net/smtp"
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
  UserName string
  Email string
  FirstName string
  LastName string
  Password string
}

type Book struct{
  Title string
  Description string
  Isbn string
  Cover_image string
  Categories string
  Available uint16
  Quantity uint16
  Author string
}

func (b Book) Testavaibility() bool {
  return b.Available == 0
}

func (b Book) Equalreturn() uint16 {
  return 0
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

var db *sql.DB
var Globalusername string
var Globalpassword string

var nouserfound bool


type HomeMessage struct {
  Valid bool
  Message string
  Username string
}

type SearchMessage struct {
  Valid bool
  Message string
  Username string
}

type AdminMessage struct {
  Valid bool
  Message string
  Username string
}

type Adminbooks struct{
  Isbn string
  Message string
  Valid bool
}

var homemessage HomeMessage

var searchmessage SearchMessage

var adminbooks Adminbooks


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
    fmt.Println("HOMEPAGE")
    fmt.Println(homemessage)
    fmt.Println("%+v\n",homemessage)
    err := templates.ExecuteTemplate(res, "homePage", homemessage)
    if err != nil {
        http.Error(res, err.Error(), http.StatusInternalServerError)
        log.Fatalln(err)
    }
    homemessage = HomeMessage{false,"", ""}
}

func Register(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
  fmt.Println("REGISTERPAGE")
fmt.Println("%+v\n",homemessage)
    err := templates.ExecuteTemplate(res, "registerPage", homemessage)
    if err != nil {
        http.Error(res, err.Error(), http.StatusInternalServerError)
        log.Fatalln(err)
    }
    homemessage = HomeMessage{false,"", ""}
    fmt.Println("%+v\n",homemessage)

}

func Registering(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {

    user := &User{
        UserName: req.FormValue("Username"),
        Email: req.FormValue("Email"),
        FirstName: req.FormValue("FirstName"),
        LastName: req.FormValue("LastName"),
        Password: req.FormValue("Password")}

        fmt.Println(user)

        stmt, err := db.Prepare("SELECT username from users where username = ?")
        if err != nil {
          log.Fatal(err)
        }
        var duser string
        err = stmt.QueryRow(user.UserName).Scan(&duser)
        if err != nil {
          if err == sql.ErrNoRows {
            stmt, err := db.Prepare("INSERT INTO users (username, email, firstname, lastname, password) VALUES (?,?,?,?,?)")
            if err != nil {
               log.Fatal(err)
            }
            result,err1 := stmt.Exec(user.UserName,user.Email,user.FirstName,user.LastName,user.Password)
            if err1 != nil {
              log.Fatal(err1)
            }
            fmt.Println(result)
            homemessage = HomeMessage{true,"User has been Registered, username ", duser}
            http.Redirect(res,req,"/Home",http.StatusFound)
          } else {
              log.Fatal(err)
          }
        } else {
          homemessage = HomeMessage{false,"User is already Registered, username ", duser}
          http.Redirect(res,req,"/Register",http.StatusFound)
        }
}



func Search(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
    err := templates.ExecuteTemplate(res, "searchPage", searchmessage)
    if err != nil {
        http.Error(res, err.Error(), http.StatusInternalServerError)
        log.Fatalln(err)
    }
    searchmessage = SearchMessage{false,"", ""}

}

func Profile(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
    err := templates.ExecuteTemplate(res, "profilePage", nil)
    if err != nil {
        http.Error(res, err.Error(), http.StatusInternalServerError)
        log.Fatalln(err)
    }
}

func Admin(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
    fmt.Println("AMDINPAGE")
    fmt.Println("%+v\n",adminbooks)
    err := templates.ExecuteTemplate(res, "adminPage", adminbooks)
    if err != nil {
        http.Error(res, err.Error(), http.StatusInternalServerError)
        log.Fatalln(err)
    }
    adminbooks = Adminbooks{"","",false}
}

func Contact(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
    err := templates.ExecuteTemplate(res, "contactPage", nil)
    if err != nil {
        http.Error(res, err.Error(), http.StatusInternalServerError)
        log.Fatalln(err)
    }
}

/* Mailing */

func ContactMail(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
fmt.Println("CONTACT MAIL")
  body := req.FormValue("Description")
  subject := req.FormValue("Subject")
  //name := req.FormValue("Name")
//  name := req.FromValue("Names")

  //fmt.Println(name)
  fmt.Println(subject)
  fmt.Println(body)

  stmt, err := db.Prepare("SELECT email,firstname,lastname,password from users where username = ?")
  if err != nil {
    log.Fatal(err)
  }
  var emailfrom string
  var firstname string
  var lastname string
  var password string
  err = stmt.QueryRow(Globalusername).Scan(&emailfrom,&firstname,&lastname,&password)
  if err != nil {
    if err == sql.ErrNoRows {
      fmt.Println(err)
    } else {
        log.Fatal(err)
    }
  }
  fmt.Println(emailfrom)

  subject = subject + " " + " from" + " " + firstname + " " + lastname


fmt.Println("set emails")
     // Send the email body.


    email := gmail.Compose(subject, body)
    email.From = emailfrom
    email.Password = password

    email.ContentType = "text/html; charset=utf-8"

    email.AddRecipient("macriad@clarkson.edu")
    err = email.Send()
    if err != nil {
        log.Fatal(err)
    }

/*
    auth := smtp.PlainAuth(
       "",
       "macriad@clarkson.edu",
       "!Entendore2012",
       "smtp.gmail.com",
   )
   // Connect to the server, authenticate, set the sender and recipient,
   // and send the email all in one step.
   err = smtp.SendMail(
       "smtp.gmail.com:587",
       auth,
       "macriad@clarkson.edu",
       []string{"macriad@clarkson.edu"},
       []byte("This is the email body."),
   )
   if err != nil {
       log.Fatal(err)
   }*/
   searchmessage = SearchMessage{true,"Email Sent", " "}
  http.Redirect(res,req,"/Search",http.StatusFound)
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
        homemessage = HomeMessage{false,"Entry for Username or Password not found", ""}
		    http.Redirect(res,req,"/Home",http.StatusFound)
	     } else {
		    log.Fatal(err)
	  }
  }
  fmt.Println(password == varify2)
  fmt.Println(varify2 == "")
  fmt.Println(password == varify2 && username == "Admin")
  if password == varify2 && username != "Admin" {
    Globalusername = username
    searchmessage =  SearchMessage{false,"", ""}
    http.Redirect(res, req, "/Search", http.StatusFound)
  }else if varify2 == "" || username == ""{
    homemessage = HomeMessage{false,"Incorrect Entry for Username or Password", ""}
    http.Redirect(res,req,"/Home",http.StatusFound)
  }else if password == varify2 && username == "Admin"{
    http.Redirect(res,req,"/Admin",http.StatusFound)
  }
}

func Logout(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
    homemessage = HomeMessage{true,"Logged Out", Globalusername}
    searchmessage = SearchMessage{false,"", ""}
    adminbooks = Adminbooks{"","",false}
    Globalusername = ""
    http.Redirect(res,req,"/Home",http.StatusFound)
}

type Check struct{
  Username string
  Isbn string
}

/* Checkin checkout */

func Checkout(res http.ResponseWriter, req *http.Request, params httprouter.Params){
fmt.Println("STARTED CHECKOUT")
  isbn := params.ByName("isbn")

  stmt, err := db.Prepare("SELECT isbn,username from checkin where isbn = ? AND username = ?")
  if err != nil {
     log.Fatal(err)
  }
  checkin := new(Check)
  err = stmt.QueryRow(isbn,Globalusername).Scan(&checkin.Username,&checkin.Isbn)
  if err != nil {
    if err == sql.ErrNoRows {
      stmt,err := db.Prepare("INSERT INTO checkin (isbn,username) VALUES (?,?)")
      if err != nil {
        log.Fatal(err)
      }
      result,err1 := stmt.Exec(&isbn,&Globalusername)
      if err1 != nil {
        log.Fatal(err1)
      }
      stmt,err = db.Prepare("UPDATE books SET available = available - 1 WHERE isbn=?")
      if err != nil {
        log.Fatal(err)
      }
      result,err = stmt.Exec(&isbn)
      if err != nil {
        log.Fatal(err)
      }
      fmt.Println(result)
      searchmessage = SearchMessage{true,"The Book is checked-out, ", Globalusername}
    } else {
        log.Fatal(err)
    }
  }else{
    searchmessage = SearchMessage{false,"The Book has been already checkout by this user,  ", Globalusername}
  }

  fmt.Println("FINISHED CHECKOUT")

  http.Redirect(res,req,"/Search",http.StatusFound)

}


func Checkin(res http.ResponseWriter, req *http.Request, params httprouter.Params){
  fmt.Println("STARTED CHECKIN")

  isbn := params.ByName("isbn")
  fmt.Println(isbn)

  stmt, err := db.Prepare("SELECT isbn,username from checkin where isbn = ? AND username = ?")
  if err != nil {
     log.Fatal(err)
  }
  checkin := new(Check)
  err = stmt.QueryRow(isbn,Globalusername).Scan(&checkin.Username,&checkin.Isbn)
  if err != nil {
    if err == sql.ErrNoRows {
      searchmessage = SearchMessage{false,"The Book is not checked-out by this user ", Globalusername}
    } else {
        log.Fatal(err)
    }
  }else{
    stmt,err := db.Prepare("UPDATE books SET available = available + 1 WHERE isbn=?")
    if err != nil {
      log.Fatal(err)
    }
    result,err1 := stmt.Exec(&isbn)
    if err1 != nil {
      log.Fatal(err1)
    }
    fmt.Println(result)
    stmt,err = db.Prepare("DELETE FROM checkin WHERE isbn = ? AND username = ?")
    if err != nil {
      log.Fatal(err)
    }
    result,err = stmt.Exec(&isbn,&Globalusername)
    if err != nil {
      log.Fatal(err)
    }
    fmt.Println(result)
    searchmessage = SearchMessage{true,"The Book is checked-in,  ", Globalusername}
  }

  fmt.Println("FINISHED CHECKIN")

  http.Redirect(res,req,"/Search",http.StatusFound)

}


/* Searching */

func BooksSearchResult(res http.ResponseWriter, req *http.Request, _ httprouter.Params)  {
    search := req.FormValue("SEARCH")
    fmt.Println(search)

    stmt, err := db.Prepare("SELECT title,description,isbn,cover_image,available,quantity,categories,author FROM books where title like CONCAT('%',?,'%') OR isbn like CONCAT('%',?,'%') OR categories like CONCAT('%',?,'%') OR author like CONCAT('%',?,'%')")
    if err != nil {
       log.Fatal(err)
    }
    defer stmt.Close()
    rows, err := stmt.Query(&search,&search,&search,&search)
    if err != nil {
      log.Fatal(err)
    }
    defer rows.Close()
    books := make([]*Book, 0)
    for rows.Next() {
        newbook := new(Book)
        err := rows.Scan(&newbook.Title,&newbook.Description,&newbook.Isbn,&newbook.Cover_image,&newbook.Available,&newbook.Quantity,&newbook.Categories,&newbook.Author)
        if err != nil {
          log.Fatal(err)
        }
        books = append(books,newbook)

    }
    if err = rows.Err(); err != nil {
      log.Fatal(err)
    }
    fmt.Println(books)
    for _ , book := range books{
      fmt.Println("%+v\n",book)
    }
    if len(books) == 0 {
      searchmessage =  SearchMessage{false,"There are no search results from search query, ", search}
      http.Redirect(res,req,"/Search",http.StatusFound)
    }

    err = templates.ExecuteTemplate(res, "searchList", books)
    if err != nil {
        http.Error(res, err.Error(), http.StatusInternalServerError)
        log.Fatalln(err)
    }

}


/** ADMINISTRATOR*/


type AdminbooksList struct{
  Books []*Book
  Valid bool
}



func getAllBooks(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {

  fmt.Println("GETALLBOOKS")
    stmt, err := db.Prepare("SELECT title,description,isbn,cover_image,available,quantity,categories,author from books")
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
	      err := rows.Scan(&newbook.Title,&newbook.Description,&newbook.Isbn,&newbook.Cover_image,&newbook.Available,&newbook.Quantity,&newbook.Categories,&newbook.Author)
        if err != nil {
          log.Fatal(err)
        }
        books = append(books,newbook)

    }
    if err = rows.Err(); err != nil {
	    log.Fatal(err)
    }
    fmt.Println(books)
    for _ , book := range books{
      fmt.Println("%+v\n",book)
    }
    response := AdminbooksList{Books: books,Valid: false}
    fmt.Println("%+v\n",response)
    //http.Redirect(res,req,"/Admin",http.StatusFound)
    err = templates.ExecuteTemplate(res, "adminList", books)
    if err != nil {
        http.Error(res, err.Error(), http.StatusInternalServerError)
        log.Fatalln(err)
    }
}


func getISBNfromAPI(isbn string) Book{

  //var url *url.URL

//  url, err := url.Parse("https://www.googleapis.com/books/v1/volumes?")

fmt.Println("HELLO ISBN")
fmt.Println(isbn)

//var apikey = "AIzaSyB6ktH7OiCAzz-hGgUbUwEVst_0mWApEA4";
  //url := fmt.Sprintf("https://www.googleapis.com/books/v1/volumes?q=isbn:%s&key%s",isbn,apikey)

  //var isbnn = "9781451648546"; // Steve Jobs book

  var url = "https://www.googleapis.com/books/v1/volumes?q=isbn:" + isbn;

  res, err := http.Get(url)
  if err != nil {
    log.Fatal(err)
  }
  defer res.Body.Close()
  temp := BookFromJSON{}
  dec := json.NewDecoder(res.Body)
  dec.Decode(&temp)
  fmt.Println("Print Temp")
  fmt.Println(temp)
  //fmt.Println("Print Volume")
  //fmt.Println(temp.Items[0].VolumeInfo)

  var realisbn string
  var authors string
  var category string
  if temp.Items[0].VolumeInfo.IndustryIdentifiers[0].Type == "ISBN_13"{
    realisbn = temp.Items[0].VolumeInfo.IndustryIdentifiers[0].Identifier
  }else if temp.Items[0].VolumeInfo.IndustryIdentifiers[1].Type == "ISBN_13" {
    realisbn = temp.Items[0].VolumeInfo.IndustryIdentifiers[1].Identifier
  }else{
    realisbn = isbn
  }

  if temp.Items[0].VolumeInfo.Authors == nil {
    authors = ""
  }else{
    authors = temp.Items[0].VolumeInfo.Authors[0]
  }

  if temp.Items[0].VolumeInfo.Categories == nil {
    category = ""
  }else{
    category = strings.Join(temp.Items[0].VolumeInfo.Categories, " ")
  }

  fmt.Println("END OF API")
  return Book{
        Title: temp.Items[0].VolumeInfo.Title,
        Description: temp.Items[0].VolumeInfo.Description,
        Isbn:  realisbn,
        Cover_image: temp.Items[0].VolumeInfo.ImageLinks.Thumbnail,
        Categories: category,
        Available: 1,
        Quantity:1,
        Author:authors}

}


func addBook(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
    fmt.Println("ADDBOOK")
    isbn := req.FormValue("ADDBOOKISBN")
    fmt.Println("PRINT ISBN")
    fmt.Println(isbn)

    if len([]rune(isbn)) != 13 {
      adminbooks = Adminbooks{isbn,"Invalid ISBN Entry. ISBN needs to be 13 digits, The ISBN entry is ",false}
    }else{

    //fmt.Println(isbn)
    //fmt.Println(res)
    //fmt.Println(req)
    bookISBN := getISBNfromAPI(isbn)
    fmt.Println("PRINT BOOK")
    fmt.Println("%+v\n",bookISBN)
    book := new(Book)
    stmt,err := db.Prepare("SELECT title,description,isbn,cover_image,available,quantity,categories,author from books WHERE isbn = ?")
    if err != nil {
	     log.Fatal(err)
    }
    err = stmt.QueryRow(bookISBN.Isbn).Scan(&book.Title,&book.Description,&book.Isbn,&book.Cover_image,&book.Available,&book.Quantity,&book.Categories,&book.Author)
    if err != nil {
      if err == sql.ErrNoRows {
        stmt,err := db.Prepare("INSERT INTO books (title,description,isbn,cover_image,available,quantity,categories,author) VALUES (?,?,?,?,?,?,?,?)")
        if err != nil {
          log.Fatal(err)
        }
        result,err1 := stmt.Exec(&bookISBN.Title,&bookISBN.Description,&bookISBN.Isbn,&bookISBN.Cover_image,&bookISBN.Available,&bookISBN.Quantity,&bookISBN.Categories,&bookISBN.Author)
        if err1 != nil {
      	  log.Fatal(err1)
        }
        fmt.Println(result)
      } else {
          log.Fatal(err)
      }
      adminbooks = Adminbooks{bookISBN.Isbn,"Inserted Book into Database, The ISBN is",true}
      fmt.Println("adminbooks")
      fmt.Println("%+v\n",adminbooks)
    }else{
      stmt,err := db.Prepare("UPDATE books SET quantity = quantity+1, available = available+1 WHERE isbn=?")
      if err != nil {
        log.Fatal(err)
      }
      result,err1 := stmt.Exec(&bookISBN.Isbn)
      if err1 != nil {
        log.Fatal(err1)
      }
      fmt.Println("result")
      fmt.Println(result)
      adminbooks = Adminbooks{bookISBN.Isbn,"Updated Quantity of Book In Database, The ISBN is ",true}
      fmt.Println("%+v\n",adminbooks)
    }
    fmt.Println("FINISHED ADD BOOK")
    }
    http.Redirect(res,req,"/Admin",http.StatusFound)
}



func deleteBook(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
    //isbn := req.FormValue("AddBookISBN")
    fmt.Println("DELETEBOOK")
    isbn :=  req.FormValue("DELETEBOOKISBN")
    fmt.Println("PRINT ISBN")
    fmt.Println(isbn)

    stmt,err := db.Prepare("SELECT title,description,isbn,cover_image,available,quantity,categories,author from books WHERE isbn = ?")
    if err != nil {
       log.Fatal(err)
    }
    defer stmt.Close()

    book := new(Book)
    err = stmt.QueryRow(&isbn).Scan(&book.Title,&book.Description,&book.Isbn,&book.Cover_image,&book.Available,&book.Quantity,&book.Categories,&book.Author)
    if err != nil {
      if err == sql.ErrNoRows {
        adminbooks = Adminbooks{isbn,"Book is Not Found in Database, input is ", false}
        fmt.Println("%+v\n",adminbooks)
      }else{

      }
    }else{
      if book.Quantity == 1{
      stmt, err := db.Prepare("DELETE FROM books WHERE isbn = ?")
      if err != nil {
  	     log.Fatal(err)
      }
      result, err1 := stmt.Exec(&isbn)
      if err1 != nil {
         log.Fatal(err1)
       }
      affected, err2 := result.RowsAffected()
      if err2 != nil {
  	     log.Fatal(err2)
      }
      fmt.Println(affected)

      adminbooks = Adminbooks{isbn,"Deleted Book From Database, The ISBN is ",true}
      fmt.Println("%+v\n",adminbooks)
      }else{
        adminbooks = Adminbooks{isbn,"Updated Quantity of Book In Database, The ISBN is ",true}
        fmt.Println("%+v\n",adminbooks)
        stmt,err := db.Prepare("UPDATE books SET quantity = quantity - 1, available = available - 1 WHERE isbn = ?")
        if err != nil {
          log.Fatal(err)
        }
        result,err1 := stmt.Exec(&isbn)
        if err1 != nil {
          log.Fatal(err1)
        }
        fmt.Println("result")
        fmt.Println(result)
      }
    }

    fmt.Println("FINISHED DELETE BOOK")
    http.Redirect(res,req,"/Admin",http.StatusFound)
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
	      err := rows.Scan(&newstudent.UserName,&newstudent.Email,&newstudent.FirstName,&newstudent.LastName,&newstudent.Password)
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
  student := req.FormValue("DeleteStudent")
  stmt, err := db.Prepare("DELETE from users WHERE username = ?")
  if err != nil {
     log.Fatal(err)
  }
  defer stmt.Close()

  result, err1 := stmt.Exec(student)
  if err1 != nil {
     log.Fatal(err1)
  }

  affected, err2 := result.RowsAffected()
  if err2 != nil {
     log.Fatal(err2)
  }
  fmt.Println(affected)
  err = templates.ExecuteTemplate(res, "adminPage", affected)
  if err != nil {
      http.Error(res, err.Error(), http.StatusInternalServerError)
      log.Fatalln(err)
  }
}

/**  GET BOOKS ISBN*/

func GetByISBN(res http.ResponseWriter, req *http.Request, param httprouter.Params) {

    fmt.Println(param)

    isbn := req.FormValue("isbn")
    fmt.Println(isbn)

    stmt,err := db.Prepare("Select * From books where isbn = ?")
    if err != nil {
  	   log.Fatal(err)
    }
    book := new(Book)
    err = stmt.QueryRow(isbn).Scan(&book.Title,&book.Description,&book.Isbn,&book.Cover_image,&book.Available,&book.Quantity)
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
    err = db.Ping()
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
    router.POST("/Emailing", ContactMail)

    router.GET("/Hello/:name", Hello)

    router.POST("/BookResults", BooksSearchResult)

    //router.POST("/checkout",checkout)
    //router.POST("/checkin",checkout)

    router.GET("/Checkin/:isbn", Checkin)
    router.GET("/Checkout/:isbn", Checkout)

    router.POST("/ListBooks", getAllBooks)
    router.POST("/ADD", addBook)
    router.POST("/DELETE", deleteBook)

    //router.POST("/AdminStudents/:student", getAllStudents)
    //router.DELETE("/AdminStudents/:student", deleteStudent)

    //router.GET("/Search/Books", GetByISBN)
    //router.POST("/Search/Books", UpdateByISBN)

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
