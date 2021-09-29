package sql

import (
	"context"
	"database/sql"
	"time"

	"github.com/pkg/errors"
	"github.com/thijzert/chesseract/chesseract"
	"github.com/thijzert/chesseract/chesseract/game"
	"github.com/thijzert/chesseract/internal/storage"

	"github.com/go-sql-driver/mysql"
)

func init() {
	storage.RegisterBackend("mysql", func(params string) (storage.Backend, error) {
		rv := &SQLBackend{}

		var err error

		var conf *mysql.Config
		conf, err = mysql.ParseDSN(params)
		if err != nil {
			return nil, errors.Wrap(err, "invalid DSN")
		}
		// Make sure DSN parameters are set:
		conf.ParseTime = true

		rv.conn, err = sql.Open("mysql", conf.FormatDSN())
		if err != nil {
			return nil, errors.Wrap(err, "error opening database")
		}

		_, err = rv.conn.Exec("SELECT 3.1415926")
		if err != nil {
			return nil, errors.Wrap(err, "error  opening database")
		}

		rv.conn.SetConnMaxLifetime(150 * time.Second)
		rv.conn.SetMaxIdleConns(10)
		rv.conn.SetMaxOpenConns(10)

		return rv, nil
	})
}

// The SQLBackend storage backend implements a storage backend using MySQL.
type SQLBackend struct {
	conn *sql.DB
}

func (d *SQLBackend) String() string {
	return "SQL"
}

func (d *SQLBackend) Close() error {
	return d.conn.Close()
}

type sqlTransactionRunning bool

// Transaction starts a transaction before running the function.
// If it returns any error, the transaction is rolled back. If it
// returns nil, it is committed.
// Transaction only returns an error if committing or rolling back fails.
func (d *SQLBackend) TransactionContext(ctx context.Context, f func(context.Context) error) error {
	if ctx.Value(sqlTransactionRunning(true)) != nil {
		return f(ctx)
	}

	var err error
	_, err = d.conn.ExecContext(ctx, "START TRANSACTION")
	if err != nil {
		return err
	}

	childCtx := context.WithValue(ctx, sqlTransactionRunning(true), true)
	err = f(childCtx)
	if err != nil {
		d.conn.ExecContext(ctx, "ROLLBACK")
		return err
	}

	_, err = d.conn.ExecContext(ctx, "COMMIT")
	return err
}

// NewSession creates a new session
func (d *SQLBackend) NewSessionContext(ctx context.Context) (storage.SessionID, storage.Session, error) {
	sid := storage.NewSessionID()
	_, err := d.conn.ExecContext(ctx, `
		INSERT INTO Session ( SessionID, PlayerID, Created, LastSeen, Inactive )
		VALUES ( ?, NULL, NOW(), NOW(), 0 )
	`, sid.String())
	return sid, storage.Session{}, err
}

// GetSession retrieves a session from the store
func (d *SQLBackend) GetSessionContext(ctx context.Context, id storage.SessionID) (storage.Session, error) {
	rv := storage.Session{}

	var strSID string
	var strPID sql.NullString

	err := d.conn.QueryRowContext(ctx, `
		SELECT SessionID, PlayerID FROM Session WHERE SessionID = ? AND Inactive = 0
	`, id.String()).Scan(&strSID, &strPID)

	if err == sql.ErrNoRows {
		// TODO: return storage.errNotPresent
		return rv, err
	} else if err != nil {
		return rv, err
	}

	d.conn.ExecContext(ctx, `UPDATE Session SET LastSeen = NOW() WHERE SessionID = ?`, strSID)

	if strPID.Valid {
		rv.PlayerID, err = storage.ParsePlayerID(strPID.String)
	}

	return rv, err
}

// StoreSession updates a modified Session in the datastore
func (d *SQLBackend) StoreSessionContext(ctx context.Context, id storage.SessionID, sess storage.Session) error {
	var strPID sql.NullString
	if !sess.PlayerID.IsEmpty() {
		strPID.Valid = true
		strPID.String = sess.PlayerID.String()
	}

	_, err := d.conn.ExecContext(ctx, `
		UPDATE Session
		SET PlayerID = ?,
			LastSeen = NOW()
		WHERE SessionID = ? AND Inactive = 0
	`, strPID, id.String())

	return err
}

// NewPlayer creates a new player
func (d *SQLBackend) NewPlayerContext(ctx context.Context) (storage.PlayerID, game.Player, error) {
	pid := storage.NewPlayerID()
	_, err := d.conn.ExecContext(ctx, `
		INSERT INTO Player ( PlayerID ) VALUES ( ? )
	`, pid.String())
	return pid, game.Player{}, err
}

// GetPlayer retrieves a player from the store
func (d *SQLBackend) GetPlayerContext(ctx context.Context, id storage.PlayerID) (game.Player, error) {
	var name, realm sql.NullString
	rv := game.Player{}
	err := d.conn.QueryRowContext(ctx, `
		SELECT Name, Realm, GenderR, GenderI, GenderJ, GenderK, ELORating
		FROM Player WHERE PlayerID = ?
	`, id.String()).Scan(&name, &realm, &rv.Gender.R, &rv.Gender.I, &rv.Gender.J, &rv.Gender.K, &rv.ELORating)

	if name.Valid {
		rv.Name = name.String
	}
	if realm.Valid {
		rv.Realm = realm.String
	}

	return rv, err
}

