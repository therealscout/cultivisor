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
