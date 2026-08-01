# README

<h3>BG</h3>

<a href="https://github.com/sajax88/bumashtina/releases/latest">Свали последната версия</a>

Тази програмка може да ви бъде полезна, ако сте самоосигуряващо се лице по смисъла на КСО.

Функции:
- генерира Декларация 1 и Декларация 6, които се качват в портала на НАП https://portal.nra.bg/
- изчислява данък за тримесечие за подаване на Декларация по чл. 55

Всичко се изчислява и се запазва локално. За повече информация вижте инструкциите долу. 

<img width="1022" height="521" alt="image" src="https://github.com/user-attachments/assets/7717fcc1-59f0-4537-8047-e0467b821845" />

<h3>EN</h3>

<a href="https://github.com/sajax88/bumashtina/releases/latest">Download the latest version</a>

This is a desktop app for all the poor souls in the Republic of Bulgaria who happen to pay their taxes and social security themselves. Only for the **self-employed** people who know what they're doing!
Check if this is your case: https://nra.bg/wps/portal/nra/osiguryavane/osiguryavam-se-sam#osigurqvam-se-sam0

App language: **BG only**!

The app can help you generate Declarations One and Six that can then be submitted to the NRA portal. You can also calculate your taxes and keep track of your general income.
The Declarations output is the same as the original NRA software gives you, but hopefully the interface is slightly better.

The data is stored locally in the user's directory. 

![dylan_moran_black_books.jpeg](dylan_moran_black_books.jpeg)

(Dylan Moran, the one and only)

## I just want to run the program

Download the latest version for your OS from [releases](https://github.com/sajax88/bumashtina/releases), unpack and run the executable file. 
- On Ubuntu 24 download **bumashtina_with_libwebkit2gtk_41** version
- On Linux you'll need to make the file executable first!
- On MacOS rename the file to **bumashtina.app**. If the app is considered unsafe, remove the quarantine attribute: `xattr -dr com.apple.quarantine ./bumashtina.app`

## I want to customize or contribute

Cool! The tech stack is Wails plus Svelte.

The following are the standard Wails commands for development. See Wails docs for more: https://wails.io/docs/introduction
The current framework version is 2.13.0.

### Run dev

```
wails dev
```

*If go is not in PATH:*

```
export PATH=$PATH:$(go env GOPATH)/bin
source ~/.bashrc
```

### Live Development

If you want to develop in a browser
and have access to your Go methods, there is also a dev server that runs on http://localhost:34115. Connect to this in your browser, and you can call your Go code from devtools.

### Building

To build a redistributable, production mode package, use `wails build`. See Wails docs for details: https://wails.io/docs/reference/cli#build 
On Ubuntu 24 you have to build with this tag:
`wails build -tags webkit2_41`

### Unit tests

Don't forget to run the unit tests before committing. Feel free to add more.


## I like your app and want to buy you a coffee/beer

You certainly can!

https://buymeacoffee.com/sajax

You can also consider buying a coffee or contributing to [Wails](https://wails.io/) or [Svelte](https://svelte.dev/) or [Lucide](https://lucide.dev/) 
