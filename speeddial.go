// Package speeddial reproduces the Master Lock Speed Dial (1500iD / ONE)
// combination mechanism: a hash from a knob sequence of any length onto one of
// the 7501 reachable states of its four disks. Mechanism from Michael Huebler,
// "The New Master Lock Speed Dial / ONE Combination Padlock: An Inside View",
// v2.0, 2009:
// https://toool.nl/images/e/e5/The_New_Master_Lock_Combination_Padlock_V2.0.pdf
package speeddial

// Move is a knob movement, and names the disk on that side of the dial. [Up],
// [Right], [Down] and [Left] are the only values; the zero value is [Up].
type Move struct {
	dial uint8 // disk index, clockwise from the top
}

// The four knob movements, clockwise.
var (
	Up    = Move{0}
	Right = Move{1}
	Down  = Move{2}
	Left  = Move{3}
)

func (m Move) String() string {
	switch m {
	case Right:
		return "right"
	case Down:
		return "down"
	case Left:
		return "left"
	default:
		return "up"
	}
}

const (
	disks          uint8 = 4
	positions            = 15 // disk stops, 24° apart
	degreesPerStep       = 360 / positions
)

// untouched is not an angle class; it marks the disk opposite the movement,
// which no pin reaches.
const untouched = 3

// Angle class each disk is driven to, by its clockwise offset from the
// movement: the facing disk lands on a multiple of 72°, its clockwise
// neighbour 24° past one, its counter-clockwise neighbour 24° short.
var classByOffset = [disks]uint8{0, 1, untouched, 2}

// State is the angle of the four combination disks, all the lock knows of the
// sequence so far. The zero value is a lock cleared by pushing the shackle in;
// equal states are sequences the lock cannot tell apart.
type State struct {
	pos [disks]uint8 // disk angle in steps of 24°, clockwise from the top
}

// Hash applies moves to a cleared lock.
func Hash(moves ...Move) State {
	var s State
	for _, m := range moves {
		s = s.Apply(m)
	}
	return s
}

// Apply moves the knob once, turning the three disks its pins reach.
func (s State) Apply(m Move) State {
	for i := range disks {
		if class := classByOffset[(i-m.dial)%disks]; class != untouched {
			s.pos[i] = advance(s.pos[i], class)
		}
	}
	return s
}

// Angles is the rotation of each combination disk, in degrees.
type Angles struct{ Top, Right, Bottom, Left int }

// Angles reports how far each disk has turned since the lock was cleared.
func (s State) Angles() Angles {
	deg := func(disk uint8) int { return int(s.pos[disk]) * degreesPerStep }
	return Angles{Top: deg(0), Right: deg(1), Bottom: deg(2), Left: deg(3)}
}

// advance turns the disk counter-clockwise to the next angle its driving pin
// can leave it at, so by 24°, 48° or 72°.
func advance(pos, class uint8) uint8 {
	step := (class+2-pos%3)%3 + 1
	return (pos + step) % positions
}
