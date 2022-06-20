package main

import (
	"math"
	"math/rand"
	"time"
)

type Boid struct {
	position Vector2D
	velocity Vector2D

	id int
}

func (b *Boid) calcAcceleration() Vector2D {

	//b点的 右上角, 左下角  范围就确定了
	upper, lowwer := b.position.AddV(viewRadius), b.position.AddV(-viewRadius)
	// lowwer, upper := b.position.AddV(viewRadius), b.position.AddV(-viewRadius)
	avgPosition, avgVelocity := Vector2D{0, 0}, Vector2D{0, 0}
	count := 0.0

	lock.Lock()
	for i := math.Max(lowwer.x, 0); i <= math.Min(upper.x, screenWidth); i++ {
		for j := math.Max(0, lowwer.y); j <= math.Min(upper.y, screenHeight); j++ {
			if otherBoidId := boidMap[int(i)][int(j)]; otherBoidId != -1 && otherBoidId != b.id {
				if dist := boids[otherBoidId].position.Distance(b.position); dist < viewRadius {
					count++
					avgVelocity = avgVelocity.Add(boids[otherBoidId].velocity)
					avgPosition = avgPosition.Add(boids[otherBoidId].position)
				}
			}
		}
	}
	lock.Unlock()

	accel := Vector2D{0, 0}
	if count > 0 {
		avgPosition, avgVelocity = avgPosition.DivisionV(count), avgVelocity.DivisionV(count)
		// accel = avgVelocity.Subtract(b.velocity).MultiplyV(adjRate)
		accelAligment := avgVelocity.Subtract(b.velocity).MultiplyV(adjRate)
		accelCohesion := avgPosition.Subtract(b.position).MultiplyV(adjRate)
		accel = accel.Add(accelAligment).Add(accelCohesion)
	}

	return accel

}

func (b *Boid) moveOne() {
	acceleration := b.calcAcceleration()
	lock.Lock()
	b.velocity = b.velocity.Add(acceleration).Limit(-1, 1)
	boidMap[int(b.position.x)][int(b.position.y)] = -1
	b.position = b.position.Add(b.velocity)
	boidMap[int(b.position.x)][int(b.position.y)] = b.id
	next := b.position.Add(b.velocity)
	if next.x > screenWidth || next.x <= 0 {
		b.velocity = Vector2D{-b.velocity.x, b.velocity.y}
	}
	if next.y > screenHeight || next.y <= 0 {
		b.velocity = Vector2D{b.velocity.x, -b.velocity.y}
	}
	lock.Unlock()
}

func (b *Boid) start() {

	for {
		b.moveOne()
		time.Sleep(5 * time.Millisecond)
	}

}

func createBoid(bid int) {

	b := Boid{
		position: Vector2D{rand.Float64() * screenWidth, rand.Float64() * screenHeight},
		// position: Vector2D{0, 0},
		velocity: Vector2D{(rand.Float64() * 2) - 1.0, (rand.Float64() * 2) - 1.0},
		id:       bid,
	}

	boids[bid] = &b

	boidMap[int(b.position.x)][int(b.position.y)] = b.id

	go b.start()
}
