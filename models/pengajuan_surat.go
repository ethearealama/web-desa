package models

import "time"

type PengajuanSurat struct {
    ID            int
    NoPengajuan   string
    Nama          string
    NIK           string
    Alamat        string
    NoHP          string
    JenisSurat    string
    Keperluan     string
    Status        string
    CatatanAdmin  string
    TanggalSelesai *time.Time
    CreatedAt     time.Time
    UpdatedAt     time.Time
}