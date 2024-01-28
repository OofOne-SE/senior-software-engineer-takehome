package dispatcher

import (
	"encoding/json"
	"net/http"
	"senior-software-engineer-takehome/models"
	"senior-software-engineer-takehome/restapi/operations/service"
	"sync"

	"github.com/ascarter/requestid"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/olahol/melody"
)

const chanSize = 100

type wrapper struct {
	event *models.WeatherUnit
}

type Dispatcher struct {
	mel      *melody.Melody
	log      func(string, ...any)
	events   chan wrapper
	isClosed bool
	mtx      sync.Mutex
}

func New(m *melody.Melody, l func(string, ...any)) *Dispatcher {
	d := Dispatcher{
		mel:    m,
		log:    l,
		events: make(chan wrapper, chanSize),
	}
	d.mel.HandleError(
		func(s *melody.Session, msg error) {
			req_id := s.MustGet("request_id").(string)
			d.log("Session closed. Error: %v. Request id: %s", msg, req_id)
			s.Close()
		},
	)
	return &d
}

func (d *Dispatcher) Dispatch() {
	for e := range d.events {
		if e.event == nil {
			continue
		}
		buf, err := json.Marshal(e.event)
		if err != nil {
			d.log("marshal: %s", err)
			continue
		}
		if err := d.mel.Broadcast(buf); err != nil {
			d.log("broadcast: %s", err)
		}
	}
}

func (d *Dispatcher) HandleWS(params service.GetWsParams) middleware.Responder {
	return middleware.ResponderFunc(
		func(rw http.ResponseWriter, _ runtime.Producer) {
			reqID, _ := requestid.FromContext(params.HTTPRequest.Context())
			err := d.mel.HandleRequestWithKeys(
				rw,
				params.HTTPRequest,
				map[string]interface{}{
					"request_id": reqID,
				})
			if err != nil {
				d.log("Error: %v. Request id: %s", err.Error(), reqID)
				http.Error(rw, err.Error(), http.StatusInternalServerError)
			}
		},
	)
}

func (d *Dispatcher) Notify(event *models.WeatherUnit) {
	d.mtx.Lock()
	defer d.mtx.Unlock()

	if d.isClosed {
		return
	}
	d.events <- wrapper{event: event}
}

func (d *Dispatcher) Close() {
	d.mtx.Lock()
	defer d.mtx.Unlock()
	d.isClosed = true
	close(d.events)
}
