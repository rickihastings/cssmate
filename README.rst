CSSMate
-------

CSSMate is a live CSS editor (reloader?). I'm not actually sure what you want to call it, but it aims to provide a no-thrills CSS reloader. It's quite similar to existing tools such as `live reload <http://livereload.com/>`_, `takana <http://usetakana.com/>`_ or the many others.

CSSMate is different in that it injects the updated stylesheet back into the page immediately without a reload, saving time and effort, it also doesn't require you to do anything fancy like running your application behind a custom server, inserting javascript into the address bar or any other nonsense!

Why
===

You're probably wondering why? There's a million out there already. I created this because I wasn't happy with the existing tools, some tools such as livereload reload the entire page. I can Ctrl+S, Alt+Tab and Ctrl+R incredibly quickly, livereload isn't cutting it for me, full reloads in complex web applications can take a little longer. Takana has a load of dependencies, nodejs? OSX? Sublime text plugin? No thanks.

Not everyone uses these, not everyone has node on their work station, people use different editors. Some of the alternatives have bulky implementations, forcing me to run my code behind a webserver they create, or injecting javascript through the address bar. I don't like this.

CSSMate was built in Go so it can be compiled to native binaries so there's no need for any dependencies, you can just drop the binary into your home directory and point it at a folder of choice. Did I mention it's completely multi-platform?

Installing
==========

There's two ways of using CSSMate, if you have the Go environment setup on your machine you can compile it stupidly easy with the following commands ::

    $ go get github.com/rickihastings/cssmate
    $ cd $GOPATH/src/github.com/rickihastings/cssmate
    $ go build -o bin/cssmate
    $ cd $GOPATH/bin

All this is providing you have the Go environment setup correctly (which I'm not going to go through here), the ``cssmate`` binary will be available in the folder you're now in.

Alternatively, you can clone the github repo and grab the pre built binaries for 64, 32 bit linux and 32 bit Windows, I never bothered doing OSX, if anyone wants it though speak up. ::

   $ git clone https://github.com/rickihastings/cssmate
   $ cd cssmate/bin

You'll see two files ::

   cssmate-386
   cssmate-amd64
   cssmate.exe

Running
=======

Once you have the binary compiled you can simply run it with the following command ::

   $ ./cssmate --path="/home/myuser/myproject/css"

CSSMate will now watch for any file changes in that folder and propogate them to any clients, before it will work though, you need to go over to your ``index.html`` or base template and insert the following script tag at the end of the body tag. ::

   <script type="text/javascript" src="http://localhost:58900/cssmate.js"></script>

You can run cssmate on a different port by passing in the ``--port`` parameter, you will need to change

Notes
=====

I built this in about 2 hours as a quick fix to one of my problem, but I'm so pleased with it I'd like to share it! It may have some small bugs, feel free to open an issue and I'll get round to fixing it. It has no tests, it could do with some, and probably will soon when I get round to it.