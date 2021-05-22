package main

import (
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/firmata"
)

func getMiliTime() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func main() {
	firmataAdaptor := firmata.NewAdaptor("/dev/ttyACM0")
	led1 := gpio.NewLedDriver(firmataAdaptor, "10")
	led2 := gpio.NewLedDriver(firmataAdaptor, "9")
	led3 := gpio.NewLedDriver(firmataAdaptor, "8")

	work := func() {
		var stime int64
		var etime int64
		var deltaTime int64
		var elapsed int64
		for {
			stime = getMiliTime()

			elapsed += deltaTime
			if elapsed >= 0 && elapsed < 500 {
				led1.On()
				led2.Off()
				led3.Off()
			} else if elapsed >= 500 && elapsed < 1000 {
				led1.Off()
				led2.On()
				led3.Off()
			} else if elapsed >= 1000 && elapsed < 1500 {
				led1.Off()
				led2.Off()
				led3.On()
			} else {
				elapsed = 0
			}

			etime = getMiliTime()
			deltaTime = etime - stime
			if deltaTime == 0 {
				time.Sleep(time.Millisecond)
				deltaTime++
			}
		}
	}

	robot := gobot.NewRobot("bot",
		[]gobot.Connection{firmataAdaptor},
		[]gobot.Device{led1, led2, led3},
		work)

	robot.Start()
}
