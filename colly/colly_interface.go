package colly

type collyType uint8

const (
	ExampleColl collyType = iota + 1
	DetailColl
)

type CollyBase struct {
	Url      string
	StoreDir string
}

type Options struct {
	Url      string
	StoreDir string
}

func (c *CollyBase) SetUrl(url string) {
	c.Url = url
}

func (c *CollyBase) SetStory(dir string) {
	c.StoreDir = dir
}

type IColly interface {
	CollyPage()
	SetUrl(str string)
	SetStory(str string)
}

func GetColly(t collyType, op Options) IColly {
	if op.Url != "" {

	}
	var c IColly
	switch t {
	case ExampleColl:
		c = &Example{}
	case DetailColl:
		c = &DetailPage{}
	default:
		panic("param not enough")
	}
	if op.Url != "" {
		c.SetUrl(op.Url)
	} else {
		panic("visit url is required")
	}
	if op.StoreDir != "" {
		c.SetStory(op.StoreDir)
	}
	return c
}
