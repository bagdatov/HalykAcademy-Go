package airport

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"golang.org/x/sync/semaphore"
)

const maxProcesses = 3

//TODO:: сделать эту структуру потокабезопастной
type runway struct {
	rw     sync.Mutex
	action func(*Plane) (<-chan *Plane, error)
}

type fleet struct {
	fl   sync.Mutex
	pool map[string]*Plane
}

type AirPort struct {
	runway

	//TODO:: сделать этот элемент потока безопасным
	fleet

	// каналы для управления посадками и взлетами
	takeoffCh chan *Plane
	landingCh chan *Plane

	// поля которые помогут вам закрыть аэропорт
	stctx context.Context
	Stop  context.CancelFunc

	// после которое говорит об завершении всех дел - программа может умирать
	Done chan struct{}
}

var errChanClosed = errors.New("channel is closed")
var errNil = errors.New("nil argument was passed")

// NewAirport создание новой аэропорта и запуска его действия в отдельной горутине
func NewAirport() *AirPort {
	stctx, stop := context.WithCancel(context.Background())

	a := &AirPort{
		fleet:     fleet{pool: make(map[string]*Plane, 10)},
		landingCh: make(chan *Plane),
		takeoffCh: make(chan *Plane),

		stctx: stctx,
		Stop:  stop,
	}
	a.action = a.runwayProcess

	//TODO:: заполнить поле planePool разными самолетами
	for {
		plane := NewPlane()
		a.fl.Lock()
		a.pool[plane.title] = plane
		a.fl.Unlock()

		if len(a.pool) == 10 {
			break
		}
	}

	go a.SendPlanes()
	// запускаем ассинхроную функцию для функционирования аэропорта
	//TODO:: как обработать ошибку от функции airportProcess?
	go func() {
		err := a.airportProcess()
		if err != nil {
			fmt.Println(err)
		}
	}()

	return a
}

func (a *AirPort) SendPlanes() error {
	a.fl.Lock()
	for _, plane := range a.pool {
		a.takeoffCh <- plane
	}
	a.fl.Unlock()
	return nil
}

// фнкция обслуживания
func (a *AirPort) airportProcess() error {
	ctx := context.Background()
	planesService := semaphore.NewWeighted(maxProcesses)

LOOP:
	for {
		select {
		// логика взлета самолета
		case plane, ok := <-a.takeoffCh:
			if !ok {
				return fmt.Errorf("takeoffCh: %w", errChanClosed)
			}

			go plane.start(ctx, a)

		case plane, ok := <-a.landingCh:
			if !ok {
				return fmt.Errorf("landingCh: %w", errChanClosed)
			}
			plane.land(ctx, a)

			// проверка на есть ли ресурсы, и если да - то я займу место в пуле
			if err := planesService.Acquire(ctx, 1); err != nil {
				return fmt.Errorf("wait for resources: %w", err)
			}

			// функиц для обслуживания конкретного самолета
			//TODO::
			go a.maintenance(planesService, plane)

			// go plane.someFuncFortakeoffCh(planesService * semaphore.Weighted)...
		// момент когда закрывается аэропорт
		case <-a.stctx.Done():
			if err := planesService.Acquire(ctx, maxProcesses); err != nil {
				return fmt.Errorf("wait for complete: %w", err)
			}
			//TODO::
			time.Sleep(3 * time.Second)
			log.Println("airport successfully finished its work")
			a.Done <- struct{}{}
			break LOOP
		}
	}
	return nil
}

func (r *runway) runwayProcess(p *Plane) (<-chan *Plane, error) {
	r.rw.Lock()
	defer r.rw.Unlock()

	if p == nil {
		return nil, fmt.Errorf("runwayProcess: %w", errNil)
	}

	ch := make(chan *Plane)
	go func() {
		ch <- p
	}()
	return ch, nil
}

func (a *AirPort) maintenance(planesService *semaphore.Weighted, p *Plane) error {
	if p == nil {
		return fmt.Errorf("maintenance: %w", errNil)
	}
	p.status = maintenance
	log.Printf("plane number: %v is under maintenence\n", p.title)

	time.Sleep(randomTime(3, 5) * time.Second)
	planesService.Release(1)

	p.status = parked
	log.Printf("plane number: %v finished maintenence and in the parking lot\n", p.title)

	if status := a.stctx.Err(); status == nil {
		go func() {
			a.takeoffCh <- p
		}()
	}
	return nil
}
