package main

import (
    "fmt"
    "strings"
    "io/ioutil"
    "net/http"
    "regexp"
    "sort"
    "strconv"
    "database/sql"
   _ "github.com/lib/pq"
)

const (
     DB_USER     = "postgres"
     DB_PASSWORD = ""
     DB_NAME     = "ab_log_db"
)

func checkErr(err error) {
        if err != nil {
            panic(err)
        }
    }



func updater(r_state int) {
        dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
            DB_USER, DB_PASSWORD, DB_NAME)
        db, err := sql.Open("postgres", dbinfo)
        checkErr(err)
        defer db.Close()

        fmt.Println("# Update relay state  values")
        var r_id int
        err = db.QueryRow("Update relays  set r_state=$1, r_last_update=current_timestamp where r_id=1  returning r_id;", r_state).Scan(&r_id)
        checkErr(err)
        fmt.Println("last updated id =", r_id)

}




//отправляем запрос
func getrelaystate(host string)(rodos_htm string){
    resp, err := http.Get(host) 
    if err != nil { 
        fmt.Println(err) 
        return
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
          fmt.Println(err)
          return
    }
    //fmt.Println(string(body))
    rodos_htm = string(body)
    return rodos_htm
}

//удаляем теги
func RemoveHtmlTag(in string) string {
	// regex to match html tag
	const pattern = `(<\/?[a-zA-A]+?[^>]*\/?>)*`
	r := regexp.MustCompile(pattern)
	groups := r.FindAllString(in, -1)
	// should replace long string first
	sort.Slice(groups, func(i, j int) bool {
		return len(groups[i]) > len(groups[j])
	})
	for _, group := range groups {
		if strings.TrimSpace(group) != "" {
			in = strings.ReplaceAll(in, group, "")
		}
	}
	return in
}

func main() {
    body := getrelaystate("http://192.168.71.117/pstat.xml")
    state, err := strconv.Atoi(strings.TrimSpace(RemoveHtmlTag(body)))
    if err != nil {
          fmt.Println(err)
          return
    }
    updater(state)
    //fmt.Println(state)
}
