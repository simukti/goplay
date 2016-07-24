// Copyright (c) 2016 - Sarjono Mukti Aji <me@simukti.net>
// Unless otherwise noted, this source code license is MIT-License

package main

import (
	"fmt"
	"github.com/nats-io/nats"
	"github.com/simukti/goplay/queue/nats/model"
	"log"
	"runtime"
	"time"
)

func main() {
	nc, ncErr := nats.Connect(nats.DefaultURL)
	if ncErr != nil {
		log.Fatalln(ncErr)
	}
	ec, ecErr := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if ecErr != nil {
		log.Fatalln(ecErr)
	}
	defer ec.Close()

	messageReceived := make(chan *model.UserEmail)
	// if worker run more than one, queue will split by num of workers
	// so every worker will NOT do the same task
	ec.BindRecvQueueChan(model.UserEmailSubject, model.UserEmailSubject, messageReceived)
	ec.Flush()

	for message := range messageReceived {
		handleUserMessage(message)
	}

	runtime.Goexit()
}

func handleUserMessage(m *model.UserEmail) {
	log.Println(fmt.Sprintf("Send %s email to user: %d", m.Type, m.UserId))
	time.Sleep(100 * time.Millisecond) // simulating process
}
