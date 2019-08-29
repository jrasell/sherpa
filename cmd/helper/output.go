package helper

import "github.com/ryanuber/columnize"

// FormatList is used by the CLI to correctly format list items with correct spacing.
func FormatList(in []string) string {
	columnConf := columnize.DefaultConfig()
	columnConf.Empty = "<none>"
	return columnize.Format(in, columnConf)
}

// FormatKV is used by the CLI to correctly format key/value items with correct indentation.
func FormatKV(in []string) string {
	columnConf := columnize.DefaultConfig()
	columnConf.Empty = "<none>"
	columnConf.Glue = " = "
	return columnize.Format(in, columnConf)
}
