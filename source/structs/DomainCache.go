package structs

import "os"
import "strings"

type DomainCache struct {
	Folder string `json:"folder"`
}

func NewDomainCache(folder string) DomainCache {

	var cache DomainCache

	if strings.HasSuffix(folder, "/") {
		folder = folder[0:len(folder)-1]
	}

	stat, err1 := os.Stat(folder)

	if err1 == nil && stat.IsDir() {

		cache.Folder = folder

	} else {

		err2 := os.MkdirAll(folder, 0750)

		if err2 == nil {
			cache.Folder = folder
		}

	}

	return cache

}
