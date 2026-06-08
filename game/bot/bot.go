package bot

import (
	"github.com/ognev-dev/goplease/game"
)

// Bot implements a simple greedy AI.
type Bot struct{}

func New() *Bot { return &Bot{} }

// TakeTurn places units and calls EndTurn on behalf of the bot player.
// Called from a goroutine; the room's internal mutex handles concurrency.
func (b *Bot) TakeTurn(room *game.Arena, p *game.Player) {
	if room.Phase != game.PlacementPhase {
		return
	}

	// TODO
}
