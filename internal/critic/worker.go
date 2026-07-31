package critic

import (
	"context"
	"fmt"
	"log"
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

// Worker polls critic jobs and posts @vibe_critic replies.
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
		httpGet:  FetchURLSnippet,
	}
}

func (w *Worker) Run(ctx context.Context) {
	log.Println("critic worker started")
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		if err := w.ProcessOne(ctx); err != nil {
			log.Printf("critic worker: %v", err)
		}
		select {
		case <-ctx.Done():
			log.Println("critic worker stopped")
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
		if post.UserID == domain.BotUserID || post.User.IsBot {
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
		if c.UserID == domain.BotUserID || c.User.IsBot {
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
		u = strings.TrimRight(u, ".,);]")
		snippet, err := w.httpGet(ctx, u)
		if err != nil {
			b.WriteString(fmt.Sprintf("- %s: (fetch failed: %v)\n", u, err))
			continue
		}
		b.WriteString(fmt.Sprintf("- %s: %s\n", u, snippet))
	}
	return b.String()
}
