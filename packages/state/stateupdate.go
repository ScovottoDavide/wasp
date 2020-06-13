package state

import (
	"fmt"
	"github.com/iotaledger/wasp/packages/sctransaction"
	"github.com/iotaledger/wasp/packages/util"
	"github.com/iotaledger/wasp/packages/variables"
	"io"
)

type stateUpdate struct {
	batchIndex uint16
	requestId  sctransaction.RequestId
	timestamp  int64
	vars       variables.Variables
}

func NewStateUpdate(reqid *sctransaction.RequestId) StateUpdate {
	var req sctransaction.RequestId
	if reqid != nil {
		req = *reqid
	}
	return &stateUpdate{
		requestId: req,
		vars:      variables.New(nil),
	}
}

func NewStateUpdateRead(r io.Reader) (StateUpdate, error) {
	ret := NewStateUpdate(nil).(*stateUpdate)
	return ret, ret.Read(r)
}

// StateUpdate

func (su *stateUpdate) String() string {
	return fmt.Sprintf("reqid: %s, ts: %d\n%s", su.requestId.String(), su.Timestamp(), su.Variables().String())
}

func (su *stateUpdate) Timestamp() int64 {
	return su.timestamp
}

func (su *stateUpdate) WithTimestamp(ts int64) StateUpdate {
	su.timestamp = ts
	return su
}

func (su *stateUpdate) RequestId() *sctransaction.RequestId {
	return &su.requestId
}

func (su *stateUpdate) Variables() variables.Variables {
	return su.vars
}

func (su *stateUpdate) Write(w io.Writer) error {
	if _, err := w.Write(su.requestId[:]); err != nil {
		return err
	}
	if err := su.vars.Write(w); err != nil {
		return err
	}
	return util.WriteUint64(w, uint64(su.timestamp))
}

func (su *stateUpdate) Read(r io.Reader) error {
	if _, err := r.Read(su.requestId[:]); err != nil {
		return err
	}
	if err := su.vars.Read(r); err != nil {
		return err
	}
	var ts uint64
	if err := util.ReadUint64(r, &ts); err != nil {
		return err
	}
	su.timestamp = int64(ts)
	return nil
}
