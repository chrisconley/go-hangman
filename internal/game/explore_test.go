package game

import (
	"testing"

	"github.com/segmentio/fasthash/fnv1a"

	"github.com/chrisconley/go-hangman/test"
)

func TestBuildPredictions(t *testing.T) {
	lines := []string{"abc", "aab", "bbb"}
	BUniverse = InitUniverse(lines)
	pw := BUniverse.PossibleWorlds
	actions := []Action{
		NewAction("a"),
		NewAction("b"),
		NewAction("c"),
	}

	actual := ExploreReactions(InitGamePlay(0.0), pw, actions, BuildTestAct(1.0))

	test.AssertAlmostEqual(actual.GetEntropy("a"), 0.9183, 4, t)
	test.AssertAlmostEqual(actual.GetEntropy("b"), 0.0000, 4, t)
	test.AssertAlmostEqual(actual.GetEntropy("c"), 0.9183, 4, t)

	test.AssertAlmostEqual(actual.GetReward("a"), 0.6667, 4, t)
	test.AssertAlmostEqual(actual.GetReward("b"), 1.0000, 4, t)
	test.AssertAlmostEqual(actual.GetReward("c"), 0.3333, 4, t)

	test.AssertAlmostEqual(actual.GetInversedMax("a"), 1.0000, 4, t)
	test.AssertAlmostEqual(actual.GetInversedMax("b"), 0.0000, 4, t)
	test.AssertAlmostEqual(actual.GetInversedMax("c"), 1.0000, 4, t)

	test.AssertAlmostEqual(actual.GetPartsCount("a"), 2.0000, 4, t)
	test.AssertAlmostEqual(actual.GetPartsCount("b"), 1.0000, 4, t)
	test.AssertAlmostEqual(actual.GetPartsCount("c"), 2.0000, 4, t)
}

func TestExploreReactionsForAction(t *testing.T) {
	lines := []string{"abc", "aab", "cab"}
	universeLength := len(lines)
	BUniverse = InitUniverse(lines)
	pw := BUniverse.PossibleWorlds

	_, counts, rewards, responseWorlds := ExploreReactionsForAction(
		InitGamePlay(0.0),
		pw.Indices(),
		NewAction("c"),
		BuildTestAct(0.222),
	)

	test.AssertAlmostEqual(counts.Get("!"), 1.0, 3, t)
	test.AssertAlmostEqual(counts.Get("c"), 2.0, 3, t)

	test.AssertAlmostEqual(rewards.Get("!"), 0.0, 3, t)
	test.AssertAlmostEqual(rewards.Get("c"), 0.444, 3, t)

	expected := NewPossibleWorlds(universeLength)
	expected.Add(1)
	responseHash := fnv1a.HashString64("!")
	test.AssertDeepEqual(responseWorlds[responseHash], expected, t)

	expected = NewPossibleWorlds(universeLength)
	expected.Add(0)
	expected.Add(2)
	responseHash = fnv1a.HashString64("c")
	test.AssertDeepEqual(responseWorlds[responseHash], expected, t)
}

func TestPredictionMarshalling(t *testing.T) {
	p1 := Prediction{
		Action:                      NewAction("a"),
		Reward:                      0.12,
		Entropy:                     0.345,
		InversedMax:                 0.56,
		PartsCount:                  12.0,
		ActionMatchesAPossibleWorld: true,
	}

	p1Bytes := p1.ToBytes()

	p2 := NewPredictionFromBytes(p1Bytes)

	test.AssertDeepEqual(p1.Action, p2.Action, t)
	test.AssertAlmostEqual(p1.Reward, p2.Reward, 4, t)
	test.AssertAlmostEqual(p1.Entropy, p2.Entropy, 4, t)
	test.AssertAlmostEqual(p1.InversedMax, p2.InversedMax, 4, t)
	test.AssertAlmostEqual(p1.PartsCount, p2.PartsCount, 4, t)
	test.AssertDeepEqual(
		p1.ActionMatchesAPossibleWorld,
		p2.ActionMatchesAPossibleWorld,
		t,
	)

}

func TestPredictionsMarshalling(t *testing.T) {
	prediction1 := Prediction{
		Action: NewAction("a"),
		Reward: 0.12,
	}
	prediction2 := Prediction{
		Action: NewAction("b"),
		Reward: 0.34,
	}

	p1 := NewPredictions()
	p1.AddPrediction("a", prediction1)
	p1.AddPrediction("b", prediction2)

	p1Bytes := p1.ToBytes()

	p2 := NewPredictionsFromBytes(p1Bytes)

	test.AssertDeepEqual(p2.Size(), 2, t)
	test.AssertDeepEqual(p2.Get("a").Action, NewAction("a"), t)
	test.AssertDeepEqual(p2.Get("b").Action, NewAction("b"), t)

}
