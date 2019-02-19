package main

import (
	"strings"
)

var (
	examp  = "4.....8.5.3..........7......2.....6.....8.4......1.......6.3.7.5..2.....1.4......"
	rows   = []string{"A", "B", "C", "D", "E", "F", "G", "H", "I"}
	cols   = []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}
	digits = "123456789"

	squares  = cross(rows, cols)       // The Array of all cell's names (81)
	unitlist = [][]string{}            // The Array of Cols, Rows and Boxes (9 column units, 9 row units and 9 box units)
	units    = map[string][][]string{} // Units (Value) per Square (Key)

	peers = map[string][]string{} // Peers, dependent squares (Value) per Square (Key)

	nassigns      = 0
	neliminations = 0
	nsearches     = 0
)

func main() {

	//___________________________________ Creating Unitlist by concatenating all units

	for _, c := range cols {
		unitlist = append(unitlist, cross(rows, []string{c}))
	}
	for _, r := range rows {
		unitlist = append(unitlist, cross([]string{r}, cols))
	}

	rrows := [][]string{{"A", "B", "C"}, {"D", "E", "F"}, {"G", "H", "I"}}
	ccols := [][]string{{"1", "2", "3"}, {"4", "5", "6"}, {"7", "8", "9"}}
	for _, rs := range rrows {
		for _, cs := range ccols {
			unitlist = append(unitlist, cross(rs, cs))
		}
	}

	//___________________________________

	for _, s := range squares {
		units[s] = [][]string{}
		for _, u := range unitlist {
			if member(s, u) {
				units[s] = append(units[s], u)
			}
		}
	}

	//___________________________________

	for _, s := range squares {
		set := []string{}
		for _, u := range units[s] {
			for _, square := range u {
				if square != s {
					set = append(set, square)
				}
			}
		}
		peers[s] = set
	}

	//___________________________________

	//fmt.Println("All cell peers", peers["I9"])
	//fmt.Println("All cell units", units["C2"])
	//fmt.Printf("len=%d cap=%d %v\n", len(squares), cap(squares), squares)

}

//___________________________________ Service method for strings concatenating

func cross(A []string, B []string) []string {
	var C []string
	for _, a := range A {
		for _, b := range B {
			C = append(C, (a + b))
		}
	}
	return C
}

//___________________________________ Service method for checking is there a string in the list

func member(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

//___________________________________ Method for Storing an initial representation of Sudoku into a map using square as key

func mapValue(grid string) map[string]string {
	symbols := digits + ".0"
	gridChars := []string{}
	m := make(map[string]string)

	for _, char := range grid {
		for _, b := range symbols {
			if char == b {
				gridChars = append(gridChars, string(char))
			}
		}
	}

	for i := 0; i < len(gridChars); i++ {
		m[squares[i]] = gridChars[i]
	}
	return m
}

//___________________________________ Method leaves only possible values for each square

func parseGrid(grid string) map[string]string {
	initValues := mapValue(grid)
	values := make(map[string]string)

	for _, s := range squares {
		values[s] = digits // Initial value in a map where each square is used as a key
	}

	for _, s := range squares {
		if strings.Contains(digits, initValues[s]) {
			assign(values, s, initValues[s])
		}
	}
	return values
}

//___________________________________ It updates the incoming values by eliminating the other values than d

func assign(values map[string]string, s string, d string) bool {
	nassigns++

	otherValues := strings.Replace(s, d, "", 1)
	for _, d2 := range otherValues {
		if !eliminate(values, s, string(d2)) {
			return false
		}
	}
	return true
}

/*
1. It removes the given value d from values[s] which is a list of potential values for s.
2. If there is no values left in s (that is we don’t have any potential value for that square), returns False
3. When there is only one potential value for s, it removes the value from all the peers of s	<== strategy (1)
4. Make sure the given value d has a place elsewhere (i.e., if no square has d as a potential value, we can not solve the puzzle)
5. Where there is only one place for the value d, remove it from the peers	<== strategy (2)
*/

func eliminate(values map[string]string, s string, d string) bool {
	neliminations++

	if strings.Contains(values[s], d) {
		return true // Already eliminated
	}
	values[s] = strings.Replace(values[s], d, "", 1)

	// <== strategy (1)
	if len(values[s]) == 0 {
		return false // Contradiction: removed last value
	} else if len(values[s]) == 1 {
		d2 := values[s]
		for _, s2 := range peers[s] {
			if !eliminate(values, s2, d2) {
				return false
			}
		}
	}

	// <== strategy (2)
	for u := range units[s] {
		dplaces := []string{}
		for i := range units[s][u] {
			sq2 := units[s][u][i]
			if strings.Contains(values[sq2], d) {
				dplaces = append(dplaces, sq2)
			}
		}
		if len(dplaces) == 0 {
			return false
		} else if len(dplaces) == 1 {
			if !assign(values, dplaces[0], d) {
				return false
			}
		}
	}

	return true
}
