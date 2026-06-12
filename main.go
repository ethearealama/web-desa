package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"go-desa/config"
	"go-desa/models"

	"github.com/gorilla/mux"
)

var templates *template.Template

func main() {
	// Koneksi database
	config.ConnectDB()
	defer config.DB.Close()

	// Template dengan function tambahan
	funcMap := template.FuncMap{
		"add": func(a, b int) int { return a + b },
		"substring": func(s string, start, end int) string {
			runes := []rune(s)
			if len(runes) < end {
				return s
			}
			return string(runes[start:end])
		},
	}

	templates = template.Must(template.New("").Funcs(funcMap).ParseGlob("views/*.html"))
	template.Must(templates.ParseGlob("views/layouts/*.html"))

	// Router
	r := mux.NewRouter()

	// Static files
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Routes
	r.HandleFunc("/", homeHandler).Methods("GET")
	r.HandleFunc("/profil", profilHandler).Methods("GET")
	r.HandleFunc("/berita", beritaHandler).Methods("GET")
	r.HandleFunc("/berita/{id}", beritaDetailHandler).Methods("GET")
	r.HandleFunc("/wisata", wisataHandler).Methods("GET")
	r.HandleFunc("/pengaduan", pengaduanHandler).Methods("GET")
	r.HandleFunc("/pengaduan/kirim", kirimPengaduanHandler).Methods("POST")
	r.HandleFunc("/bansos", bansosHandler).Methods("GET")
	r.HandleFunc("/bansos/cari", cariBansosHandler).Methods("POST")
	r.HandleFunc("/surat", suratHandler).Methods("GET")
	r.HandleFunc("/surat/kirim", kirimSuratHandler).Methods("POST")
	r.HandleFunc("/warga", wargaHandler).Methods("GET")
	r.HandleFunc("/warga/{id}", wargaDetailHandler).Methods("GET")

	fmt.Println("🌐 Server running on http://localhost:9090")
	log.Fatal(http.ListenAndServe(":9090", r))
}

// ==================== HOME ====================
func homeHandler(w http.ResponseWriter, r *http.Request) {
	stat := getStatistik()
	berita := getBeritaTerbaru()

	data := map[string]interface{}{
		"Title":         "Beranda",
		"TotalPenduduk": stat.JumlahPenduduk,
		"TotalKK":       stat.JumlahKK,
		"TotalRW":       stat.JumlahRW,
		"TotalWisata":   stat.JumlahWisata,
		"BeritaList":    berita,
	}
	templates.ExecuteTemplate(w, "index.html", data)
}

func getStatistik() models.StatistikDesa {
	var stat models.StatistikDesa
	row := config.DB.QueryRow("SELECT jumlah_penduduk, jumlah_kk, jumlah_rw, jumlah_wisata FROM statistik_desas LIMIT 1")
	err := row.Scan(&stat.JumlahPenduduk, &stat.JumlahKK, &stat.JumlahRW, &stat.JumlahWisata)
	if err != nil {
		return models.StatistikDesa{5240, 1450, 12, 8}
	}
	return stat
}

func getBeritaTerbaru() []models.Berita {
	rows, err := config.DB.Query("SELECT id, judul, slug, isi, gambar, penulis, views, created_at FROM beritas ORDER BY created_at DESC LIMIT 3")
	if err != nil {
		return []models.Berita{}
	}
	defer rows.Close()

	var beritaList []models.Berita
	for rows.Next() {
		var b models.Berita
		rows.Scan(&b.ID, &b.Judul, &b.Slug, &b.Isi, &b.Gambar, &b.Penulis, &b.Views, &b.CreatedAt)
		beritaList = append(beritaList, b)
	}
	return beritaList
}

// ==================== PROFIL ====================
func profilHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{"Title": "Profil Desa"}
	templates.ExecuteTemplate(w, "profil.html", data)
}

// ==================== BERITA ====================
func beritaHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := config.DB.Query("SELECT id, judul, slug, isi, gambar, penulis, views, created_at FROM beritas ORDER BY created_at DESC")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var beritaList []models.Berita
	for rows.Next() {
		var b models.Berita
		rows.Scan(&b.ID, &b.Judul, &b.Slug, &b.Isi, &b.Gambar, &b.Penulis, &b.Views, &b.CreatedAt)
		beritaList = append(beritaList, b)
	}

	data := map[string]interface{}{
		"Title":      "Berita Desa",
		"BeritaList": beritaList,
	}
	templates.ExecuteTemplate(w, "berita.html", data)
}

func beritaDetailHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	var b models.Berita
	row := config.DB.QueryRow("SELECT id, judul, slug, isi, gambar, penulis, views, created_at FROM beritas WHERE id = ?", id)
	err := row.Scan(&b.ID, &b.Judul, &b.Slug, &b.Isi, &b.Gambar, &b.Penulis, &b.Views, &b.CreatedAt)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	config.DB.Exec("UPDATE beritas SET views = views + 1 WHERE id = ?", id)

	data := map[string]interface{}{
		"Title":  b.Judul,
		"Berita": b,
	}
	templates.ExecuteTemplate(w, "berita_detail.html", data)
}

