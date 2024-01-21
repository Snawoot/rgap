package rgap

import (
	"context"
	"fmt"
	"log"
	"net"
)

type UDPSource struct {
	address   string
	label     string
	callback  func(string, *Announcement)
	ctx       context.Context
	ctxCancel func()
	loopDone  chan struct{}
}

func NewUDPSource(address string, label string, callback func(string, *Announcement)) *UDPSource {
	s := &UDPSource{
		address:  address,
		label:    label,
		callback: callback,
	}
	return s
}

func (s *UDPSource) Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	s.ctx = ctx
	s.ctxCancel = cancel
	s.loopDone = make(chan struct{})

	udpAddr, err := net.ResolveUDPAddr("udp", s.address)
	if err != nil {
		return fmt.Errorf("bad UDP listen address: %w", err)
	}

	var conn *net.UDPConn

	if udpAddr.IP.IsMulticast() {
		conn, err = net.ListenMulticastUDP("udp4", nil, udpAddr)
		if err != nil {
			return fmt.Errorf("UDP listen failed: %w", err)
		}
	} else {
		conn, err = net.ListenUDP("udp", udpAddr)
		if err != nil {
			return fmt.Errorf("UDP listen failed: %w", err)
		}
	}

	go func() {
		select {
		case <-ctx.Done():
		case <-s.loopDone:
		}
		conn.Close()
	}()
	go s.readLoop(conn)
	log.Printf("Started UDP source @ %s", s.address)
	return nil
}

func (s *UDPSource) Stop() error {
	s.ctxCancel()
	<-s.loopDone
	log.Printf("Stopped UDP source @ %s", s.address)
	return nil
}

func (s *UDPSource) readLoop(conn *net.UDPConn) {
	defer close(s.loopDone)
	buf := make([]byte, 4096)
	for s.ctx.Err() == nil {
		n, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			if s.ctx.Err() != nil {
				return
			}
			log.Printf("source %s: UDP read error: %w", s.label, err)
			continue
		}
		if n != announcementSize {
			continue
		}
		ann := new(Announcement)
		if err := ann.UnmarshalBinary(buf[:n]); err != nil {
			log.Println("announce unmarshaling failed: %v", err)
		}
		s.callback(s.label, ann)
	}
}
