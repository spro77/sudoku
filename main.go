package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

var (
	examp  = ".....6....59.....82....8....45........3........6..3.54...325..6.................."
	rows   = []string{"A", "B", "C", "D", "E", "F", "G", "H", "I"}
	cols   = []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}
	digits = "123456789"

	squares  = cross(rows, cols)       // The Array of all cell's names (81)
	unitlist = [][]string{}            // The Array of Cols, Rows and Boxes (9 column units, 9 row units and 9 box units)
	units    = map[string][][]string{} // Units (Value) per Square (Key)

	peers = map[string][]string{} // Peers, dependent squares (Value) per Square (Key)

)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

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

	for _, s := range squares {
		units[s] = [][]string{}
		for _, u := range unitlist {
			if member(s, u) {
				units[s] = append(units[s], u)
			}
		}
	}

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

	//___________________________________ User Prompting

	fmt.Println("From where do you want to input Sudoku initials?")
	fmt.Println("1: from file ( pls use data.txt )")
	fmt.Println("2: from prompt line")
	fmt.Println("3: from hardcoded example ( the most difficult example )")

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Your answer: ")
	char, _, err := reader.ReadRune()
	//text, err := reader.ReadString('\n')
	check(err)

	switch string(char) {
	case "1":
		buffer, err := ioutil.ReadFile("/Users/prostovsergey/go/sudoku/data.txt")
		check(err)
		display(solve(string(string(buffer))))
		fmt.Print("1 done")
	case "2":
		inputPrompt()
		fmt.Print("2 done")
	case "3":
		display(solve(string(examp)))
		fmt.Print("3 done")
	default:
		fmt.Print("def done")
	}
}

func inputPrompt() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Input you Sudoku initial, use any symbol as empty square")
	text, _ := reader.ReadString('\n')
	text = strings.Replace(text, "\n", "", -1)
	display(solve(text))
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

//___________________________________ Service method for Map duplicating

func dup(values map[string]string) map[string]string {
	newMap := make(map[string]string)
	for k, v := range values {
		newMap[k] = v
	}
	return newMap
}

//___________________________________ Vizualizing method

func display(values map[string]string) {

	for _, r := range rows {
		ind := ""
		for _, c := range cols {
			ind = r + c
			fmt.Print(values[ind] + "  ")
		}
		fmt.Println("")
	}
}

//___________________________________ Method for Storing an initial representation of Sudoku into a map using square as key

func mapValue(grid string) map[string]string {

	symbols := digits + ".0_+=,;/:\n "
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

//___________________________________ It updataes the incoming values by eliminating the other values than d

func assign(values map[string]string, s string, d string) (bool, map[string]string) {

	for _, d2 := range values[s] {
		if string(d2) != d {
			if status, _ := eliminate(values, s, string(d2)); !status {
				return false, values
			}
		}
	}
	return true, values
}

/*
1. It removes the given value d from values[s] which is a list of potential values for s.
2. If there is no values left in s (that is we don’t have any potential value for that square), returns False
3. When there is only one potential value for s, it removes the value from all the peers of s	<== strategy (1)
4. Make sure the given value d has a place elsewhere (i.e., if no square has d as a potential value, we can not solve the puzzle)
5. Where there is only one place for the value d, remove it from the peers	<== strategy (2)
*/

func eliminate(values map[string]string, s string, d string) (bool, map[string]string) {

	if !strings.Contains(values[s], d) {
		return true, values
	}

	values[s] = strings.Replace(values[s], d, "", 1)

	// <== strategy (1)
	if len(values[s]) == 0 {
		return false, values
	} else if len(values[s]) == 1 {
		d2 := values[s]
		for _, s2 := range peers[s] {
			if len(values[s2]) > 1 {
				if status, _ := eliminate(values, s2, d2); !status {
					return false, values
				}
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
			return false, values
		} else if len(dplaces) == 1 {
			if status, _ := assign(values, dplaces[0], d); !status {
				return false, values
			}
		}
	}

	return true, values
}

//___________________________________ Main iterator systematically try all possibilities until we hit one that works

func search(values map[string]string) (bool, map[string]string) {

	max, min, sq := 1, 9, ""

	for _, s := range squares {

		if len(values[s]) > max {
			max = len(values[s])
		}
		if len(values[s]) > 1 && len(values[s]) < min {
			min = len(values[s])
			sq = s
		}
	}

	if max == 1 {
		return true, values // Sudoku Solved!
	}

	for _, d := range values[sq] { // Select a square with minimum of potential values, and iterate...
		_, result := assign(dup(values), sq, string(d))
		status, values := search(result)
		if status {
			return status, values
		}
	}
	return false, values
}

func solve(grid string) map[string]string {
	status, result := search(parseGrid(grid))
	if status {
		return result
	}
	return map[string]string{}
}
