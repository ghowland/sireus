---
layout: post
title:  "A Sireus Start"
date:   2023-02-08 19:27:19 +0700
categories: sireus news
---
Hi, I'm Geoff Howland, and I'm starting this blog to write about the Sireus project's development.

Sireus is a Decision System, and sits between your Monitoring and Alerting.  It can run commands and export state
information back into your monitoring system, so you can use the information to make further decisions or execute 
commands with another system.

There will be an interactive demo coming soon where I send out a Request for Comments to get feedback before the 
design is solidified for writing the first production-ready version.

I am initially implementing Sireus against Prometheus and Alert Manager, as they are popular monitoring tools that
are easy to set up.  Later I will add more monitoring and alerting options to collect and publish to.

Sireus will also contain a Client that can be run in various privileged environments to run commands or send web requests 
in the correct places.

Check out the [Github page][sireus-gh] for more info about the project.

If you are interested in contributing, please start by creating issues around the design.  Sireus is currently in
the design RFC phase, and I want to collect opinions, feature requests and questions, so that I can fix any initial 
design flaws in this early phase of the project.

[sireus-gh]:   https://github.com/ghowland/sireus
