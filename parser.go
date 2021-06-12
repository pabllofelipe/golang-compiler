package main

import (
	"fmt"
	"github.com/fatih/color"
	"reflect"
	"strconv"
)

type Parser struct {
	*grammar

	// [symbol][state]data
	//
	// if data == 0 {
	//     null
	// } else if data > 0 {
	//     shift para simbolos terminais
	//     goto para simbolos não-terminais
	// } else if data < 0 {
	//     reduce
	// }
	table [][]int

	stateStack  []int
	symbolStack []int
}

type TreeNode struct {
	Name   string
	Leaves []*TreeNode
}

func newParser(g *grammar, table [][]int) *Parser {
	parser := new(Parser)
	parser.grammar = g
	parser.table = table
	parser.stateStack = make([]int, 0)
	parser.symbolStack = make([]int, 0)

	return parser
}

func newTreeNode(name string) *TreeNode {
	node := new(TreeNode)
	node.Name = name
	node.Leaves = make([]*TreeNode, 0)

	return node
}

func (parser *Parser) Parse(symbols []TOKEN, error_pos []Position) (*TreeNode, int, Position) {
	nodes := make([]*TreeNode, 0)
	var lastError Position
	cif := TOKEN{"$","$","Nulo"}
	symbols = append(symbols, cif)
	parser.pushState(0)
	lastSymbol := 0
	aux := TOKEN{"","",""}
	temp := 0
	counter := 0
	var temp_var []string
	for {
		symName := symbols[0].classe
		symId, exists := parser.grammar.symbolIndexByName[symName]
		if !exists {
			lastError = errorP(lastSymbol,parser,error_pos[counter], lastError)
			symbols = symbols[1:]
			continue
			//return nil, lastSymbol, error_pos[counter]
		}
		symbol := parser.grammar.symbols[symId]
		state := parser.peekState()
		//fmt.Println(symbol.name)
		inst := parser.table[symId][state]

		if inst == 0 {
			// error
			lastError = errorP(lastSymbol,parser,error_pos[counter], lastError)
			symbols = symbols[1:]
			continue
			//return nil, lastSymbol, error_pos[counter]
		} else if inst > 0 {
			// shift
			if (state == 3) || (state == 10) || (state == 68) || (state == 14){
				aux = symbols[0]
			}
			if (state == 15){
				toString := "T"+ strconv.Itoa(temp)
				temp_var = append(temp_var, toString)
				writeObjFile(temp_var[temp] + "=" + symbols[1].lexema + symbols[2].lexema + symbols[3].lexema + ";\n",1)
				writeObjFile("while(" + temp_var[temp] + ")\n",1)
				temp++
			}
			if (state == 56){
				writeObjFile("{\n",1)
			}
			if (state == 69){
				aux = symbols[0]
				if (symbols[0].classe != "id") {
					writeObjFile("}\n", 1)
				}
			}
			if (state == 43) || (state == 28) || (state == 47) || (state == 26){
				writeObjFile(";\n",1)
			}
			if (state == 13){
				toString := "T"+ strconv.Itoa(temp)
				temp_var = append(temp_var, toString)
				writeObjFile(temp_var[temp] + "=" + symbols[1].lexema + symbols[2].lexema + symbols[3].lexema + ";\n",1)
				writeObjFile("if(" + temp_var[temp] + ")\n",1)
				temp++
			}
			if (state == 62){
				writeObjFile("{\n",1)
			}
			//13
			if (state == 8){
				if (symbols[0].tipo == "literal"){
					writeObjFile("scanf(\"%s\"," + symbols[0].lexema + ")",1)
				}
				if (symbols[0].tipo == "Nulo"){
					color.Red("ERRO Semântico na linha %d coluna %d: Variável não declarada",error_pos[counter].line,error_pos[counter].column)
					lastError.line = error_pos[counter].line
					lastError.column = error_pos[counter].column
				}
				if (symbols[0].tipo == "int"){
					writeObjFile("scanf(\"%d\",&" + symbols[0].lexema + ")",1)
				}
				if (symbols[0].tipo == "double"){
					writeObjFile("scanf(\"%lf\",&" + symbols[0].lexema + ")",1)
				}
			}
			if (state == 9){
				if (symbols[0].classe == "literal") {
					if (symbols[0].lexema == ""){
						writeObjFile("printf(\"\\n\")", 1)
					}else {
						writeObjFile("printf(\""+symbols[0].lexema+"\")", 1)
					}
				}else{
				if (symbols[0].tipo == "int") {
					writeObjFile("printf(\"%d\","+symbols[0].lexema+")", 1)
				}else{
				if (symbols[0].tipo == "double") {
					writeObjFile("printf(\"%lf\","+symbols[0].lexema+")", 1)
				}else{
					if (symbols[0].tipo == "literal") {
						writeObjFile("printf(\"%s\","+symbols[0].lexema+")", 1)
					}
				}
				}
				}
			}
			if (state == 11){
				aux = symbols[0]
			}
			//19
			if (state == 6){
				if(symbols[2].lexema == ";"){
					if(aux.tipo == symbols[1].tipo){
						writeObjFile(aux.lexema + "=" + symbols[1].lexema,1)
					}else{
						color.Red("ERRO Semântico na linha %d coluna %d: Variáveis de tipos diferentes",error_pos[counter].line + 1,error_pos[counter].column)
						lastError.line = error_pos[counter].line
						lastError.column = error_pos[counter].column
					}
				}else {
					if (aux.tipo == symbols[1].tipo) {
						toString := "T" + strconv.Itoa(temp)
						temp_var = append(temp_var, toString)
						writeObjFile(temp_var[temp]+"="+symbols[1].lexema, 1)
						temp++
					} else {
						color.Red("ERRO Semântico na linha %d coluna %d: Variáveis de tipos diferentes",error_pos[counter].line,error_pos[counter].column)
						lastError.line = error_pos[counter].line
						lastError.column = error_pos[counter].column
					}
				}
			}
			//20
			if (state == 24){
				if (symbols[0].tipo == symbols[2].tipo){
					temp--
					writeObjFile(symbols[1].lexema + symbols[2].lexema + symbols[3].lexema,1)
					writeObjFile( "\n" + aux.lexema + "=" + temp_var[temp],1)
					temp++
				}
			}
			lastSymbol = symId
			parser.pushState(inst - 1)
			parser.pushSymbol(symId)
			symbols = symbols[1:]
			counter++
			nodes = append(nodes, newTreeNode(symName))

			fmt.Printf("ACTION(%d, %s) = shift %d\n", state, symbol.name, parser.peekState())
		} else if inst < 0 {
			// reduce

			proId := -inst - 1
			pro := parser.grammar.productions[proId]
			if (proId == 8) || (proId == 9) || (proId == 10){
				writeObjFile(symbols[0].tipo + " " + symbols[0].lexema,1)
				fmt.Printf("TIPO.tipo")
			}
			if(proId == 30){
				writeObjFile("}\n",1)
			}

			node := newTreeNode(parser.grammar.symbols[pro.lhs].name)
			for i := 0; i < pro.len; i++ {
				parser.popSymbol()
				parser.popState()
				node.Leaves = append(node.Leaves, nodes[len(nodes)-1])
				nodes = nodes[:len(nodes)-1]
			}
			nodes = append(nodes, node)

			parser.pushSymbol(pro.lhs)

			if proId == 0 {
				writeObjFile("\n}", 1)
				addBegFile(temp_var)
				if(lastError.line != 0) && (lastError.column != 0){
					deleteObjFile()
				}
				fmt.Println("ACCEPTED")
				break
			}
			fmt.Printf("ACTION(%d, %s) = reduce %d\n", state, symbol.name, proId)

			state = parser.peekState()
			gotoState := parser.table[pro.lhs][state]
			if gotoState == 0 {
				return nil, lastSymbol, error_pos[counter]
			}
			gotoState--
			parser.pushState(gotoState)

			fmt.Printf("ACTION(%d, %s) = goto %d\n", state, parser.grammar.symbols[pro.lhs].name, parser.peekState())

		}

		if len(symbols) == 0 {
			break
		}
	}
	aux = TOKEN{"","",""}
	return nodes[0], lastSymbol, error_pos[0]
}

