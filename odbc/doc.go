// Copyright (c) 2011, Wei guangjing <vcc.163@gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package odbc provides a ODBC database driver.

Example:

	package main

	import (
		"github.com/kch42/go-odbc/odbc"
	)

	func main() {
		conn, _ := odbc.Connect("DSN=dsn;UID=user;PWD=password")
		stmt, _ := conn.Prepare("select * from user where username = ?")
		stmt.Execute("admin")
		rows, _ := stmt.FetchAll()
		for i, row := range rows {
			println(i, row)
		}
		stmt.Close()
		conn.Close()
	}
*/
package odbc
