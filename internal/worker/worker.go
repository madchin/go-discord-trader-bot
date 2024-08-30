package worker

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/madchin/trader-bot/internal/gateway"
	"github.com/madchin/trader-bot/internal/service"
	"github.com/madchin/trader-bot/internal/storage"
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
		job, _ := scheduler.Delegate()
		if isNoWork(job) {
			time.Sleep(time.Millisecond * 100)
			continue
		}
		ctx, ctxCancel := context.WithTimeout(ctx, time.Second*5)
		worker, err := spawn(service, scheduler, factoryWorkers)
		if err != nil {
			log.Printf("worker spawn error: %v", err)
			time.Sleep(time.Second * 1)
			continue
		}
		go worker.execute(ctx, ctxCancel, job)
	}
}

func (w *worker) execute(ctx context.Context, ctxCancel context.CancelFunc, jobData *gateway.InteractionData) {
	ctx = context.WithValue(ctx, storage.DbTableDescriptorKey, storage.DbTableDescriptorValue(jobData.Command(), jobData.Interaction().GuildID))
	if err := w.exec(ctx, jobData); err != nil {
		log.Printf("error in worker with database descriptor %s %v", ctx.Value(storage.DbTableDescriptorKey), err)
	}
	ctxCancel()
}

func (w *worker) exec(ctx context.Context, jobData *gateway.InteractionData) error {
	switch jobData.Subcommand() {
	case gateway.AddSubCmdDescriptor.Descriptor():
		return w.service.Offer().Add(ctx, jobData.Interaction(), jobData.Offer())
	case gateway.ListByProductNameSubCmdDescriptor.Descriptor():
		return w.service.Offer().ListByProductName(ctx, jobData.Interaction(), jobData.Offer().Product().Name())
	case gateway.ListByVendorSubCmdDescriptor.Descriptor():
		return w.service.Offer().ListByProductName(ctx, jobData.Interaction(), jobData.Offer().Vendor().Name())
	case gateway.RemoveSubCmdDescriptor.Descriptor():
		return w.service.Offer().Remove(ctx, jobData.Interaction(), jobData.Offer())
	case gateway.UpdateCountSubCmdDescriptor.Descriptor():
		return w.service.Offer().UpdateCount(ctx, jobData.Interaction(), jobData.Offer(), jobData.UpdateOffer().Count())
	case gateway.UpdatePriceSubCmdDescriptor.Descriptor():
		return w.service.Offer().UpdatePrice(ctx, jobData.Interaction(), jobData.Offer(), jobData.UpdateOffer().Product().Price())
	}
	return errors.New("sub command happened which is not registered")
}

func isNoWork(job *gateway.InteractionData) bool {
	return job == nil
}
