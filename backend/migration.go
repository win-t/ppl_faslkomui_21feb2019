package main

var applicationID = "7a69a902ea067cdd8c08134b29417913"
var migrationStatement = []string{
	`` +
		`create table keyvalue(key text primary key, value text);` +
		`insert into keyvalue(key, value) values ('counter', '0');`,
}
