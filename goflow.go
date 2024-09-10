package goflow

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
)

var uml *bytes.Buffer

func GetUML() string {
	return uml.String()
}

func AnalyzeFunction(fileName string, funcName string) error {
	uml = &bytes.Buffer{}
	uml.WriteString("@startuml\n")
	uml.WriteString(fmt.Sprintf("start\n:%s;\n", funcName))

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, fileName, nil, parser.AllErrors)
	if err != nil {
		return err
	}

	ast.Inspect(node, func(n ast.Node) bool {
		if fn, ok := n.(*ast.FuncDecl); ok && fn.Name.Name == funcName {
			analyzeStatements(fn.Body.List)
		}
		return true
	})
	uml.WriteString("@enduml\n")
	return nil
}

// TODO select case, go routine, channel, continue
func analyzeStatements(stmts []ast.Stmt) {
	for _, stmt := range stmts {
		switch s := stmt.(type) {
		case *ast.IfStmt:
			if s.Init != nil {
				analyzeIfInitStmt(s.Init)
			}
			analyzeIfStmt(s)
		case *ast.ExprStmt:
			analyzeExprStmt(s.X)
		case *ast.AssignStmt:
			analyzeAssignStmt(s)
		case *ast.DeclStmt:
			analyzeDeclStmt(s)
		case *ast.ReturnStmt:
			analyzeReturnStmt(s)
		case *ast.RangeStmt:
			analyzeRangeStmt(s)
		case *ast.SwitchStmt:
			analyzeSwitchStmt(s)
		case *ast.ForStmt:
			analyzeForStmt(s)
		case *ast.BranchStmt:
			analyzeBranchStmt(s)
		case *ast.IncDecStmt:
			analyzeIncDecStmt(s)
		case *ast.DeferStmt:
			analyzeDeferStmt(s)
		default:
			log.Printf("Unsupported statement: %T: %+v pos:%d", s, s, s.Pos())
		}
	}
}

func analyzeDeclStmt(declStmt *ast.DeclStmt) {
	switch decl := declStmt.Decl.(type) {
	case *ast.GenDecl:
		for _, spec := range decl.Specs {
			switch spec := spec.(type) {
			case *ast.ValueSpec:
				for _, name := range spec.Names {
					uml.WriteString(":")
					uml.WriteString(name.Name)
					if spec.Type != nil {
						uml.WriteString(" ")
						analyzeExpr(spec.Type)
					}
					if len(spec.Values) > 0 && spec.Values[0] != nil {
						uml.WriteString(" = ")
						analyzeExpr(spec.Values[0])
					}
					uml.WriteString(";\n")
				}
			default:
				log.Printf("Unsupported spec type: %T, pos:%d", spec, spec.Pos())
			}
		}
	default:
		log.Printf("Unsupported decl type: %T, pos:%d", decl, decl.Pos())
	}
}

func analyzeIfStmt(ifStmt *ast.IfStmt) {
	uml.WriteString("if (")
	analyzeExpr(ifStmt.Cond)
	uml.WriteString(") then (yes)\n")

	analyzeStatements(ifStmt.Body.List)

	if ifStmt.Else != nil {
		switch elseStmt := ifStmt.Else.(type) {
		case *ast.IfStmt:
			uml.WriteString("else\n")
			analyzeIfStmt(elseStmt)
		case *ast.BlockStmt:
			uml.WriteString("else\n")
			analyzeStatements(elseStmt.List)
		default:
			log.Printf("Unsupported else statement: %T", elseStmt)
		}
	}

	uml.WriteString("endif\n")
}

func analyzeIfCondition(cond ast.Expr) {
	switch c := cond.(type) {
	case *ast.BinaryExpr:
		analyzeExpr(c.X)
		uml.WriteString(fmt.Sprintf(" %s ", c.Op))
		analyzeExpr(c.Y)
	case *ast.Ident:
		uml.WriteString(c.Name)
	default:
		log.Printf("Unsupported condition type: %T", cond)
	}
}

func analyzeElseStmt(elseStmt ast.Stmt) {
	if block, ok := elseStmt.(*ast.BlockStmt); ok {
		analyzeStatements(block.List)
	} else if ifStmt, ok := elseStmt.(*ast.IfStmt); ok {
		analyzeIfStmt(ifStmt)
	}
}

