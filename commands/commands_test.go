package commands

import (
	"context"
	"fmt"
	"testing"

	test "github.com/retro-framework/go-retro/framework/test_helper"
	"github.com/retro-framework/go-retro/framework/types"
)

type dummyCmd struct{}
type otheDummyCmd struct{ dummyCmd }

func (_ *dummyCmd) SetState(types.Aggregate) error { return nil }

func (_ *dummyCmd) Apply(context.Context, types.Aggregate, types.Depot) ([]types.Event, error) {
	return nil, nil
}

type dummyAggregate struct{}

func (_ dummyAggregate) ReactTo(types.Event) error { return nil }

func Test_Commands_Register_TwiceSameCmdRaisesError(t *testing.T) {
	assertErrEql := test.H(t).ErrEql
	err := Register(dummyAggregate{}, &dummyCmd{})
	assertErrEql(err, nil)
	err = Register(dummyAggregate{}, &dummyCmd{})
	assertErrEql(err, fmt.Errorf("can't register command *commands.dummyCmd for aggregate commands.dummyAggregate, command already registered"))
}
