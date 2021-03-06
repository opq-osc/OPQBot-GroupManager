// +build android

package androidDns

import (
	"context"
	"log"
	"net"
)

const bootstrapDNS = "8.8.8.8:53"

func SetDns() {
	log.Println("安卓设置DNS")
	var dialer net.Dialer
	net.DefaultResolver = &net.Resolver{
		PreferGo: false,
		Dial: func(context context.Context, _, _ string) (net.Conn, error) {
			conn, err := dialer.DialContext(context, "udp", bootstrapDNS)
			if err != nil {
				return nil, err
			}
			return conn, nil
		},
	}
}
