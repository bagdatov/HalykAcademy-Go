package airport

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"
)

type Plane struct {
	title string
	// взлет
	// посадка
	ctx context.Context
	// context with timeout для имитации "времени полета"
	// состояние летит, взлетает, садиться, обслуживается, припаркован
	status string
}

const (
	starting    = "plane is taking off the ground"
	flying      = "plane is on the way"
	landing     = "plane is landing"
	maintenance = "plane is fueling"
	parked      = "plane is parked at the airport"
)

func (p *Plane) start(ctx context.Context, a *AirPort) {
	p.status = starting
	log.Printf("plane number: %s is taking off\n", p.title)

	out, err := a.action(p)
	if err != nil {
		log.Println(err)
		return
	}
	<-out

	p.status = flying
	log.Printf("plane number: %s is in the sky\n", p.title)

	ctxTimeout, cancel := context.WithTimeout(ctx, randomTime(5, 15)*time.Second)

	go func() {
		for {
			if err := a.stctx.Err(); err != nil {
				cancel()
			}
		}
	}()

	<-ctxTimeout.Done()
	a.landingCh <- p

}

func (p *Plane) land(ctx context.Context, a *AirPort) {
	p.status = landing
	log.Printf("plane number: %s is landing\n", p.title)
	in, err := a.action(p)
	if err != nil {
		log.Println(err)
		return
	}
	<-in
	log.Printf("plane number: %s has landed\n", p.title)
}

// NewPlane is returning *Plane with random title and default status `parked`.
func NewPlane() *Plane {
	return &Plane{
		title:  randomName(),
		status: parked,
	}
}

// randomName is returning string formatted as `Boeing XXX`,
// where XXX is randomly generated number in between 700 and 799
func randomName() string {
	rand.Seed(time.Now().UnixNano())
	min, max := 700, 799
	num := rand.Intn(max-min) + min
	return fmt.Sprintf("Boeing %v", num)
}

// randomTime is returning randomly generated time in between min and max arguments.
func randomTime(min, max int) time.Duration {
	rand.Seed(time.Now().UnixNano())
	return time.Duration(rand.Intn(max-min) + min)
}
