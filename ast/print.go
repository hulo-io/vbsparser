// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package ast

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/hulo-io/vbsparser/token"
)

var _ Visitor = (*printer)(nil)

type printer struct {
	output io.Writer
	ident  string
}

func (p *printer) print(a ...any) (n int, err error) {
	return fmt.Fprint(p.output, a...)
}

func (p *printer) printf(format string, a ...any) (n int, err error) {
	return fmt.Fprintf(p.output, format, a...)
}

func (p *printer) println(a ...any) (n int, err error) {
	return fmt.Fprintln(p.output, a...)
}

func (p *printer) Visit(node Node) Visitor {
	switch n := node.(type) {
	case *File:
		for _, s := range n.Stmts {
			Walk(p, s)
		}
		for _, d := range n.Decls {
			Walk(p, d)
		}

	case *DimDecl:
		p.printf("%s %s", p.ident+"Dim", exprListStr(n.List))
		if n.Colon.IsValid() {
			p.printf(": Set %s = %s", exprStr(n.Set.Lhs), exprStr(n.Set.Rhs))
		}
		p.println()

	case *ReDimDecl:
		p.print("ReDim ")
		if n.Preserve.IsValid() {
			p.print("Preserve ")
		}
		p.println(exprListStr(n.List))

	case *ClassDecl:
		if n.Mod.HasPublic() {
			p.print(p.ident + "Public ")
		}
		p.printf(p.ident+"Class %s\n", exprStr(n.Name))

		temp := p.ident
		p.ident += "  "
		for _, d := range n.Decls {
			Walk(p, d)
		}

		for _, s := range n.Stmts {
			Walk(p, s)
		}
		p.ident = temp

		p.println(p.ident + "End Class")

	case *SubDecl:
		ident := p.ident
		if n.Mod.HasPublic() {
			p.print(p.ident + "Public ")
			ident = ""
		}
		p.printf(ident+"Sub %s\n", exprStr(n.Name))

		for _, s := range n.Body.List {
			temp := p.ident
			p.ident += "  "
			Walk(p, s)
			p.ident = temp
		}

		p.println(p.ident + "End Sub")

	case *FuncDecl:
		ident := p.ident
		if n.Mod.HasPublic() {
			p.print(ident + "Public ")
			ident = ""
		}
		list := []string{}
		for _, r := range n.Recv {
			if r.TokPos.IsValid() {
				list = append(list, fmt.Sprintf("%s %s", r.Tok, r.Name.Name))
			} else {
				list = append(list, r.Name.Name)
			}
		}
		p.printf(ident+"Function %s(%s)\n", exprStr(n.Name), strings.Join(list, ", "))

		for _, s := range n.Body.List {
			temp := p.ident
			p.ident += "  "
			Walk(p, s)
			p.ident = temp
		}

		p.println(p.ident + "End Function")

	case *PropertyDecl:
		ident := p.ident
		if n.Mod.HasPublic() {
			p.print(ident + "Public ")
			ident = ""
		}
		list := []string{}
		for _, r := range n.Recv {
			if r.TokPos.IsValid() {
				list = append(list, fmt.Sprintf("%s %s", r.Tok, r.Name.Name))
			} else {
				list = append(list, r.Name.Name)
			}
		}
		p.printf(ident+"Property %s %s(%s)\n", n.Tok, exprStr(n.Name), strings.Join(list, ", "))

		for _, s := range n.Body.List {
			temp := p.ident
			p.ident += "  "
			Walk(p, s)
			p.ident = temp
		}

		p.println(p.ident + "End Property")

	case *IfStmt:
		p.printf(p.ident+"If %s Then\n", exprStr(n.Cond))

		for _, s := range n.Body.List {
			temp := p.ident
			p.ident += "  "
			Walk(p, s)
			p.ident = temp
		}

		if n.ElseIf != nil {
			for _, elif := range n.ElseIf {
				p.println(p.ident+"ElseIf", exprStr(elif.Cond), "Then")
				for _, s := range elif.Body.List {
					temp := p.ident
					p.ident += "  "
					Walk(p, s)
					p.ident = temp
				}
			}
		}

		if n.Else != nil {
			p.println(p.ident + "Else")
			for _, s := range n.Else.List {
				temp := p.ident
				p.ident += "  "
				Walk(p, s)
				p.ident = temp
			}
		}

		p.println(p.ident + "End If")

	case *ExprStmt:
		p.printf(p.ident+"%s\n", exprStr(n.X))

	case *MemberStmt:
		if n.Mod.HasPublic() {
			p.print(p.ident + "Public ")
		}
		if n.Mod.HasPrivate() {
			p.print(p.ident + "Private ")
		}
		p.println(exprStr(n.Name))

	case *AssignStmt:
		modifier := ""
		switch n.Tok {
		case token.SET:
			modifier = "Set "
		case token.CONST:
			modifier = "Const "
		}
		p.printf(p.ident+"%s%s = %s\n", modifier, exprStr(n.Lhs), exprStr(n.Rhs))

	case *ForEachStmt:
		p.printf(p.ident+"For Each %s In %s\n", exprStr(n.Elem), exprStr(n.Group))
		for _, s := range n.Body.List {
			temp := p.ident
			p.ident += "  "
			Walk(p, s)
			p.ident = temp
		}
		p.print(p.ident + "Next")
		if n.Stmt != nil {
			p.print(" ")
			Walk(p, n)
		} else {
			p.println()
		}

	case *ForNextStmt:
		p.printf(p.ident+"For %s To %s", exprStr(n.Start), exprStr(n.End_))
		if n.Step != nil {
			p.printf(" Step %s\n", exprStr(n.Step))
		}
		for _, s := range n.Body.List {
			temp := p.ident
			p.ident += "  "
			Walk(p, s)
			p.ident = temp
		}
		p.println(p.ident + "Next")

	case *WhileWendStmt:
		p.printf(p.ident+"While %s\n", exprStr(n.Cond))
		for _, s := range n.Body.List {
			temp := p.ident
			p.ident += "  "
			Walk(p, s)
			p.ident = temp
		}
		p.println(p.ident + "Wend")

	case *DoLoopStmt:

	case *CallStmt:
		p.printf(p.ident+"Call %s %s\n", exprStr(n.Name), exprListStr(n.Recv))

	case *ExitStmt:
		if len(n.X) == 0 {
			p.println(p.ident + "Exit")
		} else {
			p.printf(p.ident+"Exit %s\n", n.X)
		}

	case *SelectStmt:
		p.printf(p.ident+"Select Case %s\n", exprStr(n.Var))
		for _, c := range n.Cases {
			p.printf(p.ident+"  Case %s\n", exprStr(c.Cond))
			for _, s := range c.Body.List {
				temp := p.ident
				p.ident += "    "
				Walk(p, s)
				p.ident = temp
			}
		}
		if n.Else != nil {
			p.println(p.ident + "  Case Else")
			for _, s := range n.Else.Body.List {
				temp := p.ident
				p.ident += "    "
				Walk(p, s)
				p.ident = temp
			}
		}
		p.println(p.ident + "End Select")

	case *WithStmt:
		p.printf(p.ident+"With %s\n", exprStr(n.Cond))
		for _, s := range n.Body.List {
			temp := p.ident
			p.ident += "  "
			Walk(p, s)
			p.ident = temp
		}
		p.println(p.ident + "End With")

	case *OnErrorStmt:
		if n.OnErrorGoto != nil {
			p.println(p.ident + "On Error GoTo 0")
		}
		if n.OnErrorResume != nil {
			p.println(p.ident + "On Error Resume Next")
		}
	case *StopStmt:
		p.println(p.ident + "Stop")
	case *RandomizeStmt:
		p.println(p.ident + "Randomize")
	case *OptionStmt:
		p.println(p.ident + "Option Explicit")
	}
	return nil
}

