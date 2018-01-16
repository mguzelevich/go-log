# go-log [![Build Status](https://travis-ci.org/mguzelevich/go-log.svg?branch=master)](https://travis-ci.org/mguzelevich/go-log)


go logger wrapper


# usage

init:

```
import (
	"io/ioutil"
	"os"

	"github.com/mguzelevich/go-log"
)

func main() {
	log.InitLoggers(&log.Logger{
		ioutil.Discard,
		ioutil.Discard,
		os.Stdout,
		os.Stdout,
		os.Stderr,
	})

	log.InitLoggers(&log.Logger{ Error: os.Stderr })

}
```

usage

```
import "github.com/mguzelevich/go-log"

log.Debug.Printf("some debug message")
...

log.Error.Printf("some error message %s", err)

```
*/
