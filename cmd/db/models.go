// Models package for holding structs for marchelling
package models

// Users the struct for holding Users information.
type Users struct {
	ID int,
	FirstName string `db:"first_name"`
	Lastname string `db:"last_name"`
	Username string `db:"username"`
	Email string `db:"email"`
	Active bool
	CreatedDt int `db:"created_dt"`
	ModifiedDt int `db:"modified_dt"`
}

// Books the struct for mapping books to their abbreviations.
type Books struct {
	ID int
	Name string
	Abbreviation string
}

// Scriptures the struct for holding the scripture text for each verse.
type Scriptures struct {
	ID int
	BookId int
	Chapter int
	Verse int
	Text string
}

// Comments the struct for users comments on a verse 
type Comments struct {
	ID int
	UserId int
	ScriptureId int
	Comment string
	CreatedDt int `db:"created_dt"`
	ModifiedDt int `db:"modified_dt"`
	Active bool
}