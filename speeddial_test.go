package speeddial

import (
	"fmt"
	"testing"
)

var moves = [disks]Move{Up, Right, Down, Left}

// Huebler's transition table for the upper disk, page 7.
func TestUpperDiskTable(t *testing.T) {
	table := []struct{ start, right, up, left int }{
		{336, 48, 0, 24},
		{0, 48, 72, 24},
		{24, 48, 72, 96},
		{48, 120, 72, 96},
		{72, 120, 144, 96},
		{96, 120, 144, 168},
		{120, 192, 144, 168},
		{144, 192, 216, 168},
		{168, 192, 216, 240},
		{192, 264, 216, 240},
		{216, 264, 288, 240},
		{240, 264, 288, 312},
		{264, 336, 288, 312},
		{288, 336, 0, 312},
		{312, 336, 0, 24},
	}
	for _, row := range table {
		var start State
		start.pos[Up.dial] = uint8(row.start / degreesPerStep)
		for _, c := range []struct {
			knob Move
			want int
		}{{Right, row.right}, {Up, row.up}, {Left, row.left}, {Down, row.start}} {
			if got := start.Apply(c.knob).Angles().Top; got != c.want {
				t.Errorf("upper disk at %d°, knob %v: got %d°, want %d°", row.start, c.knob, got, c.want)
			}
		}
	}
}

// Reachable states per sequence length, and their total, page 8.
func TestReachableStates(t *testing.T) {
	want := []int{1, 4, 16, 60, 168, 396, 816, 1448, 1984, 1796, 708, 104}
	seen := map[State]bool{{}: true}
	layer := []State{{}}
	for length, wantCount := range want {
		if len(layer) != wantCount {
			t.Errorf("%d moves: %d states first reached, want %d", length, len(layer), wantCount)
		}
		var next []State
		for _, s := range layer {
			for _, m := range moves {
				if n := s.Apply(m); !seen[n] {
					seen[n] = true
					next = append(next, n)
				}
			}
		}
		layer = next
	}
	if len(layer) != 0 {
		t.Errorf("%d states still unreached after %d moves", len(layer), len(want))
	}
	if len(seen) != 7501 {
		t.Errorf("got %d reachable states, want 7501", len(seen))
	}
}

// Page 7: an up from cleared leaves the upper disk on a multiple of 72°, the
// right disk 24° past one, the left disk 24° short of one, and the lower alone.
func TestAngles(t *testing.T) {
	var cleared State
	if got := cleared.Angles(); got != (Angles{}) {
		t.Errorf("cleared lock: got %+v, want every disk at 0°", got)
	}
	want := Angles{Top: 72, Right: 24, Bottom: 0, Left: 48}
	if got := Hash(Up).Angles(); got != want {
		t.Errorf("after up: got %+v, want %+v", got, want)
	}
}

// A preset combination, and what else opens on it.
func TestCollisions(t *testing.T) {
	preset := Hash(Right, Down, Down, Left)
	if preset == Hash(Left, Down, Down, Right) {
		t.Error("the preset combination opens reversed")
	}
	if preset != Hash(Right, Down, Up, Down, Left) {
		t.Error("[right down up down left] does not open the preset combination")
	}
}

func ExampleHash() {
	preset := Hash(Right, Down, Down, Left)
	fmt.Println(preset == Hash(Right, Down, Up, Down, Left)) // an extra up changes nothing
	fmt.Println(preset == Hash(Left, Down, Down, Right))     // reversed, does not open
	// Output:
	// true
	// false
}
