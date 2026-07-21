package handlers

import (
	"net/http"

	viewsconst "github.com/hilmanxcode/web-inventory-go/const"
	"github.com/hilmanxcode/web-inventory-go/utils/viewsutil"
)

func DashboardPage(w http.ResponseWriter, r *http.Request) {

	viewsutil.ShowView(viewsconst.VIEWS_DASHBOARD, nil, w)

}
