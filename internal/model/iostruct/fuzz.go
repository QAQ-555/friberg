package iostruct

type GameInfo struct {
	ID   string `json:"id" dc:"game id" bson:"_id"`
	Name string `json:"name" dc:"game name" bson:"name"`
}
type Game struct {
	ID          GuessResult   `json:"id"`
	Name        GuessResult   `json:"name"`
	ReleaseDate GuessResult   `json:"release_date"`
	Tags        []GuessResult `json:"tags"`
	Platforms   []GuessResult `json:"platforms"`
}
