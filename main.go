package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os" // Tambahan: Diperlukan untuk membaca port dari Render
	"strconv"
)

// Struktur data untuk mengirim hasil hitungan ke halaman struk/hasil
type HasilBelanja struct {
	Menu1Qty   int
	Menu2Qty   int
	Menu3Qty   int
	Menu4Qty   int
	Menu5Qty   int
	TotalHarga int
}

func main() {
	// 1. Menampilkan halaman utama kasir (index.html) saat pertama kali dibuka
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("templates/index.html")
		if err != nil {
			http.Error(w, "File index.html tidak ditemukan di folder templates", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
	})

	// 2. Handler untuk memproses tombol "Hitung Total Belanja" (POST /hitung)
	http.HandleFunc("/hitung", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		// Ambil data jumlah porsi dari input HTML (name="menu1", dst)
		qty1, _ := strconv.Atoi(r.FormValue("menu1"))
		qty2, _ := strconv.Atoi(r.FormValue("menu2"))
		qty3, _ := strconv.Atoi(r.FormValue("menu3"))
		qty4, _ := strconv.Atoi(r.FormValue("menu4"))
		qty5, _ := strconv.Atoi(r.FormValue("menu5"))

		// Definisikan harga masing-masing menu sesuai di HTML
		hargaBatagorPremium := 10000
		hargaSiomayPremium := 10000
		hargaBasoGoreng := 6000
		hargaBatagorMedium := 6000
		hargaSiomayMedium := 6000

		// Rumus Hitung Total Belanja
		total := (qty1 * hargaBatagorPremium) +
			(qty2 * hargaSiomayPremium) +
			(qty3 * hargaBasoGoreng) +
			(qty4 * hargaBatagorMedium) +
			(qty5 * hargaSiomayMedium)

		// Bungkus data hasil hitungan
		dataHasil := HasilBelanja{
			Menu1Qty:   qty1,
			Menu2Qty:   qty2,
			Menu3Qty:   qty3,
			Menu4Qty:   qty4,
			Menu5Qty:   qty5,
			TotalHarga: total,
		}

		// Tampilkan halaman struk sederhana langsung dari Go
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, `
			<!DOCTYPE html>
			<html>
			<head>
				<title>Struk Pembayaran</title>
				<style>
					body { font-family: 'Segoe UI', sans-serif; background-color: #f4f4f4; display: flex; justify-content: center; align-items: center; height: 100vh; margin: 0; }
					.struk { background: white; padding: 30px; border-radius: 10px; box-shadow: 0 4px 10px rgba(0,0,0,0.1); max-width: 400px; width: 100%%; }
					h3 { text-align: center; color: #b71c1c; margin-top: 0; }
					.garis { border-top: 2px dashed #333; margin: 15px 0; }
					.item { display: flex; justify-content: space-between; margin-bottom: 8px; font-size: 14px; }
					.total { font-size: 18px; font-weight: bold; color: #1b5e20; display: flex; justify-content: space-between; margin-top: 15px; }
					.tombol-kembali { display: block; text-align: center; background: #b71c1c; color: white; text-decoration: none; padding: 10px; margin-top: 20px; border-radius: 5px; font-weight: bold; }
				</style>
			</head>
			<body>
				<div class="struk">
					<h3>STRUK BATAGOR SAHABAT</h3>
					<div class="garis"></div>
					%s
					<div class="garis"></div>
					<div class="total">
						<span>TOTAL :</span>
						<span>Rp %s</span>
					</div>
					<a href="/" class="tombol-kembali">Kembali ke Kasir</a>
				</div>
			</body>
			</html>
		`, formatRincian(dataHasil), formatRupiah(total))
	})

	// --- PERUBAHAN UNTUK DEPLOY KE RENDER ---

	// 1. Ambil port dari sistem Render secara otomatis
	port := os.Getenv("PORT")

	// 2. Jika port kosong (artinya kamu sedang jalankan di localhost laptop), pakai default 8080
	if port == "" {
		port = "7860"
	}

	// 3. Jalankan server menggunakan variabel port dinamis
	fmt.Println("Server kasir berjalan aman di port:", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic(err)
	}
}

// Fungsi pembantu untuk menyusun teks rincian menu yang dibeli
func formatRincian(d HasilBelanja) string {
	rincian := ""
	if d.Menu1Qty > 0 {
		rincian += fmt.Sprintf("<div class='item'><span>Batagor Premium x%d</span><span>Rp %d</span></div>", d.Menu1Qty, d.Menu1Qty*10000)
	}
	if d.Menu2Qty > 0 {
		rincian += fmt.Sprintf("<div class='item'><span>Siomay Premium x%d</span><span>Rp %d</span></div>", d.Menu2Qty, d.Menu2Qty*10000)
	}
	if d.Menu3Qty > 0 {
		rincian += fmt.Sprintf("<div class='item'><span>Baso Goreng x%d</span><span>Rp %d</span></div>", d.Menu3Qty, d.Menu3Qty*6000)
	}
	if d.Menu4Qty > 0 {
		rincian += fmt.Sprintf("<div class='item'><span>Batagor Medium x%d</span><span>Rp %d</span></div>", d.Menu4Qty, d.Menu4Qty*6000)
	}
	if d.Menu5Qty > 0 {
		rincian += fmt.Sprintf("<div class='item'><span>Siomay Medium x%d</span><span>Rp %d</span></div>", d.Menu5Qty, d.Menu5Qty*6000)
	}

	if rincian == "" {
		rincian = "<div style='text-align:center; color:#999;'>Tidak ada item yang dibeli</div>"
	}
	return rincian
}

// Fungsi pembantu untuk memformat angka biasa menjadi ribuan
func formatRupiah(angka int) string {
	str := strconv.Itoa(angka)
	panjang := len(str)
	if panjang <= 3 {
		return str
	}
	hasil := ""
	hitung := 0
	for i := panjang - 1; i >= 0; i-- {
		hasil = string(str[i]) + hasil
		hitung++
		if hitung%3 == 0 && i != 0 {
			hasil = "." + hasil
		}
	}
	return hasil
}
