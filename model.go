package main

type News struct {
	Id          string
	Date        int64
	FirstName   string
	LastName    string
	Title       string
	Description string
	Body        string
	Image       bool
}

type Message struct {
	Id       string
	FullName string
	Email    string
	Phone    string
	Subject  string
	Message  string
}

type Employee struct {
	Id          string
	FirstName   string
	LastName    string
	Description string
	Title       string
	Image       bool
	Location    string
}

type NewsSorted []News

func (n NewsSorted) Len() int {
	return len(n)
}

func (n NewsSorted) Less(i, j int) bool {
	return n[i].Id > n[j].Id
}

func (n NewsSorted) Swap(i, j int) {
	n[i], n[j] = n[j], n[i]
}
