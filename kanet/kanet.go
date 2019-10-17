package kanet

import "time"

func Second(sec int) time.Duration {
	return time.Duration(sec * 1000 * 1000 * 1000)
}

func MilliSecond(ms int) time.Duration {
	return time.Duration(ms * 1000 * 1000)
}

func Init() {

}

// func main() {
// 	actor.Init()
// 	println("kkkkkkkkkkkkkk")
// 	//定时器
// 	timerSerivce := new(service.TimerService)
// 	timerSerivce.Init(1000)

// }