func (parser *Parser) peekState() int {
	return parser.stateStack[len(parser.stateStack)-1]
}

func (parser *Parser) pushState(state int) {
	parser.stateStack = append(parser.stateStack, state)
}

func (parser *Parser) popState() int {
	len := len(parser.stateStack)
	state := parser.stateStack[len-1]
	parser.stateStack = parser.stateStack[:len-1]
	return state
}

func (parser *Parser) pushSymbol(symId int) {
	parser.symbolStack = append(parser.symbolStack, symId)
}

func (parser *Parser) popSymbol() int {
	len := len(parser.symbolStack)
	symId := parser.symbolStack[len-1]
	parser.symbolStack = parser.symbolStack[:len-1]
	return symId
}

func (treeNode *TreeNode) addLeaf(leaf *TreeNode) {
	treeNode.Leaves = append(treeNode.Leaves, leaf)
}

func Distinct(arr interface{}) (reflect.Value, bool) {

	slice, ok := takeArg(arr, reflect.Slice)
	if !ok {
		return reflect.Value{}, ok
	}

	c := slice.Len()
	m := make(map[interface{}]bool)
	for i := 0; i < c; i++ {
		m[slice.Index(i).Interface()] = true
	}
	mapLen := len(m)
	out := reflect.MakeSlice(reflect.TypeOf(arr), mapLen, mapLen)
	i := 0
	for k := range m {
		v := reflect.ValueOf(k)
		o := out.Index(i)
		o.Set(v)
		i++
	}

	return out, ok
}
func Intersect(arrs ...interface{}) (reflect.Value, bool) {
	arrLength := len(arrs)
	var kind reflect.Kind
	var kindHasBeenSet bool

	tempMap := make(map[interface{}]int)
	for _, arg := range arrs {
		tempArr, ok := Distinct(arg)
		if !ok {
			return reflect.Value{}, ok
		}

		if kindHasBeenSet && tempArr.Len() > 0 && tempArr.Index(0).Kind() != kind {
			return reflect.Value{}, false
		}
		if tempArr.Len() > 0 {
			kindHasBeenSet = true
			kind = tempArr.Index(0).Kind()
		}

		c := tempArr.Len()
		for idx := 0; idx < c; idx++ {
			if _, ok := tempMap[tempArr.Index(idx).Interface()]; ok {
				tempMap[tempArr.Index(idx).Interface()]++
			} else {
				tempMap[tempArr.Index(idx).Interface()] = 1
			}
		}
	}

	numElems := 0
	for _, v := range tempMap {
		if v == arrLength {
			numElems++
		}
	}
	out := reflect.MakeSlice(reflect.TypeOf(arrs[0]), numElems, numElems)
	i := 0
	for key, val := range tempMap {
		if val == arrLength {
			v := reflect.ValueOf(key)
			o := out.Index(i)
			o.Set(v)
			i++
		}
	}

	return out, true
}

