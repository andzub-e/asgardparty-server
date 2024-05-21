package tests

import (
	"bitbucket.org/electronicjaw/asgardparty-server/internal/game"
	"encoding/json"
	"reflect"
	"testing"
)

func TestUnmarshalState(t *testing.T) {
	var state = game.State{}

	var bonusSpin = game.BonusSpin{
		Amount:        10,
		FreeSpinsLeft: 20,
		Type:          "bonus",
	}

	state.Spins = append(state.Spins, game.Spin{BonusSpin: bonusSpin, Type: bonusSpin.Type})

	baseSpin := game.BaseSpin{
		Amount: 50,
		Type:   "base",
	}

	state.Spins = append(state.Spins, game.Spin{BaseSpin: baseSpin, Type: baseSpin.Type})

	marshal, err := json.Marshal(state)

	if err != nil {
		t.Fatal(err)
	}

	t.Logf("JSON from marshal: %s\n", marshal)

	var fromJSON = game.State{}

	err = json.Unmarshal(marshal, &fromJSON)
	if err != nil {
		t.Fatal(err)
	}

	if fromJSON.Spins[0].Type != game.SpinTypeBonus || fromJSON.Spins[1].Type != game.SpinTypeBase {
		t.Fail()
	}

	if !reflect.DeepEqual(fromJSON.Spins[0].BonusSpin, bonusSpin) {
		t.Fail()
	}

	if !reflect.DeepEqual(fromJSON.Spins[1].BaseSpin, baseSpin) {
		t.Fail()
	}
}
