package formatter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/dnsimple/dnsimple-go/dnsimple"
)

type Whoami dnsimple.WhoamiResponse

func (w Whoami) FormatText(opts *Options) (io.Reader, error) {
	var data [][2]interface{}

	if user := w.Data.User; user != nil {
		data = [][2]interface{}{
			{"User", ""},
			{"  ID", user.ID},
			{"  Email", user.Email},
		}
	}

	if account := w.Data.Account; account != nil {
		data = [][2]interface{}{
			{"Accout", ""},
			{"  ID", account.ID},
			{"  Email", account.Email},
			{"  Plan identifier", account.PlanIdentifier},
			{"  Created at", account.CreatedAt},
			{"  Updated at", account.UpdatedAt},
		}
	}

	buf := new(bytes.Buffer)

	for _, i := range data {
		buf.WriteString(fmt.Sprintf("%-18s %v\n", i[0].(string)+":", i[1]))
	}

	return buf, nil
}

func (w Whoami) FormatJSON(opts *Options) (io.Reader, error) {
	return formatJSON(w, opts)
}

func (w Whoami) FormatYAML(opts *Options) (io.Reader, error) {
	return formatYAML(w, opts)
}

func (w Whoami) formatJSON(opts *Options) ([]byte, error) {
	return json.MarshalIndent(w.Data, "", "  ")
}
