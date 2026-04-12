# README

## Run dev

```
cd bumashtina
wails dev
```

*If go is not in PATH:*

```
export PATH=$PATH:$(go env GOPATH)/bin
source ~/.bashrc

## Live Development

If you want to develop in a browser
and have access to your Go methods, there is also a dev server that runs on http://localhost:34115. Connect to this in your browser, and you can call your Go code from devtools.

## Building

To build a redistributable, production mode package, use `wails build`.
