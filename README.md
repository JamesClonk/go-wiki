go-wiki
=======

My implementation of the Go example wiki from http://golang.org/doc/articles/wiki/

---

It's almost the same, except for having quite a bunch of unit tests added.
The only "big" differences are Page.Body being of type template.HTML and parsing any files found under view/* as template files.

I put a few configuration variables at the top, but I'm not sure if that's the Go way to do things. There's probably a config package or library or something to handle such configuration stuff better.

--
For my absolute first try at Go I think I'm with quite satisfied with this little project. 
Go looks like a fun language to use, and the wiki example makes as good impression.

I'll definitely have to start reading the documentation now and going through the tutorial. 
"Slices? Huh? What's that?  ..Ahh, you mean Lists?" ;-)
   
---   
   
#### Update   
Added hot_recompile.sh script.   
It will watch for file changes in the current directory through inotify and "restart" wiki.go if anything has changed.   
This kind of emulates "**Hot Reload**" as known from other web development languages / frameworks.

