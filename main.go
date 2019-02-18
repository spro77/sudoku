package main

import (
	"fmt"
)

func cross(A []string, B []string) []string {
	var C []string
	for _, a := range A {
		for _, b := range B {
			C = append(C, (a + b))
		}
	}
	return C
}

func member(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

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

		for _, unit := range units[s] {
			for _, square := range unit {
				if square != s {
					set = append(set, square)
				}
			}
		}
		peers[s] = set
	}

	//___________________________________

	fmt.Println(parseGrid(examp))

	//fmt.Println("All cell peers", peers["I9"])
	//fmt.Println("All cell units", units["C2"])
	//fmt.Printf("len=%d cap=%d %v\n", len(squares), cap(squares), squares)

}

func parseGrid(grid string) map[string]string {
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