func analyzeAssignStmt(assignStmt *ast.AssignStmt) {
	uml.WriteString(":")
	for i, lhs := range assignStmt.Lhs {
		analyzeExpr(lhs)
		if i < len(assignStmt.Lhs)-1 {
			uml.WriteString(", ")
		}
	}
	uml.WriteString(" ")
	uml.WriteString(assignStmt.Tok.String())
	uml.WriteString(" ")
	for i, rhs := range assignStmt.Rhs {
		analyzeExpr(rhs)
		if i < len(assignStmt.Rhs)-1 {
			uml.WriteString(", ")
		}
	}
	uml.WriteString(";\n")
}

func analyzeReturnStmt(returnStmt *ast.ReturnStmt) {
	uml.WriteString(":return ")
	for i, result := range returnStmt.Results {
		if i > 0 {
			uml.WriteString(", ")
		}
		analyzeExpr(result)
	}
	uml.WriteString(";\nend\n")
}

func analyzeSwitchStmt(switchStmt *ast.SwitchStmt) {
	uml.WriteString("switch (")
	if switchStmt.Tag != nil {
		analyzeExpr(switchStmt.Tag)
	}
	uml.WriteString(")\n")

	var foundDefault bool
	for _, stmt := range switchStmt.Body.List {
		caseClause, ok := stmt.(*ast.CaseClause)
		if !ok {
			continue
		}

		if len(caseClause.List) == 0 {
			foundDefault = true
			uml.WriteString("case (default)\n")
		} else {
			for _, expr := range caseClause.List {
				uml.WriteString("case (")
				analyzeExpr(expr)
				uml.WriteString(")\n")
			}
		}
		analyzeStatements(caseClause.Body)
	}
	if !foundDefault {
		uml.WriteString("case ()\n")
	}

	uml.WriteString("endswitch\n")
}

func analyzeForStmt(forStmt *ast.ForStmt) {
	if forStmt.Init != nil || forStmt.Cond != nil || forStmt.Post != nil {
		uml.WriteString("while (for ")
		if forStmt.Init != nil {
			analyzeStmtForCond(forStmt.Init)
			uml.WriteString("; ")
		} else {
			uml.WriteString("; ")
		}
		if forStmt.Cond != nil {
			analyzeExpr(forStmt.Cond)
			uml.WriteString("; ")
		} else {
			uml.WriteString("; ")
		}
		if forStmt.Post != nil {
			analyzeStmtForCond(forStmt.Post)
		}
		uml.WriteString(")\n")
	} else {
		uml.WriteString("while (for)\n")
	}

	analyzeStatements(forStmt.Body.List)

	uml.WriteString("endwhile\n")
}

func analyzeStmtForCond(stmt ast.Stmt) {
	switch s := stmt.(type) {
	case *ast.AssignStmt:
		for i, lhs := range s.Lhs {
			analyzeExpr(lhs)
			if i < len(s.Lhs)-1 {
				uml.WriteString(", ")
			}
		}
		uml.WriteString(" ")
		uml.WriteString(s.Tok.String())
		uml.WriteString(" ")
		for i, rhs := range s.Rhs {
			analyzeExpr(rhs)
			if i < len(s.Rhs)-1 {
				uml.WriteString(", ")
			}
		}
	case *ast.ExprStmt:
		analyzeExpr(s.X)
	case *ast.IncDecStmt:
		analyzeExpr(s.X)
		uml.WriteString(s.Tok.String())
	default:
		log.Printf("Unsupported statement: %T", s)
	}
}

func analyzeIncDecStmt(incDecStmt *ast.IncDecStmt) {
	uml.WriteString(":")
	analyzeExpr(incDecStmt.X)
	uml.WriteString(incDecStmt.Tok.String())
	uml.WriteString(";\n")
}

func analyzeRangeStmt(rangeStmt *ast.RangeStmt) {
	uml.WriteString("while (range ")
	analyzeExpr(rangeStmt.X)
	uml.WriteString(")\n")

	analyzeStatements(rangeStmt.Body.List)

	uml.WriteString("endwhile\n")
}

func analyzeExprStmt(expr ast.Expr) {
	switch e := expr.(type) {
	case *ast.CallExpr:
		uml.WriteString(":")
		analyzeCallExpr(e)
		uml.WriteString(";\n")
	case *ast.SelectorExpr:
		analyzeSelectorExpr(e)
	default:
		log.Printf("Unsupported expression in ExprStmt: %T", expr)
	}
}

