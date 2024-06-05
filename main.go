package main

import (
	"fmt"
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
	MAX_ASTEROIDS                   = 12
	ASTEROID_SPEED                  = 1
	ASTEROID_SIZE                   = 50.0
	ASTEROID_POINTS                 = 11
	SHIP_TIME_IN_PIECES             = 5
)

type GameState struct {
	playerShip    *PlayerShip
	asteroids     *[]Asteroid
	debug         bool
	collision     bool
	lives         int
	gameTime      float64
	destroyedTime float64
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
			points = append(points, (rand.Float32()*0.6)+0.6)

		}

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

func generateMidAsteroid(pos rl.Vector2, size float32) Asteroid {
	points := []float32{}
	for range ASTEROID_POINTS {
		points = append(points, (rand.Float32()*0.6)+0.6)
	}
	posX := pos.X
	posY := pos.Y
	orientation := rand.Float32() * (math.Pi * 2)
	directionX := float32(math.Cos(float64(orientation)))
	directionY := float32(math.Sin(float64(orientation)))
	speed := rand.Float32() * ASTEROID_SPEED
	asteroid := Asteroid{
		pos:         rl.NewVector2(posX, posY),
		speed:       speed,
		vel:         rl.Vector2Scale(rl.NewVector2(directionX, directionY), speed),
		size:        ASTEROID_SIZE / size,
		orientation: orientation,
		sizes:       points,
	}
	return asteroid
}

