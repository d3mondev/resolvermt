package multidns_test

import (
	"fmt"

	"github.com/d3mondev/multidns"
)

func Example() {
	resolvers := []string{
		"8.8.8.8",
		"9.9.9.10",
		"1.1.1.1",
	}

	domains := []string{
		"www.google.com",
		"facebook.com",
		"uber.com",
		"apple.com",
		"blogger.com",
		"instagram.com",
		"linkedin.com",
		"en.wikipedia.org",
		"cloudflare.com",
		"wordpress.org",
	}

	client := multidns.New(resolvers, 3, 10, 5)

	results := client.Resolve(domains, multidns.TypeA)

	for _, record := range results {
		fmt.Printf("%s contains %s record %s\n", record.Question, record.Type, record.Answer)
	}
}
