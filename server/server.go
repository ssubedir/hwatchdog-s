package server

import (
	"context"
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/joho/godotenv/autoload"
)

type BServer struct {
	serv http.Server
}

// Setup Http Server
func NewBServer() BServer {
	sm := mux.NewRouter()
	getR := sm.Methods(http.MethodGet).Subrouter()
	getR.HandleFunc("/", Index)
	getR.HandleFunc("/status", Status)
	sm.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("template/static")))) // serve static files

	s := http.Server{
		Addr:         ":" + os.Getenv("SERVER_PORT"), // configure the bind address
		Handler:      sm,                             // set the default handler
		ReadTimeout:  5 * time.Second,                // max time to read request from the client
		WriteTimeout: 10 * time.Second,               // max time to write response to the client
		IdleTimeout:  120 * time.Second,              // max time for connections using TCP Keep-Alive
	}

	bs := BServer{s}

	return bs

}

// Start server
func (s *BServer) Start() {
	go func() {
		log.Println("Starting server on port " + os.Getenv("SERVER_PORT"))
		err := s.serv.ListenAndServe()
		if err != nil {
			log.Fatal("Error starting orbiter service: %s\n", err)
			os.Exit(1)
		}
	}()
}

// Close Http Server
func (s *BServer) Close() {
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.serv.Shutdown(ctx)
}

func Status(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hwatchdog working.")
}

// Fetch Prices of all cpu
func Index(w http.ResponseWriter, r *http.Request) {

	q := `SELECT _price, _time
	FROM hw_cpu hc 
	LEFT JOIN hw_cpu_amazon hca 
	ON hc.AmazonId = hca. _asin
	WHERE hc.AmazonId = (SELECT hc2.AmazonId from hw_cpu hc2 where hc2.Name = ?) AND hca. _country = 1 
	ORDER BY hca. _time DESC
	LIMIT 1`

	maria, err := sql.Open("mysql", os.Getenv("USER")+":"+os.Getenv("PASSWORD")+"@("+os.Getenv("HOST")+":"+os.Getenv("PORT")+")/"+os.Getenv("DBNAME")+"?parseTime=true")
	if err != nil {
		log.Fatalln(err)
	}
	defer maria.Close()

	cpu := []CPU{}

	rows, _ := maria.Query("SELECT * FROM hw_cpu")

	for rows.Next() {

		c := CPU{}

		err = rows.Scan(&c.Id, &c.Name, &c.AmazonId)
		if err != nil {
			log.Panic(err.Error())
		}

		cpu = append(cpu, c)
	}

	data := CPUS{Cpus: cpu}

	for i, b := range data.Cpus {

		rows, _ := maria.Query(q, b.Name)

		for rows.Next() {

			err = rows.Scan(&data.Cpus[i].Price, &data.Cpus[i].Time)
			if err != nil {
				log.Panic(err.Error())
			}

		}
	}

	t, _ := template.ParseFiles("template/index.html")
	t.Execute(w, data)
}

type CPUS struct {
	Cpus []CPU
}

type PD struct {
	Price float32   `db:"_price""`
	Time  time.Time `db:"_time""`
}

type CPU struct {
	Id       int       `db:"Id"`
	Name     string    `db:"Name"`
	AmazonId string    `db:"AmazonId"`
	Price    float64   `db:"_price""`
	Time     time.Time `db:"_time""`
}
