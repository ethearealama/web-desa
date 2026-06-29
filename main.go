package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

<<<<<<< HEAD
	"go-desa/config"

=======
>>>>>>> 4a7c617d8491322f91cb76691b70156f0196a3a6
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var db *sql.DB
var templates *template.Template

// ==================== STRUKTUR DATA ====================
type Berita struct {
	ID        int
	Judul     string
	Isi       string
	Gambar    string
	Penulis   string
	Views     int
	CreatedAt time.Time
}

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

type AdminPengaduan struct {
	ID        int
	Nama      string
	NIK       string
	NoHP      string
	Judul     string
	Isi       string
	Status    string
	Balasan   string
	CreatedAt string
}

// ==================== MAIN ====================
func main() {
<<<<<<< HEAD
	// Koneksi database dari config
	config.ConnectDB()
	defer config.DB.Close()
=======
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
>>>>>>> 4a7c617d8491322f91cb76691b70156f0196a3a6

	// Template functions
	funcMap := template.FuncMap{
		"add":        func(a, b int) int { return a + b },
		"formatDate": func(t time.Time) string { return t.Format("02 Jan 2006") },
	}

	// Load semua template
	templates = template.Must(template.New("").Funcs(funcMap).ParseGlob("views/*.html"))
	template.Must(templates.ParseGlob("views/layouts/*.html"))
	template.Must(templates.ParseGlob("views/admin/dashboard/*.html"))
	template.Must(templates.ParseGlob("views/admin/wisata/*.html"))
<<<<<<< HEAD
	template.Must(templates.ParseGlob("views/admin/pengaduan/*.html"))
=======
	template.Must(templates.ParseGlob("views/admin/datawarga/*.html"))
	template.Must(templates.ParseGlob("views/admin/bansos/*.html"))
>>>>>>> 4a7c617d8491322f91cb76691b70156f0196a3a6

	// Router
	r := mux.NewRouter()

	// Static files
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// ==================== ROUTES FRONTEND ====================
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

<<<<<<< HEAD
	// ==================== ROUTES ADMIN LOGIN ====================
=======
	// Routes Admin Wisata
	r.HandleFunc("/admin/wisata", adminWisataIndexHandler).Methods("GET")
	r.HandleFunc("/admin/wisata/create", adminWisataCreateHandler).Methods("GET")
	r.HandleFunc("/admin/wisata/edit/{id}", adminWisataEditHandler).Methods("GET")

	// Routes Admin Login & Dashboard
>>>>>>> 4a7c617d8491322f91cb76691b70156f0196a3a6
	r.HandleFunc("/login", loginHandler).Methods("GET")
	r.HandleFunc("/login", doLoginHandler).Methods("POST")
	r.HandleFunc("/admin/dashboard", adminDashboardHandler).Methods("GET")
	r.HandleFunc("/logout", logoutHandler).Methods("GET")

	// ==================== ROUTES ADMIN WISATA ====================
	r.HandleFunc("/admin/wisata", adminWisataHandler).Methods("GET")
	r.HandleFunc("/admin/wisata/create", adminWisataCreateHandler).Methods("GET")
	r.HandleFunc("/admin/wisata/store", adminWisataStoreHandler).Methods("POST")
	r.HandleFunc("/admin/wisata/edit/{id}", adminWisataEditHandler).Methods("GET")
	r.HandleFunc("/admin/wisata/update/{id}", adminWisataUpdateHandler).Methods("POST")
	r.HandleFunc("/admin/wisata/delete/{id}", adminWisataDeleteHandler).Methods("POST")

	// ==================== ROUTES ADMIN PENGADUAN ====================
	r.HandleFunc("/admin/pengaduan", adminPengaduanHandler).Methods("GET")
	r.HandleFunc("/admin/pengaduan/show/{id}", adminPengaduanShowHandler).Methods("GET")
	r.HandleFunc("/admin/pengaduan/update/{id}", adminPengaduanUpdateHandler).Methods("POST")
	r.HandleFunc("/admin/pengaduan/delete/{id}", adminPengaduanDeleteHandler).Methods("POST")

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

<<<<<<< HEAD
// ==================== WISATA (FRONTEND) ====================
=======
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

>>>>>>> 4a7c617d8491322f91cb76691b70156f0196a3a6
func wisataHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, nama_wisata, deskripsi, lokasi, gambar, harga_tiket, jam_operasional, status FROM wisatas")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var wisataList []Wisata
	for rows.Next() {
		var wisataItem Wisata
		rows.Scan(&wisataItem.ID, &wisataItem.NamaWisata, &wisataItem.Deskripsi, &wisataItem.Lokasi, &wisataItem.Gambar, &wisataItem.HargaTiket, &wisataItem.JamOperasional, &wisataItem.Status)
		wisataList = append(wisataList, wisataItem)
	}

	data := map[string]interface{}{
		"Title":      "Wisata Desa",
		"WisataList": wisataList,
	}
	templates.ExecuteTemplate(w, "wisata.html", data)
}

