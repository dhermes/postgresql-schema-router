package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/auxten/postgresql-parser/pkg/sql/parser"
	"github.com/auxten/postgresql-parser/pkg/sql/sem/tree"
)

const (
	query = `
SELECT n.nspname as "Schema",
  c.relname as "Name",
  CASE c.relkind WHEN 'r' THEN 'table' WHEN 'v' THEN 'view' WHEN 'm' THEN 'materialized view' WHEN 'i' THEN 'index' WHEN 'S' THEN 'sequence' WHEN 's' THEN 'special' WHEN 'f' THEN 'foreign table' WHEN 'p' THEN 'partitioned table' WHEN 'I' THEN 'partitioned index' END as "Type",
  pg_catalog.pg_get_userbyid(c.relowner) as "Owner"
FROM pg_catalog.pg_class c
     LEFT JOIN pg_catalog.pg_namespace n ON n.oid = c.relnamespace
WHERE c.relkind IN ('r','p','')
      AND n.nspname <> 'pg_catalog'
      AND n.nspname <> 'information_schema'
      AND n.nspname !~ '^pg_toast'
  AND pg_catalog.pg_table_is_visible(c.oid)
ORDER BY 1,2;
`
)

func run() error {
	statements, err := parser.Parse(query)
	if err != nil {
		return err
	}

	if len(statements) != 1 {
		return errors.New("expected exactly 1 statement")
	}

	statement := statements[0]
	select_, ok := statement.AST.(*tree.Select)
	if !ok || select_ == nil {
		return errors.New("expected AST to be a SELECT statement")
	}
	clause, ok := select_.Select.(*tree.SelectClause)
	if !ok || clause == nil {
		return errors.New("expected SELECT clause")
	}
	tables := clause.From.Tables
	if len(tables) != 1 {
		return errors.New("expected exactly 1 table")
	}
	table, ok := tables[0].(*tree.JoinTableExpr)
	if !ok || table == nil {
		return errors.New("expected JOIN table expression")
	}
	//
	left, ok := table.Left.(*tree.AliasedTableExpr)
	if !ok || left == nil {
		return errors.New("expected left table as aliased expression")
	}
	right, ok := table.Right.(*tree.AliasedTableExpr)
	if !ok || right == nil {
		return errors.New("expected right table as aliased expression")
	}
	//
	leftName, ok := left.Expr.(*tree.TableName)
	if !ok || leftName == nil {
		return errors.New("expected left table name")
	}
	rightName, ok := left.Expr.(*tree.TableName)
	if !ok || rightName == nil {
		return errors.New("expected left table name")
	}
	fmt.Printf(" left: %s\n", leftName.FQString())
	fmt.Printf("right: %s\n", rightName.FQString())

	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