func Union(arrs ...interface{}) (reflect.Value, bool) {
	tempMap := make(map[interface{}]uint8)
	var kind reflect.Kind
	var kindHasBeenSet bool

	for _, arg := range arrs {
		tempArr, ok := Distinct(arg)
		if !ok {
			return reflect.Value{}, ok
		}
		if kindHasBeenSet && tempArr.Len() > 0 && tempArr.Index(0).Kind() != kind {
			return reflect.Value{}, false
		}
		if tempArr.Len() > 0 {
			kindHasBeenSet = true
			kind = tempArr.Index(0).Kind()
		}

		c := tempArr.Len()
		for idx := 0; idx < c; idx++ {
			tempMap[tempArr.Index(idx).Interface()] = 0
		}
	}

	mapLen := len(tempMap)
	out := reflect.MakeSlice(reflect.TypeOf(arrs[0]), mapLen, mapLen)
	i := 0
	for key := range tempMap {
		v := reflect.ValueOf(key)
		o := out.Index(i)
		o.Set(v)
		i++
	}

	return out, true
}
func Difference(arrs ...interface{}) (reflect.Value, bool) {
	tempMap := make(map[interface{}]int)
	var kind reflect.Kind
	var kindHasBeenSet bool

	for _, arg := range arrs {
		tempArr, ok := Distinct(arg)
		if !ok {
			return reflect.Value{}, ok
		}

		if kindHasBeenSet && tempArr.Len() > 0 && tempArr.Index(0).Kind() != kind {
			return reflect.Value{}, false
		}
		if tempArr.Len() > 0 {
			kindHasBeenSet = true
			kind = tempArr.Index(0).Kind()
		}

		c := tempArr.Len()
		for idx := 0; idx < c; idx++ {
			if _, ok := tempMap[tempArr.Index(idx).Interface()]; ok {
				tempMap[tempArr.Index(idx).Interface()]++
			} else {
				tempMap[tempArr.Index(idx).Interface()] = 1
			}
		}
	}

	numElems := 0
	for _, v := range tempMap {
		if v == 1 {
			numElems++
		}
	}
	out := reflect.MakeSlice(reflect.TypeOf(arrs[0]), numElems, numElems)
	i := 0
	for key, val := range tempMap {
		if val == 1 {
			v := reflect.ValueOf(key)
			o := out.Index(i)
			o.Set(v)
			i++
		}
	}

	return out, true
}

