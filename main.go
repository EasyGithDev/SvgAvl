package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	svg "github.com/ajstarks/svgo"
)

const (
	canvasWidth  = 800
	canvasHeight = 600
)

var (
	fontSize  string = "15"
	textStyle string = "text-anchor:middle;font-size:" + fontSize + "px;fill:black"
	lineStyle string = "stroke:rgb(255,0,0);stroke-width:2"
	rectStyle string = "fill:rgb(255,255,255);stroke-width:2;stroke:rgb(0,0,0)"
	dx        float64
	dy        float64
	mx        float64
	my        float64

	// Frequence format
	format = func(t *Tree) string {
		return fmt.Sprintf("%d", t.val)
	}

	// Position format
	// format = func(t *Tree) string {
	// 	return t.String()
	// }
)

// Tree position
type Point struct {
	X int
	Y int
}

func (p *Point) String() string {
	return fmt.Sprintf("(x:%d, y:%d)", p.X, p.Y)
}

func NewPoint() *Point {
	return &Point{0, 0}
}

// Tree
type Tree struct {
	*Point       // t position
	val    int   // the integer value
	left   *Tree // left child
	right  *Tree // right child
}

func (n *Tree) String() string {
	return fmt.Sprintf("%d %s", n.val, n.Point)
}

func NewTree() *Tree {
	return &Tree{NewPoint(), 0, nil, nil}
}

// Compute the tree height

// Max returns the larger of x or y.
func Max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

func Height(t *Tree) int {
	if t == nil {
		return 0
	}
	return 1 + Max(Height(t.left), Height(t.right))
}

// Display tree in pre order
func Prefixe(t *Tree) {
	if t == nil {
		return
	}
	fmt.Print(t, " ")
	Prefixe(t.left)
	Prefixe(t.right)
}

// Display tree in order
func Infixe(t *Tree) {
	if t == nil {
		return
	}
	Infixe(t.left)
	fmt.Print(t, " ")
	Infixe(t.right)
}

// Display tree in post order
func Postfixe(t *Tree) {
	if t == nil {
		return
	}
	Postfixe(t.left)
	Postfixe(t.right)
	fmt.Print(t, " ")
}

func rotateL(n *Tree) *Tree {
	tmp := n.right
	n.right = tmp.left
	tmp.left = n
	return tmp
}

func rotateR(n *Tree) *Tree {
	tmp := n.left
	n.left = tmp.right
	tmp.right = n
	return tmp
}

func avl(n *Tree) *Tree {

	fmt.Printf("node %s , left %s, right %s\n", n, n.left, n.right)

	if Height(n.left)-Height(n.right) == 2 {

		// Je fais la rotation G le sous arbre gauche
		if Height(n.left.left) < Height(n.left.right) {
			n.left = rotateL(n.left)
		}

		// Dans tous les cas je fais la rotation simple
		return rotateR(n)
	}

	if Height(n.left)-Height(n.right) == -2 {
		// Je fais la rotation inverse sur le sous arbre droit

		// Je fais la rotation D le sous arbre droit
		if Height(n.right.right) < Height(n.right.left) {
			n.right = rotateR(n.right)
		}

		// Dans tous les cas je fais la rotation simple
		return rotateL(n)
	}

	return n
}

// Insert a value in tree
func Insert(t *Tree, val int) *Tree {
	if t == nil {
		t := NewTree()
		t.val = val
		return t
	}

	if val < t.val {
		t.left = Insert(t.left, val)
	} else if val > t.val {
		t.right = Insert(t.right, val)
	}

	fmt.Println("Avl", t)
	return avl(t)
}

// Search a value in tree
func Search(t *Tree, val int) bool {

	res := false

	if t == nil {
		return res
	} else if t.val == val {
		res = true
	} else if t.val > val {
		return Search(t.left, val)
	} else {
		return Search(t.right, val)
	}

	return res
}

// Compute the position of each sub trees in the tree
func Position(t *Tree, x int, y int) int {

	if t.left != nil {
		x = Position(t.left, x, y+1)
	}

	t.X = x
	t.Y = y

	x = x + 1

	if t.right != nil {
		x = Position(t.right, x, y+1)
	}
	return x
}

// Drawing the sub trees in SVG
func Draw(t *Tree, canvas *svg.SVG) {
	if t == nil {
		return
	}

	h, _ := strconv.Atoi(fontSize)
	x1 := int(dx*float64(t.X) + mx)
	y1 := int(dy*float64(t.Y) + my)

	if t.left != nil {
		left := t.left
		x2 := int(dx*float64(left.X) + mx)
		y2 := int(dy*float64(left.Y) + my)
		// log.Printf("x1:%d y1:%d x2:%d y2:%d %s\n", x1, y1, x2, y2, lineStyle)

		canvas.Line(x1, y1+(h/2), x2, y2-h, lineStyle)
		Draw(left, canvas)
	}

	if t.right != nil {
		right := t.right
		x2 := int(dx*float64(right.X) + mx)
		y2 := int(dy*float64(right.Y) + my)
		// log.Printf("x1:%d y1:%d x2:%d y2:%d %s\n", x1, y1, x2, y2, lineStyle)

		canvas.Line(x1, y1+(h/2), x2, y2-h, lineStyle)
		Draw(right, canvas)
	}

	canvas.Text(x1, y1, format(t), textStyle)
}

// Drawing the tree in  SVG
func Display(t *Tree, w io.Writer) {

	tWidth := Position(t, 0, 0)
	tHeight := Height(t)
	fmt.Println(tHeight)
	canvas := svg.New(w)
	canvas.Start(canvasWidth, canvasHeight)
	canvas.Rect(0, 0, canvasWidth, canvasHeight, rectStyle)

	dx = float64(canvasWidth) / float64(tWidth)
	dy = float64(canvasHeight) / float64(tHeight)
	mx = dx / 2
	my = dy / 2

	Draw(t, canvas)
	canvas.End()
}

func main() {

	if len(os.Args[1:]) == 0 {
		log.Println("You must enter somme words ...")
		os.Exit(1)
	}

	// Create the t
	var t *Tree

	for _, v := range os.Args[1:] {
		val, _ := strconv.Atoi(v)
		t = Insert(t, val)
	}

	// Send result to stdout
	// Display(t, os.Stdout)

	// Display the tree on Web browser
	s := ""
	buf := bytes.NewBufferString(s)
	Display(t, buf)

	// Send the output to the client
	http.Handle("/", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "image/svg+xml")
			w.Write(buf.Bytes())
		}))
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}

}
