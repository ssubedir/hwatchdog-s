package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gocolly/colly"
	_ "github.com/joho/godotenv/autoload"
	"gopkg.in/robfig/cron.v2"
	"hwatchdog-s/server"
	"log"
	"math"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
)

type DB struct{
	User string
	Password string
	Port string
	Host string
	DBName string
}

func main() {

	log.Println("Starting hwatchdog scrapper service.")


	// start http server

	s := server.NewBServer()
	s.Start()


	// new crone
	cron := cron.New()

	// add jobs
	cron.AddFunc("@every 1h0m0s", ScrapeAmazon) // every hour
	cron.AddFunc("@every 1h0m0s", ScrapeAmazonCA) // every hour
	cron.Start()

	//Gracefully Stop
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)
	sig := <-c
	log.Println("Got signal: ", sig, ", exiting hwatchdog service.")

	cron.Stop()

	s.Close()

}

type CPU struct {
	Id       int    `db:"Id"`
	Name     string `db:"Name"`
	AmazonId string `db:"AmazonId"`
}


// AMAZON.COM
func ScrapeAmazon(){
	maria, err := sql.Open("mysql", os.Getenv("USER")+":"+os.Getenv("PASSWORD")+"@("+os.Getenv("HOST")+":"+os.Getenv("PORT")+")/"+os.Getenv("DBNAME")+"?parseTime=true")
	if err != nil {
		log.Fatalln(err)
	}
	defer maria.Close()

	// Fetch Cpu ASIN from DB
	cpu := []CPU{}


	rows, _ := maria.Query("SELECT * FROM hw_cpu")

	for rows.Next() {

		c:= CPU{}

		err = rows.Scan(&c.Id,&c.Name,&c.AmazonId)
		if err != nil {
			log.Panic(err.Error())
		}

		cpu = append(cpu,c)
	}

	log.Println(cpu)

	coll1 := colly.NewCollector()
	coll1.OnHTML("#price_inside_buybox", func(e *colly.HTMLElement) {

		a := regexp.MustCompile(`\/dp\/`)
		path := e.Request.URL.Path

		asin := a.Split(path, 2)[1]

		price, _ := ParseFloat(strings.TrimPrefix(strings.TrimSuffix(e.Text, "\n"), "\n"))

		fmt.Println("ASIN - " + asin + " USD $",price)

		// save price to database

		stmt, err := maria.Prepare("INSERT INTO hwatchdog.hw_cpu_amazon (_asin, _price, _country,_time) VALUES(?,?,?,?)")

		if err != nil {
			log.Panic(err.Error())
		}

		_, err = stmt.Exec(asin, price,1,time.Now())

		if err!= nil{
			log.Panic(err.Error())
		}

	})
	for _, c := range cpu {
		coll1.Visit("https://www.amazon.com/dp/" + c.AmazonId)
	}
}
// AMAZON.CA
func ScrapeAmazonCA(){
	maria, err := sql.Open("mysql", os.Getenv("USER")+":"+os.Getenv("PASSWORD")+"@("+os.Getenv("HOST")+":"+os.Getenv("PORT")+")/"+os.Getenv("DBNAME")+"?parseTime=true")
	if err != nil {
		log.Fatalln(err)
	}
	defer maria.Close()

	// Fetch Cpu ASIN from DB
	cpu := []CPU{}


	rows, _ := maria.Query("SELECT * FROM hw_cpu")

	for rows.Next() {

		c:= CPU{}

		err = rows.Scan(&c.Id,&c.Name,&c.AmazonId)
		if err != nil {
			log.Panic(err.Error())
		}

		cpu = append(cpu,c)
	}

	log.Println(cpu)

	coll2 := colly.NewCollector()
	coll2.OnHTML("#price_inside_buybox", func(e *colly.HTMLElement) {

		a := regexp.MustCompile(`\/dp\/`)
		path := e.Request.URL.Path

		asin := a.Split(path, 2)[1]

		price, _ := ParseFloat(strings.TrimPrefix(strings.TrimSuffix(e.Text, "\n"), "\n"))
		_ = price
		fmt.Println("ASIN - " + asin + " CAD $",price)

		stmt, err := maria.Prepare("INSERT INTO hwatchdog.hw_cpu_amazon (_asin, _price, _country,_time) VALUES(?,?,?,?)")

		if err != nil {
			log.Panic(err.Error())
		}

		_, err = stmt.Exec(asin, price,2,time.Now())

		if err!= nil{
			log.Panic(err.Error())
		}
	})
	for _, c := range cpu {
		coll2.Visit("https://www.amazon.ca/dp/" + c.AmazonId)
	}
}
// Parse Price
func ParseFloat(strr string) (float64, error) {

	var re = regexp.MustCompile(`([A-z]*\$\s*)`)
	str := re.ReplaceAllString(strr, "")

	stre:=strings.Split(str, "")

	//fmt.Println(stre[0])

	r := []rune(stre[0])
	if(!unicode.IsNumber(r[0])){
		s := []rune(str)
		res := delChar(s, 0)
		str = string(res)
	}

	//fmt.Println("{{" + str + "}}")

	val, err := strconv.ParseFloat(str, 64)
	if err == nil {
		return val, nil
	}

	//Some number may be seperated by comma, for example, 23,120,123, so remove the comma firstly
	str = strings.Replace(str, ",", "", -1)

	//Some number is specifed in scientific notation
	pos := strings.IndexAny(str, "eE")
	if pos < 0 {
		return strconv.ParseFloat(str, 64)
	}

	var baseVal float64
	var expVal int64

	baseStr := str[0:pos]
	baseVal, err = strconv.ParseFloat(baseStr, 64)
	if err != nil {
		return 0, err
	}

	expStr := str[(pos + 1):]
	expVal, err = strconv.ParseInt(expStr, 10, 64)
	if err != nil {
		return 0, err
	}

	return baseVal * math.Pow10(int(expVal)), nil
}
func IsNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}
func delChar(s []rune, index int) []rune {
	return append(s[0:index], s[index+1:]...)
}