func takeArg(arg interface{}, kind reflect.Kind) (val reflect.Value, ok bool) {
	val = reflect.ValueOf(arg)
	if val.Kind() == kind {
		ok = true
	}
	return
}

func IntersectString(args ...[]string) []string {
	arrLength := len(args)
	tempMap := make(map[string]int)
	for _, arg := range args {
		tempArr := DistinctString(arg)
		for idx := range tempArr {
			if _, ok := tempMap[tempArr[idx]]; ok {
				tempMap[tempArr[idx]]++
			} else {
				tempMap[tempArr[idx]] = 1
			}
		}
	}

	tempArray := make([]string, 0)
	for key, val := range tempMap {
		if val == arrLength {
			tempArray = append(tempArray, key)
		}
	}

	return tempArray
}

func IntersectStringArr(arr [][]string) []string {
	arrLength := len(arr)
	tempMap := make(map[string]int)
	for idx1 := range arr {
		tempArr := DistinctString(arr[idx1])
		for idx2 := range tempArr {
			if _, ok := tempMap[tempArr[idx2]]; ok {
				tempMap[tempArr[idx2]]++
			} else {
				tempMap[tempArr[idx2]] = 1
			}
		}
	}

	tempArray := make([]string, 0)
	for key, val := range tempMap {
		if val == arrLength {
			tempArray = append(tempArray, key)
		}
	}

	return tempArray
}

func UnionString(args ...[]string) []string {
	tempMap := make(map[string]uint8)
	for _, arg := range args {
		for idx := range arg {
			tempMap[arg[idx]] = 0
		}
	}

	tempArray := make([]string, 0)
	for key := range tempMap {
		tempArray = append(tempArray, key)
	}

	return tempArray
}

func UnionStringArr(arr [][]string) []string {
	tempMap := make(map[string]uint8)

	for idx1 := range arr {
		for idx2 := range arr[idx1] {
			tempMap[arr[idx1][idx2]] = 0
		}
	}

	tempArray := make([]string, 0)
	for key := range tempMap {
		tempArray = append(tempArray, key)
	}

	return tempArray
}

func DifferenceString(args ...[]string) []string {
	tempMap := make(map[string]int)
	for _, arg := range args {
		tempArr := DistinctString(arg)
		for idx := range tempArr {
			if _, ok := tempMap[tempArr[idx]]; ok {
				tempMap[tempArr[idx]]++
			} else {
				tempMap[tempArr[idx]] = 1
			}
		}
	}

	tempArray := make([]string, 0)
	for key, val := range tempMap {
		if val == 1 {
			tempArray = append(tempArray, key)
		}
	}

	return tempArray
}

func DifferenceStringArr(arr [][]string) []string {
	tempMap := make(map[string]int)
	for idx1 := range arr {
		tempArr := DistinctString(arr[idx1])
		for idx2 := range tempArr {
			if _, ok := tempMap[tempArr[idx2]]; ok {
				tempMap[tempArr[idx2]]++
			} else {
				tempMap[tempArr[idx2]] = 1
			}
		}
	}

	tempArray := make([]string, 0)
	for key, val := range tempMap {
		if val == 1 {
			tempArray = append(tempArray, key)
		}
	}

	return tempArray
}

