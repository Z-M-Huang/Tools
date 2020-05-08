package webdata

//AppCardList card list for single category
type AppCardList struct {
	Category string
	AppCards []*AppCard
}

//AppCard used to render app card in pages
type AppCard struct {
	ImageURL        string
	FontsAwesomeTag string
	Link            string
	Title           string
	Description     string
	Usage           uint64
	Liked           uint64
}
