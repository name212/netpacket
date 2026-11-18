package params

import (
	"flag"
	"fmt"
)

type (
	lookupF func(n string) *flag.Flag
	stringF func(n, def, help string) *string
	parsedF func() bool

	flagSet struct {
		set  *flag.FlagSet
		name string

		lookup  lookupF
		fString stringF
		parsed  parsedF
	}
)

func newFlagSet(set []*flag.FlagSet) (*flagSet, error) {
	var lookup lookupF = flag.Lookup
	var fString stringF = flag.String
	var parsed parsedF = flag.Parsed

	name := "global flags"

	var s *flag.FlagSet

	if len(set) > 0 {
		s = set[0]
		lookup = s.Lookup
		fString = s.String
		parsed = s.Parsed
		name = fmt.Sprintf("flag set %s", s.Name())
	}

	if parsed() {
		return nil, fmt.Errorf("%s already parsed", name)
	}

	return &flagSet{
		set:     s,
		name:    name,
		lookup:  lookup,
		fString: fString,
		parsed:  parsed,
	}, nil
}

func (f *flagSet) Lookup(n, help string) error {
	if e := flag.Lookup(n); e != nil {
		return fmt.Errorf("flag '%s' already present for %s with help '%s'", n, f.name, help)
	}

	return nil
}

func (f *flagSet) String(n, help string) *string {
	return f.fString(n, "", help)
}

func (f *flagSet) Parsed() bool {
	return f.parsed()
}

func (f *flagSet) NotParsed() error {
	if !f.set.Parsed() {
		return fmt.Errorf("extract bind address before parse %s", f.name)
	}

	return nil
}

func (f *flagSet) initFlags(names []string, destPtr *[]*string, help string) error {
	dest := *destPtr

	for _, n := range names {
		if err := f.Lookup(n, help); err != nil {
			return err
		}

		f := f.String(n, help)
		dest = append(dest, f)
	}

	return nil
}
