package main

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"io"
	"log"
	"os"
	"strconv"
	"unicode"
)

//Classes

const (
	EOF = iota
	ERRO
	id // L(L|D|_)∗
	NUM // D+(\.D+)?((E|e) (+|−)?D+)?
	Literal // “ .* “
	PT_V // ;
	AB_P // (
	FC_P // )
	OPM // + , -, *, /
	Vir // ,
	RCB // <-
	OPR // <, >, >= , <= , =, <>
	IGN
)

var tokens = []string{
	EOF: "EOF",
	ERRO: "ERRO",
	id: "id",
	NUM: "num",
	Literal: "literal",
	PT_V: "PT_V",
	AB_P: "AB_P",
	FC_P: "FC_P",
	OPM: "opm",
	Vir: "vir",
	RCB: "rcb",
	OPR: "opr",
	IGN: "IGN",

}


type Position struct {
	line   int
	column int
}

type Lexer struct {
	pos    Position
	reader *bufio.Reader
}

func NewLexer(reader io.Reader) *Lexer {
	return &Lexer{
		pos:    Position{line: 1, column: 0},
		reader: bufio.NewReader(reader),
	}
}
func isInt(val float64) bool {
	return val == float64(int(val))
}
//Tamanho do Array
const ArraySize=10

//Estrutura TOKEN
type TOKEN struct {
	classe string
	lexema string
	tipo string
}

// Palavras reservadas da linguagem mgol
func mgolkeys() []TOKEN{
	keywords := []TOKEN{
		{"inicio","inicio","Nulo"},
		{"varinicio","varinicio","Nulo"},
		{"varfim","varfim","Nulo"},
		{"escreva","escreva","Nulo"},
		{"leia","leia","Nulo"},
		{"se","se","Nulo"},
		{"entao","entao","Nulo"},
		{"fimse","fimse","Nulo"},
		{"facaAte","facaAte","Nulo"},
		{"fimFaca","fimFaca","Nulo"},
		{"fim","fim","Nulo"},
		{"int","int","int"},
		{"lit","lit","literal"},
		{"real","real","double"},
	}
	return keywords
}

//TabelaHash
type HashTable struct {
	array [ArraySize]*bucket
}

//ListaEncadeada
type bucket struct {
	head *bucketNode
}

//NóListaEncadeada
type bucketNode struct {
	key TOKEN
	next *bucketNode
}

//initHashTable
func Init() *HashTable {
	result := &HashTable{}
	for i := range result.array {
		result.array[i] = &bucket{}
	}
	for _, v := range mgolkeys(){
		result.Insert(v)
	}
	return result
}

//InserçãoTabelaHash
func (h *HashTable) Insert(key TOKEN){
	index:=hash(key)
	h.array[index].insert(key)
}

//RemoçãoTabelaHash
func (h *HashTable) Remove(key TOKEN){
	index := hash(key)
	h.array[index].remove(key)
}

//BuscaTabelaHash
func (h *HashTable) Search(key TOKEN) TOKEN{
	index := hash(key)
	return h.array[index].search(key.lexema)
}

//FunçãoHash
func hash(key TOKEN) int{
	sum:= 0
	for _,v:= range key.lexema{
		sum+=int(v)
	}
	return sum % ArraySize
}

//InserçãoListaEncadeada
func (b *bucket) insert(k TOKEN){
	newNode := &bucketNode{key: k}
	newNode.next = b.head
	b.head = newNode
}

//RemoçãoListaEncadeada
func (b *bucket) remove(k TOKEN){
	if b.head.key == k {
		b.head = b.head.next
		return
	}
	previousNode := b.head
	for previousNode.next != nil{
		if previousNode.next.key == k{
			previousNode.next = previousNode.next.next
		}
		previousNode = previousNode.next
	}
}

//BuscaListaEncadeada
func (b *bucket) search(k string) TOKEN{
	currentNode := b.head
	for currentNode != nil{
		if currentNode.key.lexema == k{
			return currentNode.key
		}
		currentNode = currentNode.next
	}
	return TOKEN{"nil","nil","nil"}
}

//Printa a tabela de simbolos
func (h *HashTable) PrintTable(){
	for _,i:=range h.array{
		curr:= i.head
		for curr != nil{
			fmt.Printf("Classe: %s, Lexema: %s, Tipo: %s \n",curr.key.classe,curr.key.lexema,curr.key.tipo)
			curr = curr.next
		}
	}
}

//Função Scanner

