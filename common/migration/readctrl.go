// The MIT License (MIT)

// Copyright (c) 2017-2020 Uber Technologies Inc.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package migration

import (
	"time"

	"github.com/uber/cadence/common/dynamicconfig"
	"github.com/uber/cadence/common/log"
	"github.com/uber/cadence/common/metrics"
)

type readerImpl[T comparable] struct {
	log                log.Logger
	scope              metrics.Scope
	backgroundTimeout  time.Duration
	rolloutCtrl        dynamicconfig.StringPropertyFn
	doResultComparison dynamicconfig.BoolPropertyFn
	comparisonFn       ComparisonFn[T]
}

func NewDualReaderWithCustomComparisonFn[T comparable](
	rolloutCtrl dynamicconfig.StringPropertyFn,
	doResultComparison dynamicconfig.BoolPropertyFn,
	log log.Logger,
	scope metrics.Scope,
	backgroundTimeout time.Duration,
	comparisonFn ComparisonFn[T],
) Reader[T] {

	return &readerImpl[T]{
		log:                log,
		scope:              scope,
		rolloutCtrl:        rolloutCtrl,
		doResultComparison: doResultComparison,
		backgroundTimeout:  backgroundTimeout,
		comparisonFn:       comparisonFn,
	}
}

func (c *readerImpl[T]) getReaderRolloutState(constraints Constraints) ReaderRolloutState {
	s := c.rolloutCtrl(
		dynamicconfig.OperationFilter(constraints.Operation),
		dynamicconfig.DomainFilter(constraints.Domain),
	)
	return ReaderRolloutState(s)
}

func (c *readerImpl[T]) shouldCompare(constraints Constraints) bool {
	return c.doResultComparison(
		dynamicconfig.OperationFilter(constraints.Operation),
		dynamicconfig.DomainFilter(constraints.Domain),
	)
}
