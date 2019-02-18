package main

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

	squares  = cross(rows, cols) // The Array of all cell's names (81)
	unitlist = [][]string{}      // The Array of Cols, Rows and Boxes (9 column units, 9 row units and 9 box units)
	units    = map[string][][]string{}

	nassigns      = 0
	neliminations = 0
	nsearches     = 0
)

func main() {

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

	for _, s := range squares {

		units[s] = [][]string{}

		for _, u := range unitlist {
			if member(s, u) {
				units[s] = append(units[s], u)
			}
		}
	}

	//fmt.Println(units)
	//fmt.Println(units["C2"])
	//fmt.Printf("len=%d cap=%d %v\n", len(squares), cap(squares), squares)

}
