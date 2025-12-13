# Raylib Game

This is documentation for our team student project in University of Gdańsk.

### Running game

Dependencies:
Of course installing Golang >=1.20 is required to *run* and *build* our game. Later we are planning to add compiled binaries as release.
Note that our game is running fine in Windows 10/11 and in Linux distros running X11. If you are using Wayland you might have to install some libs for running X11.

Steps:

1. Run server:
```bash
# in windows or in Linux running X11
go run cmd/server/main.go

# in Linux running Wayland
go run -tags "x11" cmd/server/main.go
```

2. In another terminal(s) run client(s):
```bash
# in windows or in Linux running X11
go run cmd/client/main.go

# in Linux running Wayland
go run -tags "x11" cmd/client/main.go
```

3. You are good to go!

Note, that if you have `make` installed you can just
```bash
# in one terminal
make server

# in another terminal
make client
```

### Out goal (MVP)(in Polish)
- Skrypty build / run + README (uruchomienie klienta i serwera)
- Prosty HUD - HP, stan gry, ping
- Lobby sieciowe: lista graczy + przycisk Start
- Ładowanie świata i assetów
- Podstawowy level/mapa z punktami spawnu
- Klient: sterowanie postacią (WASD + mysz)
- Ruch + kolizje z terenem i obiektami
- Modele 3D zrobione samodzielnie
- Podstawowe shadery i oświetlenie
- Podstawowe animacje postaci i synchronizacja zdarzeń animacji
- Podstawowe dźwięki
- Prosty protokół sieciowy (serializacja, entity IDs, spawn/despawn)
- Placeholder assets + fallback na brakujące zasoby
- Minimalnie jeden przeciwnik
- Podstawowy system hit detection
- Prosty scenariusz końca gry

### Ideas
- Hitscan na sprawdzanie z czym moze interaktować grać
- Wskazywanie w którą strone ma wypychać ściana/podłogo
- Mapa(słownik) pokoji z obiektami które są w środku dla bardziej wydajnej kolizji
- Enum kierunków
- Code cleanup (w main.go)
