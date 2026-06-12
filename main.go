package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var db *sql.DB
var templates *template.Template

func main() {
	// Koneksi database
	var err error
	dsn := "root:@tcp(127.0.0.1:3306)/dbdesa?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Gagal konek database:", err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal("Database tidak merespon:", err)
	}
	fmt.Println("✅ Koneksi database berhasil!")

	// Template functions
	funcMap := template.FuncMap{
		"add":        func(a, b int) int { return a + b },
		"formatDate": func(t time.Time) string { return t.Format("02 Jan 2006") },
	}

	// Load semua template
	templates = template.Must(template.New("").Funcs(funcMap).ParseGlob("views/*.html"))
	template.Must(templates.ParseGlob("views/layouts/*.html"))
	template.Must(templates.ParseGlob("views/admin/*.html"))

	// Router
	r := mux.NewRouter()

	// Static files
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Routes Frontend
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
	r.HandleFunc("/surat/tracking", suratTrackingHandler).Methods("GET")
	r.HandleFunc("/warga", wargaHandler).Methods("GET")
	r.HandleFunc("/warga/{id}", wargaDetailHandler).Methods("GET")

	// Routes Admin Login
	r.HandleFunc("/login", loginHandler).Methods("GET")
	r.HandleFunc("/login", doLoginHandler).Methods("POST")
	r.HandleFunc("/admin/dashboard", adminDashboardHandler).Methods("GET")
	r.HandleFunc("/logout", logoutHandler).Methods("GET")

	fmt.Println("🌐 Server running on http://localhost:9090")
	log.Fatal(http.ListenAndServe(":9090", r))
}

// ==================== HOME ====================
func homeHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title": "Beranda",
	}
	templates.ExecuteTemplate(w, "index.html", data)
}

func profilHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{"Title": "Profil Desa"}
	templates.ExecuteTemplate(w, "profil.html", data)
}

// ==================== BERITA ====================
type Berita struct {
	ID        int
	Judul     string
	Isi       string
	Gambar    string
	Penulis   string
	Views     int
	CreatedAt time.Time
}

func beritaHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, judul, isi, gambar, penulis, views, created_at FROM beritas ORDER BY created_at DESC")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var beritaList []Berita
	for rows.Next() {
		var b Berita
		rows.Scan(&b.ID, &b.Judul, &b.Isi, &b.Gambar, &b.Penulis, &b.Views, &b.CreatedAt)
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

	var b Berita
	err := db.QueryRow("SELECT id, judul, isi, gambar, penulis, views, created_at FROM beritas WHERE id = ?", id).
		Scan(&b.ID, &b.Judul, &b.Isi, &b.Gambar, &b.Penulis, &b.Views, &b.CreatedAt)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	db.Exec("UPDATE beritas SET views = views + 1 WHERE id = ?", id)

	data := map[string]interface{}{
		"Title":   b.Judul,
		"Berita":  b,
		"Isi":     b.Isi,
		"Gambar":  b.Gambar,
		"Penulis": b.Penulis,
		"Views":   b.Views + 1,
		"Tanggal": b.CreatedAt,
	}
	templates.ExecuteTemplate(w, "berita_detail.html", data)
}

// ==================== WISATA ====================
type Wisata struct {
	ID             int
	NamaWisata     string
	Deskripsi      string
	Lokasi         string
	Gambar         string
	HargaTiket     string
	JamOperasional string
	Status         string
}

func wisataHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, nama_wisata, deskripsi, lokasi, gambar, harga_tiket, jam_operasional, status FROM wisatas")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var wisataList []Wisata
	for rows.Next() {
		var w Wisata
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
	data := map[string]interface{}{
		"Title":   "Pengaduan",
		"success": r.URL.Query().Get("success") == "true",
	}
	templates.ExecuteTemplate(w, "pengaduan.html", data)
}

func kirimPengaduanHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	_, err := db.Exec("INSERT INTO pengaduan (nama, nik, no_hp, judul, isi, status) VALUES (?, ?, ?, ?, ?, 'baru')",
		r.FormValue("nama"), r.FormValue("nik"), r.FormValue("no_hp"), r.FormValue("judul"), r.FormValue("isi"))
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

	var nama, jenis, alamat, status string
	err := db.QueryRow("SELECT nama_penerima, jenis_bantuan, alamat, status FROM bansos WHERE nik = ?", nik).
		Scan(&nama, &jenis, &alamat, &status)

	data := map[string]interface{}{
		"Title":     "Cek Bansos",
		"Nik":       nik,
		"Nama":      nama,
		"Jenis":     jenis,
		"Alamat":    alamat,
		"Status":    status,
		"Ditemukan": err == nil,
	}
	templates.ExecuteTemplate(w, "bansos.html", data)
}

// ==================== SURAT ====================
func suratHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title":   "Layanan Surat",
		"success": r.URL.Query().Get("success") == "true",
	}
	templates.ExecuteTemplate(w, "surat.html", data)
}

func kirimSuratHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	no_pengajuan := fmt.Sprintf("SURAT/%s/%d", time.Now().Format("20060102"), time.Now().UnixNano()%1000)

	_, err := db.Exec(`INSERT INTO pengajuan_surat (no_pengajuan, nama, nik, alamat, no_hp, jenis_surat, keperluan, status) 
        VALUES (?, ?, ?, ?, ?, ?, ?, 'pending')`,
		no_pengajuan, r.FormValue("nama"), r.FormValue("nik"), r.FormValue("alamat"), r.FormValue("no_hp"), r.FormValue("jenis_surat"), r.FormValue("keperluan"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/surat?success=true", http.StatusSeeOther)
}

func suratTrackingHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{"Title": "Tracking Surat"}
	templates.ExecuteTemplate(w, "surat_tracking.html", data)
}

// ==================== WARGA ====================
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

func wargaHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, nik, nama_lengkap, alamat, no_hp, no_kk, rw, created_at, updated_at FROM warga ORDER BY id DESC")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var wargaList []Warga
	for rows.Next() {
		var w Warga
		rows.Scan(&w.ID, &w.NIK, &w.NamaLengkap, &w.Alamat, &w.NoHP, &w.NoKK, &w.RW, &w.CreatedAt, &w.UpdatedAt)
		wargaList = append(wargaList, w)
	}

	var totalKK, totalRW int
	db.QueryRow("SELECT COUNT(DISTINCT no_kk) FROM warga WHERE no_kk IS NOT NULL AND no_kk != ''").Scan(&totalKK)
	db.QueryRow("SELECT COUNT(DISTINCT rw) FROM warga WHERE rw IS NOT NULL AND rw != ''").Scan(&totalRW)

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
    id := vars["id"]

    var nama string
    err := db.QueryRow("SELECT nama_lengkap FROM warga WHERE id = ?", id).Scan(&nama)
    if err != nil {
        http.NotFound(w, r)
        return
    }

    w.Write([]byte("Nama Warga: " + nama))
}
// ==================== LOGIN ADMIN ====================
func loginHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{"error": r.URL.Query().Get("error")}
	templates.ExecuteTemplate(w, "admin/login.html", data)
}

func doLoginHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	email := r.FormValue("email")
	password := r.FormValue("password")

	if email == "admin@desa.com" && password == "password" {
		http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/login?error=Email+atau+password+salah", http.StatusSeeOther)
}

func adminDashboardHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{"Title": "Dashboard Admin"}
	templates.ExecuteTemplate(w, "admin/dashboard.html", data)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
