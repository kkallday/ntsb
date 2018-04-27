package ntsb

type flagSet interface {
	StringVar(p *string, name string, value string, usage string)
	Parse(args []string) error
}
