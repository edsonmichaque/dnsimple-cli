package printer

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/dnsimple/dnsimple-go/dnsimple"
)

type DomainList dnsimple.DomainsResponse

func (a DomainList) PrintJSON(opts *Options) (io.Reader, error) {
	return printJSON(a, opts)
}

func (a DomainList) PrintYAML(opts *Options) (io.Reader, error) {
	return printYAML(a, opts)
}

func (a DomainList) PrintTable(_ *Options) (io.Reader, error) {
	return printTable(a)
}

func (a DomainList) toJSON(opts *Options) ([]byte, error) {
	return json.MarshalIndent(a.Data, "", "  ")
}

func (a DomainList) printHeader() []string {
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

func (a DomainList) printRows() []map[string]string {
	data := make([]map[string]string, 0, len(a.Data))

	for _, k := range a.Data {
		data = append(data, map[string]string{
			"ID":            fmt.Sprintf("%d", k.ID),
			"ACCOUNT ID":    fmt.Sprintf("%d", k.AccountID),
			"REGISTRANT ID": fmt.Sprintf("%d", k.RegistrantID),
			"NAME":          k.Name,
			"UNICODE NAME":  k.UnicodeName,
			"TOKEN":         k.Token,
			"STATE":         k.State,
			"AUTO RENEW":    fmt.Sprintf("%t", k.AutoRenew),
			"PRIVATE WHOIS": fmt.Sprintf("%t", k.PrivateWhois),
			"EXPIRES AT":    k.ExpiresAt,
			"CREATED AT":    k.CreatedAt,
			"UPDATED AT":    k.UpdatedAt,
		})
	}

	return data
}