// StorePlayer updates a modified Player in the datastore
func (d *SQLBackend) StorePlayerContext(ctx context.Context, id storage.PlayerID, player game.Player) error {
	_, err := d.conn.ExecContext(ctx, `
		UPDATE Player
		SET Name = ?,
			Realm = ?,
			GenderR = ?, GenderI = ?, GenderJ = ?, GenderK = ?,
			ELORating = ?
		WHERE PlayerID = ?
	`, player.Name, player.Realm, player.Gender.R, player.Gender.I, player.Gender.J, player.Gender.K, player.ELORating, id.String())

	return err
}

func (d *SQLBackend) LookupPlayerContext(ctx context.Context, name string) (storage.PlayerID, bool, error) {
	var strPID string
	err := d.conn.QueryRowContext(ctx, `SELECT PlayerID FROM Player WHERE Name = ?`, name).Scan(&strPID)
	if err == sql.ErrNoRows {
		return storage.PlayerID{}, false, nil
	} else if err != nil {
		return storage.PlayerID{}, false, err
	}

	rv, err := storage.ParsePlayerID(strPID)
	return rv, true, err
}

// NewNonceForPlayer generates a new nonce, and assigns it to the player
// It should also invalidate any existing nonces for this player.
func (d *SQLBackend) NewNonceForPlayerContext(ctx context.Context, id storage.PlayerID) (storage.Nonce, error) {
	var err error
	if _, err = d.GetPlayerContext(ctx, id); err != nil {
		return "", err
	}

	nn := storage.NewNonce()

	_, err = d.conn.ExecContext(ctx, `DELETE FROM Nonce WHERE PlayerID = ?`, id.String())
	if err != nil {
		return "", err
	}

	_, err = d.conn.ExecContext(ctx, `INSERT INTO Nonce ( Nonce, PlayerID ) VALUES ( ?, ? )`, nn, id.String())
	if err != nil {
		return "", err
	}
	return nn, nil
}

// CheckNonce checks if the nonce exists, and is assigned to that player. A
// successful result invalidates the nonce. (Implied in the 'once' part in
// 'nonce')
func (d *SQLBackend) CheckNonceContext(ctx context.Context, id storage.PlayerID, nonce storage.Nonce) (bool, error) {
	var err error
	var nn string

	err = d.conn.QueryRowContext(ctx, `
		SELECT Nonce FROM Nonce WHERE Nonce = ? AND PlayerID = ?
	`, nonce, id.String()).Scan(&nn)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}

	_, err = d.conn.ExecContext(ctx, `DELETE FROM Nonce WHERE Nonce = ?`, nonce)
	return true, err
}

// NewGame creates a new game
func (d *SQLBackend) NewGameContext(ctx context.Context) (storage.GameID, game.Game, error) {
	g := game.Game{}
	gid := storage.NewGameID()
	_, err := d.conn.ExecContext(ctx, `INSERT INTO Match_ ( MatchID, StartTime ) VALUES ( ?, NOW() )`, gid.String())
	if err != nil {
		return storage.GameID{}, g, err
	}
	return gid, g, nil
}

