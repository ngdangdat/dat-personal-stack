package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

// GitHub API Models
type GHUser struct {
	Login     string `json:"login"`
	AvatarURL string `json:"avatar_url"`
}

type GHPRItem struct {
	ID            int64      `json:"id"`
	Number        int        `json:"number"`
	Title         string     `json:"title"`
	User          GHUser     `json:"user"`
	RepositoryURL string     `json:"repository_url"`
	HTMLURL       string     `json:"html_url"`
	State         string     `json:"state"`
	Draft         bool       `json:"draft"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	ClosedAt      *time.Time `json:"closed_at"`
	PullRequest   *struct {
		MergedAt *time.Time `json:"merged_at,omitempty"`
	} `json:"pull_request,omitempty"`
}

type GHSearchResponse struct {
	Items []GHPRItem `json:"items"`
}

type GHReview struct {
	State       string    `json:"state"`
	SubmittedAt time.Time `json:"submitted_at"`
}

type GHCommit struct {
	Commit struct {
		Author struct {
			Date time.Time `json:"date"`
		} `json:"author"`
	} `json:"commit"`
}

// Frontend PR representation
type PRDetail struct {
	ID                  int64      `json:"id"`
	Number              int        `json:"number"`
	Title               string     `json:"title"`
	Repo                string     `json:"repo"`
	Author              string     `json:"author"`
	AuthorAvatar        string     `json:"author_avatar"`
	URL                 string     `json:"url"`
	State               string     `json:"state"`
	IsDraft             bool       `json:"is_draft"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
	ClosedAt            *time.Time `json:"closed_at"`
	MergedAt            *time.Time `json:"merged_at"`
	Status              string     `json:"status"` // draft, changes_requested, approved, reviewing
	ReviewVelocityHours *float64   `json:"review_velocity_hours"`
	LeadTimeHours       *float64   `json:"lead_time_hours"`
}

type PRDashboardMetrics struct {
	AverageReviewVelocityHours float64 `json:"avg_review_velocity_hours"`
	AverageLeadTimeHours       float64 `json:"avg_lead_time_hours"`
	TotalCount                 int     `json:"total_count"`
	ApprovedCount              int     `json:"approved_count"`
	ChangesRequestedCount      int     `json:"changes_requested_count"`
	DraftCount                 int     `json:"draft_count"`
	ReviewingCount             int     `json:"reviewing_count"`
}

type PRDashboardResponse struct {
	PRs     []PRDetail         `json:"prs"`
	Metrics PRDashboardMetrics `json:"metrics"`
}

type GitHubService struct{}

func NewGitHubService() *GitHubService {
	return &GitHubService{}
}

// makeRequest is a helper to run HTTP requests against GitHub API
func (s *GitHubService) makeRequest(token, url string, target interface{}) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "token "+token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "Engineering-Assistant-Go")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("GitHub API error: status %d, body: %s", resp.StatusCode, string(body))
	}

	return json.NewDecoder(resp.Body).Decode(target)
}

