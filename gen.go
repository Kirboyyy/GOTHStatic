package main

//go:generate templ generate
//go:generate ./tailwind/tailwindcss -i tailwind/input.css -c tailwind/tailwind.config.js -o static/tailwind.css --minify