// ==================== WISATA ====================
func wisataHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := config.DB.Query("SELECT id, nama_wisata, deskripsi, lokasi, gambar, harga_tiket, jam_operasional, status FROM wisatas")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var wisataList []models.Wisata
	for rows.Next() {
		var w models.Wisata
		rows.Scan(&w.ID, &w.NamaWisata, &w.Deskripsi, &w.Lokasi, &w.Gambar, &w.HargaTiket, &w.JamOperasional, &w.Status)
		wisataList = append(wisataList, w)
	}

	data := map[string]interface{}{
		"Title":      "Wisata Desa",
		"WisataList": wisataList,
	}
	templates.ExecuteTemplate(w, "wisata.html", data)
}

// ==================== PENGADUAN ====================
func pengaduanHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{"Title": "Pengaduan"}
	templates.ExecuteTemplate(w, "pengaduan.html", data)
}

func kirimPengaduanHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	nama := r.FormValue("nama")
	nik := r.FormValue("nik")
	no_hp := r.FormValue("no_hp")
	judul := r.FormValue("judul")
	isi := r.FormValue("isi")

	_, err := config.DB.Exec("INSERT INTO pengaduan (nama, nik, no_hp, judul, isi, status) VALUES (?, ?, ?, ?, ?, 'baru')",
		nama, nik, no_hp, judul, isi)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/pengaduan?success=true", http.StatusSeeOther)
}

// ==================== BANSOS ====================
func bansosHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{"Title": "Cek Bansos"}
	templates.ExecuteTemplate(w, "bansos.html", data)
}

func cariBansosHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	nik := r.FormValue("nik")

	var b models.Bansos
	row := config.DB.QueryRow("SELECT id, nik, nama_penerima, jenis_bantuan, alamat, status FROM bansos WHERE nik = ?", nik)
	err := row.Scan(&b.ID, &b.NIK, &b.NamaPenerima, &b.JenisBantuan, &b.Alamat, &b.Status)

	data := map[string]interface{}{
		"Title": "Cek Bansos",
		"Hasil": b,
		"Cek":   err == nil,
		"Nik":   nik,
	}
	templates.ExecuteTemplate(w, "bansos.html", data)
}

// ==================== SURAT ====================
func suratHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{"Title": "Layanan Surat"}
	templates.ExecuteTemplate(w, "surat.html", data)
}

func kirimSuratHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	nama := r.FormValue("nama")
	nik := r.FormValue("nik")
	alamat := r.FormValue("alamat")
	no_hp := r.FormValue("no_hp")
	jenis_surat := r.FormValue("jenis_surat")
	keperluan := r.FormValue("keperluan")

	no_pengajuan := fmt.Sprintf("SURAT/%s/%d", time.Now().Format("20060102"), time.Now().UnixNano()%1000)

	_, err := config.DB.Exec("INSERT INTO pengajuan_surat (no_pengajuan, nama, nik, alamat, no_hp, jenis_surat, keperluan, status) VALUES (?, ?, ?, ?, ?, ?, ?, 'pending')",
		no_pengajuan, nama, nik, alamat, no_hp, jenis_surat, keperluan)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/surat?success=true&no_pengajuan="+no_pengajuan, http.StatusSeeOther)
}

// ==================== WARGA ====================
func wargaHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := config.DB.Query("SELECT id, nik, nama_lengkap, alamat, no_hp, no_kk, rw, created_at FROM warga ORDER BY id DESC")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var wargaList []models.Warga
	for rows.Next() {
		var w models.Warga
		rows.Scan(&w.ID, &w.NIK, &w.NamaLengkap, &w.Alamat, &w.NoHP, &w.NoKK, &w.RW, &w.CreatedAt)
		wargaList = append(wargaList, w)
	}

	var totalKK int
	config.DB.QueryRow("SELECT COUNT(DISTINCT no_kk) FROM warga WHERE no_kk IS NOT NULL AND no_kk != ''").Scan(&totalKK)

	var totalRW int
	config.DB.QueryRow("SELECT COUNT(DISTINCT rw) FROM warga WHERE rw IS NOT NULL AND rw != ''").Scan(&totalRW)

	data := map[string]interface{}{
		"Title":      "Data Warga",
		"WargaList":  wargaList,
		"TotalWarga": len(wargaList),
		"TotalKK":    totalKK,
		"TotalRW":    totalRW,
	}
	templates.ExecuteTemplate(w, "warga.html", data)
}

func wargaDetailHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	var w models.Warga
	row := config.DB.QueryRow("SELECT id, nik, nama_lengkap, alamat, no_hp, no_kk, rw, created_at FROM warga WHERE id = ?", id)
	err := row.Scan(&w.NIK, &w.NamaLengkap, &w.Alamat, &w.NoHP, &w.NoKK, &w.RW, &w.CreatedAt)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	data := map[string]interface{}{
		"Title": "Detail Warga",
		"Warga": w,
	}
	templates.ExecuteTemplate(w, "warga_detail.html", data)
}
