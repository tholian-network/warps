package test

import "tholian-endpoint/protocols/dns"

type SpyResolverCache struct {
	ResponseExists bool        `json:"response_exists"`
	ResponseRead   *dns.Packet `json:"response_read"`
	ResponseWrite  bool        `json:"response_write"`
	ResolvedExists string      `json:"resolved_exists"`
	ResolvedRead   string      `json:"resolved_read"`
	ResolvedWrite  string      `json:"resolved_write"`
}

func NewSpyResolverCache(response_exists bool, response_read *dns.Packet, response_write bool) SpyResolverCache {

	var cache SpyResolverCache

	cache.ResponseExists = response_exists
	cache.ResponseRead   = response_read
	cache.ResponseWrite  = response_write

	return cache

}

func (cache *SpyResolverCache) Exists(query dns.Packet) bool {

	if query.Type == "query" && len(query.Questions) > 0 {

		resolved := ""

		for q := 0; q < len(query.Questions); q++ {

			resolved += query.Questions[q].Type.String() + ":" + query.Questions[q].Name

			if q <= len(query.Questions) - 1 {
				resolved += ","
			}

		}

		cache.ResolvedExists = resolved

	}

	return cache.ResponseExists

}

func (cache *SpyResolverCache) Read(query dns.Packet) dns.Packet {

	var response dns.Packet

	if query.Type == "query" && len(query.Questions) > 0 {

		resolved := ""

		for q := 0; q < len(query.Questions); q++ {

			resolved += query.Questions[q].Type.String() + ":" + query.Questions[q].Name

			if q <= len(query.Questions) - 1 {
				resolved += ","
			}

		}

		cache.ResolvedRead = resolved

		if cache.ResponseRead != nil {
			response = *cache.ResponseRead
		}

	}

	return response

}

func (cache *SpyResolverCache) Write(query dns.Packet) bool {

	if query.Type == "query" && len(query.Questions) > 0 {

		resolved := ""

		for q := 0; q < len(query.Questions); q++ {

			resolved += query.Questions[q].Type.String() + ":" + query.Questions[q].Name

			if q <= len(query.Questions) - 1 {
				resolved += ","
			}

		}

		cache.ResolvedWrite = resolved

	}

	return cache.ResponseWrite

}

