# Product: SwiftLink (URL Shortener)

SwiftLink is a URL shortening service that lets users shorten long URLs and share them. Visiting a short link redirects to the original URL, and click counts are tracked per link.

The project is branded as **SwiftLink** in the UI, though the repository is named `url-shortener`.

## Current State

Phase 1 (Fullstack) is in progress — core shortening and redirect functionality is implemented. Phases 2–5 cover Docker, CI/CD, Kubernetes, and production monitoring respectively.

## Core Functionality

- Shorten a URL → receive a short code (6-character alphanumeric)
- Visit `/{shortCode}` → redirect to the original URL
- Click count incremented on each redirect
