package models

type Movie struct {
	TitleType string `json:"titleType"`
	Title     string `json:"primaryTitle"`
	Year      string `json:"startYear"`
	Genres    string `json:"genres"`
	Runtime   string `json:"runTimeMinutes"`
	Rating    string `json:"averageRating"`
	Votes     string `json:"numVotes"`
}

type MovieQueryParams struct {
	Genre  string
	Year   string
	Rating float32
	Amount int
}
