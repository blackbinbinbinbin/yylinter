package a	// want "golang-rule error【golang-rule-1.5.1】"

import (	//want "golang-rule error【golang-rule-1.3.1】"
	. "fmt"
)

const ABC = 1
const EFGF = 2
var d int
const EFG = 3	//want "golang-rule suggest【golang-rule-2.5.2】"
const ABC_EFG = 4	//want "golang-rule suggest【golang-rule-2.5.1】"

func main() {
	Println("yylinter testdata")
}