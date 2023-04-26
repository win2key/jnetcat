//go:build windows
// +build windows

package main

import (
	"log"

	"golang.org/x/sys/windows/svc"
)

func main() {
	isInteractive, err := svc.IsAnInteractiveSession()
	if err != nil {
		log.Fatalf("Failed to determine if session is interactive: %v", err)
	}
	if !isInteractive {
		runService()
		return
	} else {
		runNetc()
	}
}

type myService struct{}

func (m *myService) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown

	changes <- svc.Status{State: svc.StartPending}

	go runNetc()

	changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}

	for {
		select {
		case c := <-r:
			switch c.Cmd {
			case svc.Interrogate:
				changes <- c.CurrentStatus
			case svc.Stop, svc.Shutdown:
				changes <- svc.Status{State: svc.StopPending}
				os.Exit(0)
				return
			default:
				log.Printf("Unexpected control request #%d", c)
			}
		}
	}
}

func runService() {
	const serviceName = "jNetCat Service"

	err := svc.Run(serviceName, &myService{})
	if err != nil {
		log.Fatalf("Service failed: %v", err)
	}
}
