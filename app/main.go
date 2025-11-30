package main

func main() {
	InitDatabase()
	CreateViews()
	CreateTriggers()
	CreateIndexes()
	SeedDatabase()
	MainMenu()
}
