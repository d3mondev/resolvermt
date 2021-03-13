package main

import (
	"flag"
	"fmt"
	"os"
	"runtime/pprof"
	"strings"

	"github.com/d3mondev/fastdns"
)

func main() {
	var resolverFile, domainFile, cpuprofileFile string
	var retries, qps, concurrency int
	flag.StringVar(&resolverFile, "resolvers", "", "filename containing resolvers")
	flag.StringVar(&domainFile, "domains", "", "filename containing domains to resolve")
	flag.StringVar(&cpuprofileFile, "cpuprofile", "", "write cpu profile to file")
	flag.IntVar(&retries, "retry", 3, "number of times to retry a DNS query")
	flag.IntVar(&qps, "qps", 10, "maximum number of queries per second for each resolver")
	flag.IntVar(&concurrency, "concurrency", 1000, "number of concurrent DNS requests")
	flag.Parse()

	if resolverFile == "" || domainFile == "" {
		usage()
		os.Exit(1)
	}

	if cpuprofileFile != "" {
		file, err := os.Create(cpuprofileFile)

		if err != nil {
			fmt.Printf("unable to write cpu profile file %s\n\n", cpuprofileFile)
			os.Exit(1)
		}

		pprof.StartCPUProfile(file)
		defer pprof.StopCPUProfile()
	}

	content, err := os.ReadFile(resolverFile)

	if err != nil {
		fmt.Printf("unable to open file %s\n\n", resolverFile)
		usage()
		os.Exit(1)
	}

	resolvers := strings.Split(string(content), "\n")

	content, err = os.ReadFile(domainFile)

	if err != nil {
		fmt.Printf("unable to open file %s\n\n", domainFile)
		usage()
		os.Exit(1)
	}

	domains := strings.Split(string(content), "\n")

	client := fastdns.New(resolvers, retries, qps, concurrency)

	records := client.Resolve(domains, fastdns.TypeA)

	for _, record := range records {
		fmt.Printf("%s %s %s\n", record.Question, record.Type, record.Answer)
	}
}

func usage() {
	fmt.Printf("Usage: %s --resolvers <file> --domains <file> [flags]\n", os.Args[0])
	flag.PrintDefaults()
}
