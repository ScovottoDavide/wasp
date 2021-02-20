package sbtests

import (
	"github.com/iotaledger/wasp/packages/kv/kvdecoder"
	"github.com/iotaledger/wasp/packages/solo"
	"github.com/iotaledger/wasp/packages/vm/core/testcore/sbtests/sbtestsc"
	"github.com/stretchr/testify/require"
	"runtime"
	"testing"
	"time"
)

func TestCounter(t *testing.T) { run2(t, testCounter) }
func testCounter(t *testing.T, w bool) {
	_, chain := setupChain(t, nil)
	setupTestSandboxSC(t, chain, nil, w)

	req := solo.NewCallParams(SandboxSCName, sbtestsc.FuncIncCounter)
	for i := 0; i < 33; i++ {
		_, err := chain.PostRequestSync(req, nil)
		require.NoError(t, err)
	}

	ret, err := chain.CallView(SandboxSCName, sbtestsc.FuncGetCounter)
	require.NoError(t, err)

	deco := kvdecoder.New(ret, chain.Log)
	res := deco.MustGetInt64(sbtestsc.VarCounter)
	require.EqualValues(t, 33, res)
}

func TestConcurrency(t *testing.T) { run2(t, testConcurrency) }
func testConcurrency(t *testing.T, w bool) {
	//t.SkipNow()
	_, chain := setupChain(t, nil)
	setupTestSandboxSC(t, chain, nil, w)

	req := solo.NewCallParams(SandboxSCName, sbtestsc.FuncIncCounter)

	repeats := []int{100, 100, 100, 100, 100, 100, 100, 100, 100, 100}
	sum := 0
	for _, i := range repeats {
		sum += i
	}
	for r, n := range repeats {
		go func(r, n int) {
			for i := 0; i < n; i++ {
				tx := chain.RequestFromParamsToLedger(req, nil)
				chain.Env.EnqueueRequests(tx)
			}
			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			t.Logf("++++++++++++++ #%d -- alloc: %d MB, total: %d MB GC: %d",
				r, m.Alloc/(1024*1024), m.TotalAlloc/(1024*1024), m.NumGC)
		}(r, n)
	}
	time.Sleep(1 * time.Second)
	chain.WaitForEmptyBacklog(10 * time.Second)

	ret, err := chain.CallView(SandboxSCName, sbtestsc.FuncGetCounter)
	require.NoError(t, err)

	deco := kvdecoder.New(ret, chain.Log)
	res := deco.MustGetInt64(sbtestsc.VarCounter)
	require.EqualValues(t, sum, res)
}
