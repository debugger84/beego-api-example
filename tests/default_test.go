package test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"runtime"
	"path/filepath"
	_ "tournamentAPI/routers"

	"github.com/astaxie/beego"
	. "github.com/smartystreets/goconvey/convey"
	"strings"
	"fmt"
	"tournamentAPI/models"
	"encoding/json"
)

func init() {
	_, file, _, _ := runtime.Caller(1)
	apppath, _ := filepath.Abs(filepath.Dir(filepath.Join(file, ".." + string(filepath.Separator))))
	beego.TestBeegoInit(apppath)
}

// TestGet is a sample to run an endpoint test
func TestTournamentPlaying(t *testing.T) {
	var r *http.Request
	w := httptest.NewRecorder()

	r, _ = http.NewRequest("GET", "/reset", nil)
	beego.BeeApp.Handlers.ServeHTTP(w, r)
	r, _ = http.NewRequest("GET", "/fund?playerId=P1&points=300", nil)
	beego.BeeApp.Handlers.ServeHTTP(w, r)
	r, _ = http.NewRequest("GET", "/fund?playerId=P2&points=300", nil)
	beego.BeeApp.Handlers.ServeHTTP(w, r)
	r, _ = http.NewRequest("GET", "/fund?playerId=P3&points=300", nil)
	beego.BeeApp.Handlers.ServeHTTP(w, r)
	r, _ = http.NewRequest("GET", "/fund?playerId=P4&points=500", nil)
	beego.BeeApp.Handlers.ServeHTTP(w, r)
	r, _ = http.NewRequest("GET", "/fund?playerId=P5&points=1000", nil)
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	r, _ = http.NewRequest("GET", "/announceTournament?tournamentId=1&deposit=1000", nil)
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	r, _ = http.NewRequest("GET", "/joinTournament?tournamentId=1&playerId=P5", nil)
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	r, _ = http.NewRequest("GET", "/joinTournament?tournamentId=1&playerId=P1&backerId=P2&backerId=P3&backerId=P4", nil)
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	reader := strings.NewReader(`{"tournamentId": "1", "winners": [{"playerId": "P1", "prize": 2000}]}`)
	r, _ = http.NewRequest("POST", "/resultTournament", reader)
	beego.BeeApp.Handlers.ServeHTTP(w, r)


	//r, _ := http.NewRequest("GET", "/balance?playerId=P1", nil)

	//beego.BeeApp.Handlers.ServeHTTP(w, r)

	beego.Trace("testing", "Playing of tournament", fmt.Sprintf("Code[%d]\n%s", w.Code, w.Body.String()))

	Convey("Subject: Test Station Endpoint\n", t, func() {
	        Convey("Status Code Should Be 200", func() {
	                So(w.Code, ShouldEqual, 200)
	        })
	        Convey("The Result Should Not Be Empty", func() {
	                So(w.Body.Len(), ShouldBeGreaterThan, 0)
	        })
	})

	testBalance(t, "P1", 550)
	testBalance(t, "P2", 550)
	testBalance(t, "P3", 550)
	testBalance(t, "P4", 750)
	testBalance(t, "P5", 0)

}

func testBalance(t *testing.T, playerId string, expectedBalance int) {
	w1 := httptest.NewRecorder()
	r1, _ := http.NewRequest("GET", "/balance?playerId=" + playerId, nil)
	beego.BeeApp.Handlers.ServeHTTP(w1, r1)
	beego.Trace("testing", "Getting balance of P1", fmt.Sprintf("Code[%d]\n%s", w1.Code, w1.Body.String()))

	Convey("Subject: Test Station Endpoint\n", t, func() {
		Convey("Status Code Should Be 200", func() {
			So(w1.Code, ShouldEqual, 200)
		})
		Convey(fmt.Sprintf("The Balance Should Be Equal %d", expectedBalance), func() {
			var player models.Player
			json.Unmarshal(w1.Body.Bytes(), &player)
			So(player.Balance, ShouldEqual, expectedBalance)
		})
	})
}
