# Budgie

A simple task runner\
No dependencies, blazing fast

## How to install
```bash
wget -qO- https://raw.githubusercontent.com/Hand-of-Doom/budgie/main/installer/install.sh | sudo sh -s 2.0.0
```
For the current user (no sudo required)
```bash
wget -qO- https://raw.githubusercontent.com/Hand-of-Doom/budgie/main/installer/install.sh | sh -s 2.0.0
```
You can install it as a go package if you have go installed
```bash
go install github.com/Hand-of-Doom/budgie/v2@latest
```

## Examples

### Hello World
```bash
PrintHelloWorld() {
    echo "Hello world!"
}
```

### Complex example
Use the power of bash
```bash
output=$(pwd)

for i in $(seq 1 5); do
    output="$output/$i"
done

Build() {
   go build -o $output
   add() {
       echo $(($1 + $2))
   }
   echo "compilation time is $(add 2 3) ms"
}
```
### Real world example
A single page application written in lit-html and go\
[Click to go to an example](examples/real_world/tasks.budgie)
```bash
output="$(pwd)/build"
mkdir -p $output

exe_file="$output/app.bin"

BuildBackend() {
    cd ./backend
    go build -o "$exe_file"
    cp ./config.yaml $output
}

BuildFrontend() {
    cd ./frontend
    [ ! -d ./node_modules ] && npm i
    public_dir="$output/public"
    npx esbuild ./app.js --bundle --minify --outfile="$public_dir/bundle.js"
    cp -a ./static/. "$public_dir/"
}

Build() { 
    BuildBackend
    BuildFrontend
}

Run() {
    Build
    cd $output
    $exe_file
}
```

## How to use

Create a tasks.sh file with the following content
```bash
PrintHello() {
    echo -n "Hello "
}

# inline
PrintWorld() { echo "world!" }

# with dependencies
PrintHelloWorld() {
    PrintHello
    PrintWorld
}
```

Then run a line below in your terminal
```bash
$ budgie printHelloWorld
```

The output
```bash
$ Hello world!
```

If you don't want to use tasks.sh as a filename, just pass a path to the file you want as the second argument
```bash
$ budgie printHelloWorld fileYouWant
```

These commands are the same
```bash
$ budgie printHelloWorld
```
```bash
$ budgie printHelloWorld tasks.sh
```
It uses tasks.sh as the default filename if no filename is passed

## More about file parsing

### Target
```bash
TargetName() {
    echo "this target is named targetName"
}
```
Target is a function declared in a pascal case\
This is an important detail

Try run the following code snippet via `budgie targetName`
```bash
targetName() {
    :
}
```
The output
```bash
$ TargetName: command not found
```
`targetName` is a function declared in a top-level scope,\
it can be invoked from the target as follows
```bash
targetName() {
    :
}

Target() {
    targetName
}
```

### Scopes
The global scope can be accessed from the target as follows
```bash
printHello() {
    echo -n "Hello "
}

printWorld() {
    echo "world!"
}

PrintHelloWorld() {
    printHello
    printWorld
}
```
The global scope below the target cannot be accessed
```bash
PrintHelloWorld() {
    printHello
    printWorld
}

printHello() { : }

printWorld() { : }
```
The output
```bash
$ printHello: command not found
```
