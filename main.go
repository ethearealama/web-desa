package main

import (
    "database/sql"
    "fmt"
    "html/template"
    "log"
    "net/http"
    "strconv"
    "time"

    "github.com/gorilla/mux"
    _ "github.com/go-sql-driver/mysql"
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

    // Template dengan fungsi tambahan
    funcMap := template.FuncMap{
        "add":        func(a, b int) int { return a + b },
        "formatDate": func(t time.Time) string { return t.Format("02 Jan 2006") },
    }

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

    // Routes Admin
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
func beritaHandler(w http.ResponseWriter, r *http.Request) {
    rows, err := db.Query("SELECT id, judul, isi, gambar, penulis, views, created_at FROM beritas ORDER BY created_at DESC")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    type Berita struct {
        ID        int
        Judul     string
        Isi       string
        Gambar    string
        Penulis   string
        Views     int
        CreatedAt time.Time
    }

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

    var judul, isi, gambar, penulis string
    var views int
    var createdAt time.Time

    err := db.QueryRow("SELECT judul, isi, gambar, penulis, views, created_at FROM beritas WHERE id = ?", id).Scan(&judul, &isi, &gambar, &penulis, &views, &createdAt)
    if err != nil {
        http.NotFound(w, r)
        return
    }

    db.Exec("UPDATE beritas SET views = views + 1 WHERE id = ?", id)

    data := map[string]interface{}{
        "Title":   judul,
        "Judul":   judul,
        "Isi":     isi,
        "Gambar":  gambar,
        "Penulis": penulis,
        "Views":   views + 1,
        "Tanggal": createdAt,
    }
    templates.ExecuteTemplate(w, "berita_detail.html", data)
}

// ==================== WISATA ====================
func wisataHandler(w http.ResponseWriter, r *http.Request) {
    data := map[string]interface{}{"Title": "Wisata Desa"}
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

    _, err := db.Exec("INSERT INTO pengaduan (nama, nik, no_hp, judul, isi, status) VALUES (?, ?, ?, ?, ?, 'baru')",
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

    var nama, jenis, alamat, status string
    err := db.QueryRow("SELECT nama_penerima, jenis_bantuan, alamat, status FROM bansos WHERE nik = ?", nik).Scan(&nama, &jenis, &alamat, &status)

    data := map[string]interface{}{
        "Title": "Cek Bansos",
        "Nik":   nik,
        "Nama":  nama,
        "Jenis": jenis,
        "Alamat": alamat,
        "Status": status,
        "Ditemukan": err == nil,
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
    no_pengajuan := fmt.Sprintf("SURAT/%s/%d", time.Now().Format("20060102"), time.Now().UnixNano()%1000)
    _, err := db.Exec("INSERT INTO pengajuan_surat (no_pengajuan, nama, nik, alamat, no_hp, jenis_surat, keperluan, status) VALUES (?, ?, ?, ?, ?, ?, ?, 'pending')",
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
func wargaHandler(w http.ResponseWriter, r *http.Request) {
    data := map[string]interface{}{"Title": "Data Warga"}
    templates.ExecuteTemplate(w, "warga.html", data)
}

func wargaDetailHandler(w http.ResponseWriter, r *http.Request) {
    data := map[string]interface{}{"Title": "Detail Warga"}
    templates.ExecuteTemplate(w, "warga_detail.html", data)
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