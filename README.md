# element-plus

This template should help get you started developing with Vue 3 in Vite.

## Recommended IDE Setup

[VSCode](https://code.visualstudio.com/) + [Volar](https://marketplace.visualstudio.com/items?itemName=Vue.volar) (and disable Vetur) + [TypeScript Vue Plugin (Volar)](https://marketplace.visualstudio.com/items?itemName=Vue.vscode-typescript-vue-plugin).

## Customize configuration

See [Vite Configuration Reference](https://vitejs.dev/config/).

## Project Setup

```sh
go clean -modcache
rm go.sum
go get github.com/gta-spec/utils@latest
go mod tidy
```

### Compile and Hot-Reload for Development

```sh
go run .
```

### Compile and Minify for Production

```sh
go build -tags release -o myapp
go build -tags debug -o myapp
go build -tags test -o myapp
```