func analyzeCallExpr(callExpr *ast.CallExpr) {
	switch fun := callExpr.Fun.(type) {
	case *ast.Ident:
		uml.WriteString(fun.Name + "(")
		for i, arg := range callExpr.Args {
			if i > 0 {
				uml.WriteString(", ")
			}
			analyzeExpr(arg)
		}
		uml.WriteString(")")
	case *ast.SelectorExpr:
		analyzeSelectorExpr(fun)
		uml.WriteString("(")
		for i, arg := range callExpr.Args {
			if i > 0 {
				uml.WriteString(", ")
			}
			analyzeExpr(arg)
		}
		uml.WriteString(")")
	default:
		log.Printf("Unsupported function type: %T", callExpr.Fun)
	}
}

func analyzeSelectorExpr(sel *ast.SelectorExpr) {
	switch x := sel.X.(type) {
	case *ast.Ident:
		uml.WriteString(x.Name + "." + sel.Sel.Name)
	case *ast.SelectorExpr:
		analyzeSelectorExpr(x)
		uml.WriteString("." + sel.Sel.Name)
	case *ast.CallExpr:
		analyzeCallExpr(x)
		uml.WriteString("." + sel.Sel.Name)
	default:
		log.Printf("Unsupported selector type: %T", sel.X)
	}
}

func analyzeIfInitStmt(initStmt ast.Stmt) {
	if assignStmt, ok := initStmt.(*ast.AssignStmt); ok {
		analyzeAssignStmt(assignStmt)
	}
}

func analyzeTypeAssertExpr(typeAssert *ast.TypeAssertExpr) {
	analyzeExpr(typeAssert.X)
	uml.WriteString(".(")
	analyzeExpr(typeAssert.Type)
	uml.WriteString(")")
}

func analyzeBranchStmt(branchStmt *ast.BranchStmt) {
	switch branchStmt.Tok {
	case token.CONTINUE:
		uml.WriteString(":continue;\nstop\n") // Use `stop` to represent `continue`
	case token.BREAK:
		uml.WriteString(":break;\nbreak\n")
	default:
		log.Printf("Unsupported branch statement: %s", branchStmt.Tok)
	}
}

func analyzeStarExpr(starExpr *ast.StarExpr) {
	uml.WriteString("*")
	analyzeExpr(starExpr.X)
}

func analyzeArrayType(arrayType *ast.ArrayType) {
	uml.WriteString("[]")
	analyzeExpr(arrayType.Elt)
}

func analyzeUnaryExpr(unaryExpr *ast.UnaryExpr) {
	uml.WriteString(unaryExpr.Op.String())
	analyzeExpr(unaryExpr.X)
}

func analyzeBinaryExpr(binaryExpr *ast.BinaryExpr) {
	analyzeExpr(binaryExpr.X)
	uml.WriteString(fmt.Sprintf(" %s ", binaryExpr.Op))
	analyzeExpr(binaryExpr.Y)
}

func analyzeDeferStmt(deferStmt *ast.DeferStmt) {
	uml.WriteString(":defer ")
	analyzeCallExpr(deferStmt.Call)
	uml.WriteString(";\n")
}

func analyzeCompositeLit(compLit *ast.CompositeLit) {
	switch t := compLit.Type.(type) {
	case *ast.Ident:
		uml.WriteString(t.Name + "{}")
	case *ast.SelectorExpr:
		uml.WriteString(t.X.(*ast.Ident).Name + "." + t.Sel.Name + "{}")
	case *ast.ArrayType:
		analyzeArrayType(t)
	default:
		log.Printf("Unsupported composite literal type: %T", t)
	}
}

func analyzeExpr(expr ast.Expr) {
	switch e := expr.(type) {
	case *ast.BasicLit:
		uml.WriteString(e.Value)
	case *ast.Ident:
		uml.WriteString(e.Name)
	case *ast.SelectorExpr:
		analyzeSelectorExpr(e)
	case *ast.BinaryExpr:
		analyzeBinaryExpr(e)
	case *ast.CallExpr:
		analyzeCallExpr(e)
	case *ast.ArrayType:
		analyzeArrayType(e)
	case *ast.UnaryExpr:
		analyzeUnaryExpr(e)
	case *ast.TypeAssertExpr:
		analyzeTypeAssertExpr(e)
	case *ast.CompositeLit:
		analyzeCompositeLit(e)
	case *ast.StarExpr:
		analyzeStarExpr(e)
	default:
		log.Printf("Unsupported expression type: %T", expr)
	}
}
