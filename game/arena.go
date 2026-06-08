package game

import (
	"fmt"
	"math/rand"
	"sync"

	"github.com/google/uuid"
	"github.com/ognev-dev/goplease/app/ds"
)

const UnitsPerPlacementPhase = 3

// Arena holds the full state of one match.
type Arena struct {
	mu sync.Mutex

	ID         string
	Board      *Board
	Players    [2]*Player
	UnitsQueue []*Unit

	CurrentRound int
	ActivePlayer int // 0 or 1 whose turn is
	ActiveUnitID ds.ID

	CurrentTurn            int
	Phase                  RoundPhase
	UnitsPerPlacementPhase int
}

func NewArena(p1, p2 *Player) *Arena {
	return &Arena{
		ID:                     uuid.NewString(),
		Players:                [2]*Player{p1, p2},
		UnitsQueue:             []*Unit{},
		CurrentTurn:            0,
		ActivePlayer:           rand.Intn(2),
		Phase:                  PlacementPhase,
		Board:                  NewBoard(),
		UnitsPerPlacementPhase: UnitsPerPlacementPhase,
	}
}

func (a *Arena) playerByID(id ds.ID) (*Player, int) {
	for i, p := range a.Players {
		if p.ID == id {
			return p, i
		}
	}

	return nil, -1
}

func (a *Arena) ActingUnit() *Unit {
	for _, u := range a.UnitsQueue {
		if a.ActiveUnitID == u.ID {
			return u
		}
	}

	return nil
}

func (a *Arena) MarkReady(playerID ds.ID) bool {
	a.mu.Lock()
	defer a.mu.Unlock()

	_, idx := a.playerByID(playerID)
	if idx < 0 {
		return false
	}
	a.Players[idx].Ready = true

	return a.IsPlayersReady()
}

func (a *Arena) IsPlayersReady() bool {
	return a.Players[0].Ready && a.Players[1].Ready
}

func (a *Arena) IsPlayerPlacementDone(idx int) bool {
	p := a.Players[idx]
	return p.UnitsPlacedThisRound >= a.UnitsPerPlacementPhase || len(p.Units) == 0
}

func (a *Arena) PlacementActorIndex() int {
	p1 := a.Players[0].UnitsPlacedThisRound
	p2 := a.Players[1].UnitsPlacedThisRound
	if p2 < p1 {
		return 1
	}

	return 0 // tie-breaker: P1
}

func (a *Arena) PlayerByUnitID(unitID ds.ID) *Player {
	for _, u := range a.UnitsQueue {
		if u.ID == unitID {
			for _, p := range a.Players {
				if p.ID == u.OwnerID {
					return p
				}
			}
		}
	}
	return nil
}

func (a *Arena) PlayerByID(id ds.ID) (*Player, int) {
	for i, p := range a.Players {
		if p.ID == id {
			return p, i
		}
	}
	return nil, -1
}

func (a *Arena) PlaceUnitFromHandToBoard(templateID int, at HexCoord, playerID ds.ID) (u *Unit, err error) {
	player, playerIdx := a.PlayerByID(playerID)
	if player == nil {
		return nil, fmt.Errorf("player %q not found", playerID)
	}

	if a.PlacementActorIndex() != playerIdx {
		return nil, fmt.Errorf("not your turn to place")
	}

	cell, ok := a.Board.Cells[at]
	if !ok || cell == nil {
		return nil, fmt.Errorf("cell %q not found", at)
	}
	if cell.Unit != nil {
		return nil, fmt.Errorf("cell %q already has a unit", at)
	}
	if !cell.IsSafeZone || cell.SafeZonePlayer != playerIdx {
		return nil, fmt.Errorf("cell %q is not a placement zone", at)
	}

	u = player.PopUnitFromHand(templateID)
	if u == nil {
		return nil, fmt.Errorf("unit with template %d not found in hand", templateID)
	}

	u.Pos = at
	cell.Unit = u
	player.UnitsPlacedThisRound++
	a.UnitsQueue = append(a.UnitsQueue, u)

	return u, nil
}
