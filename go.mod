module main

go 1.22.8

require (
	github.com/gorilla/mux v1.8.1
	gopkg.in/yaml.v3 v3.0.1
	vanhalt.com/authservice v0.0.0-00010101000000-000000000000
)

replace vanhalt.com/authservice => ./authservice
