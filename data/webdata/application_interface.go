package webdata

//Application application interface
type Application interface {
	GetName() string
	GetDescription() string
	GetUrl() string
	GetIcon() string
}
