package main

type Token int
const (
	rules = `
		P' -> P
		P -> inicio V A
		V -> varinicio LV
		LV -> D LV
		LV -> varfim ;
		D -> TIPO L ;
		L -> id , L
		L -> id
		TIPO -> int
		TIPO -> real
		TIPO -> lit
		A -> ES A
		ES -> leia id ;
		ES -> escreva ARG ;
		ARG -> literal
		ARG -> num
		ARG -> id
		A -> CMD A
		CMD -> id rcb LD ;
		LD -> OPRD opm OPRD
		LD -> OPRD
		OPRD -> id
		OPRD -> num
		A -> COND A
		COND -> CAB CP
		CAB -> se ( EXP_R ) entao
		EXP_R -> OPRD opr OPRD
		CP -> ES CP
		CP -> CMD CP
		CP -> COND CP
		CP -> fimse
		A -> R A
		R  -> facaAte ( EXP_R ) CP_R
		CP_R -> ES CP_R
		CP_R -> CMD CP_R
		CP_R -> COND CP_R
		CP_R -> fimFaca
		A -> fim
	`
)

func main() {
	createObjFile()
	generator := NewGenerator(rules)
	parser := generator.BuildParser()
	tree, error, pos := parser.Parse(lexin())
	//
	//writeObjFile("Hello",0)
	if tree == nil {
		errorP(error,parser,pos, pos)
	}// else {
	//	printNode(tree, 0)
	//}
}
