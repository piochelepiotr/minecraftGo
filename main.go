package main

const windowWidth = 800
const windowHeight = 600

func main() {
    d := DisplayManager{windowWidth:windowWidth, windowHeight:windowHeight}
    defer d.closeDisplay()
    d.createDisplay()
	for !d.window.ShouldClose() {
        d.updateDisplay()
	}
}
