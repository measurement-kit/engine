// Package resolverlookup discovers the resolver's IP address.
package resolverlookup

import (
	"context"
	"errors"
	"net"
)

// Perform discovers the resolver's IP address.
func Perform(ctx context.Context) (string, error) {
	addrs, err := net.LookupHost("whoami.akamai.net")
	if err != nil {
		return "127.0.0.1", err
	}
	if len(addrs) != 1 {
		return "127.0.0.1", errors.New("unexpected addrs length")
	}
	return addrs[0], nil
}
