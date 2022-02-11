package docset

type Docset struct {
	Name      string
	Version   string
	CountFile int64
}

func NewDocset() *Docset {
	return &Docset{
		Name:    "AWS_JavaScript.xml",
		Version: "2.972.0",
	}
}

func (d *Docset) Collect() error {

	return nil
}
