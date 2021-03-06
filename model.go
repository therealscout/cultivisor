package main

type SmtpTemplateData struct {
	From    string
	To      string
	Subject string
	Body    string
}

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
	Question bool
}

type Employee struct {
	Id          string
	FirstName   string
	LastName    string
	Description string
	Title       string
	Image       bool
	Location    string
	Email       string
	MainPage    bool
	Facebook    string
	Twitter     string
	Linkedin    string
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

type EmployeesSorted []Employee

func (e EmployeesSorted) Len() int {
	return len(e)
}

func (e EmployeesSorted) Less(i, j int) bool {
	return e[i].Id > e[j].Id
}

func (e EmployeesSorted) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

func (n NewsSorted) Whatever(id string) (string, string) {
	var prev, next string
	for i, news := range n {
		if news.Id == id {
			if i > 0 {
				prev = n[i-1].Id
			}
			if i < len(n)-1 {
				next = n[i+1].Id
			}
			break
		}
	}
	return prev, next
}
