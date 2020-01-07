package orbiter_test

import (
	"testing"

	"github.com/caos/orbiter/internal/core/operator/test"
)

func TestOrbEmptyNodeAgent(t *testing.T) {

	t.SkipNow()

	desire, err := test.DesireDefaults(t, "v0.0.0")
	if err != nil {
		t.Error(err)
	}

	desire = desire.Chain(func(des map[string]interface{}) {
		des["nodeagent"] = make([]byte, 0)
	})

	stop := make(chan struct{})

	iterations, cleanup, _, err := test.Run(stop, "test empty nodeagent", t, `kind: orbiter.caos.ch/Orb
version: v1
`, desire)
	defer cleanup()

	if err != nil {
		return
	}

	if it := <-iterations; it.Error == nil {
		t.Fatalf("running with zero sized nodeagent did not cause an error")
	}
}