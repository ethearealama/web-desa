package models

import "time"

type Warga struct {
    ID          int
    NIK         string
    NamaLengkap string
    Alamat      string
    NoHP        string
    NoKK        string
    RW          string
    CreatedAt   time.Time
    UpdatedAt   time.Time
}