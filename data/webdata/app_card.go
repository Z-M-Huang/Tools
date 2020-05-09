package webdata

//AppCardList card list for single category
type AppCardList struct {
	Category string
	AppCards []*AppCard
}

//AppCard used to render app card in pages
type AppCard struct {
	Name            string
	TemplateName    string
	ImageURL        string
	FontsAwesomeTag string
	Link            string
	Title           string
	Description     string
	Used            bool
	AmountUsed      uint64
	Liked           bool
	AmountLiked     uint64
}
