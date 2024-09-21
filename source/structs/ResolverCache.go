package structs

import "tholian-endpoint/protocols/dns"
import "bytes"
import "encoding/json"
import "os"
import "path/filepath"
import "strings"

func countQuestionTypes(questions []dns.Question) int {

	var result int = 0

	counted := make(map[dns.Type]bool)

	for q := 0; q < len(questions); q++ {

		_, ok := counted[questions[q].Type]

		if ok == true {
			// Do Nothing
		} else {
			counted[questions[q].Type] = true
			result++
		}

	}

	return result

}

func countRecordTypes(records []dns.Record) int {

	var result int = 0

	counted := make(map[dns.Type]bool)

	for r := 0; r < len(records); r++ {

		_, ok := counted[records[r].Type]

		if ok == true {
			// Do Nothing
		} else {
			counted[records[r].Type] = true
			result++
		}

	}

	return result

}

func readDomainRecords(file string) []dns.Record {

	var records []dns.Record

	stat, err1 := os.Stat(file)

	if err1 == nil && !stat.IsDir() {

		buffer, err2 := os.ReadFile(file)

		if err2 == nil {
			json.Unmarshal(buffer, &records)
		}

	}

	return records

}

func writeDomainRecords(file string, records []dns.Record) bool {

	var result bool = false

	buffer, err1 := json.MarshalIndent(records, "", "\t")

	if err1 == nil {

		dir := filepath.Dir(file)

		stat, err2 := os.Stat(dir)

		if err2 == nil && stat.IsDir() {

			err3 := os.WriteFile(file, buffer, 0666)

			if err3 == nil {
				result = true
			}

		} else {

			err3 := os.Mkdir(dir, 0750)

			if err3 == nil {

				err4 := os.WriteFile(file, buffer, 0666)

				if err4 == nil {
					result = true
				}

			}

		}


	}

	return result

}

type ResolverCache struct {
	Folder string `json:"folder"`
}

func NewResolverCache(folder string) ResolverCache {

	var cache ResolverCache

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

func (cache *ResolverCache) Exists(query dns.Packet) bool {

	var result bool = false

	if len(query.Questions) > 0 {

		result = true

		for q := 0; q < len(query.Questions); q++ {

			question := query.Questions[q]
			resolved := question.Name + "/" + question.Type.String() + ".json"
			found := false

			if resolved != "" {

				stat, err := os.Stat(cache.Folder + "/" + resolved)

				if err == nil && !stat.IsDir() {
					found = true
				}

			}

			if found == false {
				result = false
			}

		}

	}

	return result

}

func (cache *ResolverCache) Read(query dns.Packet) dns.Packet {

	var response dns.Packet

	if query.Type == "query" && len(query.Questions) > 0 {

		response.SetIdentifier(query.Identifier)
		response.SetType("response")
		response.Flags.AuthorativeAnswer = false
		response.Flags.Truncated = false

		if query.Flags.RecursionDesired == true {
			response.Flags.RecursionAvailable = true
		}

		for q := 0; q < len(query.Questions); q++ {

			question := query.Questions[q]
			resolved := question.Name + "/" + question.Type.String() + ".json"

			if resolved != "" {

				records := readDomainRecords(cache.Folder + "/" + resolved)

				if len(records) > 0 {

					response.AddQuestion(question)

					for r := 0; r < len(records); r++ {
						response.AddAnswer(records[r])
					}

				}

			}

		}

		length_answers := countRecordTypes(response.Answers)

		if countQuestionTypes(query.Questions) == length_answers {
			response.SetResponseCode(dns.ResponseCodeNoError)
		} else if length_answers > 0 {
			response.SetResponseCode(dns.ResponseCodeNoError)
		} else if length_answers == 0 {
			response.SetResponseCode(dns.ResponseCodeNonExistDomain)
		}

	}

	return response

}

func (cache *ResolverCache) Write(response dns.Packet) bool {

	var result bool = false

	if response.Type == "response" && len(response.Answers) > 0 {

		for a := 0; a < len(response.Answers); a++ {

			record := response.Answers[a]
			resolved := record.Name + "/" + record.Type.String() + ".json"

			if resolved != "" {

				records := readDomainRecords(cache.Folder + "/" + resolved)
				changed := false

				if len(records) > 0 {

					found := false

					for r := 0; r < len(records); r++ {

						if bytes.Equal(records[r].Data, record.Data) {
							found = true
							break
						}

					}

					if found == false {
						record.TTL = uint32(0)
						records = append(records, record)
						changed = true
					}

				} else {

					record.TTL = uint32(0)
					records = append(records, record)
					changed = true

				}

				if changed == true {
					writeDomainRecords(cache.Folder + "/" + resolved, records)
					result = true
				}

			}

		}

	}

	return result

}
