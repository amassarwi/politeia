// Copyright (c) 2020 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package main

import (
	"github.com/decred/politeia/politeiawww/user"
)

const (
	// CMS events
	eventInvoiceComment      = "eventInvoiceComment"
	eventInvoiceStatusUpdate = "eventInvoiceStatusUpdate"
	eventDCCNew              = "eventDCCNew"
	eventDCCSupportOppose    = "eventDCCSupportOppose"
)

func (p *politeiawww) setupEventListenersCMS() {
	// Setup invoice comment event
	ch := make(chan interface{})
	p.events.Register(eventInvoiceComment, ch)
	go p.handleEventInvoiceComment(ch)

	// Setup invoice status update event
	ch = make(chan interface{})
	p.events.Register(eventInvoiceStatusUpdate, ch)
	go p.handleEventInvoiceStatusUpdate(ch)

	// Setup DCC new update event
	ch = make(chan interface{})
	p.events.Register(eventDCCNew, ch)
	go p.handleEventDCCNew(ch)

	// Setup DCC support/oppose event
	ch = make(chan interface{})
	p.events.Register(eventDCCSupportOppose, ch)
	go p.handleEventDCCSupportOppose(ch)
}

type dataInvoiceComment struct {
	token string // Comment token
	email string // User email
}

func (p *politeiawww) handleEventInvoiceComment(ch chan interface{}) {
	for msg := range ch {
		d, ok := msg.(dataInvoiceComment)
		if !ok {
			log.Errorf("handleEventInvoiceComment invalid msg: %v", msg)
			continue
		}

		err := p.emailInvoiceNewComment(d.email)
		if err != nil {
			log.Errorf("emailInvoiceNewComment %v: %v", err)
		}

		log.Debugf("Sent invoice comment notification %v", d.token)
	}
}

type dataInvoiceStatusUpdate struct {
	token string // Comment token
	email string // User email
}

func (p *politeiawww) handleEventInvoiceStatusUpdate(ch chan interface{}) {
	for msg := range ch {
		d, ok := msg.(dataInvoiceStatusUpdate)
		if !ok {
			log.Errorf("handleEventInvoiceStatusUpdate invalid msg: %v", msg)
			continue
		}

		err := p.emailInvoiceStatusUpdate(d.token, d.email)
		if err != nil {
			log.Errorf("emailInvoiceStatusUpdate %v: %v", err)
		}

		log.Debugf("Sent invoice status update notification %v", d.token)
	}
}

type dataDCCNew struct {
	token string // DCC token
}

func (p *politeiawww) handleEventDCCNew(ch chan interface{}) {
	for msg := range ch {
		d, ok := msg.(dataDCCNew)
		if !ok {
			log.Errorf("handleEventDCCNew invalid msg: %v", msg)
			continue
		}

		emails := make([]string, 0, 256)
		err := p.db.AllUsers(func(u *user.User) {
			// Check circumstances where we don't notify
			switch {
			case !u.Admin:
				// Only notify admin users
				return
			case u.Deactivated:
				// Never notify deactivated users
				return
			}

			emails = append(emails, u.Email)
		})
		if err != nil {
			log.Errorf("handleEventDCCNew: AllUsers: %v", err)
		}

		err = p.emailDCCSubmitted(d.token, emails)
		if err != nil {
			log.Errorf("emailDCCSubmitted %v: %v", err)
		}

		log.Debugf("Sent DCC new notification %v", d.token)
	}
}

type dataDCCSupportOppose struct {
	token string // DCC token
}

func (p *politeiawww) handleEventDCCSupportOppose(ch chan interface{}) {
	for msg := range ch {
		d, ok := msg.(dataDCCSupportOppose)
		if !ok {
			log.Errorf("handleEventDCCSupportOppose invalid msg: %v", msg)
			continue
		}

		emails := make([]string, 0, 256)
		err := p.db.AllUsers(func(u *user.User) {
			// Check circumstances where we don't notify
			switch {
			case !u.Admin:
				// Only notify admin users
				return
			case u.Deactivated:
				// Never notify deactivated users
				return
			}

			emails = append(emails, u.Email)
		})
		if err != nil {
			log.Errorf("handleEventDCCSupportOppose: AllUsers: %v", err)
		}

		err = p.emailDCCSupportOppose(d.token, emails)
		if err != nil {
			log.Errorf("emailDCCSupportOppose %v: %v", err)
		}

		log.Debugf("Sent DCC support/oppose notification %v", d.token)
	}
}
