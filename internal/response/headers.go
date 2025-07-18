package response

import (
	"fmt"

	"github.com/tiggercwh/boot-dev-httptcp/internal/headers"
)

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()
	h.Set("Content-Length", fmt.Sprintf("%d", contentLen))
	h.Set("Connection", "close")
	h.Set("Content-Type", "text/plain")
	return h
}

func GetDefaultChunkHeaders() headers.Headers {
	h := headers.NewHeaders()
	h.Set("Connection", "keep-alive")
	h.Set("Content-Type", "text/plain")
	h.Set("Transfer-Encoding", "chunked")
	return h
}
