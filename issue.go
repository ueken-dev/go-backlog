package backlog

import (
	"encoding/json"
	"strconv"
)

func validateIssueIDOrKey(issueIDOrKey string) error {
	if issueIDOrKey == "" {
		return newValidationError("issueIDOrKey must not be empty")
	}
	if issueIDOrKey == "0" {
		return newValidationError("issueIDOrKey must not be '0'")
	}
	return nil
}

// IssueService has methods for Issue.
type IssueService struct {
	method *method

	Attachment *IssueAttachmentService
}

func (s *IssueService) All(projectIDs []int, options ...*QueryOption) ([]*Issue, error) {
	validOptions := []queryType{queryOrder, queryCount, queryOffset, queryCreatedSince, queryCreatedUntil, queryUpdatedSince, queryUpdatedUntil}
	for _, option := range options {
		if err := option.validate(validOptions); err != nil {
			return nil, err
		}
	}

	query := NewQueryParams()
	for _, projectID := range projectIDs {
		if query.Get("projectId[]") == "" {
			query.Set("projectId[]", strconv.Itoa(projectID))
		} else {
			query.Add("projectId[]", strconv.Itoa(projectID))
		}
	}

	for _, option := range options {
		if err := option.set(query); err != nil {
			return nil, err
		}
	}

	resp, err := s.method.Get("issues", query)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	v := []*Issue{}
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return nil, err
	}

	return v, nil
}

func (s *IssueService) AllComments(issue *Issue) ([]*Comment, error) {
	spath := "issues/" + strconv.Itoa(issue.ID) + "/comments"
	resp, err := s.method.Get(spath, NewQueryParams())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	v := []*Comment{}
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return nil, err
	}

	return v, nil
}
