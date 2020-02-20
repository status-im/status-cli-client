# Description

This is a dead-simple CLI client for sending Status messages.

It uses the [`status-go`](https://github.com/status-im/status-go) library and is essentially a very basic example of its usage.

# Usage

```
 $ ./status-cli-client -help
Usage of ./status-cli-client:
  -addr string
    	Listening address for Whisper node thread. (default "0.0.0.0")
  -chat string
    	Name of public chat to send to. (default "whatever")
  -data string
    	Location for Status data. (default "/tmp/status-cli-client")
  -ens string
    	ENS name to send with the message.
  -key string
    	Hex private key for your Status identity.
  -message string
    	Message to send to the public channel. (default "TEST")
  -port int
    	Listening port for Whisper node thread. (default 30303)
  -timeout int
    	Timeout for message delivery in milliseconds. (default 500)
```

Example usage would be:
```bash
 $ ./status-cli-client -chat test-channel -message "Pretty cool!"
```

# Known Issues

* We're using `time.Sleep()` to let Whisper deliver the message, which is stupid
* We're using an in-memory instance of SQLite that holds things like whisper keys
