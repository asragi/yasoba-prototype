package transition

type Controller interface {
	Start()
	IsRunning() bool
}

type controller struct {
	start     func()
	isRunning func() bool
}

func (c *controller) Start() {
	if c.start != nil {
		c.start()
	}
}

func (c *controller) IsRunning() bool {
	if c.isRunning == nil {
		return false
	}
	return c.isRunning()
}

func NewController(start func(), isRunning func() bool) Controller {
	return &controller{
		start:     start,
		isRunning: isRunning,
	}
}
