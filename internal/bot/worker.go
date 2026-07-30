package bot

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/dakaii/vibegopher/internal/domain"
	"github.com/dakaii/vibegopher/internal/envvar"
	"github.com/dakaii/vibegopher/internal/repository/botjobrepo"
	"github.com/dakaii/vibegopher/internal/repository/commentrepo"
	"github.com/dakaii/vibegopher/internal/repository/postrepo"
	"github.com/google/uuid"
)

var urlRegexp = regexp.MustCompile(`https?://[^\s]+`)

type Worker struct {
	jobs     *botjobrepo.BotJobRepo
	posts    *postrepo.PostRepo
	comments *commentrepo.CommentRepo
	gemini   *GeminiClient
	interval time.Duration
	httpGet  func(ctx context.Context, url string) (string, error)
}

func NewWorker(jobs *botjobrepo.BotJobRepo, posts *postrepo.PostRepo, comments *commentrepo.CommentRepo) *Worker {
	return &Worker{
		jobs:     jobs,
		posts:    posts,
		comments: comments,
		gemini:   NewGeminiClient(envvar.GeminiAPIKey()),
		interval: 3 * time.Second,
		httpGet:  fetchURLSnippet,
	}
}

func (w *Worker) Run(ctx context.Context) {
	log.Println("bot worker started")
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		if err := w.ProcessOne(ctx); err != nil {
			log.Printf("bot worker: %v", err)
		}
		select {
		case <-ctx.Done():
			log.Println("bot worker stopped")
			return
		case <-ticker.C:
		}
	}
}

func (w *Worker) ProcessOne(ctx context.Context) error {
	job, err := w.jobs.ClaimNext()
	if err != nil {
		return err
	}
	if job == nil {
		return nil
	}

	reply, err := w.buildReply(ctx, job)
	if err != nil {
		retry := job.Attempts < 5
		_ = w.jobs.MarkFailed(job.ID, err.Error(), retry)
		return err
	}

	// Skip empty / no-op
	if strings.TrimSpace(reply) == "" {
		return w.jobs.MarkDone(job.ID)
	}

	comment := domain.Comment{
		Content: reply,
		UserID:  domain.BotUserID,
		PostID:  job.TargetID,
	}
	if job.Kind == domain.BotJobKindCommentCreated {
		parent, err := w.comments.GetCommentByID(job.TargetID)
		if err != nil {
			retry := job.Attempts < 5
			_ = w.jobs.MarkFailed(job.ID, err.Error(), retry)
			return err
		}
		comment.PostID = parent.PostID
	}

	if _, err := w.comments.CreateComment(comment); err != nil {
		retry := job.Attempts < 5
		_ = w.jobs.MarkFailed(job.ID, err.Error(), retry)
		return err
	}
	return w.jobs.MarkDone(job.ID)
}

func (w *Worker) buildReply(ctx context.Context, job *domain.BotJob) (string, error) {
	var author, content string
	var postID uuid.UUID

	switch job.Kind {
	case domain.BotJobKindPostCreated:
		post, err := w.posts.GetPostByID(job.TargetID)
		if err != nil {
			return "", err
		}
		if post.UserID == domain.BotUserID {
			return "", nil
		}
		author = post.User.Username
		content = post.Content
		postID = post.ID
	case domain.BotJobKindCommentCreated:
		c, err := w.comments.GetCommentByID(job.TargetID)
		if err != nil {
			return "", err
		}
		if c.UserID == domain.BotUserID {
			return "", nil
		}
		author = c.User.Username
		content = c.Content
		postID = c.PostID
	default:
		return "", fmt.Errorf("unknown job kind %q", job.Kind)
	}

	linkContext := w.collectLinkContext(ctx, content)
	prompt := fmt.Sprintf(
		"Kind: %s\nAuthor: @%s\nPostID: %s\nContent:\n%s\n\nLink context:\n%s\n\nWrite your single comment reply now.",
		job.Kind, author, postID, content, linkContext,
	)
	return w.gemini.GenerateComment(ctx, prompt)
}

func (w *Worker) collectLinkContext(ctx context.Context, content string) string {
	urls := urlRegexp.FindAllString(content, 2)
	if len(urls) == 0 {
		return "(none)"
	}
	var b strings.Builder
	for _, u := range urls {
		snippet, err := w.httpGet(ctx, u)
		if err != nil {
			b.WriteString(fmt.Sprintf("- %s: (fetch failed: %v)\n", u, err))
			continue
		}
		b.WriteString(fmt.Sprintf("- %s: %s\n", u, snippet))
	}
	return b.String()
}

func fetchURLSnippet(ctx context.Context, rawURL string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "VibeGopherBot/1.0")
	client := &http.Client{Timeout: 8 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("status %d", resp.StatusCode)
	}
	buf := make([]byte, 4096)
	n, _ := resp.Body.Read(buf)
	text := stripTags(string(buf[:n]))
	text = strings.Join(strings.Fields(text), " ")
	if len(text) > 500 {
		text = text[:500] + "..."
	}
	if text == "" {
		return "(empty body)", nil
	}
	return text, nil
}

func stripTags(s string) string {
	var b strings.Builder
	inTag := false
	for _, r := range s {
		switch {
		case r == '<':
			inTag = true
		case r == '>':
			inTag = false
		case !inTag:
			b.WriteRune(r)
		}
	}
	return b.String()
}
