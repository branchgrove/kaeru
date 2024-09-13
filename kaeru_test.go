package kaeru

import (
	"errors"
	"reflect"
	"regexp"
	"strings"
	"testing"
	"time"
)

type Username string
type CreatedAt struct{ time.Time }
type IsAdmin bool
type Email string

type Title string
type Body string
type Upvotes int
type Label string
type Metadata map[string]string

var usernameRegex = regexp.MustCompile("^[a-zA-Z0-9_-]{3,16}$")

func (u *Username) ParseString(s string) error {
	if !usernameRegex.MatchString(s) {
		return errors.New("Username must be between 3 and 6 characters long" +
			"only containing numbers, letters, dashes, and underscores.")
	}

	*u = Username(s)

	return nil
}

func (ca *CreatedAt) ParseString(s string) error {
	t, err := time.Parse(time.RFC3339, s)

	if err != nil {
		return err
	}

	*ca = CreatedAt{t}

	return nil
}

func (e *Email) ParseString(s string) error {
	if !strings.Contains(s, "@") {
		return errors.New("Email must contain an @ symbol")
	}

	*e = Email(s)

	return nil
}

func (t *Title) ParseString(s string) error {
	trimmed := strings.TrimSpace(s)
	if len(trimmed) < 3 || len(trimmed) > 100 {
		return errors.New("Title must be between 3 and 100 characters long")
	}
	*t = Title(trimmed)
	return nil
}

func (b *Body) ParseString(s string) error {
	trimmed := strings.TrimSpace(s)
	if len(trimmed) < 10 || len(trimmed) > 5000 {
		return errors.New("Body must be between 10 and 5000 characters long")
	}
	*b = Body(trimmed)
	return nil
}

func (u *Upvotes) ParseFloat64(f float64) error {
	*u = Upvotes(int(f))
	return nil
}

func (l *Label) ParseString(s string) error {
	trimmed := strings.TrimSpace(s)
	if len(trimmed) < 1 || len(trimmed) > 20 {
		return errors.New("Label must be between 1 and 20 characters long")
	}
	*l = Label(trimmed)
	return nil
}

type User struct {
	Username  Username
	Email     Email
	CreatedAt CreatedAt
	IsAdmin   IsAdmin
}

type Comment struct {
	Body      Body
	Metadata  Metadata
	Upvotes   Upvotes
	Commenter User
}

type Post struct {
	Title     Title
	Body      Body
	Metadata  *Metadata
	Labels    []Label
	Upvotes   Upvotes
	Poster    User
	Comments  []Comment
	AdminNote *string
}

func TestParsePost(t *testing.T) {
	// Input data
	input := map[string]any{
		"Title": "My First Post",
		"Body":  "This is the content of my first post. It's pretty exciting!",
		"Metadata": map[string]any{
			"category": "tech",
			"tags":     "golang,testing",
		},
		"Labels":  []any{"new", "featured"},
		"Upvotes": 42.0,
		"Poster": map[string]any{
			"Username":  "johndoe",
			"Email":     "john@example.com",
			"CreatedAt": "2023-09-11T10:00:00Z",
			"IsAdmin":   true,
		},
		"Comments": []any{
			map[string]any{
				"Body": "Great post! Looking forward to more.",
				"Metadata": map[string]any{
					"likes": "5",
				},
				"Upvotes": 5.0,
				"Commenter": map[string]any{
					"Username":  "janedoe",
					"Email":     "jane@example.com",
					"CreatedAt": "2023-09-10T09:00:00Z",
					"IsAdmin":   false,
				},
			},
		},
		"AdminNote": "Approved for front page",
		"FewBytes":  []byte{1, 2, 3, 4},
	}

	// Expected output
	expectedTime, _ := time.Parse(time.RFC3339, "2023-09-11T10:00:00Z")
	expectedCommentTime, _ := time.Parse(time.RFC3339, "2023-09-10T09:00:00Z")
	expected := &Post{
		Title: "My First Post",
		Body:  "This is the content of my first post. It's pretty exciting!",
		Metadata: &Metadata{
			"category": "tech",
			"tags":     "golang,testing",
		},
		Labels:  []Label{"new", "featured"},
		Upvotes: 42,
		Poster: User{
			Username:  "johndoe",
			Email:     "john@example.com",
			CreatedAt: CreatedAt{expectedTime},
			IsAdmin:   true,
		},
		Comments: []Comment{
			{
				Body: "Great post! Looking forward to more.",
				Metadata: Metadata{
					"likes": "5",
				},
				Upvotes: 5,
				Commenter: User{
					Username:  "janedoe",
					Email:     "jane@example.com",
					CreatedAt: CreatedAt{expectedCommentTime},
					IsAdmin:   false,
				},
			},
		},
		AdminNote: ptr("Approved for front page"),
	}

	// Actual output
	actual := new(Post)
	err := Parse(input, actual)

	// Check for errors
	if err != nil {
		t.Fatalf("Parse returned an error: %v", err)
	}

	// Compare actual and expected
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Parse result not as expected.\nGot: %+v\nWant: %+v", actual, expected)
	}
}

// Helper function to create a pointer to a string
func ptr(s string) *string {
	return &s
}
