module github.com/PawelZabc/ProjektZespolowy/client

go 1.25.4

require (
	github.com/chewxy/math32 v1.11.1
	github.com/gen2brain/raylib-go/raylib v0.55.1
)

require github.com/PawelZabc/ProjektZespolowy/shared v0.0.0

replace github.com/PawelZabc/ProjektZespolowy/shared => ../shared

require (
	github.com/ebitengine/purego v0.7.1 // indirect
	golang.org/x/exp v0.0.0-20240506185415-9bf2ced13842 // indirect
	golang.org/x/sys v0.20.0 // indirect
)
