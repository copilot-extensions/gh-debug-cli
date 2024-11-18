# GH Debug CLI

This tool allows you to chat with your agent locally in order to create a faster feedback loop for developers developing an extension.
Debug mode is enabled by default so that you can see clearer information around what exactly is getting parsed successfully.

The different SSE events in the [agent protocol](TODO) that the CLI gives debug output for are:
1. [errors](TODO)
2. [references](TODO)
3. [confirmations](TODO)

> Note: This tool does not handle the payload verification process. To use this tool to validate your events, please temporarily disable payload verification for local testing and re-enable when completed.

## Install the debug tool
1. Authenticate with GitHub CLI OAuth app
   ```shell
   gh auth login --web -h github.com
   ```
1. Install / upgrade extension
   ```shell
   gh extension install github.com/copilot-extensions/gh-debug-cli
   ```
1. See more info about the cli tool
   ```shell
   gh debug-cli -h
   ```

## Using the debug chat tool
1. Run the following command `gh debug-cli -h` to see the different flags that it takes in.
```
> gh debug-cli -h
This cli tool allows you to debug your agent by chatting with it locally.

Usage:
   [flags]

Flags:
  -h, --help              help for this command
      --log-level DEBUG   Log level to help debug events. Supported types are DEBUG, `TRACE`, `NONE`. `DEBUG` returns general logs. `TRACE` prints the raw http response. (default "DEBUG")
      --token string      GitHub token for chat authentication (optional)
      --url string        url to chat with your agent (default "http://localhost:8080")
      --username string   username to display in chat (default "sparklyunicorn")
```
> The token noted in the flag above is used to authenticate against the provided LLM. If you are using a different service, then this token is not needed. Generate the user-to-server token by [creating a GitHub Applicatiion](https://docs.github.com/en/apps/creating-github-apps/about-creating-github-apps/about-creating-github-apps) and then following the [using the device flow to generate a user access token](https://docs.github.com/en/apps/creating-github-apps/authenticating-with-a-github-app/generating-a-user-access-token-for-a-github-app#using-the-device-flow-to-generate-a-user-access-token) to generate the token.
2. You can alternatively set these flags as environment variables (in all caps) so you don't need to pass them in every time. The only "required" one to get this up and running is the url for your agent
```
export URL="http://localhost:8080/agent/blackbeard"
```
3. When you run the CLI, you will see any flags that were previously set in your environment variables as the output.
```
>  gh debug-cli
Setting url to http://localhost:8080/agents/blackbeard

Start typing to chat with your assistant...
sparklyunicorn: 
```
4. Type something to simulate chatting with your assistant.
```
> gh debug-cli
Setting url to http://localhost:8080/agents/blackbeard

Start typing to chat with your assistant...
sparklyunicorn: hello
assistant: Ahoy, @monalisa! A jolly good day to ye, me heartie. How can ol' Blackbeard be of service to ye today? 

Huzzah! You successfully received a message!
╔═══════════╤════════════════════════════════════════════════════════════════╗
║   Role    │                            Content                             ║
╟━━━━━━━━━━━┼━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━╢
║ assistant │ [condensed] Ahoy, @monalisa! A jolly good day to ye, me hearti ║
╟━━━━━━━━━━━┼━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━╢
║                                                        Parsed message data ║
╚═══════════╧════════════════════════════════════════════════════════════════╝
sparklyunicorn: 

```
5. To debug your SSE events, you can set up a key word that your assistant uses to send you a specific type of event. My blackbeard agent allows me to send a keyword "confirmation", and here I can see the debug output on what is parsed from the SSE event
```
> sparklyunicorn: confirmation
assistant: Arrr, @monalisa! I be ready and waitin' for yer confirmation. Be ye ready to set sail on this treacherous journey and receive a custom limerick 'bout petals? Aye or nay, let me know yer decision, and I'll be at yer service.

Huzzah! You successfully received a message!
╔═══════════╤════════════════════════════════════════════════════════════════╗
║   Role    │                            Content                             ║
╟━━━━━━━━━━━┼━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━╢
║ assistant │ [condensed] Arrr, @monalisa! I be ready and waitin' for yer co ║
╟━━━━━━━━━━━┼━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━╢
║                                                        Parsed message data ║
╚═══════════╧════════════════════════════════════════════════════════════════╝

Huzzah! You successfully received a confirmation!
╔══════════════╤═════════════════════════════════════════════════════╗
║ Key          │ Value                                               ║
╟━━━━━━━━━━━━━━┼━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━╢
║ type         │ action                                              ║
║ title        │ Be ye sure ye want a custom limerick 'bout petals ? ║
║ message      │ Arrr, this here action be irreversible, matey!      ║
║ confirmation │ map[id:123]                                         ║
╟━━━━━━━━━━━━━━┼━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━╢
║                                           Parsed confirmation data ║
╚══════════════╧═════════════════════════════════════════════════════╝

Be ye sure ye want a custom limerick 'bout petals ?
  Arrr, this here action be irreversible, matey!
Reply: [y/N]
```
6. If I got a bad confirmation, it would look something like this
```
> sparklyunicorn: bad confirmation

Alas...The following is not a valid confirmation:
 ["conf"]

assistant: Avast, @monalisa! Me apologies if I didn't quite understand yer request. Pray tell, could ye please clarify what be wrong with the confirmation? I be here to assist ye, me matey!
```
7. And if debug mode was set to false, then I would only see the confirmation prompt itself.
```
gh debug-cli --log-level none
Setting url to http://localhost:8080/agents/blackbeard

Start typing to chat with your assistant...
sparklyunicorn: confirmation
assistant: Ahoy, @monalisa! Ye be seekin' confirmation, me hearty. Are ye sure ye want a custom limerick 'bout petals? This here action be irreversible, matey!

Be ye sure ye want a custom limerick 'bout petals ?
  Arrr, this here action be irreversible, matey!
Reply: [y/N]
```
8. Currently, the supported event types for debug mode are references, errors, and confirmations! Have fun chatting with your assistant!


## Using the debug stream tool
1. to quickly parse agent response run cmd go run `main.go stream [local file name]` for example `gh-debug-cli stream test.txt`   

2. tool will take print out data packet for easy readability
```
data: {"choices":[{"delta":{"content":"A closure in JavaScript "}}],"created":1727120830,"id":"chatcmpl-AAjJW0Nz9E2Gu1P6YQMFqqmn10mdR","model":"gpt-4o-2024-05-13","system_fingerprint":"fp_80a1bad4c7"}
data: {"choices":[{"delta":{"content":"is a function that retains access "}}],"created":1727120831,"id":"chatcmpl-AAjJW0Nz9E2Gu1P6YQMFqqmn10mdR","model":"gpt-4o-2024-05-13","system_fingerprint":"fp_80a1bad4c7"}
data: {"choices":[{"delta":{"content":"to its lexical scope, even "}}],"created":1727120832,"id":"chatcmpl-AAjJW0Nz9E2Gu1P6YQMFqqmn10mdR","model":"gpt-4o-2024-05-13","system_fingerprint":"fp_80a1bad4c7"}
data: {"choices":[{"delta":{"content":"when the function is executed "}}],"created":1727120833,"id":"chatcmpl-AAjJW0Nz9E2Gu1P6YQMFqqmn10mdR","model":"gpt-4o-2024-05-13","system_fingerprint":"fp_80a1bad4c7"}
data: {"choices":[{"delta":{"content":"outside that scope. "}}],"created":1727120834,"id":"chatcmpl-AAjJW0Nz9E2Gu1P6YQMFqqmn10mdR","model":"gpt-4o-2024-05-13","system_fingerprint":"fp_80a1bad4c7"}
```