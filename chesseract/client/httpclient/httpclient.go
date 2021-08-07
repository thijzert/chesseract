package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/thijzert/chesseract/chesseract"
	"github.com/thijzert/chesseract/chesseract/client"
	"github.com/thijzert/chesseract/chesseract/game"
	"github.com/thijzert/chesseract/internal/notimplemented"
	"github.com/thijzert/chesseract/internal/storage"
	"github.com/thijzert/chesseract/web"
	http2curl "moul.io/http2curl/v2"
)

type apiError struct {
	ErrorCode     int    `json:"error_code"`
	ErrorHeadline string `json:"error"`
	ErrorMessage  string `json:"message"`
}

func (a apiError) Error() string {
	return fmt.Sprintf("[0x%02x] %s", a.ErrorCode, a.ErrorMessage)
}

type ClientConfig struct {
	ServerURI string
	Username  string

	VerboseRequestLogging io.Writer
}

// A HttpClient is an implementation of a Client that talks to a Chesseract server over HTTP
type HttpClient struct {
	client *http.Client

	baseURI string

	sessionToken string

	requestLogger *log.Logger

	gameIDs map[*game.Game]storage.GameID
}

func New(ctx context.Context, c ClientConfig) (*HttpClient, error) {
	rv := &HttpClient{
		client: &http.Client{
			Transport: &http.Transport{
				Dial: (&net.Dialer{
					Timeout: 2500 * time.Millisecond,
				}).Dial,
			},
			CheckRedirect: func(*http.Request, []*http.Request) error {
				// Do not follow redirects for any reason.
				return http.ErrUseLastResponse
			},
			Timeout: 0,
		},

		baseURI: c.ServerURI,

		gameIDs: make(map[*game.Game]storage.GameID),
	}

	if c.VerboseRequestLogging != nil {
		rv.requestLogger = log.New(c.VerboseRequestLogging, "HTTP", log.Ltime)
	}

	var sess web.NewSessionResponse
	err := rv.get(ctx, &sess, "/api/session/new", nil)
	if err != nil {
		return nil, errors.Wrap(err, "error starting client")
	}

	rv.sessionToken = sess.SessionID

	auth := web.AuthChallengeRequest{
		Username: c.Username,
	}
	var authChallenge web.AuthChallengeResponse
	err = rv.post(ctx, &authChallenge, "/api/session/auth", nil, auth)
	if err != nil {
		return nil, errors.Wrap(err, "error authenticating")
	}

	// TODO: perform authentication

	authResponse := web.AuthResponseRequest{
		Username: c.Username,
		Nonce:    authChallenge.Nonce,
		Response: "",
	}

	var authResult web.AuthResponseResponse
	err = rv.post(ctx, &authResult, "/api/session/auth/response", nil, authResponse)
	if err != nil {
		return nil, errors.Wrap(err, "error authenticating")
	}

	return rv, nil
}

func (c *HttpClient) get(ctx context.Context, rv interface{}, path string, params url.Values) error {
	return c.do(ctx, rv, "GET", path, params, nil)
}

func (c *HttpClient) post(ctx context.Context, rv interface{}, path string, params url.Values, contents interface{}) error {
	return c.do(ctx, rv, "POST", path, params, contents)
}

