package agent

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/getsentry/raven-go"
	"github.com/kilgaloon/leprechaun/api"
)

// WorkersList is default command for agents
func (d *Default) WorkersList(w http.ResponseWriter, r *http.Request) {
	resp := api.WorkersResponse{}

	if d.NumOfWorkers() < 1 {
		resp.Message = "No workers currently active!"
	}

	for name, worker := range d.GetAllWorkers() {
		startedAt := worker.StartedAt.Format(time.UnixDate)
		resp.List = append(resp.List, []string{name, startedAt, worker.WorkingOn})
	}

	w.WriteHeader(http.StatusOK)
	j, err := json.Marshal(resp)
	if err != nil {
		raven.CaptureError(err, nil)
		log.Fatal(err)
	}

	w.Write(j)
}

// KillWorker kills worker by provided name
func (d *Default) KillWorker(w http.ResponseWriter, r *http.Request) {
	resp := api.WorkersResponse{}

	worker, err := d.GetWorkerByName(r.URL.Query()["name"][0])
	if err != nil {
		resp.Message = err.Error()
	} else {
		worker.Kill()
		resp.Message = "Worker killed"
	}

	w.WriteHeader(http.StatusOK)

	j, err := json.Marshal(resp)
	if err != nil {
		raven.CaptureError(err, nil)
		log.Fatal(err)
	}

	w.Write(j)
}
