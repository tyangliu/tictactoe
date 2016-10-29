// Tictactoe game state.
import (
  fmt
)

// Board size - change this to change the size of the game board.
const boardSize = 3

/**
 * Represents a piece on a game board.
 * O - player 1 piece
 * X - player 2 piece
 * B - blank placeholder piece
 */
type Piece int
const Piece {
  O = iota
  X = iota
  B = iota
}

type Board [boardSize][boardSize]Piece

/**
 * Counts of player pieces in each row, column, and diagonal boardSize 
 * length line. If a player ever contains boardSize number of pieces 
 * in a single line, the player wins the game.
 *
 * Each board has boardSize rows and columns, and only 2 diagonals.
 */
type PlayerCounts struct {
  rows [boardSize]int
  cols [boardSize]int
  diags [2]int
}

/**
 * The result of a game move, one of:
 * - OWin - Player with piece O has won the game.
 * - XWin - Player with piece X has won the game.
 * - Tie  - Board is filled, no winner.
 * - Pending - Board is not full and no winner, keep playing.
 */
type GameResult int
const GameResult {
  OWin = iota
  XWin = iota
  Tie = iota
  Pending = iota
}

type GameState struct {
  // The boardSize * boardSize game board, each cell containing a piece 
  // (O, X, or B for blank).
  board *Board
  // The player who must make the next move, identified by their game piece
  // (O or X).
  currPiece Piece
  currPlayer string
  nextPlayer string
  // Counts of number of pieces player O has in rows, cols, and diags.
  oCounts PlayerCounts
  // Counts of number of pieces player X has in rows, cols, and diags.
  xCounts PlayerCounts
  totalPieces int
}

/**
 * Map of currently ongoing games, keyed by 'userA$$userB', where userA is 
 * lexicographically smaller than userB.
 */
currentGames map[string]*GameState

/**
 * Gets the key for the user pair, where the key is one of:
 * - "userA$$userB" if userA <= userB
 * - "userB$$userA" if userA < userB.
 *
 * This ensures that we never have two concurrent games between 
 * the same pair of users.
 */
func getUserPairKey(userA string, userB string) string {
  if userA <= userB {
    return userA + "$$" + userB
  }
  return userB + "$$" + userA
}

func initBoard(board *Board) {
  // Fill the board with blanks.
  for i := 0; i < boardSize; i++ {
    for j := 0; j < boardSize; j++ {
      board[i][j] = B
    }
  }
}

// Creates a new game between userA and userB. Overrides the previous game 
// if one already exists.
func startGame(userA string, userB string) *GameState {
  var board Board
  // Initialize board by filling with blanks.
  initBoard(&board)

  game := &GameState{board: &board, currPiece: O, currPlayer: userA}
  key := getUserPairKey(userA, userB)
  currentGames[key] = game
  return game
}

func clearGame(userA string, userB string) err {
  key := getUserPairKey(userA, userB)
  delete(currentGames, key)
  return nil
}

func getDiag(x int, y int) int {
  last := boardSize - 1
  // Top left to bottom right diagonal.
  if x == 0 && y == 0 || x == last && y == last {
    return 0
  }
  // Top right to bottom left diagonal.
  if x == last && y == 0 || x == 0 && y == last {
    return 1
  }
  // Not a diagonal.
  return -1
}

/**
 * Checks if the game is over. A game is over if either the 
 * current player has won (boardSize number of pieces in either 
 * the current row, column, or diagonal), or the board is full.
 */
func checkGameOver(game *GameState, x int, y int) GameResult {
  if game.currentPiece == O {
    diag := getDiag(x, y)
    diagWin := diag >= 0 && game.oCounts.diags[diag] == boardSize
    rowWin := game.oCounts.rows[x] == boardSize
    colWin := game.oCounts.cols[y] == boardSize

    if diagWin || rowWin || colWin {
      return OWin
    }
  } else {
    diag := getDiag(x, y)
    diagWin := diag >= 0 && game.xCounts.diags[diag] == boardSize
    rowWin := game.xCounts.rows[x] == boardSize
    colWin := game.xCounts.cols[x] == boardSize

    if diagWin || rowWin || colWin {
      return XWin
    }
  }

  // Every position is filled, but we don't have a winner, so game is a tie.
  if game.totalCount == boardSize * boardSize {
    return Tie
  }

  return Pending
}

/**
 * Makes a move by placing a piece on position (x,y) on the board if valid.
 * Returns the game result - either pending (game is not over), O or X has won, 
 * or the game is a tie.
 */
func makeMove(game *GameState, user string, x int, y int) (err, GameResult) {
  board := game.board

  if user != game.currentPlayer {
    return fmt.Errorf("It's not player %s's turn", user), Pending
  }

  if x < 0 || x >= boardSize || y < 0 || y >= boardSize {
    return fmt.Errorf("Board position %d %d is out of range.", x, y), Pending
  }

  if *board[x][y] != B {
    return fmt.Errorf("Board position %d %d is not empty.", x, y), Pending
  }

  *board[x][y] = game.currentPiece
  game.totalPieces++

  if game.currentPiece == O {
    game.oCounts.rows[x]++
    game.oCounts.cols[y]++
    diag := getDiag(x, y)
    if diag >= 0 {
      game.oCounts.diags[diag]++
    }
  } else {
    game.xCounts.rows[x]++
    game.xCounts.cols[y]++
    diag := getDiag(x, y)
    if diag >= 0 {
      game.xCounts.diags[diag]++
    }
  }

  // If game is over, we simply return the result (either a player has won 
  // or we have a tie).
  gameResult := checkGameOver(game, x, y)
  if gameResult != Pending {
    return nil, gameResult
  }

  // Change the current piece to the other one.
  if game.currentPiece == O {
    game.currentPiece = X
  } else {
    game.currentPiece = O
  }

  // Now it's nextPlayer's turn, so we swap currentPlayer and nextPlayer.
  game.currentPlayer = game.nextPlayer
  game.nextPlayer = user

  return nil, Pending
}

