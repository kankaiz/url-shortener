package main

import (
	"crypto/sha512"
	"database/sql"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

// characters used for short-urls
const (
	SYMBOLS = "0123456789abcdefghijklmnopqrsuvwxyzABCDEFGHIJKLMNOPQRSTUVXYZ"
	BASE    = uint32(len(SYMBOLS))
)

var (
	db  *sql.DB
	err error
)

func handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		short := r.URL.Path[1:]
		if short != "" {
			//urlStr := decodeURL(short)
			urlStr := getURL(short)
			log.Println(urlStr)
			http.Redirect(w, r, urlStr, http.StatusFound)
		} else {
			fmt.Fprintf(w, "<h1>Input url below</h1>"+
				"<form action=\"/save/\" method=\"POST\">"+
				//"&nbsp<textarea name=\"url\"></textarea><br><br>"+
				"&nbsp<input type=\"text\" name=\"url\"><br>"+
				"&nbsp<input type=\"submit\" value=\"Save\">"+
				"</form>")
		}
	case "POST":
		if u := r.FormValue("url"); u != "" {
			log.Println(u)
			//s := encodeURL(u)

			//validate url start with http
			rHTTP, _ := regexp.Compile("^(http|https)://")
			if !rHTTP.MatchString(u) {
				//set the url start with http as default
				u = "http://" + u
			}

			//validate the url
			_, validURLErr := http.Get(u)
			if validURLErr != nil {
				log.Println(err.Error())
				fmt.Fprintf(w, "invalid url "+u+"\n")
			} else {
				s := postURL(u)
				log.Println(s)
				fmt.Fprintf(w, "<a href=\"http://%s\">%s</a>", "brds.ht/"+s, "brds.ht/"+s)
			}

		}

		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		return
	default:
		http.Error((w), fmt.Sprintf("Unsupported method: %s", r.Method), http.StatusMethodNotAllowed)
	}
}

func getURL(short string) (url string) {
	log.Println(short)
	err = db.QueryRow("SELECT url FROM shorturl WHERE surl = $1", short).Scan(&url)

	if err == sql.ErrNoRows {
		log.Println("No Results Found")
	}

	return string(url)
}

func postURL(url string) (short string) {
	h := sha512.New()
	h.Write([]byte(url))
	bs := h.Sum(nil)

	temp := binary.BigEndian.Uint32(bs)

	short = Encode(temp)

	//insert into postgres
	_, err = db.Exec(`INSERT INTO shorturl(surl, url)
		SELECT $1,$2
		WHERE NOT EXISTS
		(SELECT surl FROM shorturl WHERE surl = $1);`, short, url)
	checkErr(err)

	return short
}

// Encode ...
func Encode(number uint32) string {
	rest := number % BASE
	result := string(SYMBOLS[rest])
	if number-rest != 0 {
		newnumber := (number - rest) / BASE
		result = Encode(newnumber) + result
	}
	return result
}

// Decode ...
func Decode(input string) uint32 {
	const floatbase = float64(BASE)
	l := len(input)
	var sum int
	for index := l - 1; index > -1; index-- {
		current := string(input[index])
		pos := strings.Index(SYMBOLS, current)
		sum = sum + (pos * int(math.Pow(floatbase, float64((l-index-1)))))
	}
	return uint32(sum)
}

//handle favourte icon request by some browsers
func handlerIcon(w http.ResponseWriter, r *http.Request) {}

func checkErr(err error) {
	if err != nil {
		log.Println(err)
		panic(err)
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3008"
	}

	connInfo := fmt.Sprintf(
		"user=%s dbname=%s password=%s host=%s port=%s sslmode=disable",
		"postgres",
		"postgres",
		os.Getenv("DB_ENV_POSTGRES_PASSWORD"),
		os.Getenv("DB_PORT_5432_TCP_ADDR"),
		os.Getenv("DB_PORT_5432_TCP_PORT"),
	)

	db, err = sql.Open("postgres", connInfo)
	checkErr(err)

	for i := 0; i < 5; i++ {
		time.Sleep(time.Duration(i) * time.Second)

		if err = db.Ping(); err == nil {
			log.Println("try to connect db")
			break
		}
		log.Println(err)
	}

	//initialise the DB table
	_, err = db.Exec(
		`create table if not exists shorturl
		(
		  surl character(10) NOT NULL,
		  url text,
		  CONSTRAINT unique_url PRIMARY KEY (surl)
		)`)
	checkErr(err)

	http.HandleFunc("/favicon.ico", handlerIcon)
	http.HandleFunc("/", handler)
	log.Println("Server started: http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
