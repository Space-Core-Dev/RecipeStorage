package declarations

import (
	"errors"
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/Space-Core-Dev/Mistaken.Utils/sequence"
	"github.com/rabbitmq/amqp091-go"
)

type RabbitDeclaration struct {
	queuePref string
	keyPref   string
	kind      string
	ch        *amqp091.Channel
	names     map[string][]string
}

func New(ch *amqp091.Channel, kind, queuePref, keyPref string, names map[string][]string) RabbitDeclaration {
	return RabbitDeclaration{
		kind:      kind,
		queuePref: queuePref,
		keyPref:   keyPref,
		ch:        ch,
		names:     names,
	}
}

func (r RabbitDeclaration) Start() error {
	err := r.areNamesOk()
	if err != nil {
		return err
	}

	err = r.checkParams()
	if err != nil {
		return err
	}

	if !r.checkKind() {
		return errors.New("kind should be \"topic, direct, fanout or headers\"")
	}

	exchanges := slices.Collect(maps.Keys(r.names))
	allQueues := slices.Collect(maps.Values(r.names))
	err = r.exchangesDeclare(exchanges)
	if err != nil {
		return err
	}

	for _, queues := range allQueues {
		err = r.queuesDeclare(queues)
		if err != nil {
			return err
		}
	}

	err = r.queuesBinding(r.names)
	if err != nil {
		return err
	}

	return nil
}

func (r RabbitDeclaration) checkParams() error {
	if r.ch == nil {
		return errors.New("channel is empty")
	}
	if strings.TrimSpace(r.queuePref) == "" {
		return errors.New("queue prefix is empty")
	}
	if strings.TrimSpace(r.keyPref) == "" {
		return errors.New("key prefix is empty")
	}
	if strings.TrimSpace(r.kind) == "" {
		return errors.New("kind is empty")
	}

	return nil
}

func (r RabbitDeclaration) checkKind() bool {
	return (r.kind == "topic" || r.kind == "direct" || r.kind == "fanout" || r.kind == "headers")
}

func (r RabbitDeclaration) areNamesOk() error {
	if len(r.names) == 0 {
		return errors.New("the names map shouldn't be empty")
	}

	if sequence.Any(slices.Collect(maps.Keys(r.names)), func(k string) bool {
		return strings.TrimSpace(k) == ""
	}) {
		return errors.New("an exchange name shouldn't be empty")
	}

	if sequence.Any(slices.Collect(maps.Values(r.names)), func(v []string) bool {
		return len(v) == 0 || sequence.Any(v, func(s string) bool {
			return strings.TrimSpace(s) == ""
		})
	}) {
		return errors.New("a queue name shouldn't be empty")
	}

	return nil
}

func (r RabbitDeclaration) exchangesDeclare(exchanges []string) error {
	for _, exchange := range exchanges {
		err := r.ch.ExchangeDeclare(exchange, r.kind, true, false, false, false, nil)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r RabbitDeclaration) queuesDeclare(queues []string) error {
	for _, name := range queues {
		_, err := r.ch.QueueDeclare(stickPref(r.queuePref, name), true, false, false, false, nil)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r RabbitDeclaration) queuesBinding(qParams map[string][]string) error {
	for exchange, names := range qParams {
		for _, name := range names {
			err := r.ch.QueueBind(stickPref(r.queuePref, name), stickPref(r.keyPref, name), exchange, false, nil)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func stickPref(pref, name string) string {
	return fmt.Sprintf("%s.%s", pref, name)
}
