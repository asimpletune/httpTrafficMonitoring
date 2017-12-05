# HTTP Log Monitoring Console Program 0.1

```
Usage:
  httpTrafficMonitoring [--time=10] [--alert=100.0] <file>

Options:
  -h --help     		Show this screen.
  --version     		Show version.
  -t --time=<seconds>	Time to refresh log stats in seconds [default: 10].
  --alert=<hits/second> Threshold to generate an alert that will persist until traffic falls below for two minutes.`
```

## Discussion

This is a CLI app written in golang, and it's intended to consume an actively written log in the Common Log Format and print some stats to the screen.

### Why golang?

I chose golang for a few reasons. First, from my personal experience, golang scales very well with your code complexity. There are plenty of other choices of languages that have good OS bindings and are more expressive, but golang strikes a good balance between expressiveness and forcing you to write code in a way that is sympatico with complexity.

Another reason why I chose golang was because of the concurrency primitives that go provides. Goroutines are very lightweight threads provided by the runtime, so you have nice semantics and with extremely low overhead. Also, channels in go are an excellent abstraction for providing concurrency. In addition, golang as a paradigm means there's support for these concepts, like the `select` keyword in golang, which allows you to multiplex on channels or block, depending on how it's written.

### Basic app architecture

The app has a flat tree structure with a few files.

```
httpTrafficMonitoring/
├── Readme.md
├── alert.go
├── controller.go
├── logEntry.go
├── main.go
├── monitor.go
├── stats.go
└── view.go
```

Although this is currently a fairly simple project, the number of files that exist is intentional and did not add much development overhead at all.

#### main.go

The first file to discuss is `main.go`. This file has a few responsibilities:

1. entry point for the binary
2. parse arguments from the command line
3. intercept OS level SIGNALs and shutdown the program.

This was the first file that I wrote, starting with just the signaling aspect, and later introducing the argument parsing using docopts. I prefer docopts over other arg parsing libraries, but it's out of the scope of the discussion as to why. For more information on that, please see [here](http://docopt.org).

#### controller.go

The next file in the discusison is our controller. While main intentionally had almost no responsibility with respect to the context of what we were solving, the controller has all the context. It receives the signal from main to shut down, starts all the relevant goroutines (monitor and alert, more on that later), receives feedback from logs, and coordinates taking that raw feedback, parsing it into structs to represent log entries, and ultimately handing it off to our view (in this case STDOUT, but it could be anything with how this code is structured).

That being said, our controller only manages the models and when to render thigns in our view. It's a go-between, and doesn't necessarily need do the work itself.

#### stats.go

The stats file can be thought of as our model. It abstracts all of the statistics that are being kept on our log file. When updates are received by the controller, it passes them along to one of the functions in the stats file to update our model.

#### view.go

The view file is what handles updating our view. It doesn't really do anything besides accept some stats and update the view. The controller is the only one who calls it.

#### monitor.go

Monitor.go is also a controller but not really the controller of our program, but rather for the file that we're interfacing with (i.e. the actual http log). Monitor is responsible for all the work pertaining to reading the logs themselves, and the raw log can be thought of as the model.

Monitor works by simply having a pointer to a file, and each successive call to read move our seek position in the file up, and then the results are passed by to the caller.

#### logEntry.go

logEntry is a model that represent a single logEntry.

#### alert.go

Part of this project is to alert the user whenever the average number of hits has exceeded our defined threshold. All of this logic is held here, however we have to signal when this has happened. Therefore, alert takes a channel called `warn`, that tells our controller when this is the case and lets the controller decide how to deal with this information.

In this case, the controller toggles a flag in the view, notifying the user of this, or logs when the flag was untoggled to the console.

The reason for having a separate alert file is because the logic for this functionality is mostly pertaining to a single feature and is therefore distinct from the rest of the application. It runs in its own goroutine and communicates with the above mentioned channel.

### Limitations

The way this app is architected allows it to be fairly flexible for adding new functionality in the future. That being said, the main drawback of the app is a lack of testing. Writing up some tests for this app would take some time and it was more than I had to spare for the project. That being said, since the app is written in such a way where everything is decoupled, it's very decoupled.

The app could be loaded as a library and therefore the test wouldn't have to interface with the CLI. Every component of the app, including the log file themselves are available as interfaces using the architecture described above. Additionally, golang provides automatic composition, so mocking up portions of the application for testing would be a cinch.

The way I tested the functionality was by manually changing the parameters and triggering conditions and then comparing what I expected with what happened. E.g. changing the alert threshold to 1 and the intervals to lower numbers (vs. 2 minutes) and waiting for system to alert and then recover.

### Future work

The main thing I would want to do in the future for an app like this is add features for more statistics. Currently, it's just tracking section statistics and overall hits, but those are just the tip of the iceberg. Adding the additional functionality wouldn't be too difficult at all, the whole app is a blueprint for how to do it.

That being said, I wouldn't personally proceed much further until I get some testing around what already exists, so that's the other side of what stands to be improved.

Lastly, besides just counting more statistics, I would open up the interface a little bit more. Instead of just reading from files I think being able to stream a log into the program would be useful. Also, having more views than just console, but perhaps a socketed view that streams data to a dashboard webpage would be really killer.
