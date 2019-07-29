# CZ19BadgeTest

A small program to plot the CPU and memory usage of your PC on the Campzone 2019 hackzone badge.

Please install the driver from https://learn.sparkfun.com/tutorials/sparkfun-serial-basic-ch340c-hookup-guide/drivers-if-you-need-them 
before running this program, use an USB cable to plugin the badge in your PC and determine the COM port used.

## Config

By default "COM1" will be used. If you are not on Windows or want to use another COM port, simply add a badge.conf next to the program
with the following content.
```
Port = "COM3"
```
COM3 can be replaced with the correct COM port. This should also work in Linux (e.g. "/dev/ttyUSB1") but that is not tested at the moment