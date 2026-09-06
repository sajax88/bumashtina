# README

<h3>BG</h3>

<b><a href="https://github.com/sajax88/bumashtina/releases/latest">Свали последната версия</a></b>

Тази програмка може да ви бъде полезна, ако сте самоосигуряващо се лице по смисъла на КСО: https://nra.bg/wps/portal/nra/osiguryavane/osiguryavam-se-sam#osigurqvam-se-sam0.

<h4>Функции:</h4>

- генерира Декларация 1 и Декларация 6, които се качват в портала на НАП https://portal.nra.bg/
- изчислява данък за тримесечие за подаване на Декларация по чл. 55
- изчислява окончателните годишни осигуровки и данъци

Всичко се изчислява и се запазва локално. За повече информация вижте инструкциите долу. 

<img width="1069" height="699" alt="image" src="https://github.com/user-attachments/assets/ababebb9-551e-4728-8642-3fba5a743464" />

## Как да пусна програмата

Последната версия се намира [тук](https://github.com/sajax88/bumashtina/releases), разопаковайте архива и пуснете изпълнимия файл.
При свалянето на новата версия просто заменяте файла с по-новия.

- Windows: **bumashtina.exe-windows-amd64.zip**
- MacOS: **bumashtina-darwin-arm64.zip** (това е универсална версия за arm64 и x86_64)
- Linux with libwebkit2gtk 4.0 (например Ubuntu 22): **bumashtina-linux-amd64.zip**
- Linux with libwebkit2gtk 4.1 (например Ubuntu 24): **bumashtina_with_libwebkit2gtk_41**

- Windows може да се оплаква, че файлът не е безопасен; трябва да разрешите да го изпълни
- На MacOS преименувайте файла в **bumashtina.app**. Ако казва, че файлът не е безопасен, трябва да изпълните тази команда: `xattr -dr com.apple.quarantine ./bumashtina.app`
- На Linux трябва да промените правата на файла, така че да стане изпълним: `chmod +x bumashtina`


<h3>EN</h3>

<b><a href="https://github.com/sajax88/bumashtina/releases/latest">Download the latest version</a></b>

This is a desktop app for all the poor souls in the Republic of Bulgaria who happen to pay their taxes and social security themselves. Only for the **self-employed** people who know what they're doing!
Check if this is your case: https://nra.bg/wps/portal/nra/osiguryavane/osiguryavam-se-sam#osigurqvam-se-sam0

App language: **BG only**!

The app can help you generate Declarations One and Six that can then be submitted to the NRA portal. You can also calculate your social security payments and taxes and keep track of your general income.
The Declarations output is the same as the original NRA software gives you, but hopefully the interface is slightly better.

The data is stored locally in the user's directory. 

![dylan_moran_black_books.jpeg](dylan_moran_black_books.jpeg)

(Dylan Moran, the one and only)

## I just want to run the program

Download the latest version for your OS from [releases](https://github.com/sajax88/bumashtina/releases), unpack and run the executable file. 
When you download a new version, just replace the old file with the new one.
- Windows: **bumashtina.exe-windows-amd64.zip**
- MacOS: **bumashtina-darwin-arm64.zip** (this is a universal version for arm64 and x86_64)
- Linux with libwebkit2gtk 4.0 (e.g. Ubuntu 22): **bumashtina-linux-amd64.zip**
- Linux with libwebkit2gtk 4.1 (e.g. Ubuntu 24): **bumashtina_with_libwebkit2gtk_41** 

- Windows might complain that the file is unsafe; you'll have to allow to run it
- On MacOS rename the file to **bumashtina.app**. If the app is considered unsafe, remove the quarantine attribute: `xattr -dr com.apple.quarantine ./bumashtina.app`
- On Linux you'll need to make the file executable: `chmod +x bumashtina`


## I want to customize or contribute

Cool! The tech stack is Wails plus Svelte.

The following are the standard Wails commands for development. See Wails docs for more: https://wails.io/docs/introduction
We use *Wails 2* and *Svelte 3*.

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

Don't forget to run the unit tests with `go test` before committing. Feel free to add more.


## I like your app and want to buy you a coffee/beer

You certainly can!

https://buymeacoffee.com/sajax

You can also consider buying a coffee or contributing to [Wails](https://wails.io/) or [Svelte](https://svelte.dev/) or [Lucide](https://lucide.dev/) 
