package checker

import (
	"sniping/helpers"
)

func Run() {
	users := helpers.LoadUsernames()
	Discord(users)
}