// GetGame retrieves a game from the store
func (d *SQLBackend) GetGameContext(ctx context.Context, id storage.GameID) (game.Game, error) {
	rv := game.Game{}

	var ruleSet string
	err := d.conn.QueryRowContext(ctx, `
		SELECT RuleSet, StartTime FROM Match_ WHERE MatchID = ?
	`, id.String()).Scan(&ruleSet, &rv.Match.StartTime)
	if err == sql.ErrNoRows {
		return rv, err
	} else if err != nil {
		return rv, err
	}

	rv.Match.RuleSet = chesseract.GetRuleSet(ruleSet)

	// Get Players
	roles := rv.Match.RuleSet.PlayerColours()
	rows, err := d.conn.QueryContext(ctx, `SELECT PlayerID, Role FROM MatchRole WHERE MatchID = ?`, id.String())
	if err != nil {
		return rv, err
	}
	type playerRole struct {
		PlayerID  string
		PlayingAs chesseract.Colour
	}
	playerIDs := make([]playerRole, 0, len(roles))
	for rows.Next() {
		var pr playerRole
		if err = rows.Scan(&pr.PlayerID, &pr.PlayingAs); err != nil {
			return rv, err
		}
		playerIDs = append(playerIDs, pr)
	}
	err = rows.Close()
	if err != nil {
		return rv, err
	}

	for _, c := range roles {
		for i, pr := range playerIDs {
			if pr.PlayingAs == c {
				pid, err := storage.ParsePlayerID(pr.PlayerID)
				if err != nil {
					return rv, err
				}
				player, err := d.GetPlayerContext(ctx, pid)
				if err != nil {
					return rv, err
				}
				rv.Players = append(rv.Players, game.MatchPlayer{
					Player:    player,
					PlayingAs: pr.PlayingAs,
				})
				playerIDs[i].PlayerID = ""
			}
		}
	}

	// Load moves
	rv.Match.Board = rv.Match.RuleSet.DefaultBoard()
	rows, err = d.conn.QueryContext(ctx, `SELECT From_, To_, Time_ FROM Move WHERE MatchID = ? ORDER BY Ordinal`, id.String())
	if err != nil {
		return rv, err
	}
	for rows.Next() {
		var sFrom, sTo string
		var seconds float64
		err = rows.Scan(&sFrom, &sTo, &seconds)
		if err != nil {
			return rv, err
		}
		p, err := rv.Match.RuleSet.ParsePosition(sFrom)
		if err != nil {
			return rv, err
		}
		q, err := rv.Match.RuleSet.ParsePosition(sTo)
		if err != nil {
			return rv, err
		}
		mv := chesseract.Move{
			From: p,
			To:   q,
			Time: time.Duration(int64(1000000.0*seconds) * int64(time.Microsecond)),
		}
		if pt, ok := rv.Match.Board.At(mv.From); ok {
			mv.PieceType = pt.PieceType
		}

		newb, err := rv.Match.RuleSet.ApplyMove(rv.Match.Board, mv)
		if err != nil {
			return rv, err
		}
		rv.Match.Board = newb
		rv.Match.Moves = append(rv.Match.Moves, mv)
	}
	err = rows.Close()
	if err != nil {
		return rv, err
	}

	return rv, nil
}

// StoreGame updates a modified Game in the datastore
func (d *SQLBackend) StoreGameContext(ctx context.Context, id storage.GameID, match game.Game) error {
	var ruleSet string
	err := d.conn.QueryRowContext(ctx, `
		SELECT RuleSet FROM Match_ WHERE MatchID = ?
	`, id.String()).Scan(&ruleSet)
	if err != nil {
		return err
	}

	finalised := false
	for _, r := range match.Result {
		if r != 0.0 {
			finalised = true
		}
	}
	_, err = d.conn.ExecContext(ctx, `
		UPDATE Match_
		SET RuleSet = ?,
			StartTime = ?,
			Finalised = ?
		WHERE MatchID = ?
	`, match.Match.RuleSet.String(), match.Match.StartTime, finalised, id.String())
	if err != nil {
		return err
	}

	playersDirty := true
	// FIXME: this might not be the case

	if playersDirty {
		_, err = d.conn.ExecContext(ctx, `DELETE FROM MatchRole WHERE MatchID = ?`, id.String())
		if err != nil {
			return err
		}

		for i, c := range match.Match.RuleSet.PlayerColours() {
			res := 0.0
			if len(match.Result) > i {
				res = match.Result[i]
			}
			for _, pl := range match.Players {
				if pl.PlayingAs == c {
					// HACK: players don't know their own ID. Should I change that?
					pid, ok, err := d.LookupPlayerContext(ctx, pl.Name)
					if err != nil {
						return nil
					}
					if ok {
						_, err = d.conn.ExecContext(ctx, `
						INSERT INTO MatchRole ( MatchID, PlayerID, Role, Result )
						VALUES ( ?, ?, ?, ? )
					`, id.String(), pid.String(), pl.PlayingAs, res)
						if err != nil {
							return err
						}
					}
				}
			}
		}
	}

	// FIXME: find a way of detecting if we can append a game rather than throwing away the whole thing every time
	_, err = d.conn.ExecContext(ctx, `DELETE FROM Move WHERE MatchID = ?`, id.String())
	if err != nil {
		return err
	}
	for i, mv := range match.Match.Moves {
		_, err := d.conn.ExecContext(ctx, `
			INSERT INTO Move ( MatchID, Ordinal, From_, To_, Time_ )
			VALUES ( ?, ?, ?, ?, ? )
		`, id.String(), i+1, mv.From.String(), mv.To.String(), mv.Time.Seconds())
		if err != nil {
			return err
		}
	}

	return nil
}

// GetActiveGames returns the GameID's of all active games in which the
// Player identified by the PlayerID is a participant
func (d *SQLBackend) GetActiveGamesContext(ctx context.Context, id storage.PlayerID) ([]storage.GameID, error) {
	rows, err := d.conn.QueryContext(ctx, `
		SELECT MatchID
		FROM MatchRole
			INNER JOIN Match_ USING ( MatchID )
		WHERE PlayerID = ? AND Finalised = 0
		GROUP BY MatchID
		ORDER BY StartTime DESC
	`, id.String())
	if err != nil {
		return nil, err
	}

	var rv []storage.GameID = nil
	for rows.Next() {
		var strGID string
		err = rows.Scan(&strGID)
		if err != nil {
			return nil, err
		}
		gid, err := storage.ParseGameID(strGID)
		if err != nil {
			return nil, err
		}
		rv = append(rv, gid)
	}

	return rv, rows.Close()
}
