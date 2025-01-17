package smGPA

import (
	"github.com/iotaledger/wasp/packages/state"
)

type blockRequestCallback interface {
	isValid() bool
	requestCompleted()
}

type blockFetcher interface {
	getCommitment() *state.L1Commitment
	notifyFetched(func(blockFetcher) (bool, error)) error // calls fun for this fetcher and each related recursively; fun for parent block is always called before fun for related block
	addCallback(blockRequestCallback)
	addRelatedFetcher(blockFetcher)
	cleanCallbacks()
}

type blockFetchers interface {
	getSize() int
	addFetcher(blockFetcher)
	takeFetcher(*state.L1Commitment) blockFetcher
	addCallback(*state.L1Commitment, blockRequestCallback) bool
	addRelatedFetcher(*state.L1Commitment, blockFetcher) bool
	getCommitments() []*state.L1Commitment
	cleanCallbacks()
}
