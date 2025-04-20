# README

## About

PolAIn is a simple "conversation with AI" interface to communicate with <https://pollinations.ai> free an open source endpoints. It is developed with [Go](https://go.dev) and the [Wails](https://wails.io/) project.

The interface uses [Vue](https://vuejs.org).

All of this make it possible to build a cross-platform (Linux, Windows and macOS) application from one source code.

## Behind the scene

Pollinations is a fascinating project which is entirely free, without the need of registration, and it offers plenty of models. It also provides an image generation endpoint on Flux1.dev and Turbo models.

PolAIn uses theses endpoints:

- text.polinations.ai/openai to post "OpenAI" compatible requests and waiting streaming response
- image.polinations.ai to create images on demand

> All of these endpoints provides a "private" option to not share your prompt and images. **It's the default in PolAIn, everything is set to private.**

Our default "system" prompt asks for the model to produce images in Markdown according to the image endpoint template. It may, sometimes, block the interface or badly generate the URL. I'm currently working on bugfixes.

I really want to thanks the Pollinations teams to offer such a service. If you want to sponsor them, please go to <https://ko-fi.com/pollinationsai>

## Help

I'm open to any help you can give:

- translation: fork the project, use the "locales/en.yaml" file as reference (or use the existing file for your language), then create a pull-request
- help on design
- help on adding more features (send files to the model, RAG, and so on)

## Note

This project wants to go to [Fyne](https://fyne.io). Fyne can make the application way faster, and I could compile the application to Android ans iOs.

The actual problem are :

- Richtext (markdown) is not interactive. So the user cannot select text, or copy the images.
- Code blocks are not highlighted

I already created the entire project and I'm waiting for updates from the Fyne project.