func (c *HttpClient) do(ctx context.Context, rv interface{}, method string, path string, params url.Values, contents interface{}) error {
	u, err := url.Parse(c.baseURI)
	if err != nil {
		return errors.Wrap(err, "invalid server URI")
	}

	u.Path = strings.TrimRight(u.Path, "/") + path

	if params != nil {
		query := u.Query()
		for k, vv := range params {
			for _, v := range vv {
				query.Add(k, v)
			}
		}
		u.RawQuery = query.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), nil)
	if err != nil {
		return errors.Wrap(err, "cannot initialise request")
	}

	if len(c.sessionToken) > 5 {
		req.Header.Set("Authorisation", "Bearer "+c.sessionToken)
	}

	var encodedBody []byte

	if contents != nil {
		if vals, ok := contents.(url.Values); ok {
			req.Header.Set("Content-Type", "application/x-form-urlencoded")
			encodedBody = []byte(vals.Encode())
		} else {
			encodedBody, err = json.Marshal(contents)
			if err != nil {
				return errors.Wrap(err, "error encoding request")
			}
			req.Header.Set("Content-Type", "application/json")
		}
		req.Body = io.NopCloser(bytes.NewReader(encodedBody))
	}

	t0 := time.Now()

	resp, err := c.client.Do(req)
	if err != nil {
		return errors.Wrap(err, "error performing request")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "error reading response")
	}

	dur := time.Since(t0)

	if c.requestLogger != nil {
		dumpBody := strings.ReplaceAll(string(body), "\n", "\n\t\t")
		// HACK: set the body again, this time for the curl dump
		req.Body = io.NopCloser(bytes.NewReader(encodedBody))
		curl, _ := http2curl.GetCurlCommand(req)
		c.requestLogger.Printf("HTTP client request (duration: %s)\n%s\n\t\t%s\n\n", dur, curl, dumpBody)
	}

	var ae apiError
	err = json.Unmarshal(body, &ae)
	if err == nil && ae.ErrorCode != 0 && ae.ErrorMessage != "" {
		return ae
	}

	err = json.Unmarshal(body, rv)
	if err != nil {
		return errors.Wrap(err, "error decoding response")
	}

	return nil
}

// Me returns the object that represents the player at the server's end
func (c *HttpClient) Me() (game.Player, error) {
	return game.Player{}, notimplemented.Error()
}

// AvailablePlayers returns the list of players available for a match
func (c *HttpClient) AvailablePlayers(context.Context) ([]game.Player, error) {
	return nil, notimplemented.Error()
}

// NewGame initialises a Game with the specified players
func (c *HttpClient) NewGame(ctx context.Context, players []game.Player) (client.GameSession, error) {
	newgame := web.NewGameRequest{
		RuleSet: "Boring2D",
	}
	for _, pl := range players {
		newgame.PlayerNames = append(newgame.PlayerNames, pl.Name)
	}
	var gameid web.NewGameResponse
	err := c.post(ctx, &gameid, "/api/game/new", nil, newgame)
	if err != nil {
		return nil, errors.Wrap(err, "error starting new game")
	}

	rv := &httpSession{
		Client: c,
		GameID: gameid.GameID,
		game:   &game.Game{},
	}

	// TODO: Query the server to see if any moves were already saved
	// TODO: Apply those locally

	return rv, nil
}

type httpSession struct {
	Client *HttpClient
	GameID string
	game   *game.Game
}

// Game returns the Game object of this session
func (s *httpSession) Game() *game.Game {
	return s.game
}

// SubmitMove submits a move by this player.
func (s *httpSession) SubmitMove(context.Context, chesseract.Move) error {
	return notimplemented.Error()
}

// NextMove waits until a move occurs, and returns it. This comprises moves
// made by all players, not just opponents. NextMove returns the move made,
// but is also assumed to have applied the move to the supplied Game.
func (s *httpSession) NextMove(context.Context) (chesseract.Move, error) {
	return chesseract.Move{}, notimplemented.Error()
}

// ProposeResult submits a possible final outcome for this game, which all
// opponents can evaluate and accept or reject. One can accept a proposed
// result by proposing the same result again.
// Proposing a nil or zero result is construed as rejecting a proposition.
func (s *httpSession) ProposeResult(context.Context, []float64) error {
	return notimplemented.Error()
}

// NextProposition waits until a result is proposed, and returns it.
func (s *httpSession) NextProposition(context.Context) ([]float64, error) {
	return nil, notimplemented.Error()
}

// GetResult retrieves the result for this game
func (s *httpSession) GetResult(context.Context) ([]float64, error) {
	return nil, notimplemented.Error()
}
