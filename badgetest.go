package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strings"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"github.com/tarm/serial"
)

func main() {
	myConfig, _ := LoadConfig()

	const B = "0x0000FFFF"
	const Z = "0"

	LABELS := []string{
		B, B, B, Z, B, B, B, Z, B, Z, B, Z,
		B, Z, Z, Z, B, B, B, Z, B, Z, B, Z,
		B, B, B, Z, B, Z, Z, Z, B, B, B, Z,
		Z, Z, Z, Z, Z, Z, Z, Z, Z, Z, Z, Z,
		B, B, B, Z, B, B, B, Z, B, B, B, Z,
		B, B, B, Z, B, B, Z, Z, B, B, B, Z,
		B, Z, B, Z, B, B, B, Z, B, Z, B, Z,
	}

	c := &serial.Config{Name: myConfig.Port, Baud: 115200}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	_, err = s.Write([]byte("\n"))
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	write("rgb.clear()", s)

	for true {
		/* The labels dont change. but keep redrawing them because the Badge can suddenly reset
		while the program runs */
		write("rgb.image(["+strings.Join(LABELS, ",")+"], (0,0), (12,7))", s)
		cpuPerc, _ := cpu.Percent(1*time.Second, false)
		drawBar(12, 0, 20, cpuPerc[0], s)

		vmem, _ := mem.VirtualMemory()
		drawBar(12, 4, 20, vmem.UsedPercent, s)
	}

}

// Draws a 3 led high percentage bar. at the given position with len being the lengt of 100%
func drawBar(x int, y int, len int, percentage float64, s *serial.Port) {
	const RED = "0xFF0000FF"
	const GREEN = "0x00FF00FF"

	str := ""

	cpuPercInt := int(math.Round(percentage / 100.0 * float64(len)))
	almostFull := len * 80 / 100.0

	for i := 0; i < cpuPercInt; i++ {
		if i == 0 {
			str = GREEN
		} else if i < almostFull {
			str += "," + GREEN
		} else {
			str += "," + RED
		}
	}

	for i := cpuPercInt; i < len; i++ {
		str = str + ",0xFF"
	}

	write(fmt.Sprintf("rgb.image(["+str+","+str+","+str+"],(%d, %d), (%d,3))", x, y, len), s)
}

func write(txt string, s *serial.Port) {
	s.Write([]byte(txt))
	s.Write([]byte("\r\n"))
}
