package server

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"fmt"
	"sync"
	"time"
	"net/rpc"
)

type Record struct {
	Name string
}

type Records []Record

type View struct {
	Num     int
	Primary string
	Backup  string
	Ack     bool
	PLast   time.Time
	BLast   time.Time
	mu		sync.RWMutex
}

type PingInfo struct {
	Num int
	Id  string
}

type Server struct {
	CurView	View
	Addr string
}

var recds map[string]Record

func (s *Server) Start(addr string) {
	http.ListenAndServe(addr, NewRouter())
}

func (s *Server) Ping(viewservice string) {
	client, err := rpc.DialHTTP("tcp", viewservice)
	if err != nil {
		fmt.Printf("failed to dial rpc server, err: %s", err)
		return
	}
	in := PingInfo{Num:s.CurView.Num, Id: s.Addr}
	var reply View
	err = client.Call("ViewService.Ping", &in, &reply)
	if err != nil {
		fmt.Printf("failed to ping rpc server, err: %s", err)
		return
	}
	s.CurView = reply
	fmt.Println(s.CurView.Num, s.CurView.Primary, s.CurView.Backup)
}

func Start(addr string, viewservice string) {
	s := Server{Addr: addr}
	recds = make(map[string]Record)
	recds["rec1"] = Record{Name: "rec1"}
	recds["rec2"] = Record{Name: "rec2"}

	go func() {
		fmt.Println("Ping")
		c := time.Tick(10 * time.Second)
		for {
				now := <-c
				fmt.Println(now)
				s.Ping(viewservice)	
		}
	} ()
	fmt.Println("Start")
	s.Start(addr)
}

func RecordDelete(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	var rec Record
	if err := json.Unmarshal(body, &rec); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}
	delete(recds, rec.Name)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode("Deleted"); err != nil {
		panic(err)
	}
}

func RecordCreate(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Request Coming")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	var rec Record
	if err := json.Unmarshal(body, &rec); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}
	recds[rec.Name] = rec
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(rec); err != nil {
		panic(err)
	}
}

func RecordShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	recordId := vars["recordId"]
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	v, ok := recds[recordId]
	if ok {
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(v); err != nil {
			panic(err)
		}
	} else {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode("Not exists"); err != nil {
			panic(err)
		}
	}
}

func RecordIndex(w http.ResponseWriter, r *http.Request) {
	records := Records{}
	for _, r := range recds {
		records = append(records, r)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(records); err != nil {
		panic(err)
	}
}
