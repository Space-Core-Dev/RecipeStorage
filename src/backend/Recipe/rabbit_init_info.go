package main

const queuePref = "queue"
const keyPref = "key"

var names map[string][]string = map[string][]string{
	"recipe": {
		"recipe.create",
		"recipe.remove",
		"recipe.update",
		"recipe.get",
	},
}
