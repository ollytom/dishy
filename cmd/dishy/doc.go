/*
Command dishy provides a command line interface to controlling and monitoring a Starlink Dishy device over the network. Its usage is:

	dishy [-a address] command

The following commands are understood:

	reboot
		Request dishy to be rebooted immediately.
	stow
		Reposition dish in a vertical orientation for easier transport.
	unstow
		Reposition the dish in the orientation prior to unstowing.

The flag -a specifies the address to connect to dishy.
By default this is the default IPv4 address and port that dishy listens on 192.168.100.1:9200.

# Examples

Unstow the device, wait 10 seconds, then reboot it:

	dishy unstow
	sleep 10
	dishy reboot

Stow the device accessible through a tunneled connection via the loopback address:

	dishy -a 127.0.0.1:9200 stow

*/
package main
