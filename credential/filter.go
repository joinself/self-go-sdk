package credential

type Filter struct {
	// TODO mock implementation to demonstrate API
}

func FilterCredentials(credentials []*Credential, filter *Filter) []*Credential {
	return nil
}

func NewFilter() *Filter {
	return &Filter{}
}

func (f *Filter) Equals(field, value string) *Filter              { return f }
func (f *Filter) NotEquals(field, value string) *Filter           { return f }
func (f *Filter) GreaterThan(field, value string) *Filter         { return f }
func (f *Filter) GreaterThanOrEquals(field, value string) *Filter { return f }
func (f *Filter) LessThan(field, value string) *Filter            { return f }
func (f *Filter) LessThanOrEquals(field, value string) *Filter    { return f }
func (f *Filter) Matches(credential *Credential) bool             { return true }
