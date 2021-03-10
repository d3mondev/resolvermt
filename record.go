package multidns

// Record describes a DNS record
type Record struct {
	Question string
	Type     RRtype
	Answer   string
}
