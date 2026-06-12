package models

import "time"

type Berita struct {
	ID        int       `json:"id"`
	Judul     string    `json:"judul"`
	Slug      string    `json:"slug"`
	Isi       string    `json:"isi"`
	Gambar    string    `json:"gambar"`
	Penulis   string    `json:"penulis"`
	Views     int       `json:"views"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
