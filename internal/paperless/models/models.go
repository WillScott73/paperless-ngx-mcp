package models

import "time"

type DocumentList struct {
	Count    int        `json:"count"`
	Next     *string    `json:"next"`
	Previous *string    `json:"previous"`
	Results  []Document `json:"results"`
}

type Document struct {
	ID                 int       `json:"id"`
	Title              string    `json:"title"`
	Content            string    `json:"content"`
	Created            time.Time `json:"created"`
	Modified           time.Time `json:"modified"`
	Added              time.Time `json:"added"`
	ArchiveSerialNumber *int      `json:"archive_serial_number"`
	DocumentType       *int      `json:"document_type"`
	Correspondent      *int      `json:"correspondent"`
	StoragePath        *int      `json:"storage_path"`
	Tags               []int     `json:"tags"`
	Notes              []Note    `json:"notes"`
}

type Note struct {
	ID      int       `json:"id"`
	Note    string    `json:"note"`
	Created time.Time `json:"created"`
	User    int       `json:"user"`
}

type TagList struct {
	Count   int   `json:"count"`
	Results []Tag `json:"results"`
}

type Tag struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

type CorrespondentList struct {
	Count   int             `json:"count"`
	Results []Correspondent `json:"results"`
}

type Correspondent struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type DocumentTypeList struct {
	Count   int            `json:"count"`
	Results []DocumentType `json:"results"`
}

type DocumentType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

