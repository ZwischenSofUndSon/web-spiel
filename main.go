package main 

import (
	"fmt"
	"net/http"
	"strconv"
	"html/template"
	"github.com/joho/godotenv"
	"os"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var database *sql.DB

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Environment didnt load")
		panic(err)
	}

	key := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", 
 os.Getenv("DB_USER"),
 os.Getenv("DB_PASS"),
 os.Getenv("DB_HOST"),
 os.Getenv("DB_PORT"),
 os.Getenv("DB_NAME"),
 )

	db, err := sql.Open("mysql", key)

	if err != nil {
		panic(err)
	} 

	database = db
	defer db.Close()
	
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))


	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
		http.ServeFile(w, r, "static/pages/index.html")
	
	})

	http.HandleFunc("/statistik", statistikHandln)
	http.HandleFunc("/activities", activitiesHandln)
	http.HandleFunc("/insertNewAktivities", insertNewAktivities)



	fmt.Println("Server is listening...")
	http.ListenAndServe(":8181", nil)
}


type Player struct {
	Health int
	Damage int
	Money  int
}

func statistikHandln(w http.ResponseWriter, r *http.Request) {
	player1 := Player{}
	rows, err := database.Query("SELECT money, health, damage FROM players_stats")
	if err != nil {
		panic(err)
		} 
		defer rows.Close()
		for rows.Next() {
			error := rows.Scan(&player1.Money, &player1.Health, &player1.Damage, )
			if error != nil {
				panic(error)
			}
		}
		
		
		tmpl, _ := template.ParseFiles("templates/statistik.html")
		tmpl.Execute(w, player1)
	}




	func activitiesHandln(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/pages/activities.html")

	}

func insertNewAktivities(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}
	schritte := r.FormValue("schritte")
	schlafindex := r.FormValue("schlafindex")
	bizeps := r.FormValue("bizeps")
	kreuzheben := r.FormValue("kreuzheben")
	crossover := r.FormValue("crossover")
	bankdrucken := r.FormValue("bankdrucken")
	schwimmen := r.FormValue("schwimmen")
	sliceAktivities := make([]string, 0) 
	sliceAktivities = append(sliceAktivities, schlafindex, kreuzheben, bizeps, crossover, bankdrucken, schwimmen, schritte)
	
    
	for _, v := range sliceAktivities{
		value := mustToInt(v)
	 _, err := database.Exec("UPDATE `players`.`players_stats` SET `money` = 'money + ?' WHERE (`id` = '1');", value)
	 if err != nil {
		fmt.Println("Error, cant insert values into database")
		panic(err)
	}
	}
	http.Redirect(w, r, "/statistik", http.StatusMovedPermanently)
}


func mustToInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		fmt.Println("Error, cant convert Sting to Int(Aktivities)")
		panic(err)
	}
	return i
}
