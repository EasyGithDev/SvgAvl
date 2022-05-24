# SvgAvl

Here is a complete GO program which reads an integer list and display an AVL tree. 

Both the insertion and the drawing routine are recursives. 

Finally, the program send the SVG associated with the tree to the browser to display.

## Install

Select or create a folder :

```sh
cd myfloder
```

Clone the project into your selected folder :

```sh
git clone git@github.com:EasyGithDev/SvgAvl.git avltree
```

Install the depencies to work with SVG :

```sh
cd svgtree
go mod init
go get github.com/ajstarks/svgo
```

## Run

You may execute the program with a short integer list as parameter :

```sh
go run main.go 8 5 4 3 1 -1
```

![alt text](../assets/avl.svg?raw=true)


If you want display the node position use the -d option like this : 

```sh
go run main.go -d=p 8 5 4 3 1 -1
```

![alt text](../assets/avl-p.svg?raw=true)

## Display the result

Open a web browser and enter the URL :

http://localhost:8000/

## Write the result

You can choose to generate a SGV file to save the result.

You must change the output like this :

```sh
go run main.go -o=stdout 8 5 4 3 1 -1 > avl.svg
```
