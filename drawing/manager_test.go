package drawing

import "testing"

func TestManagerDrawEndCallsInOrder(t *testing.T) {
	depths := []int{1, 2, 3}
	m := NewManager(depths, 8)

	target := &struct{}{}
	var calls []int

	for _, depth := range depths {
		depth := depth
		m.Draw(func(s Surface) {
			if s != target {
				t.Fatalf("unexpected surface: %#v", s)
			}
			calls = append(calls, depth)
		}, depth)
	}

	m.DrawEnd(target)

	if len(calls) != len(depths) {
		t.Fatalf("expected %d callbacks, got %d", len(depths), len(calls))
	}

	for i, depth := range depths {
		if calls[i] != depth {
			t.Fatalf("unexpected depth order at index %d: want %d, got %d", i, depth, calls[i])
		}
	}

	// ensure reset
	before := len(calls)
	m.DrawEnd(target)
	if len(calls) != before {
		t.Fatalf("callbacks should not be invoked when queue is empty; got %d, want %d", len(calls), before)
	}
}

func TestManagerResetsQueue(t *testing.T) {
	depths := []int{10}
	m := NewManager(depths, 4)
	target := &struct{}{}
	var calls int

	m.Draw(func(s Surface) {
		calls++
	}, depths[0])
	m.DrawEnd(target)

	if calls != 1 {
		t.Fatalf("expected first DrawEnd to invoke 1 callback, got %d", calls)
	}

	m.Draw(func(s Surface) {
		if s != target {
			t.Fatalf("unexpected surface: %#v", s)
		}
		calls++
	}, depths[0])
	m.DrawEnd(target)

	if calls != 2 {
		t.Fatalf("expected second DrawEnd to invoke one more callback, got %d", calls)
	}
}

func TestManagerDrawPanicsOnOverflow(t *testing.T) {
	depths := []int{42}
	m := NewManager(depths, 2)

	for i := 0; i < 2; i++ {
		m.Draw(func(Surface) {}, depths[0])
	}

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic when exceeding draw queue capacity, got none")
		}
	}()

	m.Draw(func(Surface) {}, depths[0])
}

func TestManagerPanicsOnUnknownDepth(t *testing.T) {
	m := NewManager([]int{}, 4)

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for unregistered depth, got none")
		}
	}()

	m.Draw(func(Surface) {}, 99)
}
