package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bootdotdev/learn-http-protocol/internal/request"
	"github.com/bootdotdev/learn-http-protocol/internal/response"
	"github.com/bootdotdev/learn-http-protocol/internal/server"
)

const port = 42069
const httpbin_prefix = "/httpbin/stream/"

func main() {
	server, err := server.Serve(port, handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

// and strings.TrimPrefix
func handler(w *response.Writer, req *request.Request) {
	if strings.HasPrefix(req.RequestLine.RequestTarget, httpbin_prefix) {
		fmt.Println("here!!")
		numstr := strings.TrimPrefix(req.RequestLine.RequestTarget, httpbin_prefix)
		h := response.GetDefaultChunkHeaders()
		fmt.Println(h)
		w.WriteHeaders(h)
		fmt.Println("req!!")
		res, err := http.Get(fmt.Sprintf("https://httpbin.org/stream/%s", numstr))
		if err != nil {
			handler500(w, req)
			return
		}
		fmt.Println("req made!!")
		// buf := make([]byte,1024)
		// for {
		// 	n,err:= res.Body.Read(buf)
		// 	w.WriteChunkedBody(n)
		// 	buf = buf[:0]
		// }
		buf := make([]byte, 1024)
		for {
			n, err := res.Body.Read(buf)
			if n > 0 {
				// Only send actual data read
				_, writeErr := w.WriteChunkedBody(buf[:n])
				if writeErr != nil {
					break // client disconnected or other issue
				}
			}
			if err == io.EOF {
				break
			}
			if err != nil {
				// Log or return 500 optionally
				break
			}
		}

		// Final chunk to indicate end of body
		w.WriteChunkedBodyDone()
	}
	if req.RequestLine.RequestTarget == "/yourproblem" {
		handler400(w, req)
		return
	}
	if req.RequestLine.RequestTarget == "/myproblem" {
		handler500(w, req)
		return
	}
	handler200(w, req)
	return
}

func handler400(w *response.Writer, _ *request.Request) {
	w.WriteStatusLine(response.StatusCodeBadRequest)
	body := []byte(`<html>
<head>
<title>400 Bad Request</title>
</head>
<body>
<h1>Bad Request</h1>
<p>Your request honestly kinda sucked.</p>
</body>
</html>
`)
	h := response.GetDefaultHeaders(len(body))
	h.Override("Content-Type", "text/html")
	w.WriteHeaders(h)
	w.WriteBody(body)
	return
}

func handler500(w *response.Writer, _ *request.Request) {
	w.WriteStatusLine(response.StatusCodeInternalServerError)
	body := []byte(`<html>
<head>
<title>500 Internal Server Error</title>
</head>
<body>
<h1>Internal Server Error</h1>
<p>Okay, you know what? This one is on me.</p>
</body>
</html>
`)
	h := response.GetDefaultHeaders(len(body))
	h.Override("Content-Type", "text/html")
	w.WriteHeaders(h)
	w.WriteBody(body)
}

func handler200(w *response.Writer, _ *request.Request) {
	w.WriteStatusLine(response.StatusCodeSuccess)
	body := []byte(`<html>
<head>
<title>200 OK</title>
</head>
<body>
<h1>Success!</h1>
<p>Your request was an absolute banger.</p>
</body>
</html>
`)
	h := response.GetDefaultHeaders(len(body))
	h.Override("Content-Type", "text/html")
	w.WriteHeaders(h)
	w.WriteBody(body)
	return
}
