package main

import (
	"fmt"
	"log"
	"math"
	"math/rand/v2"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	PLAYER_SHIP_SIZE                = 20
	PLAYER_SHIP_THICKNESS           = 1.5
	PLAYER_SHIP_INITIAL_ORIENTATION = 0.0
	PLAYER_SHIP_TURN_SPEED          = 0.02 * math.Pi
	PLAYER_SHIP_SPEED               = 0.3
	SCREEN_SIZE_X                   = 1024
	SCREEN_SIZE_Y                   = 768
	PROJECTILE_SPEED                = PLAYER_SHIP_SPEED + 15
	TTL_PRJECTILE                   = 45
	PROJECTILE_SIZE                 = 2.5
	MAX_SPEED                       = 5
	MAX_ASTEROIDS                   = 20
	ASTEROID_SPEED                  = 1
	ASTEROID_SIZE                   = 50.0
	ASTEROID_POINTS                 = 12
)

type GameState struct {
	playerShip *PlayerShip
	asteroids  *[]Asteroid
	debug      bool
}

type PlayerShip struct {
	pos         rl.Vector2
	orientation float32
	size        float32
	speed       float32
	vel         rl.Vector2
	projectiles *[]Projectile
}

type Projectile struct {
	pos         rl.Vector2
	speed       float32
	vel         rl.Vector2
	ttl         int
	orientation float32
	size        float32
}

type Asteroid struct {
	pos         rl.Vector2
	speed       float32
	vel         rl.Vector2
	orientation float32
	size        float32
	sizes       []float32
}

func generateAsteroids() *[]Asteroid {

	asteroids := []Asteroid{}
	positions := make(map[rl.Vector2]bool)

	for range MAX_ASTEROIDS {
		//slices.Sort(angles)
		points := []float32{}
		for range ASTEROID_POINTS {
			points = append(points, (rand.Float32()*0.7)+0.5)

		}
		/*for i := range ASTEROID_POINTS {
			log.Printf("asteroid size: %f", points[i])
		}
		log.Println("\n")*/
		cdX := rand.Float32() * SCREEN_SIZE_X
		cdY := rand.Float32() * SCREEN_SIZE_Y
		orientation := rand.Float32() * (math.Pi * 2)
		directionX := float32(math.Cos(float64(orientation)))
		directionY := float32(math.Sin(float64(orientation)))
		speed := rand.Float32() * ASTEROID_SPEED
		_, ok := positions[rl.NewVector2(cdX, cdY)]
		for ok {
			cdX := rand.Float32() * SCREEN_SIZE_X
			cdY := rand.Float32() * SCREEN_SIZE_Y
			_, ok = positions[rl.NewVector2(cdX, cdY)]
		}
		positions[rl.NewVector2(cdX, cdY)] = true
		asteroid := Asteroid{
			pos:         rl.NewVector2(cdX, cdY),
			speed:       speed,
			vel:         rl.Vector2Scale(rl.NewVector2(directionX, directionY), speed),
			size:        ASTEROID_SIZE,
			orientation: orientation,
			sizes:       points,
		}
		asteroids = append(asteroids, asteroid)
	}
	return &asteroids
}

func (s *PlayerShip) shoot() {
	circleX := s.pos.X + (s.size+10)*float32(math.Cos(float64(s.orientation)))
	circleY := s.pos.Y + (s.size+10)*float32(math.Sin(float64(s.orientation)))
	initialPosVector := rl.NewVector2(circleX, circleY)

	//Sentido y orientacion de la nave
	directionX := float32(math.Cos(float64(s.orientation)))
	directionY := float32(math.Sin(float64(s.orientation)))
	//Este vector es el sentido y orientacion de la nave
	newVector := rl.NewVector2(directionX, directionY)
	//Con el sentido y orientación de la nave se puede escalar con la rapidez para obtener la velocidad
	projectileVelocity := rl.Vector2Add(
		s.vel,
		rl.Vector2Scale(newVector, PROJECTILE_SPEED+PLAYER_SHIP_SPEED),
	)

	projectile := Projectile{
		pos:         initialPosVector,
		speed:       PROJECTILE_SPEED,
		vel:         projectileVelocity,
		ttl:         TTL_PRJECTILE,
		orientation: s.orientation,
		size:        PROJECTILE_SIZE,
	}
	*s.projectiles = append(*s.projectiles, projectile)

}

