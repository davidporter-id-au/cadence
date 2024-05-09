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

package queue

import (
	"testing"

	"github.com/golang/mock/gomock"
	"go.uber.org/goleak"

	"github.com/uber/cadence/common/persistence"
	"github.com/uber/cadence/common/reconciliation/invariant"
	"github.com/uber/cadence/service/history/config"
	"github.com/uber/cadence/service/history/execution"
	"github.com/uber/cadence/service/history/reset"
	"github.com/uber/cadence/service/history/shard"
	"github.com/uber/cadence/service/history/task"
	"github.com/uber/cadence/service/history/workflowcache"
	"github.com/uber/cadence/service/worker/archiver"
)

func TestTransferQueueProcessor_RequireStartStop(t *testing.T) {
	// some goroutine leak not from this test
	defer goleak.VerifyNone(t,
		// TODO(CDNC-8881):  TimerGate should not start background goroutine in constructor. Make it start/stoppable
		goleak.IgnoreTopFunction("github.com/uber/cadence/service/history/queue.NewLocalTimerGate.func1"),
	)
	ctrl := gomock.NewController(t)
	mockShard := shard.NewTestContext(
		t, ctrl, &persistence.ShardInfo{
			ShardID:          10,
			RangeID:          1,
			TransferAckLevel: 0,
		},
		config.NewForTest())
	defer mockShard.Finish(t)

	mockProcessor := task.NewMockProcessor(ctrl)
	mockResetter := reset.NewMockWorkflowResetter(ctrl)
	mockArchiver := &archiver.ClientMock{}
	mockInvariant := invariant.NewMockInvariant(ctrl)
	mockWorkflowCache := workflowcache.NewMockWFCache(ctrl)
	ratelimit := func(domain string) bool { return false }

	// Create a new transferQueueProcessor
	processor := NewTransferQueueProcessor(
		mockShard,
		mockShard.GetEngine(),
		mockProcessor,
		execution.NewCache(mockShard),
		mockResetter,
		mockArchiver,
		mockInvariant,
		mockWorkflowCache,
		ratelimit)
	processor.Start()
	processor.Stop()
}
