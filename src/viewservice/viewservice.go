package viewservice

import (
	"fmt"
	"net/http"
	"net/rpc"
	"time"
	"sync"
)

var views *ViewService

type PingInfo struct {
	Num int
	Id  string
}

type View struct {
	Num     int
	Primary string
	Backup  string
	Ack     bool
	PLast   time.Time
	BLast   time.Time
	mu		sync.RWMutex
}

const (
	DEADPING = time.Minute * 2
)

func (v View) ViewOutOfDate() bool {
	if ((v.Primary != "") && time.Now().Sub(v.PLast) > DEADPING) || ((v.Backup != "") && time.Now().Sub(v.BLast) > DEADPING) {
		return true
	}
	return false
}

func (v View) String() string {
	return fmt.Sprintf("View #: %d, Primary: %s, Backup: %s, Ack: %v", v.Num, v.Primary, v.Backup, v.Ack)
}

func (v *View) SetPrimary(id string) {
	v.Primary = id
	v.PLast = time.Now()
	v.Num += 1
}

func (v *View) SetBackup(id string) {
	v.Num += 1
	v.Backup = id
	v.BLast = time.Now()
}

type ViewService struct {
	curView View
	pServer string
}

func GetView() View {
	return views.curView
}

func (v *ViewService) Ping(args *PingInfo, reply *View) error {
	if v.curView.Num == 0 {
		v.curView.SetPrimary(args.Id)
		*reply = v.curView
		return nil
	}

	if args.Id == v.curView.Primary {
		v.curView.Ack = true
		v.curView.PLast = time.Now()
		if v.pServer != "" && (v.curView.Backup == "") {
			v.curView.SetBackup(v.pServer)
			v.pServer = ""
		}
		*reply = v.curView
		return nil
	}

	if args.Id == v.curView.Backup {
		v.curView.BLast = time.Now()
		*reply = v.curView
		return nil
	}

	if v.curView.Backup == "" && v.curView.Ack {
		v.curView.SetBackup(args.Id)
		*reply = v.curView
	} else {
		v.pServer = args.Id
		*reply = v.curView

	}

	return nil
}

func StartViewService(addr string) {
	views = &ViewService{curView: View{Num: 0, Primary: "", Backup: "", Ack: false, PLast: time.Now(), BLast: time.Now()}}
	go func() {
		rpc.Register(views)
		rpc.HandleHTTP()
		err := http.ListenAndServe(addr, nil)
		if err != nil {
			fmt.Println(err.Error())
		}
	}()
	go views.update()
	fmt.Println("Listen done")
}

func (v *ViewService) update() {
	for {
		if !v.curView.ViewOutOfDate() {
			continue
		}
		v.curView.mu.Lock()
		
		if time.Now().Sub(v.curView.PLast) > DEADPING {
			v.curView.Primary = ""
			v.curView.PLast = time.Now()
		}

		if time.Now().Sub(v.curView.BLast) > DEADPING {
			v.curView.Backup = ""
			v.curView.BLast = time.Now()
		} else {
			v.curView.Primary = v.curView.Backup
			v.curView.Backup = ""
			v.curView.Ack = false
			v.curView.BLast = time.Now()
			v.curView.PLast = time.Now()
		}

		if v.pServer != "" {
			if v.curView.Primary == "" {
				v.curView.Primary = v.pServer
				v.curView.Ack = false
				v.curView.PLast = time.Now()
			} else if v.curView.Backup == "" {
				v.curView.Backup = v.pServer
				v.curView.BLast = time.Now()
			}
			v.pServer = ""
		}
		v.curView.Num += 1
		v.curView.mu.Unlock()
		fmt.Println("Update done")
		fmt.Println(v.curView)
	}
}
