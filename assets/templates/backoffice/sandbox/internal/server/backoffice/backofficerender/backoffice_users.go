package backofficerender

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
	"{{.Module}}/sandbox/internal/databases/backoffice_db"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficeauth"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficeusers"
)

// BackofficeUsersPage is what backoffice/backoffice_users.html is rendered with.
type BackofficeUsersPage struct {
	Viewer Viewer
	// Notice is shown above the list, its Text "" for none.
	Notice Notice
	Users  []BackofficeUserRow
	// Search, Role and Limit refill the filter form, and every page link
	// carries them.
	Search string
	Role   string
	Limit  int
	// Limits are the page sizes the filter form offers.
	Limits []LimitOption
	Total  int
	Page   int
	Pages  int
	// From and To number the first and the last user shown, counted from 1;
	// both 0 when the page is empty.
	From     int
	To       int
	HasPrev  bool
	PrevPage int
	HasNext  bool
	NextPage int
}

// BackofficeUserRow is one user of the list.
type BackofficeUserRow struct {
	Id       string
	Username string
	Email    string
	Role     string
	IsRoot   bool
	// IsSelf marks the signed-in user, who may not remove themselves.
	IsSelf bool
}

// Notice is the outcome of the last add, edit or remove.
type Notice struct {
	Text string
	// Kind is "ok" or "error", and styles the banner.
	Kind string
}

// LimitOption is one page size of the filter form.
type LimitOption struct {
	Value    int
	Selected bool
}

// pageSizes are the page sizes the filter form always offers.
var pageSizes = []int{10, 20, 50, 100}

// noticeOf is how the list page words a backofficeusers notice code; an
// unknown code shows nothing.
func noticeOf(sandbox *api.Sandbox, code string) Notice {
	switch code {
	case backofficeusers.NoticeAdded:
		return Notice{Text: "User added.", Kind: "ok"}
	case backofficeusers.NoticeUpdated:
		return Notice{Text: "User updated.", Kind: "ok"}
	case backofficeusers.NoticePasswordChanged:
		return Notice{Text: "User updated. Their sessions were ended and their API tokens revoked.", Kind: "ok"}
	case backofficeusers.NoticeRemoved:
		return Notice{Text: "User removed.", Kind: "ok"}
	case backofficeusers.NoticeNotFound:
		return Notice{Text: "That user no longer exists.", Kind: "error"}
	case backofficeusers.NoticeSelf:
		return Notice{Text: "You cannot remove your own account.", Kind: "error"}
	case backofficeusers.NoticeNotRoot:
		return Notice{Text: "You are no longer a root user, so nothing was removed.", Kind: "error"}
	}
	return Notice{}
}

// limitOptions are the page sizes offered, the one shown included even when
// it is not a standard one.
func limitOptions(sandbox *api.Sandbox, limit int) []LimitOption {
	options := []LimitOption{}
	listed := false
	for _, size := range pageSizes {
		if !listed && limit < size {
			options = append(options, LimitOption{Value: limit, Selected: true})
			listed = true
		}
		if size == limit {
			listed = true
		}
		options = append(options, LimitOption{Value: size, Selected: size == limit})
	}
	if !listed {
		options = append(options, LimitOption{Value: limit, Selected: true})
	}
	return options
}

// RenderUsersPage answers the backoffice user list page for user: listing,
// with the notice code notice above it.
func RenderUsersPage(sandbox *api.Sandbox, response *serverdeps.Response, user *backoffice_db.BackofficeUserRecord, listing backofficeusers.Listing, notice string) error {
	rows := []BackofficeUserRow{}
	for _, item := range listing.Users {
		rows = append(rows, BackofficeUserRow{
			Id:       sandbox.Deps.StringsDeps.FormatInt(item.Id, 10),
			Username: item.Username,
			Email:    item.Email,
			Role:     backofficeauth.RoleName(sandbox, backofficeauth.Role(item.Role)),
			IsRoot:   backofficeauth.Role(item.Role) == backofficeauth.RoleRoot,
			IsSelf:   item.Id == user.Id,
		})
	}

	page := BackofficeUsersPage{
		Viewer:   viewerOf(sandbox, user),
		Notice:   noticeOf(sandbox, notice),
		Users:    rows,
		Search:   listing.Search,
		Role:     listing.Role,
		Limit:    listing.Limit,
		Limits:   limitOptions(sandbox, listing.Limit),
		Total:    listing.Total,
		Page:     listing.Page,
		Pages:    listing.Pages,
		HasPrev:  listing.Page > 1,
		PrevPage: listing.Page - 1,
		HasNext:  listing.Page < listing.Pages,
		NextPage: listing.Page + 1,
	}
	if len(rows) > 0 {
		page.From = (listing.Page-1)*listing.Limit + 1
		page.To = page.From + len(rows) - 1
	}
	return RenderHTML(sandbox, response, api.StatusOK, "backoffice/backoffice_users.html", page)
}
