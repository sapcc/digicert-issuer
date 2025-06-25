// SPDX-FileCopyrightText: 2025 SAP SE or an SAP affiliate company
// SPDX-License-Identifier: Apache-2.0

// SPDX-FileCopyrightText: 2022 SAP SE or an SAP affiliate company
// SPDX-License-Identifier: Apache-2.0

package version

import (
	"bytes"
	"runtime"
	"strings"
	"text/template"
)

var (
	Version   string
	Revision  string
	Branch    string
	BuildDate string
	GoVersion = runtime.Version()
)

var versionInfoTmpl = `
{{.program}}, version {{.version}} (branch: {{.branch}}, revision: {{.revision}})
  build date:       {{.buildDate}}
  go version:       {{.goVersion}}
`

// Print returns version information.
func Print(program string) string {
	m := map[string]string{
		"program":   program,
		"version":   Version,
		"revision":  Revision,
		"branch":    Branch,
		"buildDate": BuildDate,
		"goVersion": GoVersion,
	}
	t := template.Must(template.New("version").Parse(versionInfoTmpl))

	var buf bytes.Buffer
	if err := t.ExecuteTemplate(&buf, "version", m); err != nil {
		panic(err)
	}
	return strings.TrimSpace(buf.String())
}
