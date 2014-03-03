// Copyright (c) 2012, Wei guangjing <vcc.163@gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package driver provides a driver for database/sql.

Example:

	package main

	import (
	   "database/sql"
	   _ "github.com/waterfield/go-odbc/odbc/driver"
	   "fmt"
	)

	func main() {
	   db, err := sql.Open("odbc", "DSN=test;")
	   defer db.Close()

	   stmt, err :=	db.Prepare("select name from table")
	   defer stmt.Close()

	   rows, err :=	stmt.Query()
	   defer rows.Close()

	   for rows.Next() {
	       var name string

	       _ = rows.Scan(&name)
	       fmt.Println(name)
	   }
	}
*/
package driver

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"io"
	"odbc"
)

func init() {
	d := &Driver{}
	sql.Register("odbc", d)
}

type Driver struct {
}

func (d *Driver) Open(dsn string) (driver.Conn, error) {
	c, err := odbc.Connect(dsn)
	if err != nil {
		return nil, err
	}
	conn := &conn{c: c}
	return conn, nil
}

func (d *Driver) Close() error {
	return nil
}

type conn struct {
	c *odbc.Connection
	t *tx
}

func (c *conn) Prepare(query string) (driver.Stmt, error) {
	st, err := c.c.Prepare(query)
	if err != nil {
		return nil, err
	}

	stmt := &stmt{st: st}
	return stmt, nil
}

func (c *conn) Begin() (driver.Tx, error) {
	if err := c.c.AutoCommit(false); err != nil {
		return nil, err
	}

	return &tx{c: c}, nil
}

func (c *conn) Close() error {
	if c.c != nil {
		return c.c.Close()
	}
	return nil
}

type tx struct {
	c *conn
}

func (t *tx) Commit() error {
	err := t.c.c.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (t *tx) Rollback() error {
	err := t.c.c.Rollback()
	if err != nil {
		return err
	}
	return nil
}

type stmt struct {
	st *odbc.Statement
}

func (s *stmt) Exec(args []driver.Value) (driver.Result, error) {
	if err := s.st.Execute2(args); err != nil {
		return nil, err
	}

	rowsAffected, err := s.st.RowsAffected()
	if err != nil {
		return nil, err
	}

	r := &result{rowsAffected: int64(rowsAffected)}
	return r, nil
}

func (s *stmt) NumInput() int {
	return s.st.NumParams()
}

func (s *stmt) Query(args []driver.Value) (driver.Rows, error) {
	if err := s.st.Execute2(args); err != nil {
		return nil, err
	}
	rows := &rows{s: s}
	return rows, nil
}

func (s *stmt) Close() error {
	s.st.Close()
	return nil
}

type result struct {
	rowsAffected int64
}

func (r *result) LastInsertId() (int64, error) {
	return 0, errors.New("not supported")
}

func (r *result) RowsAffected() (int64, error) {
	return r.rowsAffected, nil
}

type rows struct {
	s *stmt
}

func (r *rows) Columns() []string {
	c, err := r.s.st.NumFields()
	if err != nil {
		return nil
	}
	columns := make([]string, c)
	for i, _ := range columns {
		f, err := r.s.st.FieldMetadata(i + 1)
		if err != nil {
			return nil
		}
		columns[i] = f.Name
	}
	return columns
}

func (r *rows) Close() error {
	err := r.s.Close()
	if err != nil {
		return err
	}
	return nil
}

func (r *rows) Next(dest []driver.Value) error {
	eof, err := r.s.st.FetchOne2(dest)
	if err != nil {
		return err
	}
	if eof {
		return io.EOF
	}
	return nil
}
