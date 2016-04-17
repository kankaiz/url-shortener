package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"hash/crc32"
	"log"
	"math"
	"net/http"
	"os"
	"strings"
)

const (
	// characters used for short-urls
	SYMBOLS = "0123456789abcdefghijklmnopqrsuvwxyzABCDEFGHIJKLMNOPQRSTUVXYZ"
	BASE    = uint32(len(SYMBOLS))

	//DB parameters
	DB_NAME    = "broadsheet"
	DB_USER    = "iamkk"
	DATASOURCE = "user=" + DB_USER + " dbname=" + DB_NAME + " sslmode=disable"
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
			fmt.Fprintf(w, "no url to shorten"+"\n")
		}
	case "POST":
		if u := r.FormValue("url"); u != "" {
			log.Println(u)
			//s := encodeURL(u)
			s := postURL(u)
			log.Println(s)

			fmt.Fprintf(w, "192.168.15.107:3008/"+s+"\n")
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
	db, err := sql.Open("postgres", DATASOURCE)
	checkErr(err)
	err = db.QueryRow("SELECT url FROM shorturl WHERE surl = $1", short).Scan(&url)

	if err == sql.ErrNoRows {
		log.Fatal("No Results Found")
	}

	return string(url)
}

func postURL(url string) (short string) {
	temp := crc32.ChecksumIEEE([]byte(url))
	short = Encode(temp)

	//insert into postgres
	db, err := sql.Open("postgres", DATASOURCE)
	checkErr(err)
	_, err = db.Exec(`INSERT INTO shorturl(surl, url)
		SELECT $1,$2
		WHERE NOT EXISTS
		(SELECT surl FROM shorturl WHERE surl = $1);`, short, url)
	checkErr(err)

	return short
}

func Encode(number uint32) string {
	rest := number % BASE
	result := string(SYMBOLS[rest])
	if number-rest != 0 {
		newnumber := (number - rest) / BASE
		result = Encode(newnumber) + result
	}
	return result
}

func Decode(input string) uint32 {
	const floatbase = float64(BASE)
	l := len(input)
	var sum int = 0
	for index := l - 1; index > -1; index -= 1 {
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
		log.Fatal(err)
		panic(err)
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3008"
	}
	http.HandleFunc("/favicon.ico", handlerIcon)
	http.HandleFunc("/", handler)
	log.Println("Server started: http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
