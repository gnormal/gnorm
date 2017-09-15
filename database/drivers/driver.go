package drivers

// Templates is an interface for gnorm templates. The templates vary depending
// on the active driver.
type Templates interface {
	TplNames() []string
	Tpl(name string) ([]byte, error)
}
