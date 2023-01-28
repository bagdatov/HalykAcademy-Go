package airport

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

type Plane struct {
	title  int
	status string
}

// flying функция для полета самолета, в конце она отправляет самолет на посадку
func (p *Plane) flying(a *Airport) {
	p.status = "fly"
	fmt.Printf("%s plane %v: is %s\n", getCurrentTime(), p.title, p.status)

	r := rand.Intn(10)
	if r < 3 {
		r = 3
	}

	//TODO:: логика полета.
	// Полет либо должен закончится по таймауту, либо если аэропорт скажет садить - мы закрываемся
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r)*time.Second)
	defer cancel()

	if a.isClose() {
		fmt.Printf("%s plane %v: received signal of airport closure\n", getCurrentTime(), p.title)
		cancel()
	}

	<-ctx.Done()

	fmt.Printf("%s plane %v: is finished flight and going to land\n", getCurrentTime(), p.title)

	//TODO:: самолет нужно отправить на посадку
	a.landingCh <- p
}

// servicing функция обслуживания самолета, в конце она отправляет самолет обратно на взлет
func (p *Plane) servicing(a *Airport) {

	if !a.isClose() {
		fmt.Printf("%s plane %v: is starting servicing\n", getCurrentTime(), p.title)
		p.status = "on service"

		r := rand.Intn(3)
		if r < 1 {
			r = 1
		}
		//TODO:: логика обслуживания самолета.
		// Обслуивание либо должено закончится по таймауту, либо если аэропорт скажет заканчивай - мы закрываемся
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r)*time.Second)
		defer cancel()

		<-ctx.Done()
		fmt.Printf("%s plane %v: is finished servicing\n", getCurrentTime(), p.title)
	}

	p.status = "parking"

	//TODO:: самолет нужно отправить на попытку взлета
	a.takeoffCh <- p
}