func (h *HashTable)scanner(lexer *Lexer, erro *int, last_t TOKEN) (TOKEN, Position, TOKEN) {
	t:=TOKEN{"","","Nulo"}
	pos, tok, lit := lexer.Lex()
	t.lexema = lit
	t.classe = lit
	if tok == ERRO{
		fmt.Printf("Classe: ERRO%v, Lexema: %s, Tipo: %s \n",*erro,t.lexema,t.tipo)
		color.Red("ERRO%v - Caractere inválido na linguagem, linha: %d, coluna %d \n",*erro,pos.line,pos.column)
		*erro++
		return t, pos, t
	}else if tok == EOF {
		//fmt.Printf("Classe: %s, Lexema: %s, Tipo: %s \n",t.classe,t.lexema,t.tipo)
		return t, pos, t
	}else if tok == Literal {
		t.classe = "literal"
		//fmt.Printf("Classe: %s, Lexema: %s, Tipo: %s \n",t.classe,t.lexema,t.tipo)
		return t, pos, t
	}else if tok == OPR {
		t.classe = "opr"
		//fmt.Printf("Classe: %s, Lexema: %s, Tipo: %s \n",t.classe,t.lexema,t.tipo)
		return t, pos, t
	}else if t.classe == "" {
		t.classe = "literal"
		//fmt.Printf("Classe: %s, Lexema: %s, Tipo: %s \n",t.classe,t.lexema,t.tipo)
		return t, pos, t
	}else if tok == OPM {
		t.classe = "opm"
		//fmt.Printf("Classe: %s, Lexema: %s, Tipo: %s \n",t.classe,t.lexema,t.tipo)
		return t, pos, t
	}else if tok == RCB {
		t.classe = "rcb"
		//fmt.Printf("Classe: %s, Lexema: %s, Tipo: %s \n",t.classe,t.lexema,t.tipo)
		return t, pos, t
	}else if tok == NUM {
		t.classe = "num"
		s, err := strconv.ParseFloat(t.lexema, 32)
		if err == nil {
			if (isInt(s)){
				t.tipo = "int"
			}else{
				t.tipo = "double"
			}
			//fmt.Printf("Classe: %s, Lexema: %s, Tipo: %s \n", t.classe, t.lexema, t.tipo)
			return t, pos, t
		}
	}else if tok == id {
		if (last_t.lexema == "real")||(last_t.lexema == "int")||(last_t.lexema == "lit") {
			t.classe = tokens[tok]
		}
		if t.lexema == "lit"{
			t.classe = "lit"
			t.tipo = "literal"
			last_t = t
		}
		if t.lexema == "int"{
			t.classe = "int"
			t.tipo = "int"
			last_t = t
		}
		if t.lexema == "real"{
			t.classe = "real"
			t.tipo = "double"
			last_t = t
		}
		res := h.Search(t)
		if res.lexema != "nil"{
			//fmt.Printf("Classe: %s, Lexema: %s, Tipo: %s \n",t.classe,t.lexema,t.tipo)
			return res, pos, res
		}else{
			res := h.Search(t)
			if res.lexema != "nil"{
				//fmt.Printf("Classe: %s, Lexema: %s, Tipo: %s \n", t.classe, t.lexema, t.tipo)
				return res, pos, res
			}else {
				t.tipo = last_t.tipo
				t.classe = tokens[tok]
				h.Insert(t)
				//TODO: adicionar no print a estrutura TOKEN
				//fmt.Printf("Classe: %s, Lexema: %s, Tipo: Nulo \n", tokens[tok], lit)
				return t, pos, t
			}
		}
	}else{
		//fmt.Printf("Classe: %s, Lexema: %s, Tipo: Nulo \n", tokens[tok], lit)
		return t, pos, t
	}
	//fmt.Printf("%d:%d\t%s\t%s\n", pos.line, pos.column, tok, lit)

	//for _, each_ln := range text {
	//	//fmt.Println(each_ln)
	//	res := strings.Fields(each_ln)
	//	for _,each_word := range res{
	//		t = TOKEN{each_word,each_word,"Nulo"}
	//		if(h.Search(t)){
	//			fmt.Printf("Classe: %s, Lexema: %s, Tipo: %s \n",t.classe,t.lexema,t.tipo)
	//		}
	//
	//	}
	//	//fmt.Println(res)
	//}
	return t, pos, t
}

