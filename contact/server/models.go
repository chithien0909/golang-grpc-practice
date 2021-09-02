package main

type ContactInfo struct {
	ID 				int 	`orm:"auto"`
	PhoneNumber 	string 	`orm:"size(15)"`
	Name 			string
}