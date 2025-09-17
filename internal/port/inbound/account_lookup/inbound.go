package accountlookup

import "net/http"

type UserSearchAdapter interface {
	SearchUser(w http.ResponseWriter, r *http.Request)
}
