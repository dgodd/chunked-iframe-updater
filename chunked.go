package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		fmt.Fprint(res, "<html><head><title>Chunked updated iframe</title></head><body><h1>Stream</h1><iframe src='/data' style='border:0;width:100%;height:3em;' scrolling='yes'></iframe><p>Footer</p><p><a href='/chunked.go'>Download source code</a></p></body></html>")
	})
	mux.HandleFunc("/chunked.go", func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "text/plain")
		f, err := os.Open("chunked.go")
		if err != nil {
			res.Header().Set("Status", "500")
			return
		}
		defer f.Close()
		io.Copy(res, f)
	})
	mux.HandleFunc("/data", func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Transfer-Encoding", "chunked")

		fmt.Fprint(res, "<html><head><style>body { margin:0;padding:0; } body > div { position: absolute; top: 0; min-height:10em; width: 100%; background:white; }</style></head><body>")

		fmt.Fprint(res, "<div>foo</div>")
		if f, ok := res.(http.Flusher); ok {
			f.Flush()
		}
		time.Sleep(3 * time.Second)
		fmt.Fprint(res, "<div>bar bar bar</div>")
		if f, ok := res.(http.Flusher); ok {
			f.Flush()
		}
		time.Sleep(3 * time.Second)
		fmt.Fprint(res, "<div>More data</div>")
		if f, ok := res.(http.Flusher); ok {
			f.Flush()
		}
		time.Sleep(3 * time.Second)
		// fmt.Fprint(res, "<script>alert('Last Data')</script>")
		fmt.Fprint(res, "<div>Last data</div>")
	})

	// listen and serve using `ServeMux`
	http.ListenAndServe(":9000", mux)

}
