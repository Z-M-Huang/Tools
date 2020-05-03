package webdata

//AppCardList card list for single category
type AppCardList struct {
	Category string
	AppCards []*AppCard
}

//AppCard used to render app card in pages
type AppCard struct {
	ImageURL    string
	Link        string
	Title       string
	Description string
	Up          uint64
	Saved       uint64
}
