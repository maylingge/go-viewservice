package viewservice

import (
	"net/rpc"
	"testing"
	"time"
)

func TestViewService(t *testing.T) {
	StartViewService()
	<-time.After(5 * time.Second)

	cases := []struct {
		in   PingInfo
		want View
	}{
		{in: PingInfo{Num: 0, Id: "One"}, want: View{Num: 1, Primary: "One", Backup: "", Ack: false}},
		{in: PingInfo{Num: 0, Id: "Second"}, want: View{Num: 1, Primary: "One", Backup: "", Ack: false}},
		{in: PingInfo{Num: 1, Id: "One"}, want: View{Num: 2, Primary: "One", Backup: "Second", Ack: true}},
		{in: PingInfo{Num: 1, Id: "Second"}, want: View{Num: 2, Primary: "One", Backup: "Second", Ack: true}},

	}

	client, err := rpc.DialHTTP("tcp", ":1234")
	if err != nil {
		t.Fatalf("failed to dial rpc server, err: %s", err)
	}
	var reply View
	for _, c := range cases {
		err := client.Call("ViewService.Ping", &c.in, &reply)
		if err != nil {
			t.Fatalf("failed to ping rpc server, err: %s", err)
		}
		if reply.String() != c.want.String() {
			t.Fatalf("In(%s), Reply(%s), Want(%s)", c.in, reply, c.want)
		}
	}

	<-time.After(1*time.Minute)
	reply = View{}
	err = client.Call("ViewService.Ping", &PingInfo{Num: 2, Id: "Second"}, &reply)
	want := View{Num: 2, Primary: "One", Backup: "Second", Ack: true}
	if reply.String() != want.String() {
		t.Fatalf("Reply(%s), Want(%s)", reply, want)
	}
	
	<-time.After(90*time.Second)
	reply = View{}
	err = client.Call("ViewService.Ping", &PingInfo{Num: 2, Id: "Second"}, &reply)
	want = View{Num: 3, Primary: "Second", Backup: "", Ack: true}
	if reply.String() != want.String() {
		t.Fatalf("Reply(%s), Want(%s)", reply, want)
	}

	reply = View{}
	err = client.Call("ViewService.Ping", &PingInfo{Num: 0, Id: "One"}, &reply)
	want = View{Num: 4, Primary: "Second", Backup: "One", Ack: true}
	if reply.String() != want.String() {
		t.Fatalf("Reply(%s), Want(%s)", reply, want)
	}
}

