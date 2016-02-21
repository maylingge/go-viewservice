package viewservice

import (
	"fmt"
	"net/http"
	"net/rpc"
	"time"
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
}

const (
	DEADPING = time.Minute * 2
)

func (v View) ViewOutOfDate() bool {
	if time.Now().Sub(v.PLast) > DEADPING || time.Now().Sub(v.BLast) > DEADPING {
		return true
	}
	return false
}

func (v View) String() string {
	return fmt.Sprintf("View #: %d, Primary: %s, Backup: %s", v.Num, v.Primary, v.Backup)
}

type ViewService struct {
	curView View
}

func GetView() View {
	return views.curView
}

func (v *ViewService) Ping(args *PingInfo, reply *View) error {
	fmt.Println(args.Id, args.Num)
	if v.curView.Num != 0 && args.Num == v.curView.Num {
		if args.Id == v.curView.Primary {
			v.curView.Ack = true
			v.curView.PLast = time.Now()
		} else if args.Id == v.curView.Backup {
			v.curView.BLast = time.Now()
		}
	} else if v.curView.Primary == "" {
		v.curView.Ack = false
		v.curView.Num = v.curView.Num + 1
		v.curView.Primary = args.Id
		v.curView.PLast = time.Now()
	} else if v.curView.Backup == "" {
		if v.curView.Ack {
			v.curView.Backup = args.Id
			v.curView.Num = v.curView.Num + 1
			v.curView.BLast = time.Now()
		}
	}
	*reply = v.curView
	return nil
}

func StartViewService() {
	views =	&ViewService{curView: View{Num: 0, Primary: "", Backup: "", Ack: false, PLast: time.Now(), BLast: time.Now()}}
	go func() {
		rpc.Register(views)
		rpc.HandleHTTP()
		err := http.ListenAndServe(":1234", nil)
		if err != nil {
			fmt.Println(err.Error())
		}
	}()
	go views.update()
	fmt.Println("Listen done")
}

func (v *ViewService) update() {
	if !v.curView.ViewOutOfDate() {
		return
	}

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
}
