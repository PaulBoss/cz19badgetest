package main

import (
	"log"
	"math"
	"os"
	"strconv"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"github.com/tarm/serial"
)

func main() {
	myConfig, _ := LoadConfig()

	scrbuf := make([]int, 8*32)

	const B = 0x0000FFFF
	const Z = 0

	LABELS := []int{
		B, B, B, Z, B, B, B, Z, B, Z, B, Z,
		B, Z, Z, Z, B, B, B, Z, B, Z, B, Z,
		B, B, B, Z, B, Z, Z, Z, B, B, B, Z,
		Z, Z, Z, Z, Z, Z, Z, Z, Z, Z, Z, Z,
		B, B, B, Z, B, B, B, Z, B, B, B, Z,
		B, B, B, Z, B, B, Z, Z, B, B, B, Z,
		B, Z, B, Z, B, B, B, Z, B, Z, B, Z,
	}

	for y := 0; y < 7; y++ {
		for x := 0; x < 12; x++ {
			scrbuf[x+y*32] = LABELS[x+y*12]
		}
	}

	c := &serial.Config{Name: myConfig.Port, Baud: 115200, ReadTimeout: time.Millisecond * 100}
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
	write("rgb.disablecomp()", s)

	for true {
		cpuPerc, _ := cpu.Percent(1*time.Second, false)

		drawBar(12, 0, 20, cpuPerc[0], scrbuf)
		drawBar(12, 1, 20, cpuPerc[0], scrbuf)
		drawBar(12, 2, 20, cpuPerc[0], scrbuf)

		vmem, _ := mem.VirtualMemory()
		drawBar(12, 4, 20, vmem.UsedPercent, scrbuf)
		drawBar(12, 5, 20, vmem.UsedPercent, scrbuf)
		drawBar(12, 6, 20, vmem.UsedPercent, scrbuf)

		dumpScreenBuf(scrbuf, s)
	}
}

// Dump the screenbuffer to the display with a python command
func dumpScreenBuf(buf []int, s *serial.Port) {
	str := "rgb.frame(["

	first := true

	for _, char := range buf {
		if first {
			str += strconv.Itoa(char)
			first = false
		} else {
			str += "," + strconv.Itoa(char)
		}

	}

	str += "])"
	write(str, s)
}

// Draws a 3 led high percentage bar. at the given position with len being the length of 100%
func drawBar(x int, y int, len int, percentage float64, s []int) {
	const RED = 0xFF0000FF
	const GREEN = 0x00FF00FF

	cpuPercInt := int(math.Round(percentage / 100.0 * float64(len)))
	almostFull := len * 80 / 100.0

	for i := x; i < x+cpuPercInt; i++ {
		if i < almostFull+x {
			s[i+y*32] = GREEN
		} else {
			s[i+y*32] = RED
		}
	}

	for i := x + cpuPercInt; i < x+len; i++ {
		s[i+y*32] = 0
	}

}

// Write to the display
func write(txt string, s *serial.Port) {
	s.Write([]byte(txt))
	s.Write([]byte("\r\n"))
}
