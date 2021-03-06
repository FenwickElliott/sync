package sync

import (
	"encoding/json"
	"net/http"
	"strings"

	"gopkg.in/mgo.v2/bson"
)

func print(w http.ResponseWriter, r *http.Request) {
	// TODO: redirects (also functionalize them)
	var all []bson.M
	err = c.Find(nil).All(&all)
	check(err)

	err := json.NewEncoder(w).Encode(all)
	check(err)
}

func in(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	partner := r.FormValue("partner")
	partnerCookie := r.FormValue("cookie")
	nativeCookie := getOrSetCookie(&w, r)

	err = insert(nativeCookie.Value, partner, partnerCookie)
	check(err)

	if service.Redirect != "" {
		var res bson.M
		c.FindId(nativeCookie.Value).One(&res)

		str := service.Redirect + "/forward?"
		for k, v := range res {
			str += k + "=" + v.(string) + "&"
		}
		str += "back=" + service.Address
		str = strings.Replace(str, "_id", service.Name, -1)
		http.Redirect(w, r, str, 307)
	}
}

func out(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	partner := r.FormValue("partner")
	partnerCookie := r.FormValue("cookie")
	if partner == service.Name {
		partner = "_id"
	}

	var res bson.M
	err = c.Find(bson.M{partner: partnerCookie}).One(&res)
	if err != nil && err.Error() == "not found" {
		http.Error(w, "Cookie not found", 404)
		return
	}
	check(err)

	err := json.NewEncoder(w).Encode(res)
	check(err)
}

func forward(w http.ResponseWriter, r *http.Request) {
	nativeCookie := getOrSetCookie(&w, r)

	for k, v := range r.Form {
		err = insert(nativeCookie.Value, k, v[0])
		check(err)
	}

	str := r.FormValue("back") + "/back?partner=" + service.Name + "&cookie=" + nativeCookie.Value
	http.Redirect(w, r, str, 307)
}

func back(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	partner := r.FormValue("partner")
	partnerCookie := r.FormValue("cookie")

	nativeCookie := getOrSetCookie(&w, r)
	err = insert(nativeCookie.Value, partner, partnerCookie)
	check(err)
}