func (p *Projectile) drawProjectile() {
	rl.DrawCircleV(p.pos, PROJECTILE_SIZE, rl.White)
}

func (s *PlayerShip) drawProjectiles() {
	for _, p := range *s.projectiles {
		p.drawProjectile()
	}
}

func (s *PlayerShip) moveProjectiles() {
	for i := range *s.projectiles {
		(*s.projectiles)[i].pos = rl.Vector2Add((*s.projectiles)[i].pos, (*s.projectiles)[i].vel)
		resetPosition(&(*s.projectiles)[i].pos)
		(*s.projectiles)[i].ttl -= 1
	}
}

func (s *PlayerShip) removeProjectiles() {
	for i, p := range *s.projectiles {
		if p.ttl < 1 {
			removeItem(s.projectiles, i)
		}
	}
}

func removeItem[T any](slice *[]T, index int) {
	if index < 0 || index >= len(*slice) {
		return
	}
	copy((*slice)[index:], (*slice)[index+1:])
	*slice = (*slice)[:len(*slice)-1]
}

func getDirection(orientation float32) rl.Vector2 {
	circleX := float32(math.Cos(float64(orientation)))
	circleY := float32(math.Sin(float64(orientation)))

	newVector := rl.NewVector2(circleX, circleY)
	return newVector

}

func (s *PlayerShip) drawShip() {
	verticalDirection := rl.Vector2Scale(getDirection(s.orientation), s.size)
	horizontalDirection := rl.Vector2Scale(getDirection(s.orientation+math.Pi*0.5), s.size)

	points := []rl.Vector2{
		rl.Vector2Add(s.pos, verticalDirection),
		rl.Vector2Subtract(rl.Vector2Subtract(s.pos, verticalDirection), horizontalDirection),
		s.pos,
		rl.Vector2Add(rl.Vector2Subtract(s.pos, verticalDirection), horizontalDirection),
		rl.Vector2Add(s.pos, verticalDirection),
	}

	for i := range points {
		rl.DrawLineV(
			points[i],
			points[(i+1)%len(points)],
			rl.White,
		)
	}
}

func (g *GameState) input() {
	if rl.IsKeyPressed(rl.KeyF1) {
		g.debug = !g.debug
	}
	if rl.IsKeyDown(rl.KeyD) {
		newOrientation := g.playerShip.orientation + PLAYER_SHIP_TURN_SPEED
		if newOrientation >= 2*math.Pi {
			g.playerShip.orientation = 0.0
		} else if newOrientation <= -2*math.Pi {
			g.playerShip.orientation = 0.0
		} else {
			g.playerShip.orientation = newOrientation
		}

		log.Printf("D key pressed")
		log.Printf("Ship orientation: %f", g.playerShip.orientation)

	}
	if rl.IsKeyDown(rl.KeyA) {
		newOrientation := g.playerShip.orientation - PLAYER_SHIP_TURN_SPEED
		if newOrientation >= 2*math.Pi {
			g.playerShip.orientation = 0.0
		} else if newOrientation <= -2*math.Pi {
			g.playerShip.orientation = 0.0
		} else {
			g.playerShip.orientation = newOrientation
		}
		log.Printf("A key pressed")
		log.Printf("Ship orientation: %f", g.playerShip.orientation)
	}

	//Sentido y orientacion de la nave
	directionX := float32(math.Cos(float64(g.playerShip.orientation)))
	directionY := float32(math.Sin(float64(g.playerShip.orientation)))

	//Este vector es el sentido y orientacion de la nave
	newVector := rl.NewVector2(directionX, directionY)

	if rl.IsKeyDown(rl.KeyW) {
		//Agregarle la rapidez
		g.playerShip.vel = rl.Vector2Add(
			g.playerShip.vel,
			rl.Vector2Scale(newVector, g.playerShip.speed),
		)
	}

	if rl.IsKeyDown(rl.KeyS) {
		//Agregarle la rapidez
		g.playerShip.vel = rl.Vector2Subtract(
			g.playerShip.vel,
			rl.Vector2Scale(newVector, g.playerShip.speed),
		)
	}

	if rl.IsKeyPressed(rl.KeySpace) {
		g.playerShip.shoot()
	}

	g.playerShip.pos = rl.Vector2Add(g.playerShip.pos, g.playerShip.vel)

}

