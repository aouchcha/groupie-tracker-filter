package functions

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

type Err struct {
	Message string
	Title   string
	Code    int
}

var Error Err

func ServeStyle(w http.ResponseWriter, r *http.Request) {
	v := http.StripPrefix("/styles/", http.FileServer(http.Dir("./styles")))
	tmpl1, err2 := template.ParseFiles("templates/errors.html")
	if err2 != nil {
		http.Error(w, "Error 500", http.StatusInternalServerError)
		return
	}
	if r.URL.Path == "/styles/" {
		ChooseError(w, 403)
		tmpl1.Execute(w, Error)
		return
	}
	v.ServeHTTP(w, r)
}

func FirstPage(w http.ResponseWriter, r *http.Request) {

	tmpl, err := template.ParseFiles("templates/welcome.html")
	tmpl1, err2 := template.ParseFiles("templates/errors.html")

	if err != nil || err2 != nil {
		if err2 != nil {
			http.Error(w, "Error 500", http.StatusInternalServerError)
			return
		} else {
			ChooseError(w, 500)
			tmpl1.Execute(w, Error)
			return
		}
	}
	if r.URL.Path != "/" {
		ChooseError(w, 404)
		tmpl1.Execute(w, Error)
		return
	}
	if r.Method != http.MethodGet {
		ChooseError(w, 405)
		tmpl1.Execute(w, Error)
		return
	}
	tmpl.Execute(w, artists)
}

func SuggestHandler(w http.ResponseWriter, r *http.Request) {
	input := r.URL.Query().Get("userinput")

	suggestions := getSuggestions(input)
	w.Header().Set("Content-Type", "text/plain")
	for _, item := range suggestions {
		w.Write([]byte(item + "\n"))
	}
}

func getSuggestions(input string) []string {
	var suggestions []string
	input = strings.ToLower(input)
	for i := range artists {
		if strings.HasPrefix(strings.ToLower(artists[i].Name), input) {
			suggestions = append(suggestions, artists[i].Name+"-> Band")
		}
		if strings.HasPrefix(strings.ToLower(artists[i].FirstAlbum), input) {
			suggestions = append(suggestions, artists[i].FirstAlbum+"-> First Album Date")
		}
		if strings.HasPrefix(strings.ToLower(strconv.Itoa(artists[i].CreationDate)), input) {
			suggestions = append(suggestions, strconv.Itoa(artists[i].CreationDate)+"-> Creation Date")
		}
		for j := range artists[i].Members {
			if strings.HasPrefix(strings.ToLower(artists[i].Members[j]), input) {
				suggestions = append(suggestions, artists[i].Members[j]+"->Member")
				break
			}
		}
	}
	for i := range locals.Index {
		for j := range locals.Index[i].Location {
			if strings.Contains(strings.ToLower(locals.Index[i].Location[j]), input) {
				suggestions = append(suggestions, locals.Index[i].Location[j]+"->Location")
			}
		}
	}
	if suggestions == nil {
		suggestions = append(suggestions, "There is no data like that")
	}
	return suggestions
}

func OtherPages(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/details.html")
	tmpl1, err2 := template.ParseFiles("templates/errors.html")

	if err != nil || err2 != nil {
		if err2 != nil {
			http.Error(w, "Error 500", http.StatusInternalServerError)
			return
		} else {
			ChooseError(w, 500)
			tmpl1.Execute(w, Error)
			return
		}
	}

	if r.URL.Path != "/artist" {
		ChooseError(w, 404)
		tmpl1.Execute(w, Error)
		return
	}
	max := artists[len(artists)-1].ID
	url := r.URL.Query().Get("ID")
	index, err := strconv.Atoi(string(url))
	if err != nil || index < 1 || index > max {
		ChooseError(w, 404)
		tmpl1.Execute(w, Error)
		return
	}
	index -= 1
	if r.Method != http.MethodGet {
		ChooseError(w, 405)
		tmpl1.Execute(w, Error)
		return
	}

	artistinfos := struct {
		ID            int
		Name          string
		Image         string
		Members       []string
		CreationDate  int
		FirstAlbum    string
		Localisations []string
		Relations     map[string][]string
		Dates         []string
	}{
		ID:            artists[index].ID,
		Name:          artists[index].Name,
		Image:         artists[index].Image,
		Members:       artists[index].Members,
		CreationDate:  artists[index].CreationDate,
		FirstAlbum:    artists[index].FirstAlbum,
		Localisations: locals.Index[index].Location,
		Relations:     rel.Index[index].DateLocations,
		Dates:         dat.Index[index].Date,
	}
	tmpl.Execute(w, artistinfos)
}