// FetchPRs fetches and computes PR statistics for a given query type ("reviewing" or "mine")
func (s *GitHubService) FetchPRs(token, queryType, username string) (*PRDashboardResponse, error) {
	var q string
	if queryType == "reviewing" {
		// PRs assigned to user for review, showing open ones
		q = fmt.Sprintf("is:pr is:open review-requested:%s", username)
	} else {
		// PRs authored by user (open & closed so we can compute lead time metrics from recently closed)
		q = fmt.Sprintf("is:pr author:%s", username)
	}

	url := fmt.Sprintf("https://api.github.com/search/issues?q=%s&sort=updated&order=desc&per_page=30", strings.ReplaceAll(q, " ", "+"))
	var searchResp GHSearchResponse
	if err := s.makeRequest(token, url, &searchResp); err != nil {
		return nil, err
	}

	prs := make([]PRDetail, len(searchResp.Items))
	var wg sync.WaitGroup
	// Limit concurrency to 5 requests at a time to stay safe on rate limit
	sem := make(chan struct{}, 5)

	for i, item := range searchResp.Items {
		// Extract repo name from repository_url (e.g. "https://api.github.com/repos/owner/name")
		parts := strings.Split(item.RepositoryURL, "/repos/")
		repoName := ""
		if len(parts) > 1 {
			repoName = parts[1]
		}

		prs[i] = PRDetail{
			ID:           item.ID,
			Number:       item.Number,
			Title:        item.Title,
			Repo:         repoName,
			Author:       item.User.Login,
			AuthorAvatar: item.User.AvatarURL,
			URL:          item.HTMLURL,
			State:        item.State,
			IsDraft:      item.Draft,
			CreatedAt:    item.CreatedAt,
			UpdatedAt:    item.UpdatedAt,
			ClosedAt:     item.ClosedAt,
			Status:       "reviewing", // default status
		}

		if item.PullRequest != nil && item.PullRequest.MergedAt != nil {
			prs[i].State = "merged"
			prs[i].MergedAt = item.PullRequest.MergedAt
		}

		// Spawn concurrent workers to fetch commits/reviews for metrics calculation
		wg.Add(1)
		go func(idx int, rName string, prNum int) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			s.enrichPRDetails(token, rName, prNum, &prs[idx])
		}(i, repoName, item.Number)
	}

	wg.Wait()

	// Calculate overall metrics
	var totalVel, totalLead float64
	var velCount, leadCount int
	var approved, changesReq, drafts, reviewing int

	for _, pr := range prs {
		if pr.IsDraft {
			drafts++
		} else {
			switch pr.Status {
			case "approved":
				approved++
			case "changes_requested":
				changesReq++
			case "reviewing":
				reviewing++
			}
		}

		if pr.ReviewVelocityHours != nil {
			totalVel += *pr.ReviewVelocityHours
			velCount++
		}
		if pr.LeadTimeHours != nil {
			totalLead += *pr.LeadTimeHours
			leadCount++
		}
	}

	metrics := PRDashboardMetrics{
		TotalCount:            len(prs),
		ApprovedCount:         approved,
		ChangesRequestedCount: changesReq,
		DraftCount:            drafts,
		ReviewingCount:        reviewing,
	}

	if velCount > 0 {
		metrics.AverageReviewVelocityHours = totalVel / float64(velCount)
	}
	if leadCount > 0 {
		metrics.AverageLeadTimeHours = totalLead / float64(leadCount)
	}

	return &PRDashboardResponse{
		PRs:     prs,
		Metrics: metrics,
	}, nil
}

// enrichPRDetails fetches commits and reviews from GitHub to compute metrics
func (s *GitHubService) enrichPRDetails(token string, repo string, prNum int, pr *PRDetail) {
	if repo == "" || prNum == 0 {
		return
	}

	// 1. Fetch Reviews to calculate Review Velocity and overall Status
	reviewsUrl := fmt.Sprintf("https://api.github.com/repos/%s/pulls/%d/reviews", repo, prNum)
	var reviews []GHReview
	if err := s.makeRequest(token, reviewsUrl, &reviews); err == nil {
		var firstReview *GHReview
		// Determine final review status: approved, changes_requested, or reviewing
		hasApproval := false
		hasChangesReq := false

		for _, r := range reviews {
			// Save the first review timestamp
			if firstReview == nil || r.SubmittedAt.Before(firstReview.SubmittedAt) {
				cp := r
				firstReview = &cp
			}

			if r.State == "APPROVED" {
				hasApproval = true
			} else if r.State == "CHANGES_REQUESTED" {
				hasChangesReq = true
			}
		}

		// Update overall PR status badge based on reviews
		if pr.IsDraft {
			pr.Status = "draft"
		} else if hasChangesReq {
			pr.Status = "changes_requested"
		} else if hasApproval {
			pr.Status = "approved"
		} else {
			pr.Status = "reviewing"
		}

		// Calculate Review Velocity: time between creation and first review
		if firstReview != nil {
			vel := firstReview.SubmittedAt.Sub(pr.CreatedAt).Hours()
			if vel < 0 {
				vel = 0
			}
			pr.ReviewVelocityHours = &vel
		}
	}

	// 2. Fetch Commits to calculate Lead Time
	commitsUrl := fmt.Sprintf("https://api.github.com/repos/%s/pulls/%d/commits?per_page=1", repo, prNum)
	var commits []GHCommit
	if err := s.makeRequest(token, commitsUrl, &commits); err == nil && len(commits) > 0 {
		firstCommitTime := commits[0].Commit.Author.Date

		// If merged, calculate lead time from first commit to merge time
		if pr.State == "merged" && pr.MergedAt != nil {
			lt := pr.MergedAt.Sub(firstCommitTime).Hours()
			if lt < 0 {
				lt = 0
			}
			pr.LeadTimeHours = &lt
		} else if pr.State == "closed" && pr.ClosedAt != nil {
			// If closed but not merged, use closed_at as end time
			lt := pr.ClosedAt.Sub(firstCommitTime).Hours()
			if lt < 0 {
				lt = 0
			}
			pr.LeadTimeHours = &lt
		}
	}
}
