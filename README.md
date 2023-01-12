# How to overcomplicate [Advent of Code](https://adventofcode.com/)

This repo is meant as a self exercise for some topics and drastically overcomplicates some things in order
to have more fun while developing solutions for the Advent of Code.

So what does this mean?

1. Let's make it and enterprise app, so __WE NEED A SERVER + A WEB UI!__
2. New coding challenges are activated each day. We want to implement them
   but also they need to seemlessly integratable into our 24/7 running server
   -> __WE NEED A DYNAMIC PLUGIN SYSTEM__
3. We want to be able to program solutions in most modern programming languages -> __THE PLUGIN SYSTEM MUST SUPPORT ALL THE PROGRAMMING__
4. Don't overcomplicate the whole web JS world but instead __USE PLAIN HTML + CSS TO ACHIEVE EVERYTHING IF POSSIBLE!__

## Server + Web UI

The server is written in Go since I like it. It should serve the web UI via template rendered HTML files (of cource embedded within the app).
Since we want to avoid JS as much as we can -> use form submits and the rendering for updates.

## Plugin system

I like what [the guys at HashiCorp have developed](https://github.com/hashicorp/go-plugin) and it supports gRPC so almost any programming language can be used.
Let's further enhance upon that by allowing plugins to be added via the web UI.
