package formatter

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/dnsimple/dnsimple-go/dnsimple"
)

type Filter struct {
	Fields []string
}

type AccountList dnsimple.AccountsResponse

func (a AccountList) FormatJSON(opts *Options) (io.Reader, error) {
	return formatJSON(a, opts)
}

func (a AccountList) FormatYAML(opts *Options) (io.Reader, error) {
	return formatYAML(a, opts)
}

func (a AccountList) FormatTable(_ *Options) (io.Reader, error) {
	return formatTable(a)
}

func (a AccountList) formatJSON(opts *Options) ([]byte, error) {
	return json.MarshalIndent(a.Data, "", "  ")
}

func (a AccountList) formatHeader() []string {
	return []string{
		"ID",
		"EMAIL",
		"PLAN IDENTIFIER",
		"CREATED AT",
		"UPDATED AT",
	}
}

func (a AccountList) formatRows() []map[string]string {
	data := make([]map[string]string, 0, len(a.Data))

	for _, k := range a.Data {
		data = append(data, map[string]string{
			"ID":              fmt.Sprintf("%d", k.ID),
			"EMAIL":           k.Email,
			"PLAN IDENTIFIER": k.PlanIdentifier,
			"CREATED AT":      k.CreatedAt,
			"UPDATED AT":      k.UpdatedAt,
		})
	}

	return data
}
