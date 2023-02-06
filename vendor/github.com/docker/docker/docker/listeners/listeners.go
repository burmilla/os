package listeners

import (
	"crypto/tls"
	"net"

	"github.com/burmilla/go-connections-old/sockets"
	"github.com/sirupsen/logrus"
)

func initTCPSocket(addr string, tlsConfig *tls.Config) (l net.Listener, err error) {
	if tlsConfig == nil || tlsConfig.ClientAuth != tls.RequireAndVerifyClientCert {
		logrus.Warn("/!\\ DON'T BIND ON ANY IP ADDRESS WITHOUT setting -tlsverify IF YOU DON'T KNOW WHAT YOU'RE DOING /!\\")
	}
	if l, err = sockets.NewTCPSocket(addr, tlsConfig); err != nil {
		return nil, err
	}
	return
}
