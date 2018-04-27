package ntsb

type flagSet interface {
	StringVar(p *string, name string, value string, usage string)
	IntVar(p *int, name string, value int, usage string)
	Parse(args []string) error
}
