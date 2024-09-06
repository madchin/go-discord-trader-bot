package worker

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/madchin/trader-bot/internal/gateway"
	"github.com/madchin/trader-bot/internal/gateway/command"
	"github.com/madchin/trader-bot/internal/service"
	"github.com/madchin/trader-bot/internal/storage"
)

// we assume its concurrent safe
type scheduler interface {
	Delegate() (gateway.Job, error)
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
			time.Sleep(time.Millisecond * 100)
			continue
		}
		go worker.execute(ctx, ctxCancel, job)
	}
}

func (w *worker) execute(ctx context.Context, ctxCancel context.CancelFunc, job gateway.Job) {
	ctx = context.WithValue(
		ctx,
		"offer",
		storage.TableWithGuildIdSuffix(job.Data().Metadata.Subcommand(), job.Data().Metadata.Interaction().GuildID),
	)
	ctx = context.WithValue(
		ctx,
		"item",
		storage.TableWithGuildIdSuffix("item", job.Data().Metadata.Interaction().GuildID),
	)
	if isOfferCommand(job) {
		if err := w.execOffer(ctx, job); err != nil {
			log.Printf("error in worker %v", err)
		}
	}
	if isItemRegistrarCommand(job) {
		if err := w.execItemRegistrar(ctx, job); err != nil {
			log.Printf("error in worker %v", err)
		}
	}
	//ctxCancel()
}

func (w *worker) execOffer(ctx context.Context, job gateway.Job) error {
	switch job.Data().Metadata.Action() {
	case command.Offer.Action.Add.Descriptor():
		return w.service.Offer.Add(ctx, job.Data().Metadata.Interaction(), job.Data().OfferEvent.VendorOffer())
	case command.Offer.Action.ListByProductName.Descriptor():
		return w.service.Offer.ListByProductName(ctx, job.Data().Metadata.Interaction(), job.Data().OfferEvent.VendorOffer().Product.Name())
	case command.Offer.Action.ListByVendor.Descriptor():
		return w.service.Offer.ListByVendor(ctx, job.Data().Metadata.Interaction(), job.Data().OfferEvent.VendorOffer().VendorIdentity())
	case command.Offer.Action.Remove.Descriptor():
		return w.service.Offer.Remove(ctx, job.Data().Metadata.Interaction(), job.Data().OfferEvent.VendorOffer())
	case command.Offer.Action.UpdateCount.Descriptor():
		return w.service.Offer.UpdateCount(ctx, job.Data().Metadata.Interaction(), job.Data().OfferEvent.VendorOffer())
	case command.Offer.Action.UpdatePrice.Descriptor():
		return w.service.Offer.UpdatePrice(ctx, job.Data().Metadata.Interaction(), job.Data().OfferEvent.VendorOffer(), job.Data().OfferEvent.UpdatePrice())
	}
	return errors.New("sub command happened which is not registered")
}

func (w *worker) execItemRegistrar(ctx context.Context, job gateway.Job) error {
	switch job.Data().Metadata.Subcommand() {
	case command.ItemRegistrar.SubCommand.Add.Descriptor():
		return w.service.ItemRegistrar.Add(ctx, job.Data().Metadata.Interaction(), job.Data().ItemRegistrarEvent.Item())
	case command.ItemRegistrar.SubCommand.Remove.Descriptor():
		return w.service.ItemRegistrar.Remove(ctx, job.Data().Metadata.Interaction(), job.Data().ItemRegistrarEvent.Item())
	case command.ItemRegistrar.SubCommand.List.Descriptor():
		return w.service.ItemRegistrar.List(ctx, job.Data().Metadata.Interaction())
	}
	return errors.New("sub command happened which is not registered")
}

func isNoWork(job gateway.Job) bool {
	return job == nil
}

func isOfferCommand(job gateway.Job) bool {
	return command.Offer.Command.Descriptor() == job.Data().Metadata.Command()
}

func isItemRegistrarCommand(job gateway.Job) bool {
	return command.ItemRegistrar.Command.Descriptor() == job.Data().Metadata.Command()
}
