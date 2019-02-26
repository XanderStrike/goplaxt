package plexhooks

type Account struct {
	Id    int
	Thumb string
	Title string
}

type Server struct {
	Title string
	Uuid  string
}

type Player struct {
	Local         bool
	PublicAddress string
	Title         string
	Uuid          string
}

type GenericItem struct {
	Id     int
	Filter string
	Tag    string
	Role   string
	Thumb  string
	Count  int
}

type Genre GenericItem
type Director GenericItem
type Writer GenericItem
type Producer GenericItem
type Country GenericItem
type Similar GenericItem
type Role GenericItem

type Metadata struct {
	LibrarySectionType    string
	RatingKey             string
	Key                   string
	ParentRatingKey       string
	GrandparentRatingKey  string
	Guid                  string
	LibrarySectionTitle   string
	LibrarySectionID      int
	LibrarySectionKey     string
	Studio                string
	Type                  string
	Title                 string
	TitleSort             string
	GrandparentKey        string
	ParentKey             string
	GrandparentTitle      string
	ParentTitle           string
	ContentRating         string
	Summary               string
	Tagline               string
	Index                 int
	ParentIndex           int
	Rating                float32
	RatingCount           int
	AudienceRating        float32
	ViewOffset            int
	ViewCount             int
	LastViewedAt          int
	Year                  int
	Thumb                 string
	Art                   string
	Duration              int
	ParentThumb           string
	GrandparentThumb      string
	GrandparentArt        string
	GrandparentTheme      string
	OriginallyAvailableAt string
	AddedAt               int
	UpdatedAt             int
	AudienceRatingImage   string
	PrimaryExtraKey       string
	RatingImage           string

	Genre    []Genre
	Director []Director
	Writer   []Writer
	Producer []Producer
	Country  []Country

	Similar []Similar

	Role []Role
}

type PlexResponse struct {
	Event    string
	User     bool
	Owner    bool
	Account  Account
	Metadata Metadata
	Server   Server
	Player   Player
}
