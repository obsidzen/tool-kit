package tuikit

import (
	"fmt"
	"testing"
)

func tailWith(lines int, height int) TailModel {
	tail := TailModel{Height: height}
	for i := 0; i < lines; i++ {
		tail = tail.Append(fmt.Sprintf("line %d", i))
	}
	return tail
}

func TestTailShowsTheEndWhileFollowing(t *testing.T) {
	tail := tailWith(100, 10)

	visible, below := tail.Window()

	if !tail.AtBottom() {
		t.Fatalf("a fresh tail should follow the output")
	}
	if below != 0 {
		t.Fatalf("below = %d, want 0 while following", below)
	}
	if visible[len(visible)-1] != "line 99" {
		t.Fatalf("last visible = %q, want the newest line", visible[len(visible)-1])
	}
}

func TestScrollingUpHoldsPositionWhileOutputContinues(t *testing.T) {
	tail := tailWith(100, 10).ScrollUp(20)
	before, _ := tail.Window()

	for i := 0; i < 30; i++ {
		tail = tail.Append("more")
	}
	after, below := tail.Window()

	if after[0] != before[0] {
		t.Fatalf("view moved: %q -> %q; scrolled-up output must stay put", before[0], after[0])
	}
	if below == 0 {
		t.Fatalf("below = 0, want the newly hidden lines to be reported")
	}
}

func TestScrollingBackToTheBottomFollowsAgain(t *testing.T) {
	tail := tailWith(100, 10).ScrollUp(20).ScrollBottom()

	visible, below := tail.Window()

	if !tail.AtBottom() || below != 0 {
		t.Fatalf("AtBottom = %v, below = %d; want to be following again", tail.AtBottom(), below)
	}
	if visible[len(visible)-1] != "line 99" {
		t.Fatalf("last visible = %q, want the newest line", visible[len(visible)-1])
	}
}

func TestScrollingStopsAtTheOldestLine(t *testing.T) {
	tail := tailWith(30, 10).ScrollUp(1000)

	visible, _ := tail.Window()

	if visible[0] != "line 0" {
		t.Fatalf("first visible = %q, want the oldest line", visible[0])
	}
}

func TestDroppingOldLinesKeepsTheViewOnTheSameContent(t *testing.T) {
	tail := TailModel{Height: 10, Max: 50}
	for i := 0; i < 50; i++ {
		tail = tail.Append(fmt.Sprintf("line %d", i))
	}
	tail = tail.ScrollUp(20)
	before, _ := tail.Window()

	// Past Max every append drops one from the front, which would slide the
	// window across the log if the offset were left alone.
	for i := 50; i < 70; i++ {
		tail = tail.Append(fmt.Sprintf("line %d", i))
	}
	after, _ := tail.Window()

	if after[0] != before[0] {
		t.Fatalf("view drifted: %q -> %q", before[0], after[0])
	}
}

func TestScrollKeysAreConsumedInTailMode(t *testing.T) {
	tail := tailWith(100, 10)

	for _, key := range []string{"up", "k", "down", "j", "pgup", "pgdown", "home", "end"} {
		if _, ok := scrollTail(tail, key); !ok {
			t.Fatalf("key %q was not consumed; it would close the log instead", key)
		}
	}
	if _, ok := scrollTail(tail, "x"); ok {
		t.Fatalf("an unrelated key must fall through to the return behaviour")
	}
}
