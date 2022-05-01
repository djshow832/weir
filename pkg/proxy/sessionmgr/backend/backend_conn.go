package backend

import (
	"errors"
	"net"
	"time"

	pnet "github.com/djshow832/weir/pkg/proxy/net"
	"github.com/pingcap/tidb/util/arena"
)

const (
	DialTimeout = 5 * time.Second
)

type connectionPhase byte

type BackendConnection interface {
	Connect() error
	PacketIO() *pnet.PacketIO
	Close() error
}

type BackendConnectionImpl struct {
	pkt        *pnet.PacketIO // a helper to read and write data in packet format.
	alloc      arena.Allocator
	phase      connectionPhase
	capability uint32
	address    string
}

func NewBackendConnectionImpl(address string) *BackendConnectionImpl {
	return &BackendConnectionImpl{
		address: address,
		alloc:   arena.NewAllocator(32 * 1024),
	}
}

func (bc *BackendConnectionImpl) Connect() error {
	cn, err := net.DialTimeout("tcp", bc.address, DialTimeout)
	if err != nil {
		return errors.New("dial backend error")
	}

	bufReadConn := pnet.NewBufferedReadConn(cn)
	pkt := pnet.NewPacketIO(bufReadConn)
	bc.pkt = pkt
	return nil
}

func (bc *BackendConnectionImpl) PacketIO() *pnet.PacketIO {
	return bc.pkt
}

func (bc *BackendConnectionImpl) Close() error {
	return bc.pkt.Close()
}
