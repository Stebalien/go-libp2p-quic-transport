package libp2pquic

import (
	"net"

	ic "github.com/libp2p/go-libp2p-crypto"
	peer "github.com/libp2p/go-libp2p-peer"
	tpt "github.com/libp2p/go-libp2p-transport"
	smux "github.com/libp2p/go-stream-muxer"
	quic "github.com/lucas-clemente/quic-go"
	ma "github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr-net"
)

type conn struct {
	sess      quic.Session
	transport tpt.Transport

	localPeer      peer.ID
	privKey        ic.PrivKey
	localMultiaddr ma.Multiaddr

	remotePeerID    peer.ID
	remotePubKey    ic.PubKey
	remoteMultiaddr ma.Multiaddr
}

var _ tpt.Conn = &conn{}

func (c *conn) Close() error {
	return c.sess.Close(nil)
}

// IsClosed returns whether a connection is fully closed.
func (c *conn) IsClosed() bool {
	return c.sess.Context().Err() != nil
}

// OpenStream creates a new stream.
func (c *conn) OpenStream() (smux.Stream, error) {
	qstr, err := c.sess.OpenStreamSync()
	return &stream{Stream: qstr}, err
}

// AcceptStream accepts a stream opened by the other side.
func (c *conn) AcceptStream() (smux.Stream, error) {
	qstr, err := c.sess.AcceptStream()
	return &stream{Stream: qstr}, err
}

// LocalPeer returns our peer ID
func (c *conn) LocalPeer() peer.ID {
	return c.localPeer
}

// LocalPrivateKey returns our private key
func (c *conn) LocalPrivateKey() ic.PrivKey {
	return c.privKey
}

// RemotePeer returns the peer ID of the remote peer.
func (c *conn) RemotePeer() peer.ID {
	return c.remotePeerID
}

// RemotePublicKey returns the public key of the remote peer.
func (c *conn) RemotePublicKey() ic.PubKey {
	return c.remotePubKey
}

// LocalMultiaddr returns the local Multiaddr associated
func (c *conn) LocalMultiaddr() ma.Multiaddr {
	return c.localMultiaddr
}

// RemoteMultiaddr returns the remote Multiaddr associated
func (c *conn) RemoteMultiaddr() ma.Multiaddr {
	return c.remoteMultiaddr
}

func (c *conn) Transport() tpt.Transport {
	return c.transport
}

// TODO: there must be a better way to do this
func quicMultiaddr(na net.Addr) (ma.Multiaddr, error) {
	udpMA, err := manet.FromNetAddr(na)
	if err != nil {
		return nil, err
	}
	quicMA, err := ma.NewMultiaddr(udpMA.String() + "/quic")
	if err != nil {
		return nil, err
	}
	return quicMA, nil
}
