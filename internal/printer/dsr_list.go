package printer

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/dnsimple/dnsimple-go/dnsimple"
)

type DSRList dnsimple.DelegationSignerRecordsResponse

func (a DSRList) PrintJSON(opts *Options) (io.Reader, error) {
	return printJSON(a, opts)
}

func (a DSRList) PrintYAML(opts *Options) (io.Reader, error) {
	return printYAML(a, opts)
}

func (a DSRList) PrintTable(_ *Options) (io.Reader, error) {
	return printTable(a)
}

func (a DSRList) toJSON(opts *Options) ([]byte, error) {
	return json.MarshalIndent(a.Data, "", "  ")
}

func (a DSRList) printHeader() []string {
	return []string{
		"ID",
		"DOMAIN ID",
		"ALGORITHM",
		"DIGEST",
		"DIGEST TYPE",
		"KEYTAG",
		"PUBLIC KEY",
		"CREATED AT",
		"UPDATED AT",
	}
}

func (a DSRList) printRows() []map[string]string {
	data := make([]map[string]string, 0, len(a.Data))

	dsr := a.Data

	const txtLen = 10

	for i := range dsr {
		data = append(data, map[string]string{
			"ID":          fmt.Sprintf("%d", dsr[i].ID),
			"DOMAIN ID":   fmt.Sprintf("%d", dsr[i].DomainID),
			"ALGORITHM":   dsr[i].Algorithm,
			"DIGEST":      truncate(dsr[i].Digest, txtLen),
			"DIGEST TYPE": dsr[i].DigestType,
			"KEYTAG":      dsr[i].Keytag,
			"PUBLIC KEY":  truncate(dsr[i].PublicKey, txtLen),
			"CREATED AT":  dsr[i].CreatedAt,
			"UPDATED AT":  dsr[i].UpdatedAt,
		})
	}

	return data
}

func truncate(s string, length int) string {
	if len(s) <= length {
		return s
	}

	return s[:length] + "..."
}
