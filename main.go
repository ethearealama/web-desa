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
	dsn := "root:@tcp(localhost:3306)/dbdesa?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Gagal konek database:", err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal("Database tidak merespon:", err)
	}

	// Template functions
	funcMap := template.FuncMap{
		"add":        func(a, b int) int { return a + b },
		"formatDate": func(t time.Time) string { return t.Format("02 Jan 2006") },
	}

	templates = template.Must(template.New("").Funcs(funcMap).ParseGlob("views/*.html"))
	template.Must(templates.ParseGlob("views/layouts/*.html"))
	template.Must(templates.ParseGlob("views/admin/dashboard/*.html"))
	template.Must(templates.ParseGlob("views/admin/wisata/*.html"))
	template.Must(templates.ParseGlob("views/admin/datawarga/*.html"))
	template.Must(templates.ParseGlob("views/admin/bansos/*.html"))

	// Router
	r := mux.NewRouter()

	// Static files
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Routes Frontend
	r.HandleFunc("/", homeHandler).Methods("GET")
	r.HandleFunc("/admin", adminDashboardHandler).Methods("GET")
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

	// Routes Admin Wisata
	r.HandleFunc("/admin/wisata", adminWisataIndexHandler).Methods("GET")
	r.HandleFunc("/admin/wisata/create", adminWisataCreateHandler).Methods("GET")
	r.HandleFunc("/admin/wisata/edit/{id}", adminWisataEditHandler).Methods("GET")

	// Routes Admin Login & Dashboard
	r.HandleFunc("/login", loginHandler).Methods("GET")
	r.HandleFunc("/login", doLoginHandler).Methods("POST")
	r.HandleFunc("/admin/dashboard", adminDashboardHandler).Methods("GET")
	r.HandleFunc("/logout", logoutHandler).Methods("GET")

	fmt.Println("🌐 Server running on http://localhost:9090")
	log.Fatal(http.ListenAndServe(":9090", r))
}

// ==================== HOME ====================
func homeHandler(w http.ResponseWriter, r *http.Request) {
	var totalPenduduk, totalKK, totalRW, totalWisata int

	db.QueryRow("SELECT COUNT(*) FROM warga").Scan(&totalPenduduk)
	db.QueryRow("SELECT COUNT(DISTINCT no_kk) FROM warga WHERE no_kk IS NOT NULL AND no_kk != ''").Scan(&totalKK)
	db.QueryRow("SELECT COUNT(DISTINCT rw) FROM warga WHERE rw IS NOT NULL AND rw != ''").Scan(&totalRW)
	db.QueryRow("SELECT COUNT(*) FROM wisatas").Scan(&totalWisata)

	data := map[string]interface{}{
		"Title":         "Beranda",
		"TotalPenduduk": totalPenduduk,
		"TotalKK":       totalKK,
		"TotalRW":       totalRW,
		"TotalWisata":   totalWisata,
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

// ==================== ADMIN WISATA ====================
func adminWisataIndexHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "admin/wisata/index.html", nil)
}

func adminWisataCreateHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "admin/wisata/create.html", nil)
}

func adminWisataEditHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "admin/wisata/edit.html", nil)
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
		var wg Warga
		rows.Scan(&wg.ID, &wg.NIK, &wg.NamaLengkap, &wg.Alamat, &wg.NoHP, &wg.NoKK, &wg.RW, &wg.CreatedAt, &wg.UpdatedAt)
		wargaList = append(wargaList, wg)
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
	id, _ := strconv.Atoi(vars["id"])

	var wg Warga
	query := "SELECT id, nik, nama_lengkap, alamat, no_hp, no_kk, rw, created_at, updated_at FROM warga WHERE id = ?"
	err := db.QueryRow(query, id).Scan(
		&wg.ID,
		&wg.NIK,
		&wg.NamaLengkap,
		&wg.Alamat,
		&wg.NoHP,
		&wg.NoKK,
		&wg.RW,
		&wg.CreatedAt,
		&wg.UpdatedAt,
	)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	data := map[string]interface{}{
		"Title": "Detail Warga",
		"Warga": wg,
	}
	templates.ExecuteTemplate(w, "warga_detail.html", data)
}

// ==================== LOGIN ADMIN ====================
func loginHandler(w http.ResponseWriter, r *http.Request) {
	errorMsg := r.URL.Query().Get("error")
	html := `
    <!DOCTYPE html>
    <html>
    <head>
        <title>Login Admin</title>
        <style>
            body { font-family: Arial; display: flex; justify-content: center; align-items: center; height: 100vh; background: #f0f2f5; }
            .box { background: white; padding: 30px; border-radius: 10px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); width: 350px; text-align: center; }
            input { width: 100%; padding: 10px; margin: 5px 0 20px; border: 1px solid #ccc; border-radius: 5px; }
            button { background: #2e7d32; color: white; padding: 10px; border: none; border-radius: 5px; width: 100%; cursor: pointer; }
            button:hover { background: #1b5e20; }
            .error { color: red; margin-bottom: 10px; }
            h2 { color: #2e7d32; }
        </style>
    </head>
    <body>
        <div class="box">
            <h2>Login Admin Desa</h2>
            <form method="POST" action="/login">
                <input type="email" name="email" placeholder="Email" required>
                <input type="password" name="password" placeholder="Password" required>
                <button type="submit">Login</button>
            </form>
        </div>
    </body>
    </html>
    `
	if errorMsg != "" {
		html = `<div style="color:red; text-align:center; margin-bottom:10px;">` + errorMsg + `</div>` + html
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
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
	data := map[string]interface{}{
		"Title": "Dashboard Admin",
	}
	err := templates.ExecuteTemplate(w, "dashboard.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
