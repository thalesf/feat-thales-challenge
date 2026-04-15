# Niche.com Full-Stack Coding Exercise

## Overview

Welcome to the Niche.com full-stack coding exercise! This exercise is designed to assess your skills in building a full-stack application with Go, TypeScript, and React. You'll be implementing a college review search and autocomplete system using real college review data from Niche.com.

For this exercise, you'll need to:
1. Implement the back-end functionality for loading reviews, autocompleting college names, and retrieving reviews for a specific college.
2. Create a React front-end that implements an autocomplete feature and displays college reviews.

## Guidelines
1.  Feel free to ask questions about the exercise, if you need clarity on approach, requirements, or anything else.  Asking meaningful questions can only reflect well on you. 
2.  Please submit your solution within 7 days. 
3.  Feel free to use AI tooling if you believe it will be useful.  Be prepared to discuss your approach to using this tooling, and document your approach as described in the [Submission](#submission) section. 

## Submission

Your submission should include the following:
- A branch containing your changes
- A new file in that branch that documents your approach to solving this problem.
  - Call out key decisions you made, and why
- If you used AI tooling of any sort, please include your approach to working with those tools, as well as specific prompts you used, in this documentation. 

To submit your solution, create a pull request against the `main` branch in this repository and email your recruiter.

## Requirements

### Back-End Requirements

The back-end is a Go service with some incomplete functionality. Your task is to:

1. Implement the `loadReviews` function in `reviews.go`
   - Process the review data from the provided CSV file.
   - Implement `ReviewsData` as an in-memory data structure to support the following operations:
     - Retrieve reviews for a given college
     - Support autocomplete for college names
   - You do not need to persist this data - it's fine to read it from the CSV each time the service starts.
2. Implement `handleGetReviews` in `server.go`
   - The endpoint should accept a college url (as specified in the source data).
   - The endpoint should return reviews for the specified entity.
2. Implement `handleAutocomplete` in `server.go`
   - The endpoint should accept a query parameter and return matching college names.
   - Return any matching colleges including their Name and url.

### Front-End Requirements

Using React and TypeScript, implement:

1. An autocomplete component for searching college names
   - The component should display matching college names as the user types
   - The component should support keyboard navigation (arrow keys, enter for selection)
   - Visual styling with CSS

2. A review display component
   - When a college is selected from autocomplete, fetch and display its reviews
   - Format the reviews in a clean, readable way
   - Visual styling with CSS

## Project Structure

### Back-End (Go)
- `back-end/main.go` - Entry point for the Go service
- `back-end/server.go` - HTTP server implementation with endpoints
- `back-end/reviews.go` - Reviews data processing
- `back-end/data/niche_reviews.csv` - A CSV file containing college reviews data

### Front-End (React/TypeScript)
- `front-end/` - Contains the React application

## Getting Started

### Back-End
1. Navigate to the `back-end` directory
2. Run `go mod tidy` to install dependencies
3. Run `go run .` to start the server
4. The server will be available at `http://localhost:8080`

### Front-End
1. Navigate to the `front-end` directory
2. Run `npm install` to install dependencies
3. Run `npm run dev` to start the development server
4. The application will be available at `http://localhost:3000`

## Evaluation Criteria

Your submission will be evaluated based on:
1. **Functionality** - Does it correctly implement all the required features?
2. **Code Quality** - Is the code well-structured, readable, and maintainable?
3. **Performance** - Does the service start quickly? Are page interactions fast?
4. **UI/UX** - Is the interface intuitive and responsive?
5. **Technical Decisions** - Can you justify your technical choices?

## Submission Guidelines

1. Please document your approach in a `SOLUTION.md` file, explaining:
   - Your implementation strategy
   - Any trade-offs or assumptions you made
   - How you used AI tools (if applicable)
   - Any challenges you encountered
   - How your solution could be improved with more time

2. Include instructions for running your solution

3. If you use AI tools like ChatGPT, GitHub Copilot, Claude Code or other similar tools, please document:
   - Which parts of your solution used AI assistance
   - The specific prompts you used
   - How you verified and modified the AI-generated code
   - Where you found the AI tools most and least helpful and why

We encourage you to ask questions about any requirements that may be unclear. Feel free to make reasonable assumptions where necessary, but please document these assumptions in your submission.

## Time Expectation

We expect this exercise to take approximately 2-3 hours to complete. Don't worry if you don't finish everything - focus on demonstrating your best work within a reasonable timeframe.

Good luck!
