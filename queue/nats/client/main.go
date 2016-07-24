// Copyright (c) 2016 - Sarjono Mukti Aji <me@simukti.net>
// Unless otherwise noted, this source code license is MIT-License

package main

import (
	"github.com/nats-io/nats"
	"github.com/simukti/goplay/queue/nats/model"
	"log"
)

// make sure to run worker before sending message
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
	// fire and forget
	// see: http://nats.io/documentation/faq/#gmd
	userEmail := make(chan *model.UserEmail)
	ec.BindSendChan(model.UserEmailSubject, userEmail)
	ec.Flush()

	for i := 1; i <= 100; i++ {
		content := &model.UserEmail{
			UserId: uint32(i),
			Type:   "registration",
		}

		userEmail <- content
	}
}
