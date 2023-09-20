# This is my personal wiki.


It is inspired by wiki's, weblogs, and digital gardens like Devine Lu Linvega's [XIIVV](https://wiki.xxiivv.com/site/home.html), Chase McCoy's [chasem.co](https://chasem.co/), Gwern Branwen's [gwern.net](https://gwern.net/) and so many [more](/site/links). 

If you would like to go deeper on digital gardens, check out Maggie Appleton's post on the [history and ethos of digital gardens.](https://maggieappleton.com/garden-history)


## How it started
 
 I gave myself the challenge of building a simple static site generator that would act as my personal wiki and my portfolio. I wanted it to be extremely bare bones, with as few dependencies as possible.

At the time of creating this, I was studying Go and I discovered the html/template package and Blackfriday. I thought these tools would work well for what I wanted to do.

## The bones

So, this website is written in Go. I have two html templates, one for the site pages and one for the journal entries. The site pages and journal entries are written in Markdown and converted to HTML using Blackfriday a Markdown processor implemented in Go.

Find the code here [iamseeley/ts-wiki](https://github.com/iamseeley/ts-wiki).

## How it's going

I'm working on making a project template so other people can try it out!

> "It's a living document that outlines where one has 
> been, and a tool that advises where one could go."

> <cite>--- Devine Lu Linvega</cite>

