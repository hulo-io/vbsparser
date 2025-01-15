// Copyright 2025 The Hulo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package ast_test

import (
	"testing"

	"github.com/hulo-io/vbsparser/ast"
	"github.com/hulo-io/vbsparser/token"
	"github.com/stretchr/testify/assert"
)

func TestAST(t *testing.T) {
	testset := []struct {
		node     ast.Node
		expected string
	}{
		{
			&ast.File{
				Decls: []ast.Decl{
					&ast.DimDecl{Vars: []ast.Expr{&ast.Ident{Name: "A"}}},
				},
				Stmts: []ast.Stmt{
					&ast.AssignStmt{
						Lhs: &ast.Ident{Name: "A"},
						Rhs: &ast.CallExpr{
							Func: &ast.Ident{Name: "Array"},
							Recv: []ast.Expr{&ast.Ident{Name: "10"}, &ast.Ident{Name: "20"}, &ast.Ident{Name: "30"}},
						},
					},
				},
			}, `Dim A
A = Array(10,20,30)`},
		{&ast.File{Decls: []ast.Decl{
			&ast.DimDecl{Vars: []ast.Expr{&ast.IndexExpr{X: &ast.Ident{Name: "Names"}, Index: &ast.Ident{Name: "9"}}}},
			&ast.DimDecl{Vars: []ast.Expr{&ast.IndexListExpr{X: &ast.Ident{Name: "Names"}, Indices: []ast.Expr{&ast.Ident{Name: "10"}, &ast.Ident{Name: "10"}, &ast.Ident{Name: "10"}}}}},
			&ast.DimDecl{Vars: []ast.Expr{&ast.Ident{Name: "MyVar"}, &ast.Ident{Name: "MyNum"}}}}}, `Dim Names(9)
Dim Names(10, 10, 10)
Dim MyVar, MyNum`},
		{&ast.BlockStmt{List: []ast.Stmt{
			&ast.OnErrorStmt{OnErrorResume: &ast.OnErrorResume{}},
			&ast.ExprStmt{
				X: &ast.CallExpr{
					Func: &ast.SelectorExpr{X: &ast.Ident{Name: "Err"}, Sel: &ast.Ident{Name: "Raise"}},
					Recv: []ast.Expr{&ast.BasicLit{Kind: token.INTEGER, Value: "6"}},
				}},
			&ast.ExprStmt{
				Doc: &ast.CommentGroup{List: []*ast.Comment{{Tok: token.APOSTROPHE, Text: "Clear the error"}}},
				X: &ast.CallExpr{
					Func: &ast.Ident{Name: "MsgBox"},
					Recv: []ast.Expr{&ast.BinaryExpr{
						X:  &ast.Ident{Name: `"Error # "`},
						Op: token.BITAND,
						Y: &ast.BinaryExpr{
							X: &ast.CallExpr{
								Func: &ast.Ident{Name: "CStr"},
								Recv: []ast.Expr{&ast.SelectorExpr{X: &ast.Ident{Name: "Err"}, Sel: &ast.Ident{Name: "Number"}}},
							},
							Op: token.BITAND,
							Y: &ast.BinaryExpr{
								X:  &ast.BasicLit{Kind: token.STRING, Value: " "},
								Op: token.BITAND,
								Y: &ast.SelectorExpr{
									X:   &ast.Ident{Name: "Err"},
									Sel: &ast.Ident{Name: "Description"},
								},
							},
						}}}}},
		}}, `On Error Resume Next
Err.Raise 6   ' Raise an overflow error.
MsgBox ("Error # " & CStr(Err.Number) & " " & Err.Description)
Err.Clear      ' Clear the error`},
	}
	for _, tt := range testset {
		assert.Equal(t, tt.expected, tt.node)
	}
}