func drawGameOverScreen() {
	rl.DrawTextPro(rl.GetFontDefault(), "Game Over", rl.Vector2{
		X: SCREEN_SIZE_X / 2,
		Y: SCREEN_SIZE_Y / 2,
	}, rl.Vector2Scale(rl.MeasureTextEx(rl.GetFontDefault(), "Game Over", 100.0, 1.0), 0.5), 0.0, 100.0, 1.0, rl.White)
	rl.DrawTextPro(rl.GetFontDefault(), "Press Enter to try again", rl.Vector2{
		X: SCREEN_SIZE_X / 2,
		Y: SCREEN_SIZE_Y/2 + 100,
	}, rl.Vector2Scale(rl.MeasureTextEx(rl.GetFontDefault(), "Press Enter to try again", 50.0, 1.0), 0.5), 0.0, 50.0, 1.0, rl.White)

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

func (s *PlayerShip) drawShipExplosion() {

	verticalDirection := rl.Vector2Add(rl.Vector2Scale(getDirection(s.orientation), s.size), s.pos)
	horizontalDirection := rl.Vector2Add(rl.Vector2Scale(getDirection(s.orientation+math.Pi*0.7), s.size), s.pos)
	opVerticalDirection := rl.Vector2Add(rl.Vector2Scale(getDirection(s.orientation+math.Pi*0.25), s.size), s.pos)
	opHorizontalDirection := rl.Vector2Add(rl.Vector2Scale(getDirection(math.Pi+s.orientation+math.Pi*0.1), s.size), s.pos)

	rl.DrawLineV(rl.Vector2AddValue(s.pos, 25), verticalDirection, rl.White)
	rl.DrawLineV(rl.Vector2AddValue(s.pos, 50), horizontalDirection, rl.White)
	rl.DrawLineV(rl.Vector2AddValue(s.pos, -40), opVerticalDirection, rl.White)
	rl.DrawLineV(rl.Vector2AddValue(s.pos, 35), opHorizontalDirection, rl.White)
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

func drawLife(pos rl.Vector2, size float32, orientation float32) {
	verticalDirection := rl.Vector2Scale(getDirection(orientation), size)
	horizontalDirection := rl.Vector2Scale(getDirection(orientation+math.Pi*0.5), size)

	points := []rl.Vector2{
		rl.Vector2Add(pos, verticalDirection),
		rl.Vector2Subtract(rl.Vector2Subtract(pos, verticalDirection), horizontalDirection),
		pos,
		rl.Vector2Add(rl.Vector2Subtract(pos, verticalDirection), horizontalDirection),
		rl.Vector2Add(pos, verticalDirection),
	}

	for i := range points {
		rl.DrawLineV(
			points[i],
			points[(i+1)%len(points)],
			rl.White,
		)
	}

}

func (s *PlayerShip) getShipPoints() []rl.Vector2 {
	verticalDirection := rl.Vector2Scale(getDirection(s.orientation), s.size)
	horizontalDirection := rl.Vector2Scale(getDirection(s.orientation+math.Pi*0.5), s.size)

	points := []rl.Vector2{
		rl.Vector2Add(s.pos, verticalDirection),
		rl.Vector2Subtract(rl.Vector2Subtract(s.pos, verticalDirection), horizontalDirection),
		s.pos,
		rl.Vector2Add(rl.Vector2Subtract(s.pos, verticalDirection), horizontalDirection),
		rl.Vector2Add(s.pos, verticalDirection),
	}
	return points
}

func (g *GameState) getAsteroidsPoints() [][]rl.Vector2 {
	asteroidsPoints := [][]rl.Vector2{}
	for _, a := range *g.asteroids {
		asteroidPoints := []rl.Vector2{
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
		}
		asteroidsPoints = append(asteroidsPoints, asteroidPoints)
	}
	return asteroidsPoints
}

func (g *GameState) checkColissions() bool {
	shipPoints := g.playerShip.getShipPoints()
	asteroidsPoints := g.getAsteroidsPoints()

	for i := range shipPoints {
		for j := range *g.asteroids {
			collisionPoint := rl.Vector2{}
			if rl.CheckCollisionLines(shipPoints[i], shipPoints[(i+1)%len(shipPoints)], asteroidsPoints[j][0], asteroidsPoints[j][1], &collisionPoint) ||
				rl.CheckCollisionLines(shipPoints[i], shipPoints[(i+1)%len(shipPoints)], asteroidsPoints[j][1], asteroidsPoints[j][2], &collisionPoint) ||
				rl.CheckCollisionLines(shipPoints[i], shipPoints[(i+1)%len(shipPoints)], asteroidsPoints[j][2], asteroidsPoints[j][3], &collisionPoint) ||
				rl.CheckCollisionLines(shipPoints[i], shipPoints[(i+1)%len(shipPoints)], asteroidsPoints[j][3], asteroidsPoints[j][4], &collisionPoint) ||
				rl.CheckCollisionLines(shipPoints[i], shipPoints[(i+1)%len(shipPoints)], asteroidsPoints[j][4], asteroidsPoints[j][5], &collisionPoint) ||
				rl.CheckCollisionLines(shipPoints[i], shipPoints[(i+1)%len(shipPoints)], asteroidsPoints[j][5], asteroidsPoints[j][6], &collisionPoint) ||
				rl.CheckCollisionLines(shipPoints[i], shipPoints[(i+1)%len(shipPoints)], asteroidsPoints[j][6], asteroidsPoints[j][7], &collisionPoint) ||
				rl.CheckCollisionLines(shipPoints[i], shipPoints[(i+1)%len(shipPoints)], asteroidsPoints[j][7], asteroidsPoints[j][8], &collisionPoint) ||
				rl.CheckCollisionLines(shipPoints[i], shipPoints[(i+1)%len(shipPoints)], asteroidsPoints[j][8], asteroidsPoints[j][9], &collisionPoint) ||
				rl.CheckCollisionLines(shipPoints[i], shipPoints[(i+1)%len(shipPoints)], asteroidsPoints[j][9], asteroidsPoints[j][10], &collisionPoint) {
				return true
			}
		}
	}
	return false
}

func (g *GameState) checkProjectileCollisions() {

	for i, p := range *g.playerShip.projectiles {
		for j, a := range *g.asteroids {
			asteroidsPoints := g.getAsteroidsPoints()
			if j < len(*g.asteroids) {
				if rl.CheckCollisionPointLine(p.pos, asteroidsPoints[j][0], asteroidsPoints[j][1], 20) ||
					rl.CheckCollisionPointLine(p.pos, asteroidsPoints[j][1], asteroidsPoints[j][2], 20) ||
					rl.CheckCollisionPointLine(p.pos, asteroidsPoints[j][2], asteroidsPoints[j][3], 20) ||
					rl.CheckCollisionPointLine(p.pos, asteroidsPoints[j][3], asteroidsPoints[j][4], 20) ||
					rl.CheckCollisionPointLine(p.pos, asteroidsPoints[j][4], asteroidsPoints[j][5], 20) ||
					rl.CheckCollisionPointLine(p.pos, asteroidsPoints[j][5], asteroidsPoints[j][6], 20) ||
					rl.CheckCollisionPointLine(p.pos, asteroidsPoints[j][6], asteroidsPoints[j][7], 20) ||
					rl.CheckCollisionPointLine(p.pos, asteroidsPoints[j][7], asteroidsPoints[j][8], 20) ||
					rl.CheckCollisionPointLine(p.pos, asteroidsPoints[j][8], asteroidsPoints[j][9], 20) ||
					rl.CheckCollisionPointLine(p.pos, asteroidsPoints[j][9], asteroidsPoints[j][10], 20) {
					removeItem(g.asteroids, j)
					removeItem(g.playerShip.projectiles, i)
					if a.size == ASTEROID_SIZE {
						*g.asteroids = append(*g.asteroids, generateMidAsteroid(a.pos, 2.0))
						*g.asteroids = append(*g.asteroids, generateMidAsteroid(a.pos, 2.0))
					}
					if a.size == ASTEROID_SIZE/2 {
						*g.asteroids = append(*g.asteroids, generateMidAsteroid(a.pos, 4.0))
						*g.asteroids = append(*g.asteroids, generateMidAsteroid(a.pos, 4.0))
					}

				}

			}

		}

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
	g.gameTime = rl.GetTime()
	if !g.collision && g.lives > 0 {
		g.input()
	}
	if g.lives > 0 {
		g.playerShip.pos = *resetPosition(&g.playerShip.pos)
		g.playerShip.moveProjectiles()
		g.playerShip.removeProjectiles()
		g.checkProjectileCollisions()

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

		if g.collision {
			g.destroyedTime -= 0.1
			if g.destroyedTime < 0 {
				g.restartGame()
			}
		} else {
			if g.checkColissions() {
				g.collision = true
			}
		}
	}
	g.moveAsteroids()

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
		if g.collision {
			rl.DrawTextEx(rl.GetFontDefault(), fmt.Sprintf("Collision: %v", g.collision), rl.Vector2{
				X: 210,
				Y: 10,
			}, 10.0, 1.0, rl.Red)
		} else {
			rl.DrawTextEx(rl.GetFontDefault(), fmt.Sprintf("Collision: %v", g.collision), rl.Vector2{
				X: 210,
				Y: 10,
			}, 10.0, 1.0, rl.White)
		}
		rl.DrawTextEx(rl.GetFontDefault(), fmt.Sprintf("Gametime: %f", g.gameTime), rl.Vector2{
			X: 300,
			Y: 10,
		}, 10.0, 1.0, rl.White)

	}
	if g.collision {
		g.playerShip.drawShipExplosion()

	} else if g.lives > 0 {
		g.playerShip.drawShip()
	}

	g.playerShip.drawProjectiles()
	g.drawAsteroids()
	for i := range g.lives {
		drawLife(rl.Vector2{
			X: 25 + 45*float32(i),
			Y: SCREEN_SIZE_Y - 25,
		}, 20, math.Pi+math.Pi*0.5)
	}
	if g.lives <= 0 {
		drawGameOverScreen()
	}
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
		asteroids:     generateAsteroids(),
		debug:         true,
		collision:     false,
		lives:         3,
		gameTime:      0,
		destroyedTime: SHIP_TIME_IN_PIECES,
	}
	return &gState
}

func (g *GameState) restartGame() {
	g.playerShip.pos = rl.Vector2{
		X: SCREEN_SIZE_X / 2,
		Y: SCREEN_SIZE_Y / 2,
	}
	g.playerShip.orientation = -math.Pi / 2
	g.playerShip.vel = rl.Vector2{
		X: 0,
		Y: 0,
	}
	g.collision = false
	g.lives = g.lives - 1
	g.destroyedTime = SHIP_TIME_IN_PIECES

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
