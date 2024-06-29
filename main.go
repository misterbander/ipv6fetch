package main

import (
	"errors"
	"log"
	"net"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		addr, err := getIPV6Address()
		if err != nil {
			log.Panicln("Can't get ipv6 address: " + err.Error())
		}

		w.Write([]byte(addr))
		w.WriteHeader(http.StatusOK)
	})

	log.Println("🌐 Server started at port 25566!")
	http.ListenAndServe(":25566", r)
}

func getIPV6Address() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, iface := range interfaces {
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			ipnet, ok := addr.(*net.IPNet)
			if ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() == nil && !ipnet.IP.IsLinkLocalUnicast() {
				return ipnet.IP.String(), nil
			}
		}
	}

	return "", errors.New("no ipv6 found")
}
