package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/go-vgo/robotgo"
	"github.com/gocolly/colly"
	hook "github.com/robotn/gohook"
)

var secondsToStart = flag.Int("start", 10, "Segundos a esperar antes de empezar a escribir")
var msWriteLatency = flag.Int("speed", 100, "Velocidad de escritura, minimo 5")

func main() {
	// parseamos las flags (opciones a la hora de ejecutar la aplicacion desde una consola)
	flag.Parse()
	// le damos un valor al Seed del paquete random para conseguir numeros aleatorios
	rand.Seed(time.Now().UnixNano())
	// iniciamos la escucha de la tecla ESC para salir
	go events()

	// colly es una libreria para leer paginas web
	c := colly.NewCollector()
	// mw-content-text es el ID del bloque central de texto de la página de wikipedia
	c.OnHTML("#mw-content-text", func(e *colly.HTMLElement) {
		// por si acaso, si el texto esta vacío, terminamos la funcion
		if e.Text == "" {
			return
		}
		// dividimos el texto letra a letra
		text := strings.Split(e.Text, "")
		for { // bucle infinito
			for _, s := range text { // Por cada letra
				robotgo.TypeStr(s)            // escribimos la letra
				if s == "\n" || s == "\r\n" { // si es un salto de linea, le damos a la tecla Enter
					robotgo.KeyTap("enter")
				}
				robotgo.MilliSleep(*msWriteLatency) // esperamos unos milisegundos
			}
		}
	})

	fmt.Println("Presiona ESC para terminar")
	// escribimos por pantalla los segundos que quedan para empezar
	for *secondsToStart > 0 {
		fmt.Println("Empiezo a escribir en", *secondsToStart, "segundos!")
		timer1 := time.NewTimer(time.Second) // creamos un timer en Go
		<-timer1.C                           // esperamos a que el timer nos diga que sigamos
		*secondsToStart--
	}
	// enviamos el collector de colly a leer esta URL
	c.Visit("https://es.wikipedia.org/wiki/Special:Random") // esta URL es una URL especial de wikipedia que devuelve un articulo aleatorio
}

func events() {
	robotgo.EventHook(hook.KeyDown, []string{"esc"}, func(e hook.Event) {
		robotgo.EventEnd()
		os.Exit(1)
	})

	s := robotgo.EventStart()
	<-robotgo.EventProcess(s)
}
