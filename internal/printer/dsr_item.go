package printer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/dnsimple/dnsimple-go/dnsimple"
)

/*
 * 	ID         int64  `json:"id,omitempty"`
	DomainID   int64  `json:"domain_id,omitempty"`
	Algorithm  string `json:"algorithm"`
	Digest     string `json:"digest,omitempty"`
	DigestType string `json:"digest_type,omitempty"`
	Keytag     string `json:"keytag,omitempty"`
	PublicKey  string `json:"public_key,omitempty"`
	CreatedAt  string `json:"created_at,omitempty"`
	UpdatedAt  string `json:"updated_at,omitempty"`

*/

type DSRItem dnsimple.DelegationSignerRecordResponse

func (d DSRItem) PrintText(opts *Options) (io.Reader, error) {
	keys := []string{
		"id",
		"domain_id",
		"algorithm",
		"digest",
		"digest_type",
		"keytag",
		"public_key",
		"created_at",
		"updated_at",
	}

	const txtLen = 8

	values := map[string]interface{}{
		"id":          d.Data.ID,
		"domain_id":   d.Data.DomainID,
		"algorithm":   d.Data.Algorithm,
		"digest":      truncate(d.Data.Digest, txtLen),
		"digest_type": d.Data.DigestType,
		"keytag":      d.Data.Keytag,
		"public_key":  truncate(d.Data.PublicKey, txtLen),
		"created_at":  d.Data.CreatedAt,
		"updated_at":  d.Data.UpdatedAt,
	}

	titles := map[string]string{
		"id":          "ID",
		"domain_id":   "Domain ID",
		"algorithm":   "Algorithm",
		"digest":      "Digest",
		"digest_type": "Digest type",
		"keytag":      "Keytag",
		"public_key":  "Public key",
		"created_at":  "Created at",
		"updated_at":  "Updated at",
	}

	buf := new(bytes.Buffer)
	for _, v := range keys {
		buf.WriteString(fmt.Sprintf("%-20s%v\n", titles[v]+":", values[v]))
	}

	return buf, nil
}

func (d DSRItem) PrintJSON(opts *Options) (io.Reader, error) {
	return printJSON(d, opts)
}

func (d DSRItem) PrintYAML(opts *Options) (io.Reader, error) {
	return printYAML(d, opts)
}

func (d DSRItem) toJSON(opts *Options) ([]byte, error) {
	return json.MarshalIndent(d.Data, "", "  ")
}