func (a *Asteroid) drawAsteroid() {

	//rl.DrawCircleV(a.pos, 2, rl.Blue)
	//rl.DrawLineV(a.pos, rl.Vector2Add(a.pos, verticalDirection), rl.Pink)
	//rl.DrawLineV(a.pos, rl.Vector2Add(a.pos, horizontalDirection), rl.Red)
	//points = append(points, a.pos)
	//rl.Vector2Add(rl.Vector2Scale(getDirection(a.anglePoints[i]), a.sizePoints[i]), a.pos)
	types := [][]rl.Vector2{
		{
			rl.Vector2Add(rl.Vector2Scale(getDirection(a.orientation), a.size*a.sizes[0]), a.pos),
			rl.Vector2Add(rl.Vector2Scale(getDirection(a.orientation+(math.Pi*2)), a.size*a.sizes[1]), a.pos),
			rl.Vector2Add(rl.Vector2Scale(getDirection(a.orientation+1.1*(math.Pi*2)), a.size*a.sizes[2]), a.pos),
			rl.Vector2Add(rl.Vector2Scale(getDirection(a.orientation+1.2*(math.Pi*2)), a.size*a.sizes[3]), a.pos),
			rl.Vector2Add(rl.Vector2Scale(getDirection(a.orientation+1.3*(math.Pi*2)), a.size*a.sizes[4]), a.pos),
			rl.Vector2Add(rl.Vector2Scale(getDirection(a.orientation+1.4*(math.Pi*2)), a.size*a.sizes[5]), a.pos),
			rl.Vector2Add(rl.Vector2Scale(getDirection(a.orientation+1.5*(math.Pi*2)), a.size*a.sizes[6]), a.pos),
			rl.Vector2Add(rl.Vector2Scale(getDirection(a.orientation+1.6*(math.Pi*2)), a.size*a.sizes[7]), a.pos),
			rl.Vector2Add(rl.Vector2Scale(getDirection(a.orientation+1.7*(math.Pi*2)), a.size*a.sizes[8]), a.pos),
			rl.Vector2Add(rl.Vector2Scale(getDirection(a.orientation+1.8*(math.Pi*2)), a.size*a.sizes[9]), a.pos),
			rl.Vector2Add(rl.Vector2Scale(getDirection(a.orientation+1.9*(math.Pi*2)), a.size*a.sizes[10]), a.pos),
		},
	}
	//points := types[rand.IntN(len(types)-1)]
	points := types[0]

	for i := range points {
		rl.DrawLineV(
			//a.pos,
			points[i],
			points[(i+1)%len(points)],
			rl.White,
		)
	}
}

func (g *GameState) drawAsteroids() {
	for _, p := range *g.asteroids {
		p.drawAsteroid()
	}
}

func (g *GameState) moveAsteroids() {
	for i := range *g.asteroids {
		(*g.asteroids)[i].pos = rl.Vector2Add((*g.asteroids)[i].pos, (*g.asteroids)[i].vel)
		resetPosition(&(*g.asteroids)[i].pos)
	}
}

