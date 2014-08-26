package main

import (
	"fmt"
	"net/http"
	"os"
	"github.com/yvasiyarov/gorelic"
)

func main() {

	agent := newRelicAgent()
	http.HandleFunc("/", agent.WrapHTTPHandlerFunc(handler))

	fmt.Println("listening...")
	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		panic(err)
	}
}

// configures new relic agent
func newRelicAgent() (*gorelic.Agent) {
	agent := gorelic.NewAgent()
	agent.Verbose = true
	agent.NewrelicName = "fmgr-api"
	agent.CollectHTTPStat = true
	agent.NewrelicLicense = "4ea69fc41f601a44712b071ab214352bd00087d4"
	agent.Run()
	return agent
}

func handler(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(res, "hello, world")
}
