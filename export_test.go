package imosrpc

import (
	"net"
)

var BuildQuery = buildQuery
var ParseQuery = parseQuery

func SetHttpListener(listener net.Listener) {
	httpListener = listener
}