<<<<<<< HEAD
// ==================== PENGADUAN (FRONTEND) ====================
=======
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
>>>>>>> 4a7c617d8491322f91cb76691b70156f0196a3a6
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

// ==================== WARGA (FRONTEND) ====================
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
<<<<<<< HEAD
	html := `
    <!DOCTYPE html>
    <html>
    <head>
        <title>Dashboard Admin</title>
        <script src="https://cdn.tailwindcss.com"></script>
    </head>
    <body class="bg-gray-100">
        <nav class="bg-green-700 text-white shadow-lg">
            <div class="container mx-auto px-4 py-3 flex justify-between">
                <span class="font-bold">Admin Panel - Desa Sukaindah</span>
                <a href="/logout" class="bg-red-600 px-4 py-1 rounded">Logout</a>
            </div>
        </nav>
        <div class="container mx-auto px-4 py-8">
            <div class="bg-white p-6 rounded shadow">
                <h1 class="text-2xl font-bold text-green-800">Selamat Datang, Admin!</h1>
                <p>Ini dashboard admin Desa Sukaindah.</p>
                <div class="grid md:grid-cols-4 gap-4 mt-6">
                    <a href="/admin/wisata" class="bg-blue-100 p-4 rounded text-center hover:bg-blue-200">🏞️ Wisata</a>
                    <a href="/admin/pengaduan" class="bg-red-100 p-4 rounded text-center hover:bg-red-200">📢 Pengaduan</a>
                    <a href="/admin/berita" class="bg-green-100 p-4 rounded text-center hover:bg-green-200">📰 Berita</a>
                    <a href="/admin/bansos" class="bg-yellow-100 p-4 rounded text-center hover:bg-yellow-200">🎁 Bansos</a>
                </div>
            </div>
        </div>
    </body>
    </html>
    `
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
=======
	data := map[string]interface{}{
		"Title": "Dashboard Admin",
	}
	err := templates.ExecuteTemplate(w, "dashboard.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
>>>>>>> 4a7c617d8491322f91cb76691b70156f0196a3a6
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// ==================== ADMIN WISATA ====================
func adminWisataHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := config.DB.Query("SELECT id, nama_wisata, lokasi, harga_tiket, status FROM wisatas ORDER BY id DESC")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var wisataList []Wisata
	for rows.Next() {
		var wisataItem Wisata
		rows.Scan(&wisataItem.ID, &wisataItem.NamaWisata, &wisataItem.Lokasi, &wisataItem.HargaTiket, &wisataItem.Status)
		wisataList = append(wisataList, wisataItem)
	}

	data := map[string]interface{}{
		"Title":      "Kelola Wisata",
		"WisataList": wisataList,
	}
	templates.ExecuteTemplate(w, "admin/wisata/index.html", data)
}

func adminWisataCreateHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{"Title": "Tambah Wisata"}
	templates.ExecuteTemplate(w, "admin/wisata/create.html", data)
}

func adminWisataStoreHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	nama := r.FormValue("nama_wisata")
	deskripsi := r.FormValue("deskripsi")
	lokasi := r.FormValue("lokasi")
	harga := r.FormValue("harga_tiket")
	jam := r.FormValue("jam_operasional")
	status := r.FormValue("status")

	_, err := config.DB.Exec(`INSERT INTO wisatas (nama_wisata, deskripsi, lokasi, harga_tiket, jam_operasional, status) 
        VALUES (?, ?, ?, ?, ?, ?)`,
		nama, deskripsi, lokasi, harga, jam, status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/admin/wisata", http.StatusSeeOther)
}

func adminWisataEditHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	var wisataItem Wisata
	err := config.DB.QueryRow("SELECT id, nama_wisata, deskripsi, lokasi, harga_tiket, jam_operasional, status, gambar FROM wisatas WHERE id = ?", id).
		Scan(&wisataItem.ID, &wisataItem.NamaWisata, &wisataItem.Deskripsi, &wisataItem.Lokasi, &wisataItem.HargaTiket, &wisataItem.JamOperasional, &wisataItem.Status, &wisataItem.Gambar)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	data := map[string]interface{}{
		"Title":  "Edit Wisata",
		"Wisata": wisataItem,
	}
	templates.ExecuteTemplate(w, "admin/wisata/edit.html", data)
}

func adminWisataUpdateHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	r.ParseForm()
	nama := r.FormValue("nama_wisata")
	deskripsi := r.FormValue("deskripsi")
	lokasi := r.FormValue("lokasi")
	harga := r.FormValue("harga_tiket")
	jam := r.FormValue("jam_operasional")
	status := r.FormValue("status")

	_, err := config.DB.Exec(`UPDATE wisatas SET 
        nama_wisata = ?, deskripsi = ?, lokasi = ?, harga_tiket = ?, jam_operasional = ?, status = ? 
        WHERE id = ?`,
		nama, deskripsi, lokasi, harga, jam, status, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/admin/wisata", http.StatusSeeOther)
}

func adminWisataDeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	_, err := config.DB.Exec("DELETE FROM wisatas WHERE id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/admin/wisata", http.StatusSeeOther)
}

// ==================== ADMIN PENGADUAN ====================
func adminPengaduanHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := config.DB.Query("SELECT id, nama, nik, no_hp, judul, isi, status, balasan, created_at FROM pengaduan ORDER BY id DESC")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var list []AdminPengaduan
	var pending, processed, done int
	for rows.Next() {
		var p AdminPengaduan
		rows.Scan(&p.ID, &p.Nama, &p.NIK, &p.NoHP, &p.Judul, &p.Isi, &p.Status, &p.Balasan, &p.CreatedAt)
		list = append(list, p)
		switch p.Status {
		case "baru", "dibaca":
			pending++
		case "ditindaklanjuti":
			processed++
		case "selesai":
			done++
		}
	}

	data := map[string]interface{}{
		"Title":          "Kelola Pengaduan",
		"PengaduanList":  list,
		"PendingCount":   pending,
		"ProcessedCount": processed,
		"DoneCount":      done,
	}
	templates.ExecuteTemplate(w, "admin/pengaduan/index.html", data)
}

func adminPengaduanShowHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	var p AdminPengaduan
	err := config.DB.QueryRow("SELECT id, nama, nik, no_hp, judul, isi, status, balasan, created_at FROM pengaduan WHERE id = ?", id).
		Scan(&p.ID, &p.Nama, &p.NIK, &p.NoHP, &p.Judul, &p.Isi, &p.Status, &p.Balasan, &p.CreatedAt)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	data := map[string]interface{}{
		"Title":     "Detail Pengaduan",
		"Pengaduan": p,
	}
	templates.ExecuteTemplate(w, "admin/pengaduan/show.html", data)
}

func adminPengaduanUpdateHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	r.ParseForm()
	status := r.FormValue("status")
	balasan := r.FormValue("balasan")

	_, err := config.DB.Exec("UPDATE pengaduan SET status = ?, balasan = ? WHERE id = ?", status, balasan, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/admin/pengaduan", http.StatusSeeOther)
}

func adminPengaduanDeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	_, err := config.DB.Exec("DELETE FROM pengaduan WHERE id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/admin/pengaduan", http.StatusSeeOther)
}