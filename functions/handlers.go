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
	tmpl, err := template.ParseFiles("./templates/search.html")
	tmpl1, err1 := template.ParseFiles("templates/errors.html")
	if err != nil || err1 != nil {
		if err1 != nil {
			http.Error(w, "Internal Server Error", 500)
			return
		} else {
			ChooseError(w, 500)
			tmpl1.Execute(w, Error)
			return
		}
	}
	if r.Method != http.MethodPost {
		ChooseError(w, 405)
		tmpl1.Execute(w, Error)
		return
	}
	var types, text string
	types = r.FormValue("typessearch")
	text = strings.ToLower(r.FormValue("search"))
	temp := strings.Split(text, "->")
	text = temp[0]
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

	if len(ss) == 0 {
		ChooseError(w, 1000)
		tmpl1.Execute(w, Error)
		return
	}

	tmpl.Execute(w, ss)
}

func FilterHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./templates/search.html")
	tmpl1, err1 := template.ParseFiles("templates/errors.html")
	if err != nil || err1 != nil {
		if err1 != nil {
			http.Error(w, "Internal Server Error", 500)
			return
		} else {
			ChooseError(w, 500)
			tmpl1.Execute(w, Error)
			return
		}
	}
	if r.Method != http.MethodPost {
		ChooseError(w, 405)
		tmpl1.Execute(w, Error)
		return
	}
	err = r.ParseForm()
	if err != nil {
		ChooseError(w, 500)
		tmpl1.Execute(w, Error)
		return
	}
	members_checked := r.Form["member-radio"]
	Dates_range := r.Form["creation"]
	type_filter := r.FormValue("type_filtre")
	text := strings.ToLower(r.FormValue("filter_input"))
	min_date, err1 := strconv.Atoi(Dates_range[0])
	max_date, err2 := strconv.Atoi(Dates_range[1])
	if err1 != nil || err2 != nil {
		ChooseError(w, 400)
		tmpl1.Execute(w, Error)
		return
	}
	var min_mem, max_mem, check int
	if len(members_checked) == 0 && (min_date == 1958 && max_date == 2018) && len(text) == 0 {
		ChooseError(w, 400)
		tmpl1.Execute(w, Error)
		return
	} else {
		if len(members_checked) == 0 {
			min_mem = 1
			max_mem = 8
			check = 1
		} else if len(members_checked) == 1 {
			temp, err3 := strconv.Atoi(members_checked[0])
			if err3 != nil {
				ChooseError(w, 400)
				tmpl1.Execute(w, Error)
				return
			}
			min_mem = temp
			max_mem = temp
			check = 2
		} else {
			temp1, err3 := strconv.Atoi(members_checked[0])
			temp2, err4 := strconv.Atoi(members_checked[len(members_checked)-1])
			if err3 != nil || err4 != nil {
				ChooseError(w, 400)
				tmpl1.Execute(w, Error)
				return
			}
			min_mem = temp1
			max_mem = temp2
			check = 2
		}
	}
	if len(text) != 0 {
		check = 3
	}

	var ff []Artist

	if check <= 2 && check > 0 {
		for i := range artists {
			if type_filter == "creation" {
				if artists[i].CreationDate >= min_date && artists[i].CreationDate <= max_date {
					if check == 2 {
						if len(artists[i].Members) >= min_mem && len(artists[i].Members) <= max_mem {
							ff = append(ff, artists[i])
						}
					} else {
						ff = append(ff, artists[i])
					}
				}
			} else if type_filter == "firstalbum" {
				sli := strings.Split(artists[i].FirstAlbum, "-")
				temp, err := strconv.Atoi(sli[len(sli)-1])
				if err != nil {
					ChooseError(w, 400)
					tmpl1.Execute(w, Error)
					return
				}
				if temp >= min_date && temp <= max_date {
					if check == 2 {
						if len(artists[i].Members) >= min_mem && len(artists[i].Members) <= max_mem {
							ff = append(ff, artists[i])
						}
					} else {
						ff = append(ff, artists[i])
					}
				}
			} else if type_filter == "both" {
				sli := strings.Split(artists[i].FirstAlbum, "-")
				temp, err := strconv.Atoi(sli[len(sli)-1])
				if err != nil {
					ChooseError(w, 400)
					tmpl1.Execute(w, Error)
					return
				}
				if temp >= min_date && temp <= max_date && artists[i].CreationDate >= min_date && artists[i].CreationDate <= max_date {
					if check == 2 {
						if len(artists[i].Members) >= min_mem && len(artists[i].Members) <= max_mem {
							ff = append(ff, artists[i])
						}
					} else {
						ff = append(ff, artists[i])
					}
				}
			} else {
				if len(artists[i].Members) >= min_mem && len(artists[i].Members) <= max_mem {
					ff = append(ff, artists[i])
				}
			}
		}
	} else if check == 3 {
		sli := strings.Split(text, " ")
		if len(sli) == 2 {
			text = sli[0]
		} else if len(sli) == 3 {
			text = sli[1]
		}
		var tf string
		for _, char := range text {
			if char >= 'a' && char <= 'z' {
				tf += string(char)
			}
		}
		text = tf

		for i := range locals.Index {
			for j := range locals.Index[i].Location {
				if strings.Contains(strings.ToLower(locals.Index[i].Location[j]), text) {
					fmt.Println("1")
					fmt.Println(artists[locals.Index[i].Id-1].Name)
					if artists[locals.Index[i].Id-1].CreationDate >= min_date && artists[locals.Index[i].Id-1].CreationDate <= max_date {
						fmt.Println("2")
						fmt.Println(artists[locals.Index[i].Id-1].Name)

						if len(artists[locals.Index[i].Id-1].Members) >= min_mem && len(artists[locals.Index[i].Id-1].Members) <= max_mem {
							fmt.Println("3")
							fmt.Println(artists[locals.Index[i].Id-1].Name)

							ff = append(ff, artists[locals.Index[i].Id-1])
						}
					}
				}
			}
		}

	}

	tmpl.Execute(w, ff)
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