func Print(node Node) {
	Walk(&printer{ident: "", output: os.Stdout}, node)
}

func String(node Node) string {
	buf := &strings.Builder{}
	Walk(&printer{ident: "", output: buf}, node)
	return buf.String()
}

func exprStr(e Expr) string {
	switch e := e.(type) {
	case *Ident:
		return e.Name
	case *BasicLit:
		if e.Kind == token.STRING {
			return fmt.Sprintf(`"%s"`, e.Value)
		}
		return e.Value
	case *SelectorExpr:
		return fmt.Sprintf("%s.%s", exprStr(e.X), exprStr(e.Sel))
	case *BinaryExpr:
		return fmt.Sprintf("%s %s %s", exprStr(e.X), e.Op, exprStr(e.Y))
	case *CallExpr:
		return fmt.Sprintf("%s(%s)", exprStr(e.Func), exprListStr(e.Recv))
	case *IndexExpr:
		return fmt.Sprintf("%s(%s)", exprStr(e.X), exprStr(e.Index))
	case *IndexListExpr:
		return fmt.Sprintf("%s(%s)", exprStr(e.X), exprListStr(e.Indices))
	}
	return ""
}

func exprListStr(list []Expr) string {
	res := []string{}
	for _, e := range list {
		res = append(res, exprStr(e))
	}
	return strings.Join(res, ", ")
}