func DistinctString(arg []string) []string {
	tempMap := make(map[string]uint8)

	for idx := range arg {
		tempMap[arg[idx]] = 0
	}

	tempArray := make([]string, 0)
	for key := range tempMap {
		tempArray = append(tempArray, key)
	}
	return tempArray
}

func IntersectUint64(args ...[]uint64) []uint64 {
	arrLength := len(args)
	tempMap := make(map[uint64]int)
	for _, arg := range args {
		tempArr := DistinctUint64(arg)
		for idx := range tempArr {
			if _, ok := tempMap[tempArr[idx]]; ok {
				tempMap[tempArr[idx]]++
			} else {
				tempMap[tempArr[idx]] = 1
			}
		}
	}

	tempArray := make([]uint64, 0)
	for key, val := range tempMap {
		if val == arrLength {
			tempArray = append(tempArray, key)
		}
	}

	return tempArray
}

func DistinctIntersectUint64(args ...[]uint64) []uint64 {
	arrLength := len(args)
	tempMap := make(map[uint64]int)
	for _, arg := range args {
		for idx := range arg {
			if _, ok := tempMap[arg[idx]]; ok {
				tempMap[arg[idx]]++
			} else {
				tempMap[arg[idx]] = 1
			}
		}
	}

	tempArray := make([]uint64, 0)
	for key, val := range tempMap {
		if val == arrLength {
			tempArray = append(tempArray, key)
		}
	}

	return tempArray
}

func sortedIntersectUintHelper(a1 []uint64, a2 []uint64) []uint64 {
	intersection := make([]uint64, 0)
	n1 := len(a1)
	n2 := len(a2)
	i := 0
	j := 0
	for i < n1 && j < n2 {
		switch {
		case a1[i] > a2[j]:
			j++
		case a2[j] > a1[i]:
			i++
		default:
			intersection = append(intersection, a1[i])
			i++
			j++
		}
	}
	return intersection
}

func SortedIntersectUint64(args ...[]uint64) []uint64 {
	tempIntersection := args[0]
	argsLen := len(args)

	for k := 1; k < argsLen; k++ {
		switch len(tempIntersection) {
		case 0:
			return tempIntersection

		default:
			tempIntersection = sortedIntersectUintHelper(tempIntersection, args[k])
		}
	}

	return tempIntersection
}

func IntersectUint64Arr(arr [][]uint64) []uint64 {
	arrLength := len(arr)
	tempMap := make(map[uint64]int)
	for idx1 := range arr {
		tempArr := DistinctUint64(arr[idx1])
		for idx2 := range tempArr {
			if _, ok := tempMap[tempArr[idx2]]; ok {
				tempMap[tempArr[idx2]]++
			} else {
				tempMap[tempArr[idx2]] = 1
			}
		}
	}

	tempArray := make([]uint64, 0)
	for key, val := range tempMap {
		if val == arrLength {
			tempArray = append(tempArray, key)
		}
	}

	return tempArray
}

func SortedIntersectUint64Arr(arr [][]uint64) []uint64 {
	tempIntersection := arr[0]
	argsLen := len(arr)

	for k := 1; k < argsLen; k++ {
		switch len(tempIntersection) {
		case 0:
			return tempIntersection

		default:
			tempIntersection = sortedIntersectUintHelper(tempIntersection, arr[k])
		}
	}

	return tempIntersection
}

func DistinctIntersectUint64Arr(arr [][]uint64) []uint64 {
	arrLength := len(arr)
	tempMap := make(map[uint64]int)
	for idx1 := range arr {
		for idx2 := range arr[idx1] {
			if _, ok := tempMap[arr[idx1][idx2]]; ok {
				tempMap[arr[idx1][idx2]]++
			} else {
				tempMap[arr[idx1][idx2]] = 1
			}
		}
	}

	tempArray := make([]uint64, 0)
	for key, val := range tempMap {
		if val == arrLength {
			tempArray = append(tempArray, key)
		}
	}

	return tempArray
}

func UnionUint64(args ...[]uint64) []uint64 {
	tempMap := make(map[uint64]uint8)

	for _, arg := range args {
		for idx := range arg {
			tempMap[arg[idx]] = 0
		}
	}

	tempArray := make([]uint64, 0)
	for key := range tempMap {
		tempArray = append(tempArray, key)
	}

	return tempArray
}

