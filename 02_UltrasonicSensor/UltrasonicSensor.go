package main

import (
	"fmt"
	"time"
	"errors"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/firmata"
)

func pulseIn(driver *gpio.DirectPinDriver, state int) (int64, error) {
	var stime time.Time
	var etime time.Time

	for {
		val, err := driver.DigitalRead()
		stime = time.Now()
		if err != nil {
			fmt.Println(err)
			break
		}
		if val == 0 { continue }
		break
	}
	for {
		val, err := driver.DigitalRead()
		etime = time.Now()
		if err != nil {
			fmt.Println(err)
			break
		}
		if val == state {
			duration := etime.Sub(stime)
			if duration.Microseconds() > 1000000 {
				return 0, errors.New("over max time(1sec)")
			}
			continue
		}
		break
	}

	duration := etime.Sub(stime)
	return duration.Microseconds(), nil
}

func main() {
	firmataAdaptor := firmata.NewAdaptor("/dev/ttyACM0")
	ultrasonic_trig := gpio.NewDirectPinDriver(firmataAdaptor, "8")
	ultrasonic_echo := gpio.NewDirectPinDriver(firmataAdaptor, "10")

	work := func() {
		gobot.Every(1 * time.Second, func() {
			ultrasonic_trig.DigitalWrite(byte(0))
			time.Sleep(2 * time.Microsecond)
			ultrasonic_trig.DigitalWrite(byte(1))
			time.Sleep(10 * time.Microsecond)
			ultrasonic_trig.DigitalWrite(byte(0))

			duration, err := pulseIn(ultrasonic_echo, 1)
			if err != nil {
				fmt.Println(err)
			}

			distance := duration * 17 / 1000
			fmt.Printf("distance: %d\n", distance);
		})
	}

	robot := gobot.NewRobot("bot",
		[]gobot.Connection{firmataAdaptor},
		[]gobot.Device{ultrasonic_trig, ultrasonic_echo},
		work)

	robot.Start()
}
