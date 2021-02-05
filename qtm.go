package main

import (
	"github.com/lqhoang99/tcp-quic-tools/client"
	"github.com/lqhoang99/tcp-quic-tools/server"
	"github.com/lqhoang99/tcp-quic-tools/util/cli"
	"log"
)

func main() {
	opt := cli.ParseOptions()

	if opt.IsServerMode {
		log.Println("Tool started in SERVER mode")

		s, err := server.NewServer(opt)
		if err != nil {
			log.Fatalln(err.Error())
		}

		wg, err := s.Listen(&opt.Address)
		if err != nil {
			log.Fatalf("Server could not start listening: %s", err.Error())
		}

		wg.Wait() // Wait for server termination
	} else {
		log.Println("Tool started in CLIENT mode")

		c, err := client.NewClient(opt)
		if err != nil {
			log.Fatalln(err.Error())
		}

		if opt.Duration > -1 {
			// Send for the set duration to the server
			sentBytes, err := c.SendDuration(opt.Duration, opt.BufferSize)
			if err != nil {
				log.Fatalf("Encountered error when trying to send for %d ns to the server. Error: %s", opt.Duration.Nanoseconds(), err.Error())
			}

			log.Printf("Sent %d bytes in %d nanoseconds", sentBytes, opt.Duration.Nanoseconds())
			mb := float64(opt.Bytes*8)/float64(1000000)
			log.Printf("%f Mbit/s", mb/(opt.Duration.Seconds()))
		}else if opt.Bytes > -1 {
			// Send the set amount of bytes to the server
			time, err := c.SendBytes(opt.Bytes)
			if err != nil {
				log.Fatalf("Encountered error when trying to send %d bytes to the server. Error: %s", opt.Bytes, err.Error())
			}

			log.Printf("Sent %d bytes in %d nanoseconds", opt.Bytes, time.Nanoseconds())
			mb := float64(opt.Bytes*8)/float64(1000000)
			log.Printf("%f Mbit/s", mb/time.Seconds())
		} else {
			log.Fatalf("You need to either set --bytes or --duration to measure throughput")
		}

		err = c.Cleanup()
		if err != nil {
			log.Fatalf("Could not cleanup client properly. Error: %s", err.Error())
		}
	}
}
