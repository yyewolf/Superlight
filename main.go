package main

import (
	"fmt"
	"math"
	"time"
)

var lastBrightness float64

func main() {
	err := LoadConfig()
	if err != nil {
		err = InitialConfig()
		if err != nil {
			fmt.Println("Error creating config file:", err)
			// Waits for input
			fmt.Scanln()
			return
		}
		fmt.Println("Please edit superlight.json and run again")
		// Waits for input
		fmt.Scanln()
		return
	}

	err = getToken()
	if err != nil {
		fmt.Println("Error getting token:", err)
		// Waits for input
		fmt.Scanln()
		return
	}

	fmt.Println("Successfully connected to Tuya API.\nPress Ctrl+C to exit")

	for {
		brightness := calcBrightness()
		if brightness < 5 {
			brightness = -100
		} else {
			// normalize brightness to be between 0 and 100
			brightness = (brightness - 5) * 1000 / 990
		}

		if lastBrightness == brightness {
			time.Sleep(10 * time.Millisecond)
			continue
		}

		if math.Abs(brightness-lastBrightness) < 2 {
			time.Sleep(10 * time.Millisecond)
			continue
		}

		if config.Debug {
			fmt.Println("Brightness:", brightness)
		}

		sendDeviceBrightness(brightness)
		lastBrightness = brightness
		time.Sleep(20 * time.Millisecond)
	}
}
