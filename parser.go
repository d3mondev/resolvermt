package fastdns

import (
	"strings"

	"github.com/miekg/dns"
)

type msgParser struct{}

func (s *msgParser) Parse(query string, msg *dns.Msg) []Record {
	records := []Record{}

	for _, answer := range msg.Answer {
		switch t := answer.(type) {
		case *dns.A:
			rec := Record{
				Question: query,
				Type:     TypeA,
				Answer:   t.A.String(),
			}
			records = append(records, rec)
		case *dns.AAAA:
			rec := Record{
				Question: query,
				Type:     TypeAAAA,
				Answer:   t.AAAA.String(),
			}
			records = append(records, rec)
		case *dns.CNAME:
			rec := Record{
				Question: query,
				Type:     TypeCNAME,
				Answer:   strings.TrimSuffix(t.Target, "."),
			}
			records = append(records, rec)
		case *dns.TXT:
			for _, txt := range t.Txt {
				rec := Record{
					Question: query,
					Type:     TypeTXT,
					Answer:   txt,
				}
				records = append(records, rec)
			}
		case *dns.MX:
			rec := Record{
				Question: query,
				Type:     TypeMX,
				Answer:   strings.TrimSuffix(t.Mx, "."),
			}
			records = append(records, rec)
		case *dns.NS:
			rec := Record{
				Question: query,
				Type:     TypeNS,
				Answer:   strings.TrimSuffix(t.Ns, "."),
			}
			records = append(records, rec)
		}
	}

	return records
}
