package chesseract

import (
	"encoding/json"
	"time"

	"github.com/pkg/errors"
)

// The matchJsonProxy struct is a JSON proxy for the Match struct
type matchJsonProxy struct {
	RuleSet   string          `json:"ruleset"`
	Board     boardJsonProxy  `json:"board"`
	StartTime string          `json:"start_time,omitempty"`
	Moves     []moveJsonProxy `json:"moves"`
}

// The positionJsonProxy is a JSON proxy for the Position interface
type positionJsonProxy string

// The pieceJsonProxy struct is a JSON proxy for the Piece struct
type pieceJsonProxy struct {
	PieceType PieceType         `json:"type"`
	Colour    Colour            `json:"colour"`
	Position  positionJsonProxy `json:"position"`
}

// The boardJsonProxy struct is a JSON proxy for the Board struct
type boardJsonProxy struct {
	Pieces []pieceJsonProxy `json:"pieces"`
	Turn   Colour           `json:"turn"`
}

// The moveJsonProxy struct is a JSON proxy for the Move struct
type moveJsonProxy struct {
	PieceType PieceType         `json:"type"`
	From      positionJsonProxy `json:"from"`
	To        positionJsonProxy `json:"to"`
	Time      string            `json:"time,omitempty"`
}

func (m Move) MarshalJSON() ([]byte, error) {
	proxy := moveJsonProxy{
		PieceType: m.PieceType,
		From:      positionJsonProxy(m.From.String()),
		To:        positionJsonProxy(m.To.String()),
		Time:      m.Time.String(),
	}
	return json.Marshal(proxy)
}

func (m Match) MarshalJSON() ([]byte, error) {
	proxy := matchJsonProxy{
		RuleSet: m.RuleSet.String(),
		Board: boardJsonProxy{
			Turn: m.Board.Turn,
		},
		StartTime: m.StartTime.Format(time.RFC3339),
	}

	for _, pc := range m.Board.Pieces {
		proxy.Board.Pieces = append(proxy.Board.Pieces, pieceJsonProxy{
			PieceType: pc.PieceType,
			Colour:    pc.Colour,
			Position:  positionJsonProxy(pc.Position.String()),
		})
	}

	for _, mv := range m.Moves {
		proxy.Moves = append(proxy.Moves, moveJsonProxy{
			PieceType: mv.PieceType,
			From:      positionJsonProxy(mv.From.String()),
			To:        positionJsonProxy(mv.To.String()),
			Time:      mv.Time.String(),
		})
	}

	return json.Marshal(proxy)
}

func (m *Match) UnmarshalJSON(buf []byte) error {
	proxy := matchJsonProxy{}
	err := json.Unmarshal(buf, &proxy)
	if err != nil {
		return err
	}

	m.RuleSet = GetRuleSet(proxy.RuleSet)
	if m.RuleSet == nil {
		return errors.New("unknown rule set")
	}

	m.StartTime, err = time.Parse(time.RFC3339, proxy.StartTime)
	if err != nil {
		return errors.Wrap(err, "error decoding match")
	}

	m.Board = Board{
		Pieces: nil,
		Turn:   proxy.Board.Turn,
	}
	for _, pc := range proxy.Board.Pieces {
		pos, err := m.RuleSet.ParsePosition(string(pc.Position))
		if err != nil {
			return errors.Wrap(err, "error decoding match")
		}
		m.Board.Pieces = append(m.Board.Pieces, Piece{
			PieceType: pc.PieceType,
			Colour:    pc.Colour,
			Position:  pos,
		})
	}

	m.Moves = nil
	for _, mv := range proxy.Moves {
		from, err := m.RuleSet.ParsePosition(string(mv.From))
		if err != nil {
			return errors.Wrap(err, "error decoding match")
		}
		to, err := m.RuleSet.ParsePosition(string(mv.To))
		if err != nil {
			return errors.Wrap(err, "error decoding match")
		}
		dur, err := time.ParseDuration(mv.Time)
		if err != nil {
			return errors.Wrap(err, "error decoding match")
		}
		m.Moves = append(m.Moves, Move{
			PieceType: mv.PieceType,
			From:      from,
			To:        to,
			Time:      dur,
		})
	}

	return nil
}
