package rollbar

import (
	"encoding/json"
	"strconv"
)

// Invite represents nested invites from the ListInvitesResponse.
type Invite struct {
	ID           int    `json:"id"`
	FromUserID   int    `json:"from_user_id"`
	TeamID       int    `json:"team_id"`
	ToEmail      string `json:"to_email"`
	Status       string `json:"status"`
	DateCreated  int    `json:"date_created"`
	DateRedeemed int    `json:"date_redeemed"`
}

// ListInvitesResponse represents the ListInvites response.
type ListInvitesResponse struct {
	Error  int `json:"err"`
	Result []Invite
}

// ListInvites lists all the invites.
func (c *Client) ListInvites(teamID int) ([]Invite, error) {
	var invites []Invite
	// Invitation call has pagination.
	// There's a feature request to expire the invitations after some time.
	// Looping until we get an empty invitations list [].
	// Page=0 and page=1 return the same result.
	for i := 1; ; i++ {
		var data ListInvitesResponse
		var dataResponses []Invite

		pageNum := i
		queryParams := map[string]string{"page": strconv.Itoa(pageNum)}

		bytes, err := c.getWithQueryParams(queryParams, "team", strconv.Itoa(teamID), "invites")
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(bytes, &data)
		dataResponses = data.Result
		invites = append(invites, dataResponses...)

		if err != nil {
			return nil, err
		}

		if len(data.Result) == 0 {
			break
		}
	}
	return invites, nil
}
