package sql

import (
	"context"
)

func (d *SQLBackend) InitialiseContext(ctx context.Context) error {
	var err error
	_, err = d.conn.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS Player (
			PlayerID   CHAR(33)     CHARSET ASCII   NOT NULL DEFAULT '',
			Name       VARCHAR(50)  CHARSET UTF8MB4     NULL,
			Realm      VARCHAR(50)  CHARSET UTF8MB4     NULL,
			GenderR    DOUBLE                       NOT NULL DEFAULT 0.0,
			GenderI    DOUBLE                       NOT NULL DEFAULT 0.0,
			GenderJ    DOUBLE                       NOT NULL DEFAULT 0.0,
			GenderK    DOUBLE                       NOT NULL DEFAULT 0.0,
			ELORating  DECIMAL(7,2)                 NOT NULL DEFAULT 100.00,
			Pubkey     CHAR(64)     CHARSET ASCII   NOT NULL DEFAULT '',
			PRIMARY KEY ( PlayerID ),
			UNIQUE KEY uq_player_realm ( Name, Realm )
		) ENGINE=InnoDB
	`)
	if err != nil {
		return err
	}

	_, err = d.conn.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS Match_ (
			MatchID    CHAR(33)     CHARSET ASCII   NOT NULL DEFAULT '',
			RuleSet    CHAR(15)     CHARSET UTF8MB4 NOT NULL DEFAULT '',
			StartTime  DATETIME                     NOT NULL,
			Finalised  TINYINT(1)                   NOT NULL DEFAULT 0,
			PRIMARY KEY ( MatchID )
		) ENGINE=InnoDB
	`)
	if err != nil {
		return err
	}

	_, err = d.conn.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS MatchRole (
			MatchID    CHAR(33)     CHARSET ASCII   NOT NULL,
			PlayerID   CHAR(33)     CHARSET ASCII   NOT NULL,
			Role       INT                          NOT NULL,
			Result     DECIMAL(8,6)                 NOT NULL DEFAULT 0.000000,
			PRIMARY KEY ( MatchID, PlayerID ),
			FOREIGN KEY ( MatchID ) REFERENCES Match_(MatchID) ON UPDATE CASCADE ON DELETE RESTRICT,
			FOREIGN KEY ( PlayerID ) REFERENCES Player(PlayerID) ON UPDATE CASCADE ON DELETE RESTRICT
		) ENGINE=InnoDB
	`)
	if err != nil {
		return err
	}

	_, err = d.conn.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS Move (
			MatchID    CHAR(33)     CHARSET ASCII   NOT NULL,
			Ordinal    INT                          NOT NULL,
			From_      CHAR(8)      CHARSET UTF8MB4 NOT NULL DEFAULT '',
			To_        CHAR(8)      CHARSET UTF8MB4 NOT NULL DEFAULT '',
			Time_      DECIMAL(9,3)                 NOT NULL DEFAULT 0.000,
			PRIMARY KEY ( MatchID, Ordinal ),
			FOREIGN KEY ( MatchID ) REFERENCES Match_(MatchID) ON UPDATE CASCADE ON DELETE RESTRICT
		) ENGINE=InnoDB
	`)
	if err != nil {
		return err
	}

	_, err = d.conn.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS Session (
			SessionID  CHAR(68)     CHARSET ASCII   NOT NULL,
			PlayerID   CHAR(33)     CHARSET ASCII   NULL,
			Created    DATETIME                     NOT NULL,
			LastSeen   DATETIME                     NOT NULL,
			Inactive   TINYINT(1)                   NOT NULL DEFAULT 0,
			PRIMARY KEY ( SessionID ),
			FOREIGN KEY ( PlayerID ) REFERENCES Player(PlayerID) ON UPDATE CASCADE ON DELETE CASCADE
		) ENGINE=InnoDB
	`)
	if err != nil {
		return err
	}

	_, err = d.conn.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS Nonce (
			Nonce      CHAR(68)     CHARSET ASCII   NOT NULL,
			PlayerID   CHAR(33)     CHARSET ASCII   NOT NULL,
			PRIMARY KEY ( Nonce ),
			FOREIGN KEY ( PlayerID ) REFERENCES Player(PlayerID) ON UPDATE CASCADE ON DELETE CASCADE
		) ENGINE=InnoDB
	`)
	if err != nil {
		return err
	}

	return nil
}