func SearchPage(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("./templates/search.html")
	tmpl1, _ := template.ParseFiles("templates/errors.html")
	if r.Method != http.MethodPost {
		ChooseError(w, 405)
		tmpl1.Execute(w, Error)
		return
	}
	project := r.FormValue("project")
	var types, years_range, text string
	var membernumber []string
	var min, max int
	var check = false
	// fmt.Println(project)
	types = r.FormValue("typessearch")
	if project == "search" {
		text = strings.ToLower(r.FormValue("search"))
		temp := strings.Split(text, "->")
		text = temp[0]
	} else if project == "filtre" {
		check = true
		years_range = r.FormValue("years-range")
		membernumber = r.Form["member-radio"]
		if len(membernumber) == 0 && years_range == "1012" {
			ChooseError(w, 400)
			tmpl1.Execute(w, Error)
			return
		} else if len(membernumber) != 0 && years_range == "1012" {
			if len(membernumber) == 1 {
				max, _ = strconv.Atoi(membernumber[0])
				min = max
			} else {
				min, _ = strconv.Atoi(membernumber[0])
				max, _ = strconv.Atoi(membernumber[len(membernumber)-1])
			}
			years_range = "1000"
			types = "Band"
		} else if len(membernumber) == 0 && years_range != "1012" {
			min = 1
			max = 8
		} else {
			if len(membernumber) == 1 {
				max, _ = strconv.Atoi(membernumber[0])
				min = max
			} else {
				min, _ = strconv.Atoi(membernumber[0])
				max, _ = strconv.Atoi(membernumber[len(membernumber)-1])
			}
		}
		text = years_range
	}

	fmt.Println(min, max)
	fmt.Println(check)
	fmt.Println(years_range)
	fmt.Println(membernumber)
	// if err != nil {
	// 	ChooseError(w, 400)
	// 	tmpl1.Execute(w, Error)
	// 	return
	// }
	// fmt.Println(years_range)
	// fmt.Println(check)
	// fmt.Println(doublecheck)
	fmt.Println(types)
	fmt.Println(text)

	if text == "" {
		fmt.Println("hna")
		ChooseError(w, 400)
		tmpl1.Execute(w, Error)
		return
	}

	var ss []Artist
	fmt.Println(ss)
	if check {
		for i := range artists {
			if len(artists[i].Members) >= min && len(artists[i].Members) <= max {
				fmt.Println(len(artists[i].Members))
				if types == "firstalbum" {
					temp := strings.Split(artists[i].FirstAlbum, "-")
					val, _ := strconv.Atoi(temp[len(temp)-1])
					val2, _ := strconv.Atoi(text)
					if val >= val2 {
						ss = append(ss, artists[i])
					}
				} else if types == "creation" {
					val, _ := strconv.Atoi(text)
					if val <= artists[i].CreationDate {
						ss = append(ss, artists[i])
					}
				} else {
					ss = append(ss, artists[i])
				}
			}
		}
	} else {
		if types == "Band" {
			for i := range artists {
				if strings.HasPrefix(strings.ToLower(artists[i].Name), text) {
					ss = append(ss, artists[i])
				}
			}
		} else if types == "firstalbum" || types == "creation" {
			if types == "firstalbum" {
				for i := range artists {
					temp := strings.Split(artists[i].FirstAlbum, "-")
					val, _ := strconv.Atoi(temp[len(temp)-1])
					val2, _ := strconv.Atoi(text)

					fmt.Println("first album check false")

					if val >= val2 {
						ss = append(ss, artists[i])
					}

				}
			}
			if types == "creation" {
				val, _ := strconv.Atoi(text)
				for i := range artists {

					fmt.Println(val)
					if val <= artists[i].CreationDate {
						ss = append(ss, artists[i])
					}

				}

			}
		} else if types == "Members" {
			for i := range artists {

				for j := range artists[i].Members {
					if strings.HasPrefix(strings.ToLower(artists[i].Members[j]), text) {
						ss = append(ss, artists[i])
					}
				}

			}
		} else if types == "location" {
			for i := range locals.Index {
				for j := range locals.Index[i].Location {
					if strings.Contains(strings.ToLower(locals.Index[i].Location[j]), text) {
						if len(ss) == 0 {
							ss = append(ss, artists[locals.Index[i].Id-1])
						} else {
							var checkrepitition bool
							for k := range ss {
								if ss[k].ID == locals.Index[i].Id {
									checkrepitition = true
								} else {
									checkrepitition = false
								}
							}
							if !checkrepitition {
								ss = append(ss, artists[locals.Index[i].Id-1])
							}
						}

					}
				}
			}
		}
	}

	if len(ss) == 0 {
		ChooseError(w, 1000)
		tmpl1.Execute(w, Error)
		return
	}
	fmt.Println(ss)

	tmpl.Execute(w, ss)
}

func ChooseError(w http.ResponseWriter, code int) {
	if code == 404 || code == 0 {
		Error.Title = "Error 404"
		Error.Message = "The page web doesn't exist\nError 404"
		Error.Code = 404
		w.WriteHeader(404)
	} else if code == 405 {
		Error.Title = "Error 405"
		Error.Message = "The method is not alloweded\nError 405"
		Error.Code = code
		w.WriteHeader(code)
	} else if code == 400 {
		Error.Title = "Error 400"
		Error.Message = "Bad Request\nError 400"
		Error.Code = code
		w.WriteHeader(code)
	} else if code == 500 {
		Error.Title = "Error 500"
		Error.Message = "Internal Server Error\nError 500"
		Error.Code = code
		w.WriteHeader(code)
	} else if code == 403 {
		Error.Title = "Error 403"
		Error.Message = "This page web is forbidden\nError 403"
		Error.Code = code
		w.WriteHeader(code)
	}
}
