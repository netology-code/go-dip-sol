package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

type EventLogger struct {
	eventsChan chan string
	done       chan struct{}
	file       *os.File
}

func NewEventLogger(filePath string) *EventLogger {
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	return &EventLogger{
		eventsChan: make(chan string, 100),
		done:       make(chan struct{}),
		file:       f,
	}
}

func (el *EventLogger) LogEvent(event string) {
	select {
	case el.eventsChan <- event:
	default:
		log.Printf("Warning: event logger channel is full")
	}
}

func (el *EventLogger) Start() {
	go el.worker()
}

func (el *EventLogger) worker() {
	for event := range el.eventsChan {
		time.Sleep(1 * time.Second)
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		line := fmt.Sprintf("[%s] %s\n", timestamp, event)

		if _, err := el.file.WriteString(line); err != nil {
			log.Printf("Failed to write to log file: %v", err)
		}
	}

	if err := el.file.Close(); err != nil {
		log.Printf("Failed to close log file: %v", err)
	}

	el.done <- struct{}{}
}

func (el *EventLogger) Stop() {
	close(el.eventsChan)
	<-el.done
	log.Println("Event logger stopped gracefully")
}
