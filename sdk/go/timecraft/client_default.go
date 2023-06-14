//go:build !wasip1

package timecraft

import (
	"context"
	"crypto/tls"
	"net"

	"github.com/stealthrocket/timecraft/internal/timecraft"
)

func dialContext(ctx context.Context, network, addr string, _ *tls.Config) (net.Conn, error) {
	var d net.Dialer
	return d.DialContext(ctx, "unix", timecraft.Socket)
}
