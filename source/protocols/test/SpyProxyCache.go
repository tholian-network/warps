package test

import "tholian-endpoint/protocols/http"

type SpyProxyCache struct {
	ResponseExists  bool         `json:"response_exists"`
	ResponseRead    *http.Packet `json:"response_read"`
	ResponseWrite   bool         `json:"response_write"`
	RequestedExists string       `json:"requested_exists"`
	RequestedRead   string       `json:"requested_read"`
	RequestedWrite  string       `json:"requested_write"`
}

func NewSpyProxyCache(response_exists bool, response_read *http.Packet, response_write bool) SpyProxyCache {

	var cache SpyProxyCache

	cache.ResponseExists = response_exists
	cache.ResponseRead = response_read
	cache.ResponseWrite = response_write

	return cache

}

func (cache *SpyProxyCache) Exists(request http.Packet) bool {

	if request.Type == "request" && request.URL != nil {
		cache.RequestedExists = request.URL.String()
	}

	return cache.ResponseExists

}

func (cache *SpyProxyCache) Read(request http.Packet) http.Packet {

	var response http.Packet

	if request.Type == "request" && request.URL != nil {

		if cache.ResponseRead != nil {
			response = *cache.ResponseRead
		}

		cache.RequestedRead = request.URL.String()

	}

	return response

}

func (cache *SpyProxyCache) Write(request http.Packet) bool {

	if request.Type == "request" && request.URL != nil {
		cache.RequestedWrite = request.URL.String()
	}

	return cache.ResponseWrite

}

