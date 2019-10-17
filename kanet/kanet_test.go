package kanet

import (
	"fmt"
	"kanet/base"
	"kanet/game"
	pb "kanet/proto/pbgo"
	"os"
	"os/signal"
	"reflect"
	"strconv"
	"syscall"
	"testing"

	"github.com/golang/protobuf/proto"
)

type User struct {
	ID      int
	Account string
}

func (u *User) String() string {
	return "id:" + strconv.Itoa(u.ID) + ", account:" + u.Account
}

func funa(u *User) {
	println("funa addr:", u)
	u.ID++
}

func TestKanet(t *testing.T) {
	// var i interface{}
	// i = 42
	// println(reflect.ValueOf(i).Type().String())
	// buf := new(bytes.Buffer)
	// encoder := gob.NewEncoder(buf)
	// user := new(User)
	// user.ID = 10000
	// user.Account = "kanon lee"
	// encoder.Encode(user)

	// uv := reflect.ValueOf(user)
	// uu := uv.Elem().Interface().(User)
	// println("====", uu.String())

	// println("buf.len", buf.Bytes())

	// buf2 := bytes.NewBuffer(buf.Bytes())
	// // var u2 User
	// decoder := gob.NewDecoder(buf2)
	// val := reflect.New(reflect.TypeOf(user))
	// err := decoder.DecodeValue(val)
	// // err := decoder.Decode(&u2)
	// if err != nil {
	// 	log.Printf("err [%s]", err)
	// }

	// v := val.Elem().Elem().Interface().(User)
	// println(v.String())
	// k := reflect.TypeOf(funa)
	// println(k.In(1).Kind().String())

	// actors[42].Send(actor.CallIO{

	// })

	// funa(&u)

	// ts := new(service.TestService)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)
	<-c
}

func TestProtobuf(t *testing.T) {
	data := &pb.Data{
		Key:    "lee",
		StrVal: "kanon",
		IntVal: 42,
	}
	buf, err := proto.Marshal(data)
	if err != nil {
		println(err.Error())
		return
	}
	newData := pb.Data{}
	typ := reflect.TypeOf(newData)
	val := reflect.New(typ)

	packet := val.Interface().(proto.Message)

	err = proto.Unmarshal(buf, packet)
	if err != nil {
		println(err.Error())
		return
	}
	newData = val.Elem().Interface().(pb.Data)
	fmt.Println(newData.String())

	u1 := User{
		ID:      1,
		Account: "kkk",
	}
	u2 := u1
	u2.ID = 2
	println(u1.ID)
}

func TestSomething(t *testing.T) {
	buf := make([]byte, 0, 5)
	buf = append(buf, []byte{1, 2}...)
	fmt.Println(buf)
}

func TestClientPacket(t *testing.T) {
	ret := pb.Data{
		Key:    "lee",
		StrVal: "kanon",
	}
	buf, err := base.ClientPack(game.MTest["test"].ProtoID, &ret, game.SystemError["success"])
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("pack buf:", buf)
	buf = buf[2:]
	protoID := base.BytesToInt(buf[3:5])
	data, err := base.ClientUnPack(protoID, buf[5:])
	if err != nil {
		fmt.Println(err)
	}
	ret = data.(pb.Data)
	fmt.Println(ret.String())
}
