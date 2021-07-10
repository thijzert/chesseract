Use cases
=========
![Use case overview](../.resources/doc_usecases-mermaid.svg)



Use case: play match
--------------------
**Primary actor**: player  
**Stakeholders and interests**
- Player: wants to pass the time and get better at 4D chess.

**Preconditions**: player is authenticated at a server, and has set up a new match with at least one opponent.  
**Postconditions**: none  

### Main success scenario
1. The player moves one piece on the board
2. Their opponent moves one piece on the board
3. Steps 1-2 are repeated until the board is in a state of checkmate.
4. The system records the result of the match, which is reflected in each player's stats.

### Extensions
#### Player forfeits
3a. Steps 1-2 are repeated until the board is in a state of checkmate, or the player forfeits the match.

#### Player and opponent agree to a draw
3b. Steps 1-2 are repeated until the board is in a state of checkmate, or all players agree to a draw.


Use case: set up match
----------------------
**Primary actor**: player  
**Stakeholders and interests**:
- Player: wants to play a match of 4D chess

**Preconditions**: player is authenticated at a server  
**Postconditions**: player has set up a new match with at least one opponent.

### Main success scenario
1. The player enters the server lobby.
2. The system marks that player as available to play a match, and displays a list of other available players.
3. The player picks a suitable opponent from the list.
4. The system messages that opponent, asking if they're willing to play a match against the player
5. If the opponent accepts, a new match between the two is started.

### Extensions
#### Opponent finds player
3a. If an opponent selects the player from their lobby list before the player does so:
1. The system displays a message indicating the opponent wishes to play a match against the player, and asks the player to accept or reject this offer.
2. If the player accepts, a new match between the two is started.

