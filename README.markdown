# go-hangman

go-hangman is a Go engine that plays guessing games like hangman and Mastermind. On each turn it looks at every guess it could make and asks how that guess would split the answers that are still possible. It scores each guess on four measures:

- **Reward**: the expected reward of the guess. In hangman that means the chance the letter is in the word; in Mastermind it's the expected black/white score.
- **Entropy**: the information the guess gives, in bits.
- **MiniMax**: how small the largest group of remaining answers would be (Knuth's 1976 rule).
- **MostParts**: how many different responses the guess can produce.

You decide how much each measure counts, from one measure alone to any mix of the four. The engine then plays a whole dictionary that way and reports how many guesses and how many misses it took, so strategies can be compared on the same games. It can also try every mix at a chosen step size and report which one needs the fewest guesses or the fewest misses.

## Build and test

Tested with Go 1.27.2 on Linux. go.mod declares `go 1.15`.

```
go build ./...
make test
```

## Play

`cmd/player` plays one game against every answer in a dictionary and prints the totals.

```
cat fixtures/nine-letter-words.txt | go run ./cmd/player -game hangman -motivation.entropy 1.0
```

The parts of that command:

- `cat fixtures/nine-letter-words.txt |` is the dictionary, one answer per line on stdin. Every line is a possible answer, and the engine plays one game against each line.
- `-game hangman` picks the rules: `hangman` (guess a letter, see where it appears) or `mastermind` (guess a code, get black and white pegs). Required.
- `-motivation.entropy 1.0` sets the weights. There are four flags, `-motivation.reward`, `-motivation.entropy`, `-motivation.minimax` and `-motivation.mostparts`. Each defaults to 0, and they must sum to exactly 1.0.
- `-limit N` (optional) stops after the first N lines of the file, so you can try a run without playing all 25,011 games. The engine still treats every line in the file as a word the answer could be.

It runs for about 6 seconds and prints one line of JSON. It is laid out here for reading, with the averages rounded and the per-game list cut short:

```
GameStats: {
  "GameCount": 25011,
  "TotalTurns": 146356,
  "AverageTurns": 5.85,
  "MaxTurns": 13,
  "TotalReward": 103686,
  "AverageReward": 4.15,
  "TotalPain": 42670,
  "AveragePain": 1.71,
  "TurnCounts": [3, 5, 6, 7, 7, ...]
}
```

- `GameCount`: games played, one per answer.
- `TotalTurns`, `AverageTurns`, `MaxTurns`: guesses over all games, per game, and in the longest game.
- `TotalReward` and `TotalPain`: in hangman, the guessed letters that were in the word and the ones that were not, so the two add up to `TotalTurns`.
- `TurnCounts`: the number of guesses for each answer, in input order.

### Games and dictionaries

| Dictionary | Use with | Answers | What it is |
|---|---|---|---|
| `fixtures/nine-letter-words.txt` | `-game hangman` | 25,011 | nine-letter English words |
| `fixtures/commercial-mastermind.txt` | `-game mastermind` | 1,296 | standard Mastermind: 4 pegs, 6 colours (A–F) |
| `fixtures/mastermind-ABCD-3.txt` | `-game mastermind` | 64 | 3 pegs, 4 colours |
| `fixtures/mastermind-ABCDEFG-4.txt` | `-game mastermind` | 2,401 | 4 pegs, 7 colours |
| `fixtures/mastermind-ABCDEFGH-5.txt` | `-game mastermind` | 32,768 | 5 pegs, 8 colours |

You can also pipe in your own list, one answer per line.

### More examples

A blend of three measures, the best one found for nine-letter hangman (146,258 turns):

```
cat fixtures/nine-letter-words.txt | go run ./cmd/player -game hangman \
    -motivation.entropy 0.7 -motivation.minimax 0.15 -motivation.mostparts 0.15
```

Standard Mastermind with MostParts alone (1,296 games, 5,668 turns, at most 6 per game):

```
cat fixtures/commercial-mastermind.txt | go run ./cmd/player -game mastermind -motivation.mostparts 1.0
```

The smallest game, for a quick try (64 games, 208 turns, at most 4 per game):

```
cat fixtures/mastermind-ABCD-3.txt | go run ./cmd/player -game mastermind -motivation.entropy 1.0
```

## Search the weights

`cmd/optimize` tries every combination of the four weights at a given step size. For each combination it plays all the games and prints one line of JSON. It ends with the best combination for total turns and for total pain.

```
cat fixtures/nine-letter-words.txt | go run ./cmd/optimize -game hangman -trial.stepSize 0.5
```

A step of 0.5 gives 10 combinations and takes about 40 seconds on 4 CPUs. The results below used a step of 0.05, which gives 1,771 combinations and takes roughly 177 times as long.

## Results in this repository

"Turns" means the total number of guesses over all games.

**Hangman, 25,011 nine-letter words** (`hangman-results-oct26.csv`, step 0.05)

| Weights | Total turns | Max turns |
|---|---|---|
| Entropy 0.70, MiniMax 0.15, MostParts 0.15 (best) | 146,258 | 13 |
| Entropy alone | 146,356 | 13 |
| MostParts alone | 148,401 | 15 |
| MiniMax alone | 149,500 | 12 |
| Reward alone | 170,676 | 14 |

The fewest misses (total pain 34,695) is a tie between Reward 0.75 / Entropy 0.25 (156,126 turns) and Reward 0.70 / Entropy 0.30 (155,451 turns).

The CSV has one row per weight combination. The full output of that run also listed every game's turn count; it isn't in the repository, and `cmd/optimize` regenerates it with `-trial.stepSize 0.05`.

**Mastermind, 1,296 codes (4 pegs, 6 colours)** (`mastermind-results.json`, step 0.05)

| Weights | Total turns | Max turns |
|---|---|---|
| MostParts alone (best) | 5,668 | 6 |
| Entropy alone | 5,723 | 6 |
| MiniMax alone | 5,801 | 5 |
| Best with at most 5 turns: Entropy 0.05, MiniMax 0.70, MostParts 0.25 | 5,677 | 5 |

## Optimizations

The slow part is scoring every guess against every remaining answer on every turn. The engine cuts that work down in four ways.

- **Remaining answers as a bitmap.** The set of answers still possible is a compressed bitmap (`roaring64`) with a 64-bit hash, so a set is cheap to copy, store and compare.
- **Two caches shared by all games.** Many games pass through the same situation. The engine keeps the scores of every guess for a situation (the prediction cache) and the answers left after a guess and its response (the observe cache), each under a hash. In the nine-letter run above, 146,356 turns needed 30,139 scoring passes; the rest were cache hits. Both caches use `freecache`, which keeps them out of the garbage collector's way. `-cacheSize` sets the prediction cache in MB (default 1000).
- **Fewer Mastermind guesses to score.** Before the first guess, AABB, CCDD and EEFF are the same guess with the colours renamed, so only one of them needs scoring. The engine uses that idea on every turn for the colours that haven't been played yet and scores one guess per pattern: 5 opening guesses in place of 1,296. This is a shortcut, not an exact reduction. From the third guess on it can skip a guess that isn't truly equivalent to the one it keeps, so a total can differ slightly from scoring every code. With MiniMax alone, the 1,296 games of standard Mastermind take 5,800 guesses with the shortcut and 5,801 without it; the difference is tie-order noise. The code is `GetActions` in `internal/mastermind/agreement.go`.
- **Weight combinations in parallel.** `cmd/optimize` runs one combination per CPU at a time.

`-profiling cpu` or `-profiling mem` writes a profile to `/tmp/profiling`, or to the path given with `-profiling.profilePath`.

`symmetries.markdown` (2022) tests a tempting way to prune further: call two guesses equivalent whenever they split the remaining codes into groups of the same sizes. Such guesses get the same score on all four measures, so it looks safe to keep just one. After an opening ABCD, the second guesses fall into 57 families that really are equivalent. The size rule finds exactly those 57 when the first response is two pegs. When it is one peg, it finds only 53: four pairs of different families happen to give the same group sizes. With three pegs, two pairs do. Two guesses with the same group sizes are equally good for this turn but can leave different positions behind, so the rule would change how games play out. The note lists every case, and the rule isn't used in play.

## Compared with published results

The single-measure rules are the ones studied in the Mastermind literature, and their totals here agree with the published ones to within one guess.

| Game | Rule | Published | This repository |
|---|---|---|---|
| 4 pegs, 6 colours | Most parts | 5,668 (Kooi 2005, Ville 2013) | 5,668 |
| 4 pegs, 6 colours | Entropy | 5,723 (Ville 2013); 5,722 (Kooi 2005) | 5,723 |
| 4 pegs, 6 colours | Largest group, Knuth's rule (MiniMax here) | 5,801 (Ville 2013) | 5,801 in `mastermind-results.json`; 5,800 with the guess pruning above |
| 4 pegs, 7 colours | Entropy | 11,382 (Ville 2013) | 11,381 |
| 4 pegs, 7 colours | Most parts | 11,388 (Ville 2013) | 11,388 |
| 4 pegs, 7 colours | Largest group | 11,613 (Ville 2013) | 11,613 |

The 7-colour figures are from `mastermind4x7-noreward-results.json`. The one-guess differences come from which guess is played when several score the same: Kooi's own tables give 5,722 for Entropy with one tie order and 5,723 with the reverse.

These rules look one guess ahead. The published optimal strategy for 4 pegs and 6 colours needs 5,625 guesses (Koyama and Lai 1993). This engine does not search for optimal strategies.

What this engine adds:

- **Blends.** Any weighted mix of the four measures can be played, and `cmd/optimize` searches the mixes. On standard Mastermind, Entropy 0.05 / MiniMax 0.70 / MostParts 0.25 needs 5,677 guesses and never more than 5 in a game. That is 124 fewer than Knuth's rule, with the same 5-guess limit. With 7 colours, Entropy 0.85 / MostParts 0.15 needs 11,366, fewer than any single measure.
- **Reward and pain.** A fourth measure, the expected immediate reward of a guess, and a second total, pain, so a strategy can be judged on misses as well as on guesses.
- **Hangman.** The same engine and the same four measures play a letter-guessing game. Nothing in the scoring is specific to Mastermind.

Full citations are in [REFERENCES.md](REFERENCES.md).

## License

MIT. See [LICENSE](LICENSE).
