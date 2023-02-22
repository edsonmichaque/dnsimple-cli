package formatter

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/dnsimple/dnsimple-go/dnsimple"
)

type DomainList dnsimple.DomainsResponse

func (a DomainList) FormatJSON(opts *Options) (io.Reader, error) {
	return formatJSON(a, opts)
}

func (a DomainList) FormatYAML(opts *Options) (io.Reader, error) {
	return formatYAML(a, opts)
}

func (a DomainList) FormatTable(_ *Options) (io.Reader, error) {
	return formatTable(a)
}

func (a DomainList) formatJSON(opts *Options) ([]byte, error) {
	return json.MarshalIndent(a.Data, "", "  ")
}

func (a DomainList) formatHeader() []string {
	return []string{
		"ID",
		"ACCOUNT ID",
		"REGISTRANT ID",
		"NAME",
		"UNICODE NAME",
		"TOKEN",
		"STATE",
		"AUTO RENEW",
		"PRIVATE WHOIS",
		"EXPIRES AT",
		"CREATED AT",
		"UPDATED AT",
	}
}

func (a DomainList) formatRows() []map[string]string {
	data := make([]map[string]string, 0, len(a.Data))

	domains := a.Data

	for i := range domains {
		data = append(data, map[string]string{
			"ID":            fmt.Sprintf("%d", domains[i].ID),
			"ACCOUNT ID":    fmt.Sprintf("%d", domains[i].AccountID),
			"REGISTRANT ID": fmt.Sprintf("%d", domains[i].RegistrantID),
			"NAME":          domains[i].Name,
			"UNICODE NAME":  domains[i].UnicodeName,
			"TOKEN":         domains[i].Token,
			"STATE":         domains[i].State,
			"AUTO RENEW":    fmt.Sprintf("%t", domains[i].AutoRenew),
			"PRIVATE WHOIS": fmt.Sprintf("%t", domains[i].PrivateWhois),
			"EXPIRES AT":    domains[i].ExpiresAt,
			"CREATED AT":    domains[i].CreatedAt,
			"UPDATED AT":    domains[i].UpdatedAt,
		})
	}

	return data
}
