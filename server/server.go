package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"path"
	"path/filepath"
	"strings"
)

type fileHandler struct {
	root string
}

func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP
}

func main() {
	var port string

	fmt.Print("Enter port:")
	fmt.Scanf("%s\n", &port)
	mux := http.NewServeMux()

	staticHandler := fileServer("../website")

	mux.Handle("/", http.StripPrefix("/", staticHandler))
	log.Printf("Listening on: %s\n", GetOutboundIP().String())
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

func fileServer(root string) http.Handler {
	return &fileHandler{root}
}

func (f *fileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	upath := r.URL.Path
	if !strings.HasPrefix(upath, "/") {
		upath = "/" + upath
		if upath == "/" {
			// upath += "index.html"
		}
		r.URL.Path = upath
	}

	name := filepath.Join(f.root, path.Clean(upath))
	log.Printf("fileHandler.ServeHTTP: path=%s", name)
	http.ServeFile(w, r, name)
}