func (l *Lexer) Lex() (Position, Token, string) {
	// Continua até retornar um token
	for {
		l.pos.column++
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				return l.pos, EOF, "EOF"
			}
			panic(err)
		}
		switch r {
		case '\n'://salto de linha
			l.resetPosition()
		case ';'://Ponto e vírgula
			return l.pos, PT_V, ";"
		case '+'://Adição
			return l.pos, OPM, "+"
		case '-'://subtração
			return l.pos, OPM, "-"
		case '*'://Multiplicação
			return l.pos, OPM, "*"
		case '/'://Divisão
			return l.pos, OPM, "/"
		case '('://Abre parenteses
			return l.pos,AB_P,"("
		case ')'://Fecha parenteses
			return l.pos,FC_P,")"
		case '.'://Valor fracionado
			lit:= string(r)
			l.pos.column++
			r, _, _ = l.reader.ReadRune()
			if unicode.IsDigit(r) {
				l.backup()
				lit = lit + l.lexInt()
				return l.pos, NUM, lit
			}else{
				return l.pos, ERRO,lit
			}
		case '{'://Comentário
			lit := l.lexIdent() + " "
			r, _, _ := l.reader.ReadRune()
			l.backup()
			for r != '}'{
				r, _, _ = l.reader.ReadRune()
				lit = lit + l.lexIdent()
				l.pos.column++
			}
			r, _, _ = l.reader.ReadRune()
			lit = lit + l.lexIdent()
			l.pos.column++
			return l.pos,IGN,lit
		case '"'://Constante Literal
			lit := ""
			r, _, _ := l.reader.ReadRune()
			for r != '"'{
				if(r != 92) {
					lit = lit + string(r)
					l.pos.column++
				}else{
					l.pos.column++
					r, _, _ = l.reader.ReadRune()
				}
				r, _, _ = l.reader.ReadRune()
			}
			if lit != "" {
				return l.pos, Literal, lit
			}else{
				return l.pos, IGN, lit
			}
		case '='://Operador relacional de igualdade
			return l.pos, OPR, "="
		case '<'://Verificação adicional
			r, _, err := l.reader.ReadRune()
			if err != nil {
				if err == io.EOF {
					return l.pos, EOF, ""
				}
				panic(err)
			}
			switch r {
			case '-'://Atribuição
				l.pos.column++
				return l.pos, RCB, "<-"
			case '>'://Operador relacional ?
				l.pos.column++
				return l.pos, OPR, "<>"
			case '='://Operador relacional de menor ou igual
				l.pos.column++
				return l.pos, OPR, "<="
			default://Operador relacional de menor
				return l.pos, OPR, "<"
			}
		case '>'://Verificação adicional
			switch r {
			case '='://Operador relacional de maior ou igual
				l.pos.column++
				return l.pos, OPR, ">="
			default://Operador relacional de maior
				return l.pos, OPR, ">"
			}
		default:
			if unicode.IsSpace(r) {
				continue // Ignora TAB, espaço
			}else if unicode.IsDigit(r) {
				// Constante numérica
				startPos := l.pos
				l.backup()
				lit := l.lexInt()
				r, _, _ = l.reader.ReadRune()
				if r == '.'{
					lit = lit + string(r)
					l.pos.column++
					r, _, _ = l.reader.ReadRune()
					l.backup()
					lit = lit + l.lexInt()
				}
				return startPos, NUM, lit
			}else if unicode.IsLetter(r) || r == '_' {
				// id
				startPos := l.pos
				l.backup()
				lit := ""
				lit = lit + l.lexInt()
				U:=false
				if r =='_' {
					U = true
				}
				L:=unicode.IsLetter(r)
				D:=unicode.IsDigit(r)
				l.pos.column++
				r, _, _ = l.reader.ReadRune()
				for L || D || U{
					L = unicode.IsLetter(r)
					D = unicode.IsDigit(r)
					if L {
						l.backup()
						lit = lit + l.lexIdent()
					}
					if D {
						l.backup()
						lit = lit + l.lexInt()
					}
					if r == '_'{
						U = true
						lit = lit + string(r)
					}else {
						U = false
					}
					r, _, _ = l.reader.ReadRune()
				}
				if r!=0 {
					l.backup()
				}
				return startPos, id, lit
			} else {
				return l.pos, ERRO, string(r)
			}
		}
	}
}

func (l *Lexer) resetPosition() {
	l.pos.line++
	l.pos.column = 0
}

func (l *Lexer) backup() {
	if err := l.reader.UnreadRune(); err != nil {
		panic(err)
	}

	l.pos.column--
}

func (l *Lexer) lexInt() string {
	var lit string
	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				// at the end of the int
				return lit
			}
		}

		l.pos.column++
		if unicode.IsDigit(r) {
			lit = lit + string(r)
		} else {
			l.backup()
			return lit
		}
	}
}

func (l *Lexer) lexIdent() string {
	var lit string
	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				// at the end of the identifier
				return lit
			}
		}

		l.pos.column++
		if unicode.IsLetter(r) {
			lit = lit + string(r)
		} else {
			l.backup()
			return lit
		}
	}
}

func lex_init()(*HashTable,*Lexer,*int){
	var erro *int = new(int)
	*erro = 1
	argsWithProg := os.Args[1]
	file, err := os.Open(argsWithProg)
	if err != nil {
		log.Fatalf("failed to open")
	}
	lexer := NewLexer(file)
	result:= Init()
	result.PrintTable()
	return result, lexer, erro
}

func lexin() ([]TOKEN, []Position){
	lex_table, lex_file, lex_erro := lex_init()
	var res []TOKEN
	var res_pos []Position
	last_t := TOKEN{"","",""}
	tok := TOKEN{"","",""}
	pos := Position{0,0}
	get_tok := ""
	for get_tok != "EOF" {
			tok, pos, last_t = lex_table.scanner(lex_file, lex_erro, last_t)
		get_tok = tok.classe
		if get_tok == "EOF" {
			break
		}
		res = append(res, tok) // Invoca a função scanner, tendo como parametro o endereço do arquivo fonte
		res_pos = append(res_pos,pos)
	}
	return res, res_pos
}