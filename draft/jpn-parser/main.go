package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

func greet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World! %s", time.Now())
}

func main() {
	// http.HandleFunc("/", greet)
	// http.ListenAndServe(":8080", nil)
	line := sentences[0]
	lines := strings.Split(line, "\n")

	tk := []TK{}

	for _, v := range lines {
		tokens := strings.Split(strings.TrimSpace(v), " ")
		for _, t := range tokens {
			switch {
			case strings.HasPrefix(t, "a="):
				tk = append(tk, TK{Type: Property, Param1: Accent, Param2: strings.TrimPrefix(t, "a=")})
			case t == "p":
				tk = append(tk, TK{Type: BinaryOperator, Param1: Pronounce})
			case t == "v":
				tk = append(tk, TK{Type: BinaryOperator, Param1: Verb})
			default:
				tk = append(tk, TK{Type: Word, Param1: t})
			}
		}
		tk = append(tk, TK{Type: NewLine})
	}

	for _, v := range tk {
		fmt.Printf("%s %q %q\n", v.Type, v.Param1, v.Param2)
		// switch v.Type {
		// case Word:
		// 	fmt.Printf("%s %q\n", v.Type, v.Param1)
		// case Property:
		// 	fmt.Printf("%s %q %q\n", v.Type, v.Param1, v.Param2)
		// case NewLine:
		// 	fmt.Printf("%s\n", v.Type)
		// }
	}
	fmt.Println(tk)
}

type Token struct {
	Word      string
	Pronounce string
	Accent    string
	Suffix    string
	Verb      string
}

type TK struct {
	Type   string
	Param1 string
	Param2 string
}

const (
	Word           = "WORD"
	Property       = "PROPERTY"
	NewLine        = "NEWLINE"
	BinaryOperator = "BINOP"

	Accent    = "Accent"
	Pronounce = "Pronounce"
	Verb      = "Verb"
)

var sentences []string = []string{
	`さあ　
どこに a=a
試 p た し 斬 p ぎ り a=h
行 p い こう v いく a=k2
かなァ`,
	`マンション に a=a
広 #p ひろ い a=k2
駐車場 #p ちゅうしゃじょう が a=h
あります 。`,
	`友達 #p ともだち が a=h
綺麗 #p きれい な a=a
葉書 #p はがき を a=h
送 #p おく って b=おくる a=h くれました 。`,
	`それじゃ a=h
折角 p せっかく a=h
だし a=h
金魚 p きんぎょ すくい
でも
しよう a=k2
かなぁ`,
	`ビール は a=a
苦 p にが い a=k2
ので 、
余 p あま り a=a
飲 p の v のむ みません a=k1 。`,
	`まー
視聴者 p しちょうしゃ a=n2 でも
反応 p はんのう a=h が
真っ二 p まっぷた つ a=n2,n3 に
分 p わ かれる a=k3 だろうし`,
	`お 姉 p ねえ さん a=n2`,
	`取 p と り 掛 p か かります`,
	`文 p ぶん を a=a
成 p な り 立 p た たせる a=h
基本的 p きほんてき な a=h
成分 p せいぶん で a=a
ある。`,
	`こと に a=o
述語 p じゅつご は、 a=h
文 p ぶん を a=a
まとめる
重要 p じゅうよう な a=h
役割 p やくわり を a=o
果 たす。a=k2`,
	`自分 p じぶん の a=h
国 p く に a=o
について a=n2
作文 p さくぶん を a=h
書 p か きましょう v かくa=k1 。`,
}

// 自分[じぶん;h] の 国[くに;h] について 作文[さくぶん;h] を 書[か,かく;k1]きましょう 。