func UnionUint64Arr(arr [][]uint64) []uint64 {
	tempMap := make(map[uint64]uint8)

	for idx1 := range arr {
		for idx2 := range arr[idx1] {
			tempMap[arr[idx1][idx2]] = 0
		}
	}

	tempArray := make([]uint64, 0)
	for key := range tempMap {
		tempArray = append(tempArray, key)
	}

	return tempArray
}

func DifferenceUint64(args ...[]uint64) []uint64 {
	tempMap := make(map[uint64]int)
	for _, arg := range args {
		tempArr := DistinctUint64(arg)
		for idx := range tempArr {
			if _, ok := tempMap[tempArr[idx]]; ok {
				tempMap[tempArr[idx]]++
			} else {
				tempMap[tempArr[idx]] = 1
			}
		}
	}

	tempArray := make([]uint64, 0)
	for key, val := range tempMap {
		if val == 1 {
			tempArray = append(tempArray, key)
		}
	}

	return tempArray
}

func DifferenceUint64Arr(arr [][]uint64) []uint64 {
	tempMap := make(map[uint64]int)
	for idx1 := range arr {
		tempArr := DistinctUint64(arr[idx1])
		for idx2 := range tempArr {
			if _, ok := tempMap[tempArr[idx2]]; ok {
				tempMap[tempArr[idx2]]++
			} else {
				tempMap[tempArr[idx2]] = 1
			}
		}
	}

	tempArray := make([]uint64, 0)
	for key, val := range tempMap {
		if val == 1 {
			tempArray = append(tempArray, key)
		}
	}

	return tempArray
}

func DistinctUint64(arg []uint64) []uint64 {
	tempMap := make(map[uint64]uint8)

	for idx := range arg {
		tempMap[arg[idx]] = 0
	}

	tempArray := make([]uint64, 0)
	for key := range tempMap {
		tempArray = append(tempArray, key)
	}
	return tempArray
}

func errorP(id int, p *Parser, error_pos Position, lastError Position) Position {
	if (lastError == error_pos){
		return lastError
	}
	stop := 0
	var nextIds []int
	IdNext := 0
	IDNEXT:
	if id == 11 {
		color.Red("ERRO na linha %d coluna %d: TOKEN's ESPERADOS: [id], [varfim], [escreva], [leia], [se], [entao], [fimse], [fim]\n",error_pos.line,error_pos.column)
		lastError = error_pos
		return lastError
	}else{
		for _, prod := range p.grammar.productions {
			for i, rhsId := range prod.rhs {
				if id == rhsId && (len(prod.rhs)-1) > i {
					stop = 1
				} else if id == rhsId && (len(prod.rhs)-1) == i{
					IdNext = prod.lhs
				}else if stop == 1 {
					nextIds = append(nextIds, rhsId)
					stop = 0
				}
			}
		}
	}
	if IdNext != 0{
		id = IdNext
		IdNext = 0
		goto IDNEXT
	}
	DERIVATION:
	for i, lhsId := range nextIds{
		if !p.grammar.symbols[lhsId].terminal{
			for _, prod := range p.grammar.productions {
					if lhsId == prod.lhs {
						if id == prod.rhs[0]{
							copy(nextIds[i:],nextIds[i+1:])
							nextIds[len(nextIds)-1] = 0
							nextIds = nextIds[:len(nextIds)-1]
						}else{
							nextIds[i] = prod.rhs[0]
							i++
							nextIds = append(nextIds,0)
						}
					}
			}
		}
	}
	for _, lhsId := range nextIds {
		if !p.grammar.symbols[lhsId].terminal{
			goto DERIVATION
		}
	}

	z, ok := Union(nextIds,nextIds)
	if !ok {
		fmt.Println("Não foi possivel fazer a uniao")
	}
	nextIds, ok = z.Interface().([]int)
	if !ok {
		fmt.Println("Nao foi possivel converter em slice")
	}

	color.Red("ERRO na linha %d coluna %d: TOKEN's ESPERADOS:",error_pos.line,error_pos.column)
	for i:=0; i < len(nextIds);i++{
		if(nextIds[i]!=0){
			color.Red("[%s] ",p.grammar.symbols[nextIds[i]].name)
		}
	}
	fmt.Printf("\n")
	lastError = error_pos
	return lastError
}