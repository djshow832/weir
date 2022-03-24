// Copyright 2013 The Go-MySQL-Driver Authors. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at http://mozilla.org/MPL/2.0/.

// The MIT License (MIT)
//
// Copyright (c) 2014 wandoulabs
// Copyright (c) 2014 siddontang
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// Copyright 2015 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package net

import (
	"bufio"
	"io"
	"time"

	"github.com/pingcap/errors"
	"github.com/pingcap/tidb/errno"
	"github.com/pingcap/tidb/parser/mysql"
	"github.com/pingcap/tidb/parser/terror"
	"github.com/pingcap/tidb/util/dbterror"
)

const defaultWriterSize = 16 * 1024

var (
	errInvalidSequence = dbterror.ClassServer.NewStd(errno.ErrInvalidSequence)
)

// packetIO is a helper to read and write data in packet format.
type PacketIO struct {
	bufReadConn *BufferedReadConn
	bufWriter   *bufio.Writer
	sequence    uint8
	readTimeout time.Duration
}

func NewPacketIO(bufReadConn *BufferedReadConn) *PacketIO {
	p := &PacketIO{sequence: 0}
	p.SetBufferedReadConn(bufReadConn)
	return p
}

func (p *PacketIO) SetBufferedReadConn(bufReadConn *BufferedReadConn) {
	p.bufReadConn = bufReadConn
	p.bufWriter = bufio.NewWriterSize(bufReadConn, defaultWriterSize)
}

func (p *PacketIO) SetReadTimeout(timeout time.Duration) {
	p.readTimeout = timeout
}

func (p *PacketIO) ResetSequence() {
	p.sequence = 0
}

func (p *PacketIO) ReadOnePacket() ([]byte, error) {
	var header [4]byte

	if _, err := io.ReadFull(p.bufReadConn, header[:]); err != nil {
		return nil, errors.Trace(err)
	}

	sequence := header[3]
	if sequence != p.sequence {
		return nil, errInvalidSequence.GenWithStack("invalid sequence %d != %d", sequence, p.sequence)
	}

	p.sequence++

	length := int(uint32(header[0]) | uint32(header[1])<<8 | uint32(header[2])<<16)

	data := make([]byte, length)
	if _, err := io.ReadFull(p.bufReadConn, data); err != nil {
		return nil, errors.Trace(err)
	}
	return data, nil
}

// ReadPacket reads data and removes the header
func (p *PacketIO) ReadPacket() ([]byte, error) {
	data, err := p.ReadOnePacket()
	if err != nil {
		return nil, errors.Trace(err)
	}

	if len(data) < mysql.MaxPayloadLen {
		return data, nil
	}

	// handle multi-packet
	for {
		buf, err := p.ReadOnePacket()
		if err != nil {
			return nil, errors.Trace(err)
		}

		data = append(data, buf...)

		if len(buf) < mysql.MaxPayloadLen {
			break
		}
	}

	return data, nil
}

// writePacket writes data without a header
func (p *PacketIO) WritePacket(data []byte) error {
	length := len(data)
	header := make([]byte, 4, 4)

	for length >= mysql.MaxPayloadLen {
		header[0] = 0xff
		header[1] = 0xff
		header[2] = 0xff
		header[3] = p.sequence

		if n, err := p.bufWriter.Write(header); err != nil || n != 4 {
			return errors.Trace(mysql.ErrBadConn)
		}
		if n, err := p.bufWriter.Write(data[:mysql.MaxPayloadLen]); err != nil {
			return errors.Trace(mysql.ErrBadConn)
		} else if n != mysql.MaxPayloadLen {
			return errors.Trace(mysql.ErrBadConn)
		} else {
			p.sequence++
			length -= mysql.MaxPayloadLen
			data = data[mysql.MaxPayloadLen:]
		}
	}

	header[0] = byte(length)
	header[1] = byte(length >> 8)
	header[2] = byte(length >> 16)
	header[3] = p.sequence

	if n, err := p.bufWriter.Write(header); err != nil || n != 4 {
		return errors.Trace(mysql.ErrBadConn)
	}
	if n, err := p.bufWriter.Write(data); err != nil {
		terror.Log(errors.Trace(err))
		return errors.Trace(mysql.ErrBadConn)
	} else if n != len(data) {
		return errors.Trace(mysql.ErrBadConn)
	} else {
		p.sequence++
		return nil
	}
}

func (p *PacketIO) Flush() error {
	err := p.bufWriter.Flush()
	if err != nil {
		return errors.Trace(err)
	}
	return err
}
