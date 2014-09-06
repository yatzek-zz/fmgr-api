package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"fmt"
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
	newrelic := newRelicAgent()
	http.HandleFunc("/players", newrelic.WrapHTTPHandlerFunc(playersHandleFunc))

	port := os.Getenv("PORT")
	fmt.Println("Listening on port:" + port)
	err := http.ListenAndServe(":" + port, nil)
	if err != nil {
		panic(err)
	}
}

func playersHandleFunc(res http.ResponseWriter, req *http.Request) {
	jzon, err := json.Marshal(getPlayers())
	checkErr(err, "Serialization of players into json failed")
	fmt.Fprintln(res, string(jzon))
}

func getPlayers() []Player {
	dbmap := initDb()
	defer dbmap.Db.Close()

	var players []Player
	_, err := dbmap.Select(&players, "select id, name, surname, email from players order by id")
	checkErr(err, "Select failed")
	return players
}

func initDb() *gorp.DbMap {
	db, err := sql.Open("postgres", connectionUrl())
	checkErr(err, "sql.Open failed")

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
