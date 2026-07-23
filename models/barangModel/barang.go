package barangmodel

type barang struct {
	Sku        string
	NamaBarang string
	Kategori   string
	Stock      string
	Image      string
}

var BarangColumn = barang{
	Sku:        "sku",
	NamaBarang: "nama_barang",
	Kategori:   "kategori",
	Stock:      "stock",
	Image:      "image",
}
