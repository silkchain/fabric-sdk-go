// +build deprecated

/*
Copyright SecureKey Technologies Inc., Unchain B.V. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package peer

import (
	"crypto/x509"

	"github.com/hyperledger/fabric-sdk-go/pkg/context/api/core"
	"github.com/hyperledger/fabric-sdk-go/pkg/context/api/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config/urlutil"
	"google.golang.org/grpc/keepalive"
)

// NewPeerTLSFromCert constructs a Peer given its endpoint configuration settings.
// url is the URL with format of "host:port".
// certificate is ...
// serverNameOverride is passed to NewClientTLSFromCert in grpc/credentials.
// Deprecated: use peer.New() instead
func NewPeerTLSFromCert(url string, certPath string, serverHostOverride string, config core.Config) (*Peer, error) {
	var certificate *x509.Certificate
	var err error

	if urlutil.IsTLSEnabled(url) {
		certConfig := core.TLSConfig{Path: certPath}
		certificate, err = certConfig.TLSCert()

		if err != nil {
			return nil, err
		}
	}
	var kap keepalive.ClientParameters

	// TODO: config is declaring TLS but cert & serverHostOverride is being passed-in...
	endorseRequest := peerEndorserRequest{
		target:             url,
		certificate:        certificate,
		serverHostOverride: serverHostOverride,
		dialBlocking:       connBlocking,
		config:             config,
		kap:                kap,
		failFast:           false,
		allowInsecure:      true,
	}
	conn, err := newPeerEndorser(&endorseRequest)

	if err != nil {
		return nil, err
	}

	return NewPeerFromProcessor(url, conn, config)
}

// NewPeerFromConfig constructs a Peer from given peer configuration and global configuration setting.
// Deprecated: use peer.New() instead
func NewPeerFromConfig(peerCfg *core.NetworkPeer, config core.Config) (*Peer, error) {

	serverHostOverride := ""
	if str, ok := peerCfg.GRPCOptions["ssl-target-name-override"].(string); ok {
		serverHostOverride = str
	}

	allowInsecure := false
	if boolVal, ok := peerCfg.GRPCOptions["allow-insecure"].(bool); ok {
		allowInsecure = !urlutil.HasProtocol(peerCfg.URL) && boolVal
	}

	var certificate *x509.Certificate
	var err error
	kap := getKeepAliveOptions(peerCfg)
	failFast := getFailFast(peerCfg)
	if urlutil.IsTLSEnabled(peerCfg.URL) {
		certificate, err = peerCfg.TLSCACerts.TLSCert()

		if err != nil {
			return nil, err
		}
	}

	endorseRequest := peerEndorserRequest{
		target:             peerCfg.URL,
		certificate:        certificate,
		serverHostOverride: serverHostOverride,
		dialBlocking:       connBlocking,
		config:             config,
		kap:                kap,
		failFast:           failFast,
		allowInsecure:      allowInsecure,
	}
	conn, err := newPeerEndorser(&endorseRequest)

	if err != nil {
		return nil, err
	}

	newPeer, err := NewPeerFromProcessor(peerCfg.URL, conn, config)
	if err != nil {
		return nil, err
	}

	// TODO: Remove upon making peer interface immutable
	newPeer.SetMSPID(peerCfg.MspID)

	return newPeer, nil
}

// NewPeer constructs a Peer given its endpoint configuration settings.
// url is the URL with format of "host:port".
// Deprecated: use peer.New() instead
func NewPeer(url string, config core.Config) (*Peer, error) {
	var kap keepalive.ClientParameters
	endorseRequest := peerEndorserRequest{
		target:             url,
		certificate:        nil,
		serverHostOverride: "",
		dialBlocking:       connBlocking,
		config:             config,
		kap:                kap,
		failFast:           false,
		allowInsecure:      true,
	}
	conn, err := newPeerEndorser(&endorseRequest)
	if err != nil {
		return nil, err
	}

	return NewPeerFromProcessor(url, conn, config)
}

// NewPeerFromProcessor constructs a Peer with a ProposalProcessor to simulate transactions.
// Deprecated: use peer.New() instead
func NewPeerFromProcessor(url string, processor fab.ProposalProcessor, config core.Config) (*Peer, error) {
	return &Peer{url: url, processor: processor}, nil
}
