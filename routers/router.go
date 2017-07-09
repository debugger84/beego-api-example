// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"tournamentAPI/controllers"

	"github.com/astaxie/beego"
)

func init() {

	beego.Router("/balance", &controllers.BalanceController{}, "get:Balance")
	beego.Router("/take", &controllers.BalanceController{}, "get:Charge")
	beego.Router("/fund", &controllers.BalanceController{}, "get:Fund")

	beego.Router("/announceTournament", &controllers.TournamentController{}, "get:AnnounceTournament")
	beego.Router("/joinTournament", &controllers.TournamentController{}, "get:JoinTournament")
	beego.Router("/resultTournament", &controllers.TournamentController{}, "post:ResultTournament")

	beego.Router("/reset", &controllers.SystemController{}, "get:Reset")
}
