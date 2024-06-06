## Asteroids-Go
The old school videogame Asteroids made in Go without any assets. All done with raylib's drawing functions and math.

### Prerequisites

* Go v1.22.2
* [Raylib 5.0](https://www.raylib.com/)

I recommend installing Raylib's bindings for the Go programming language using the [pure go installation](https://github.com/gen2brain/raylib-go/tree/master?tab=readme-ov-file#purego-without-cgo-ie-cgo_enabled0). It's extremely easy to set up and start working.

### Installation
#### On Windows
1. Clone the repo 
```sh
   git clone https://github.com/Rodolfato/asteroids-go
```
2. Run the batch script
```sh
   .\build_and_run.bat
```
3. Play the game!

### Controls
* `A` and `D` to rotate the ship
* `W` to move forward, `S` to move backwards
* `Space` to shot projectiles

### Editing the game

You can change any of the global variables at the start of the `main.go` file to change the games starting settings. After any changes that you've made run the `build_and_run.bat` to test the game.

If you happen to run in to this repository feel free to download, check it out and use it as a learning tool for raylib or Go!

```go
const (
	PLAYER_SHIP_SIZE                = 20
	PLAYER_SHIP_THICKNESS           = 1.5
	PLAYER_SHIP_INITIAL_ORIENTATION = math.Pi + (math.Pi * 0.5)
	PLAYER_SHIP_TURN_SPEED          = 0.02 * math.Pi
	PLAYER_SHIP_SPEED               = 0.3
	SCREEN_SIZE_X                   = 1024
	SCREEN_SIZE_Y                   = 768
	PROJECTILE_SPEED                = PLAYER_SHIP_SPEED + 15
	TTL_PRJECTILE                   = 45
	PROJECTILE_SIZE                 = 2.5
	MAX_SPEED                       = 5
	MAX_ASTEROIDS                   = 12
	ASTEROID_SPEED                  = 1
	ASTEROID_SIZE                   = 50.0
	ASTEROID_POINTS                 = 11
	SHIP_TIME_IN_PIECES             = 5
	LIVES                           = 3
)
```

### To do
- [ ] Organize the code for better understanding
- [ ] Add alien ship enemies
- [ ] Add Linux installation instructions
- [ ] Add a configuration file for easy customization
