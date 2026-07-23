package httputil

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"

	"github.com/hilmanxcode/web-inventory-go/sessions"
)

func RedirectWithError(w http.ResponseWriter, r *http.Request, errorMessage string, path string) {
	c, errCookie := r.Cookie("session_token")

	if errCookie != nil {
		sessions.SetSession(sessions.Session{
			ErrorMessages: []string{errorMessage},
		}, w)
		http.Redirect(w, r, path, http.StatusFound)
		return
	}

	val, err := sessions.GetSession(c.Value)
	if err != nil {
		sessions.SetSession(sessions.Session{
			ErrorMessages: []string{"Invalid Session", errorMessage},
		}, w)
		http.Redirect(w, r, path, http.StatusFound)
		return
	}

	val.ErrorMessages = []string{errorMessage}
	sessions.UpdateSession(c.Value, val)
	http.Redirect(w, r, path, http.StatusFound)
}

func UploadBase64Handler(w http.ResponseWriter, r *http.Request, imageForm string) string {
	// 1. Batasi ukuran memori (misal 10MB)
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		RedirectWithError(w, r, "Ukuran file terlalu besar", "/dashboard")
		return ""
	}

	// 2. Ambil file dari form
	file, _, err := r.FormFile(imageForm)
	if err != nil {
		RedirectWithError(w, r, "Gambar Wajib di isi", "/dashboard")
		return ""
	}
	defer file.Close()

	// 3. Baca seluruh isi file ke dalam memori (slice of byte)
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		RedirectWithError(w, r, "Gagal membaca isi file", "/dashboard")
		return ""
	}

	// 4. Deteksi tipe MIME asli gambar (misal: "image/jpeg" atau "image/png")
	// Ini penting agar string Base64 nantinya bisa langsung dirender di tag <img src="...">
	mimeType := http.DetectContentType(fileBytes)

	// 5. Konversi byte menjadi string Base64
	base64Encoded := base64.StdEncoding.EncodeToString(fileBytes)

	// 6. Gabungkan tipe MIME dengan string Base64 (Format Data URI standar)
	base64Image := fmt.Sprintf("data:%s;base64,%s", mimeType, base64Encoded)

	// SELESAI! Variabel `base64Image` sekarang berisi string seperti:
	// data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAA...
	// Anda bisa menyimpan string ini ke kolom text di database Anda.

	// Untuk simulasi, kita render langsung string tersebut ke dalam tag HTML
	return base64Image
}