func resetPosition(position *rl.Vector2) *rl.Vector2 {

	position.X = float32(
		math.Mod(
			float64(position.X), float64(SCREEN_SIZE_X)))

	position.Y = float32(
		math.Mod(
			float64(position.Y), float64(SCREEN_SIZE_Y)))

	if position.X <= 0 {
		position.X = SCREEN_SIZE_X
	}
	if position.Y <= 0 {
		position.Y = SCREEN_SIZE_Y
	}
	return position
}

func (g *GameState) update() {
	g.input()
	g.playerShip.pos = *resetPosition(&g.playerShip.pos)
	g.playerShip.moveProjectiles()
	g.playerShip.removeProjectiles()
	g.moveAsteroids()
	g.playerShip.pos = rl.Vector2Add(g.playerShip.pos, g.playerShip.vel)
	if g.playerShip.vel.X > MAX_SPEED {
		g.playerShip.vel.X = MAX_SPEED
	}
	if g.playerShip.vel.Y > MAX_SPEED {
		g.playerShip.vel.Y = MAX_SPEED
	}
	if g.playerShip.vel.X < -MAX_SPEED {
		g.playerShip.vel.X = -MAX_SPEED
	}
	if g.playerShip.vel.Y < -MAX_SPEED {
		g.playerShip.vel.Y = -MAX_SPEED
	}
}

func (g *GameState) render() {

	rl.BeginDrawing()
	g.draw()
	rl.EndDrawing()
}

func (g *GameState) draw() {
	rl.ClearBackground(rl.Black)
	if g.debug {
		rl.DrawTextEx(rl.GetFontDefault(), fmt.Sprintf("Ship position: (%f, %f)", g.playerShip.pos.X, g.playerShip.pos.Y), rl.Vector2{
			X: 10,
			Y: 10,
		}, 10.0, 1.0, rl.White)
		rl.DrawTextEx(rl.GetFontDefault(), fmt.Sprintf("Velocity: (%f, %f)", g.playerShip.vel.X, g.playerShip.vel.Y), rl.Vector2{
			X: 10,
			Y: 30,
		}, 10.0, 1.0, rl.White)

		for i, p := range *g.playerShip.projectiles {
			rl.DrawTextEx(rl.GetFontDefault(), fmt.Sprintf("P(%f, %f)", p.pos.X, p.pos.Y), rl.Vector2{
				X: 150,
				Y: 50 + 10*float32(i),
			}, 10.0, 1.0, rl.White)
		}

		for i, a := range *g.asteroids {
			rl.DrawTextEx(rl.GetFontDefault(), fmt.Sprintf("A(%f, %f)", a.pos.X, a.pos.Y), rl.Vector2{
				X: 10,
				Y: 50 + 10*float32(i),
			}, 10.0, 1.0, rl.White)
		}
	}
	g.playerShip.drawShip()
	g.playerShip.drawProjectiles()
	g.drawAsteroids()
}

func initGame() *GameState {
	gState := GameState{
		playerShip: &PlayerShip{
			pos: rl.Vector2{
				X: SCREEN_SIZE_X / 2,
				Y: SCREEN_SIZE_Y / 2,
			},
			size:        PLAYER_SHIP_SIZE,
			orientation: PLAYER_SHIP_INITIAL_ORIENTATION,
			speed:       PLAYER_SHIP_SPEED,
			vel: rl.Vector2{
				X: 0,
				Y: 0,
			},
			projectiles: &[]Projectile{},
		},
		asteroids: generateAsteroids(),
		debug:     true,
	}
	return &gState
}

func main() {
	rl.InitWindow(SCREEN_SIZE_X, SCREEN_SIZE_Y, "Rokas espasiales")

	defer rl.CloseWindow()
	gState := initGame()
	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		gState.update()
		gState.render()
	}
}
