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
		"add":        func(a, b int) int { return a + b },
		"formatDate": func(t time.Time) string { return t.Format("02 Jan 2006") },
	}

	templates = template.Must(template.New("").Funcs(funcMap).ParseGlob("views/*.html"))
	template.Must(templates.ParseGlob("views/layouts/*.html"))
	template.Must(templates.ParseGlob("views/admin/*.html"))

	// Router
	r := mux.NewRouter()
	// Auth routes
	r.HandleFunc("/login", loginHandler).Methods("GET")
	r.HandleFunc("/login", doLoginHandler).Methods("POST")
	r.HandleFunc("/admin/dashboard", adminDashboardHandler).Methods("GET")
	r.HandleFunc("/logout", logoutHandler).Methods("GET")

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
	r.HandleFunc("/admin/wisata/create", adminWisataCreateHandler).Methods("GET")

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

	// Hitung dari tabel warga
	config.DB.QueryRow("SELECT COUNT(*) FROM warga").Scan(&stat.JumlahPenduduk)
	config.DB.QueryRow("SELECT COUNT(DISTINCT no_kk) FROM warga WHERE no_kk IS NOT NULL AND no_kk != ''").Scan(&stat.JumlahKK)
	config.DB.QueryRow("SELECT COUNT(DISTINCT rw) FROM warga WHERE rw IS NOT NULL AND rw != ''").Scan(&stat.JumlahRW)
	config.DB.QueryRow("SELECT COUNT(*) FROM wisatas").Scan(&stat.JumlahWisata)

	return stat
}

func getBeritaTerbaru() []models.Berita {
	rows, err := config.DB.Query("SELECT id, judul, slug, isi, gambar, penulis, views, created_at, updated_at FROM beritas ORDER BY created_at DESC LIMIT 3")
	if err != nil {
		return []models.Berita{}
	}
	defer rows.Close()

	var beritaList []models.Berita
	for rows.Next() {
		var b models.Berita
		rows.Scan(&b.ID, &b.Judul, &b.Slug, &b.Isi, &b.Gambar, &b.Penulis, &b.Views, &b.CreatedAt, &b.UpdatedAt)
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
	rows, err := config.DB.Query("SELECT id, judul, slug, isi, gambar, penulis, views, created_at, updated_at FROM beritas ORDER BY created_at DESC")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var beritaList []models.Berita
	for rows.Next() {
		var b models.Berita
		rows.Scan(&b.ID, &b.Judul, &b.Slug, &b.Isi, &b.Gambar, &b.Penulis, &b.Views, &b.CreatedAt, &b.UpdatedAt)
		beritaList = append(beritaList, b)
	}

	data := map[string]interface{}{
		"Title":      "Berita Desa",
		"BeritaList": beritaList,
	}
	templates.ExecuteTemplate(w, "berita.html", data)
}

// BERITA
func beritaDetailHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	var b models.Berita
	row := config.DB.QueryRow("SELECT id, judul, slug, isi, gambar, penulis, views, created_at, updated_at FROM beritas WHERE id = ?", id)
	err := row.Scan(&b.ID, &b.Judul, &b.Slug, &b.Isi, &b.Gambar, &b.Penulis, &b.Views, &b.CreatedAt, &b.UpdatedAt)
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
		var wisata models.Wisata
		rows.Scan(&wisata.ID, &wisata.NamaWisata, &wisata.Deskripsi, &wisata.Lokasi, &wisata.Gambar, &wisata.HargaTiket, &wisata.JamOperasional, &wisata.Status)
		wisataList = append(wisataList, wisata)
	}

	data := map[string]interface{}{
		"Title":      "Wisata Desa",
		"WisataList": wisataList,
	}
	templates.ExecuteTemplate(w, "wisata.html", data)
}

// ADMIN
func adminWisataCreateHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "admin/create_wisata.html", nil)
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

	var bansos models.Bansos
	row := config.DB.QueryRow("SELECT id, nik, nama_penerima, jenis_bantuan, alamat, status FROM bansos WHERE nik = ?", nik)
	err := row.Scan(&bansos.ID, &bansos.NIK, &bansos.NamaPenerima, &bansos.JenisBantuan, &bansos.Alamat, &bansos.Status)

	data := map[string]interface{}{
		"Title": "Cek Bansos",
		"Hasil": bansos,
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
	rows, err := config.DB.Query("SELECT id, nik, nama_lengkap, alamat, no_hp, no_kk, rw, created_at, updated_at FROM warga ORDER BY id DESC")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var wargaList []models.Warga
	for rows.Next() {
		var warga models.Warga
		rows.Scan(&warga.ID, &warga.NIK, &warga.NamaLengkap, &warga.Alamat, &warga.NoHP, &warga.NoKK, &warga.RW, &warga.CreatedAt, &warga.UpdatedAt)
		wargaList = append(wargaList, warga)
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

	var warga models.Warga
	row := config.DB.QueryRow("SELECT id, nik, nama_lengkap, alamat, no_hp, no_kk, rw, created_at, updated_at FROM warga WHERE id = ?", id)
	err := row.Scan(&warga.ID, &warga.NIK, &warga.NamaLengkap, &warga.Alamat, &warga.NoHP, &warga.NoKK, &warga.RW, &warga.CreatedAt, &warga.UpdatedAt)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	data := map[string]interface{}{
		"Title": "Detail Warga",
		"Warga": warga,
	}
	templates.ExecuteTemplate(w, "warga_detail.html", data)
}

// ==================== LOGIN ADMIN ====================
func loginHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"error": r.URL.Query().Get("error"),
	}
	templates.ExecuteTemplate(w, "admin/login.html", data)
}

func doLoginHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	email := r.FormValue("email")
	password := r.FormValue("password")

	if email == "admin@desa.com" && password == "password" {
		// Redirect ke dashboard admin
		http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/login?error=Email+atau+password+salah", http.StatusSeeOther)
}

func adminDashboardHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "admin/dashboard.html", nil)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
