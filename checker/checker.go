package checker

import (
	"sniping/helpers"
)

func Run() {
	users := helpers.LoadUsernames()
	net := helpers.InitConfig()
	Discord(users, net)
}