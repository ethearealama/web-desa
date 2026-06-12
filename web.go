package main

import (
	"fmt"
	"html/template"
	"net/http"
)

var templates *template.Template

func main() {
	// Load semua template HTML
	templates = template.Must(template.ParseGlob("views/*.html"))
	template.Must(templates.ParseGlob("views/layouts/*.html"))

	// Static files (CSS, images)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Routes
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/profil", profilHandler)
	http.HandleFunc("/berita", beritaHandler)
	http.HandleFunc("/wisata", wisataHandler)
	http.HandleFunc("/pengaduan", pengaduanHandler)
	http.HandleFunc("/bansos", bansosHandler)
	http.HandleFunc("/surat", suratHandler)

	fmt.Println("Server running on http://localhost:8082")
	http.ListenAndServe(":8082", nil)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title":         "Beranda",
		"TotalPenduduk": 5240,
		"TotalKK":       1450,
		"TotalRW":       12,
		"TotalWisata":   8,
	}
	templates.ExecuteTemplate(w, "index.html", data)
}

func profilHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title": "Profil Desa",
	}
	templates.ExecuteTemplate(w, "profil.html", data)
}

func beritaHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title": "Berita Desa",
	}
	templates.ExecuteTemplate(w, "berita.html", data)
}

func wisataHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title": "Wisata Desa",
	}
	templates.ExecuteTemplate(w, "wisata.html", data)
}

func pengaduanHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title": "Pengaduan",
	}
	templates.ExecuteTemplate(w, "pengaduan.html", data)
}

func bansosHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title": "Cek Bansos",
	}
	templates.ExecuteTemplate(w, "bansos.html", data)
}

func suratHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title": "Layanan Surat",
	}
	templates.ExecuteTemplate(w, "surat.html", data)
}
