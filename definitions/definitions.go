package definitions

type Member struct {
	EmailAddress string `datastore:"emailAddress"`
	FirstName    string `datastore:"firstname"`
	LastName     string `datastore:"lastname"`
	ID           int64  `datastore:"id"`
	TeamID       int64  `datastore:"teamID"`
}

type Team struct {
	TeamName    string   `datastore:"name"`
	TeamID      int64    `datastore:"id"`
	Description string   `datastore:"description"`
	Members     []Member `datastore:"members"`
}
