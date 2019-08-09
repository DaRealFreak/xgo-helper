## xgo-helper

Wrapper for docker image cross for compiling applications local projects and fixing the local build on windows.  
This project was primarily written for a small build pipeline for a CGO dependency dependant project where goreleaser
is not viable, so the possibility for you requiring this helper is probably quite unlikely.

### Usage
```
Usage:
  xgo_helper [flags]

Flags:
      --dest string      Destination folder to put binaries in (empty = current)
  -h, --help             help for app
      --image string     Docker Image used (default "diebaumchen/xgo")
      --module string    Module name for local compilation (empty = external git repository)
      --out string       Prefix to use for output naming
      --pkg string       Package of main.go
      --source string    Repository source (branch/tag/commit hash)
      --targets string   Build targets
```

The helper will send external repository requests to the original xgo module, so you can use all additional xgo arguments as well.  
Only the arguments above are implemented for local builds though.