package worker

import (
	"context"
	"log"
	"time"

	"github.com/madchin/trader-bot/internal/domain/offer"
	"github.com/madchin/trader-bot/internal/gateway"
	"github.com/madchin/trader-bot/internal/service"
)

// we assume its concurrent safe
type scheduler interface {
	Delegate() (*gateway.InteractionData, error)
}

type worker struct {
	service   *service.Service
	scheduler scheduler
}

// Worker spawner, cant spawn more workers than specified in factoryWorkers config
func spawn(service *service.Service, scheduler scheduler, factoryWorkers *factoryWorkers) (*worker, error) {
	if factoryWorkers.actualCount == factoryWorkers.maximumCount {
		return nil, ErrWorkerSpawnFactoryIsFull
	}
	factoryWorkers.actualCount++
	return &worker{service, scheduler}, nil
}

func Spawner(ctx context.Context, service *service.Service, scheduler scheduler, factoryWorkers *factoryWorkers) {
	for {
		ctx, _ := context.WithTimeout(ctx, time.Second*5)
		job, _ := scheduler.Delegate()
		if job == nil {
			time.Sleep(time.Millisecond * 100)
			continue
		}
		worker, err := spawn(service, scheduler, factoryWorkers)
		if err != nil {
			log.Printf("worker spawn error: %v", err)
			time.Sleep(time.Second * 1)
			continue
		}
		go worker.execute(ctx, job)
	}
}

func (w *worker) execute(ctx context.Context, jobData *gateway.InteractionData) {
	off := jobData.Offer()
	log.Printf("executing worker job...")
	if off.Type() == offer.Buy {
		ctx = context.WithValue(ctx, "dbTableDescriptor", "buy"+"_"+jobData.Interaction().GuildID)
		log.Printf("executing buy offer with job data %v", jobData)
		w.executeBuyOffer(ctx, jobData)
	}
	if off.Type() == offer.Sell {
		ctx = context.WithValue(ctx, "dbTableDescriptor", "sell"+"_"+jobData.Interaction().GuildID)
		log.Printf("executing sell offer with job data %v", jobData)
		w.executeSellOffer(ctx, jobData)
	}
}

func (w *worker) executeBuyOffer(ctx context.Context, jobData *gateway.InteractionData) {
	switch jobData.Offer().Action() {
	case offer.Add:
		err := w.service.BuyService().Add(ctx, jobData.Interaction(), jobData.Offer())
		if err != nil {
			log.Printf("error in worker buy add: %v", err)
		}
	case offer.List:
		err := w.service.BuyService().List(ctx, jobData.Interaction(), jobData.Offer())
		if err != nil {
			log.Printf("error in worker buy list: %v", err)
		}
	case offer.Remove:
		err := w.service.BuyService().Remove(ctx, jobData.Interaction(), jobData.Offer())
		if err != nil {
			log.Printf("error in worker buy remove: %v", err)
		}
	case offer.Update:
		err := w.service.BuyService().Update(ctx, jobData.Interaction(), jobData.Offer(), jobData.Offer())
		if err != nil {
			log.Printf("error in worker buy update: %v", err)
		}
	}
}

func (w *worker) executeSellOffer(ctx context.Context, jobData *gateway.InteractionData) {
	switch jobData.Offer().Action() {
	case offer.Add:
		err := w.service.SellService().Add(ctx, jobData.Interaction(), jobData.Offer())
		if err != nil {
			log.Printf("error in worker sell add: %v", err)
		}
	case offer.List:
		err := w.service.SellService().List(ctx, jobData.Interaction(), jobData.Offer())
		if err != nil {
			log.Printf("error in worker sell list: %v", err)
		}
	case offer.Remove:
		err := w.service.SellService().Remove(ctx, jobData.Interaction(), jobData.Offer())
		if err != nil {
			log.Printf("error in worker sell remove: %v", err)
		}
	case offer.Update:
		err := w.service.SellService().Update(ctx, jobData.Interaction(), jobData.Offer(), jobData.Offer())
		if err != nil {
			log.Printf("error in worker sell update: %v", err)
		}
	}
}
