package main

import (
	"context"
	"log"
	"math/rand"
	"time"
)

func main() {
	p1 := newPlayer("p1")
	p2 := newPlayer("p2")
	t := table{players: [2]player{p1, p2}}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	ctx = context.WithValue(ctx, "key", "value")
	t.play(ctx)
}

func newPlayer(name string) player {
	return player{name: name, in: make(chan ball), out: make(chan ball)}
}

type ball struct {
	hit int
}

type player struct {
	name string
	// in is a channel used for two-way communication.
	in  chan ball
	out chan ball
}

func (p player) play() {
	for b := range p.in {
		took := time.Duration(250+rand.Intn(100)) * time.Millisecond
		time.Sleep(took)
		log.Printf("%s was playing for %s", p.name, took)
		b.hit++
		p.out <- b
	}
}

type table struct {
	players [2]player
}

// local state can be safely mutated, since the access is serialized by channel operations
func (t *table) play(ctx context.Context) {
	if value, ok := ctx.Value("key").(string); ok {
		log.Printf("value for key is: %q\n", value)
	}

	go t.players[0].play()
	go t.players[1].play()

	b := ball{}

	for {
		select {
		// There is NO GUARANTEE that this case will be handled before others.
		// The order of cases does not matter. When more than once the operation can be performed, one is picked at random.
		case <-ctx.Done():
			log.Printf("game was stopped: %s", ctx.Err())
			return
		// When there is more than one consumer of the channel, one is picked at random.
		// Usually reading from and writing to the channel from the same goroutine is a no-no.
		// It greatly increases the chances of a deadlock or a bug.
		// That's why player receives the ball from channel in and 'returns' it through channel out.
		case b = <-t.players[0].out:
			log.Printf("ping %d\n", b.hit)
			// Mutating the local state is thread-safe, since there is only one goroutine that can do it at once.
			// Even if more than one case of the select block does not block, the second one can be executed after the first case was handled.
			t.players[0], t.players[1] = t.players[1], t.players[0]
		case t.players[0].in <- b:
			log.Printf("pong %d\n", b.hit)
		}
	}
}
