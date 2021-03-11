package multidns

// Record contains a DNS record.
type Record struct {
	Question string
	Type     RRtype
	Answer   string
}
