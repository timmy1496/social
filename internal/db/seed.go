package db

import (
	"context"
	"fmt"
	"github.com/timmy1496/social/internal/store"
	"log"
	"math/rand"
)

var usernames = []string{
	"Liam", "Emma", "Noah", "Olivia", "William", "Ava", "James", "Isabella", "Oliver", "Sophia",
	"Benjamin", "Mia", "Elijah", "Charlotte", "Lucas", "Amelia", "Mason", "Harper", "Logan", "Evelyn",
	"Alexander", "Abigail", "Ethan", "Emily", "Jacob", "Elizabeth", "Michael", "Mila", "Daniel", "Ella",
	"Henry", "Avery", "Jackson", "Sofia", "Sebastian", "Camila", "Aiden", "Aria", "Matthew", "Scarlett",
	"Samuel", "Victoria", "David", "Luna", "Joseph", "Grace", "Carter", "Chloe", "Owen", "Penelope",
}

var titles = []string{
	"Beyond the Basics",
	"Why It Still Matters",
	"A Shortcut to Clarity",
	"Lessons I Learned Too Late",
	"The Hidden Costs of Waiting",
	"Fast, Cheap, and Wrong",
	"What No One Tells You",
	"Do Less, Achieve More",
	"The Problem You’re Not Seeing",
	"When Simple Isn’t Easy",
	"Rethinking the Obvious",
	"Small Wins, Big Gains",
	"How It Broke (and Why)",
	"What Changed My Mind",
	"You’re Doing It Backwards",
	"One Rule That Works",
	"The Power of Default Settings",
	"What We Overlook Every Day",
	"Fixing the Wrong Thing",
	"It’s Not About the Tools",
}

var contents = []string{
	"Discover how small habits lead to big results.",
	"A simple trick to improve focus and clarity.",
	"Why doing less can help you achieve more.",
	"Hidden patterns in everyday decisions.",
	"A lesson I learned the hard way.",
	"The key to making smarter choices.",
	"When speed becomes your worst enemy.",
	"What most people miss about motivation.",
	"How to rethink failure as a strategy.",
	"One mindset shift that changes everything.",
	"Why routines matter more than goals.",
	"The unexpected power of boredom.",
	"What silence can teach you about progress.",
	"How to break out of autopilot mode.",
	"Small wins that compound over time.",
	"What I wish I knew a year ago.",
	"The problem isn’t the task — it’s the timing.",
	"A framework for thinking clearly.",
	"Why simple questions matter most.",
	"The overlooked value of starting ugly.",
}

var tags = []string{
	"productivity",
	"technology",
	"lifestyle",
	"startups",
	"mental-health",
	"design",
	"leadership",
	"self-improvement",
	"inspiration",
	"career",
	"coding",
	"remote-work",
	"learning",
	"motivation",
	"ai",
	"minimalism",
	"growth",
	"writing",
	"tools",
	"time-management",
}

var comments = []string{
	"Great point, I never thought about it that way!",
	"Can you explain this a bit more?",
	"Thanks for sharing this insight.",
	"I totally agree with you.",
	"This changed how I approach the problem.",
	"Interesting perspective, but I’m not sure I agree.",
	"I’ve had the same experience.",
	"This made my day, thank you!",
	"What would you recommend as a next step?",
	"Couldn’t have said it better myself.",
	"This is underrated advice.",
	"Loved the simplicity of your explanation.",
	"Bookmarking this for later!",
	"I respectfully disagree — here’s why...",
	"That’s exactly what I needed to hear.",
	"You nailed it.",
	"More people need to know this.",
	"How would you apply this in real life?",
	"Clear and concise, well done.",
	"Can you link to more resources on this?",
}

func Seed(store *store.Storage) {
	ctx := context.Background()

	users := generateUsers(100)
	for _, user := range users {
		if err := store.Users.Create(ctx, user); err != nil {
			log.Println("Error seeding user", err)
			return
		}
	}

	posts := generatePosts(200, users)
	for _, post := range posts {
		if err := store.Posts.Create(ctx, post); err != nil {
			log.Println("Error seeding post", err)
			return
		}
	}

	comments := generateComments(500, users, posts)
	for _, comment := range comments {
		if err := store.Comments.Create(ctx, comment); err != nil {
			log.Println("Error seeding comment", err)
			return
		}
	}

	log.Println("Seeding complete")
}

func generateUsers(num int) []*store.User {
	users := make([]*store.User, num)

	for i := 0; i < num; i++ {
		userName := usernames[i%len(usernames)] + fmt.Sprintf("%d", i)

		users[i] = &store.User{
			Username: userName,
			Email:    userName + "@example.com",
			Password: "123",
		}
	}

	return users
}

func generatePosts(num int, users []*store.User) []*store.Post {
	posts := make([]*store.Post, num)

	for i := 0; i < num; i++ {
		user := users[rand.Intn(len(users))]

		posts[i] = &store.Post{
			UserID:  user.ID,
			Title:   titles[rand.Intn(len(titles))],
			Content: contents[rand.Intn(len(contents))],
			Tags: []string{
				tags[rand.Intn(len(tags))],
				tags[rand.Intn(len(tags))],
			},
		}
	}

	return posts
}

func generateComments(num int, users []*store.User, posts []*store.Post) []*store.Comment {
	cms := make([]*store.Comment, num)

	for i := 0; i < num; i++ {
		cms[i] = &store.Comment{
			PostID:  posts[rand.Intn(len(posts))].ID,
			UserID:  users[rand.Intn(len(users))].ID,
			Content: comments[rand.Intn(len(comments))],
		}
	}

	return cms
}
