package text

import (
	"bytes"
	"fmt"
	"strconv"
)

// Book A struct for holding book info from scripture.
// This structure is based on the Digital Bible Platform /library/books API.
type Book struct {
	DamID            string `json:"dam_id"`
	BookID           string `json:"book_id"`
	BookName         string `json:"book_name"`
	BookOrder        string `json:"book_order"`
	NumberOfChapters int    `json:"number_of_chapters,string"`
	Chapters         string `json:"chapters"`
}

/*
   "dam_id": "ENGESVO2ET",
   "book_id": "Gen",
   "book_name": "Genesis",
   "book_order": "1",
   "number_of_chapters": "50",
   "chapters": "1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23,24,25,26,27,28,29,30,31,32,33,34,35,36,37,38,39,40,41,42,43,44,45,46,47,48,49,50"
*/
// Verse A struct for holding a verse of scripture.
type Verse struct {
	ID           string `json:"id"`
	BookName     string `json:"book_name"`
	BookID       string `json:"book_id"`
	BookOrder    int    `json:"book_order,string"`
	ChapterID    int    `json:"chapter_id,string"`
	ChapterTitle string `json:"chapter_title"`
	Verse        int    `json:"verse_id,string"`
	VerseText    string `json:"verse_text"`
	ParagraphNbr int    `json:"paragraph_number,string"`
}

// Stringer method for Verse.
func (v Verse) String() string {
	buf := bytes.Buffer{}
	buf.WriteString(fmt.Sprintf("\nBookName: %s\n", v.BookName))
	buf.WriteString(fmt.Sprintf("\tBookId: %s\n", v.BookID))
	buf.WriteString(fmt.Sprintf("\tChapter: %d\n", v.ChapterID))
	buf.WriteString(fmt.Sprintf("\tVerse: %d\n", v.Verse))
	buf.WriteString(fmt.Sprintf("\tText: %s", v.VerseText[:25]))
	return buf.String()
}

// GetID Getter for Verse unique identifier.
// Uses BookName_ChapterID_Verse as the ID
func (v *Verse) GetID() string {
	if v.ID == "" {
		v.ID = v.BookName + "_" + strconv.Itoa(v.ChapterID) + "_" + strconv.Itoa(v.Verse)
	}
	return v.ID
}

// Collection A struct for a collection of verses, ie. multiple verses, chapter or entire books.
type Collection struct {
	Verses []*Verse
}
