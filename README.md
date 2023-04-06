# Budgie

A simple task runner\
No dependencies, blazing fast

## How to install
```bash
sudo sh -c "$(wget -qO- https://raw.githubusercontent.com/Hand-of-Doom/budgie/main/installer/install.sh)" 1.0.1
```
You can install it as a go package if you have go installed
```bash
go install github.com/Hand-of-Doom/budgie@latest
```

## Examples

### Hello World
```bash
target printHelloWorld:
    echo "Hello world!"
end
```

### Complex example
Use the power of bash
```bash
output=$(pwd)

for i in $(seq 1 5); do
    output="$output/$i"
done

target build:
   go build -o $output
   add() {
       echo $(($1 + $2))
   }
   echo "compilation time is $(add 2 3) ms"
end
```
### Real world example
A single page application written in lit-html and go\
[Click to go to an example](examples/real_world/tasks.budgie)
```bash
output="$(pwd)/build"
mkdir -p $output

exe_file="$output/app.bin"

target buildBackend:
    cd ./backend
    go build -o "$exe_file"
    cp ./config.yaml $output
end

target buildFrontend:
    cd ./frontend
    [ ! -d ./node_modules ] && npm i
    public_dir="$output/public"
    npx esbuild ./app.js --bundle --minify --outfile="$public_dir/bundle.js"
    cp -a ./static/. "$public_dir/"
end

target build(buildBackend buildFrontend):end

target run(build):
    cd $output
    $exe_file
end
```

## How to use

Create a tasks.budgie file with the following content
```bash
target printHello:
    echo -n "Hello "
end

# inline
target printWorld: echo "world!" end

# with dependencies
target printHelloWorld(printHello printWorld):end
```

Then run a line below in your terminal
```bash
$ budgie printHelloWorld
```

How it works
1. Reads the targets from the tasks.budgie file in your working directory
2. Finds a target named printHelloWorld
3. Executes it and the targets depends on it

The output
```bash
$ Hello world!
```

If you don't want to use tasks.budgie as a filename, just pass the path to the file you want as the second argument
```bash
$ budgie printHelloWorld anotherFile.budgie
```

These commands are the same
```bash
$ budgie printHelloWorld
```
```bash
$ budgie printHelloWorld tasks.budgie
```
It uses tasks.budgie as the default filename if no filename is passed

## More about budgie files

### Target
```bash
target targetName:
    echo "this target is named targetName"
end
```
You can use all characters in the target names except spaces\
These names are valid: target-name, $&*targetName

You can declare a target that depends on other targets\
They will run after each other
```bash
target targetName1: echo "the first line displayed" end
target targetName2: echo "the second line displayed" end

target print(targetName1 targetName2):end
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

target printHelloWorld:
    printHello
    printWorld
end
```
The global scope below the target cannot be accessed
```bash
target printHelloWorld:
    printHello
    printWorld
end

printHello() {}

printWorld() {}
```
The output
```bash
printHello: command not found
```
