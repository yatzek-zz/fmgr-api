package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/coopernurse/gorp"
	_ "github.com/lib/pq"
	"github.com/yvasiyarov/gorelic"
)

type Player struct {
	Id      int64  `json:"id"`
	Name    string `json:"name"`
	Surname string `json:"surname"`
	Email   string `json:"email"`
}


func main() {
	dbmap := initDb()
	defer dbmap.Db.Close()

	newrelic := newRelicAgent()

	dbHandler := &DbHandler{dbmap: dbmap}

	http.HandleFunc("/availability", newrelic.WrapHTTPHandlerFunc(dbHandler.availabilityFunc))
	http.HandleFunc("/players", newrelic.WrapHTTPHandlerFunc(dbHandler.playersHandleFunc))

	port := os.Getenv("PORT")
	fmt.Println("Listening on port:" + port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic(err)
	}
}

type DbHandler struct {
	dbmap *gorp.DbMap
}

func (h *DbHandler) playersHandleFunc(res http.ResponseWriter, _ *http.Request) {
	jzon, err := json.Marshal(h.getPlayers())
	checkErr(err, "Serialization of players into json failed")
	fmt.Fprintln(res, string(jzon))
}

func (h *DbHandler) availabilityFunc(res http.ResponseWriter, _ *http.Request) {
	error := h.dbmap.Db.Ping()
	checkErr(error, "Error connecting to database")
	fmt.Fprintln(res, string("OK"))
}

func (h *DbHandler) getPlayers() []Player {
	var players []Player
	_, err := h.dbmap.Select(&players, "select id, name, surname, email from players order by id")
	checkErr(err, "Select failed")
	return players
}

func initDb() *gorp.DbMap {
	db, err := sql.Open("postgres", connectionUrl())
	checkErr(err, "sql.Open failed")

	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)

	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}
	dbmap.AddTableWithName(Player{}, "players").SetKeys(true, "Id")

	return dbmap
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}

func connectionUrl() string {
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		panic("DATABASE_URL environment param na specified, exiting!")
	}
	return url + "?sslmode=disable"
}

// configures new relic agent
func newRelicAgent() *gorelic.Agent {
	agent := gorelic.NewAgent()
	agent.Verbose = true
	agent.NewrelicName = "fmgr-api"
	agent.CollectHTTPStat = true
	agent.NewrelicLicense = "6250f7427b4873ef4ece6aba345e4801aa690ec8"
	agent.Run()
	return agent
}
