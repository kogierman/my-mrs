# My Merge Requests CLI lister for Gitlab

## Why

I simply lacked overview of all my open merge requests in Gitlab which Github has ;)

## Installation and usage

Requirements:
- `go` 

Set `GITLAB_RO_TOKEN` environment variable to your [Gitlab Personal Access Token](https://docs.gitlab.com/ee/user/profile/personal_access_tokens.html#create-a-personal-access-token) (required scope `read_api`).

Then you can just run it:

```
$ go run .

Your merge requests:
=====
1: chore: release some service [ review: ✘ ] (https://gitlab.com/XDorg/XDproject/-/merge_requests/69)
2: feat: add new feature [ review: ✔ ] (https://gitlab.com/XDorg/XDproject/-/merge_requests/420)
3: fix: fix nasty bug [ review: ✔ ] (https://gitlab.com/XDorg/XDproject/-/merge_requests/2137)
``` 

Optional flags:

- `-t` - Gitlab token, overrides `GITLAB_RO_TOKEN`
- `-a` - print all merge requests, also merged and closed (merged are marked green, closed - red)


## License
![WTFPL](http://www.wtfpl.net/wp-content/uploads/2012/12/wtfpl-badge-1.png)

DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE

Version 2, December 2004 


Copyright (C) 2004 Sam Hocevar <sam@hocevar.net> 


Everyone is permitted to copy and distribute verbatim or modified copies of this license document, and changing it is allowed as long as the name is changed. 

DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE 

TERMS AND CONDITIONS FOR COPYING, DISTRIBUTION AND MODIFICATION 

0. You just DO WHAT THE FUCK YOU WANT TO.