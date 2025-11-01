package transition

type Option struct {
	waitForComplete bool
}

type OptionPort func(string) *Option

func NewOption(wait bool) *Option {
	return &Option{
		waitForComplete: wait,
	}
}

func (o *Option) WaitForComplete() bool {
	if o == nil {
		return false
	}
	return o.waitForComplete
}
