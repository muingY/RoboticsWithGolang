package main

import (
	"fmt"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/firmata"
	"gobot.io/x/gobot/platforms/keyboard"
)

func main() {
	firmataAdaptor := firmata.NewAdaptor("/dev/ttyACM0")
	servo := gpio.NewServoDriver(firmataAdaptor, "8")
	keyListener := keyboard.NewDriver()

	work := func() {
		var angle uint8
		var inputTogle bool

		/* key contral */
		keyListener.On(keyboard.Key, func(data interface{}) {
			key := data.(keyboard.KeyEvent)

			if key.Key == keyboard.ArrowLeft {
				if angle > 0 {
					inputTogle = true
					angle--
				}
			} else if key.Key == keyboard.ArrowRight {
				if angle < 180 {
					inputTogle = true
					angle++
				}
			} else if key.Key == keyboard.ArrowUp {
				angle = uint8(gobot.Rand(180))
				fmt.Printf("random angle: %d\n", angle)
				inputTogle = true
			}
		})

		/* servo contral */
		for {
			if inputTogle {
				servo.Move(angle)
				fmt.Printf("angle = %d\n", angle)
				inputTogle = false
			}
		}
	}

	robot := gobot.NewRobot("servo bot",
		[]gobot.Connection{firmataAdaptor},
		[]gobot.Device{servo, keyListener},
		work)

	robot.Start()
}
