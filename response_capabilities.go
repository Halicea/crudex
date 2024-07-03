package crudex

type ResponseCapabilities struct {
	UI     bool
	API    bool
	Layout bool
}

func NewResponseCapabilities() *ResponseCapabilities {
	return &ResponseCapabilities{
		UI:     true,
		API:    true,
		Layout: true,
	}
}

func (self *ResponseCapabilities) HasUI() bool {
	return self.UI
}

func (self *ResponseCapabilities) HasAPI() bool {
	return self.API
}
func (self *ResponseCapabilities) EnableLayoutOnNonHxRequest() bool {
	return self.Layout
}
