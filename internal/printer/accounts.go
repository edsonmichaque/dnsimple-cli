package printer

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

func (a AccountList) PrintJSON(opts *Options) (io.Reader, error) {
	return printJSON(a, opts)
}

func (a AccountList) PrintYAML(opts *Options) (io.Reader, error) {
	return printYAML(a, opts)
}

func (a AccountList) PrintTable(_ *Options) (io.Reader, error) {
	return printTable(a)
}

func (a AccountList) toJSON(opts *Options) ([]byte, error) {
	return json.MarshalIndent(a.Data, "", "  ")
}

func (a AccountList) printHeader() []string {
	return []string{
		"ID",
		"EMAIL",
		"PLAN IDENTIFIER",
		"CREATED AT",
		"UPDATED AT",
	}
}

func (a AccountList) printRows() []map[string]string {
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
