/***** BEGIN LICENSE BLOCK *****
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this file,
# You can obtain one at http://mozilla.org/MPL/2.0/.
#
# The Initial Developer of the Original Code is the Mozilla Foundation.
# Portions created by the Initial Developer are Copyright (C) 2012
# the Initial Developer. All Rights Reserved.
#
# Contributor(s):
#   Rob Miller (rmiller@mozilla.com)
#
# ***** END LICENSE BLOCK *****/
package hekagrater

import (
	"log"
	"net"
)

type Input interface {
	Start(pipeline func(*[]byte, chan<- *[]byte))
}

type UdpInput struct {
	listener *net.UDPConn
}

func NewUdpInput(addrStr string) *UdpInput {
	address, err := net.ResolveUDPAddr("udp", addrStr)
	if err != nil {
		log.Printf("ResolveUDPAddr failed: %s\n", err.Error())
		return nil
	}
	listener, err := net.ListenUDP("udp", address)
	if err != nil {
		log.Printf("ListenUDP failed: %s\n", err.Error())
		return nil
	}
	return &UdpInput{listener}
}

func (self *UdpInput) Start(pipeline func(*[]byte, chan<- *[]byte)) {
	go func() {
		recycleChan := make(chan *[]byte, 1000)
		for {
			var inData []byte
			select {
			case inDataPtr := <-recycleChan:
				inData = *inDataPtr
			default:
				inData = make([]byte, 65536)
			}
			n, _, error := (*self.listener).ReadFrom(inData)
			if error != nil {
				continue
			}
			inData = inData[:n]
			go pipeline(&inData, recycleChan)
		}
	}()
}
