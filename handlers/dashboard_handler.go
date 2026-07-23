package handlers

import (
	"net/http"

	viewsconst "github.com/hilmanxcode/web-inventory-go/const"
	"github.com/hilmanxcode/web-inventory-go/sessions"
	"github.com/hilmanxcode/web-inventory-go/utils/viewsutil"
)

func DashboardPage(w http.ResponseWriter, r *http.Request) {

	var data = map[string]any{
		"currentPage": viewsconst.VIEWS_DASHBOARD,
		"csrf_token":  sessions.GetCSRFToken(w, r),
	}

	viewsutil.ShowView(viewsconst.VIEWS_DASHBOARD, data, w)

}

func MasterBarang(w http.ResponseWriter, r *http.Request) {

	csrfToken := sessions.GetCSRFToken(w, r)

	var data = map[string]any{
		"csrf_token":  csrfToken,
		"currentPage": viewsconst.VIEWS_MASTER_BARANG,
	}

	viewsutil.ShowView(viewsconst.VIEWS_MASTER_BARANG, data, w)
}

func BarangMasuk(w http.ResponseWriter, r *http.Request) {
	var data = map[string]any{
		"currentPage": viewsconst.VIEWS_BARANG_MASUK,
	}

	viewsutil.ShowView(viewsconst.VIEWS_BARANG_MASUK, data, w)
}

func BarangKeluar(w http.ResponseWriter, r *http.Request) {
	var data = map[string]any{
		"currentPage": viewsconst.VIEWS_BARANG_KELUAR,
	}

	viewsutil.ShowView(viewsconst.VIEWS_BARANG_KELUAR, data, w)
}

func LaporanStok(w http.ResponseWriter, r *http.Request) {
	var data = map[string]any{
		"currentPage": viewsconst.VIEWS_LAPORAN_STOK,
	}

	viewsutil.ShowView(viewsconst.VIEWS_LAPORAN_STOK, data, w)
}
