# speeddial

Master Lock Speed Dial (1500iD) finite-state hash function as a little Go library.

The lock takes a knob sequence of any length and hashes it into the angles of
four internal disks, so it checks states, not sequences. There are 7501 of
them, all reachable in 11 moves or fewer, which is why many different sequences
open the same lock: one preset to right, down, down, left also opens on right,
down, up, down, left.

```go
preset := speeddial.Hash(speeddial.Right, speeddial.Down, speeddial.Down, speeddial.Left)
open := speeddial.Hash(attempt...) == preset
```

Mechanism from Michael Huebler, [*The New Master Lock Speed Dial / ONE
Combination Padlock: An Inside View*](https://toool.nl/images/e/e5/The_New_Master_Lock_Combination_Padlock_V2.0.pdf) (v2.0, 2009).
