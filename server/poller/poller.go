package poller

import (
	"context"
	"sync"
	"time"

	gmg "github.com/brandenc40/green-mountain-grill"
	"github.com/brandenc40/green-mountain-grill/server/respository"
	"github.com/brandenc40/green-mountain-grill/server/respository/mapper"
	"github.com/brandenc40/green-mountain-grill/server/respository/model"
	"github.com/google/uuid"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Provide(New)

type Params struct {
	fx.In

	GrillClient gmg.Client
	Repository  respository.Repository
	Logger      *zap.Logger
	Lifecycle   fx.Lifecycle
}

// New -
func New(p Params) *Poller {
	poll := &Poller{
		grillClient: p.GrillClient,
		logger:      p.Logger,
		repo:        p.Repository,
		subscribers: make(map[uuid.UUID]chan *model.GrillState),
	}
	p.Lifecycle.Append(
		fx.Hook{
			OnStop: func(ctx context.Context) error {
				p.Logger.Info("Stopping grill polling")
				poll.StopPolling()
				time.Sleep(time.Millisecond) // time to ensure all is shutdown properly
				return nil
			},
		},
	)
	return poll
}

// Poller -
type Poller struct {
	grillClient     gmg.Client
	logger          *zap.Logger
	repo            respository.Repository
	currentSession  uuid.UUID
	stopChan        chan bool
	isPolling       bool
	subscribers     map[uuid.UUID]chan *model.GrillState
	subscriberMutex sync.Mutex
}

// NewSession -
func (p *Poller) NewSession() {
	p.currentSession = uuid.New()
}

// SetSession -
func (p *Poller) SetSession(sessionUUID uuid.UUID) {
	p.currentSession = sessionUUID
}

// CurrentSession -
func (p *Poller) CurrentSession() uuid.UUID {
	return p.currentSession
}

// StopPolling -
func (p *Poller) StopPolling() {
	if p.IsPolling() {
		p.stopChan <- true
	}
}

// IsPolling -
func (p *Poller) IsPolling() bool {
	return p.isPolling
}

// StartPolling -
func (p *Poller) StartPolling(interval time.Duration) error {
	_, err := p.grillClient.GetState()
	if err != nil {
		return err
	}
	if p.IsPolling() {
		p.StopPolling()
	}
	if p.currentSession == uuid.Nil {
		p.NewSession()
	}
	go p.pollGrill(interval)
	return nil
}

// Subscribers - The current number of subscribers
func (p *Poller) Subscribers() int {
	return len(p.subscribers)
}

// Subscribe - Subscribe to all polls of the current grill state, the state structs will be sent through the channel
// when they are polled from the grill.
//
// `defer unsubscribe()` should be used to ensure channels are closed after the subscriber is finished subscribing
func (p *Poller) Subscribe() (channel chan *model.GrillState, unsubscribe func()) {
	p.subscriberMutex.Lock()
	defer p.subscriberMutex.Unlock()
	u := uuid.New()
	channel = make(chan *model.GrillState)
	p.subscribers[u] = channel
	unsubscribe = func() {
		p.subscriberMutex.Lock()
		defer p.subscriberMutex.Unlock()
		delete(p.subscribers, u)
		close(channel)
	}
	return
}

func (p *Poller) pollGrill(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	p.isPolling = true
	for p.isPolling {
		select {
		case <-ticker.C:
			s, err := p.grillClient.GetState()
			if err != nil {
				p.logger.Error("error getting polling grill state", zap.Error(err))
				p.isPolling = false
			}
			m := mapper.GrillStateEntityToModel(s, p.currentSession)
			p.broadcastToSubscribers(m)
			if err := p.repo.InsertStateData(m); err != nil {
				p.logger.Error("error inserting state data", zap.Error(err))
				p.isPolling = false
			}
		case <-p.stopChan:
			p.isPolling = false
		}
	}
}

func (p *Poller) broadcastToSubscribers(model *model.GrillState) {
	p.subscriberMutex.Lock()
	defer p.subscriberMutex.Unlock()
	for _, subChan := range p.subscribers {
		subChan <- model
	}
}
