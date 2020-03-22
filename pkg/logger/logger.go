package logger

import (
	"fmt"
	"log"
	"os"
)

const (
	er    string = "\x1b[31;1mERROR: \x1b[0m"
	fatal string = "\x1b[31;1mFATAL: \x1b[0m"
)

//Info вывод информационного лога в консоль
func Info(v ...interface{}) {
	l := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	if err := l.Output(2, fmt.Sprintln(v...)); err != nil {
		fmt.Println("ERROR: не удалось вывести сообщение в лог. Сообщение:", fmt.Sprintf("%v", v...))
	}
}

//Fatal вывод  ошибки типа Fatal в консоль. Программа завершает работу
func Fatal(v ...interface{}) {
	l := log.New(os.Stderr, fatal, log.Ldate|log.Ltime|log.Lshortfile)
	if err := l.Output(2, fmt.Sprintln(v...)); err != nil {
		fmt.Println("ERROR: не удалось вывести сообщение в лог. Сообщение:", fmt.Sprintf("%v", v...))
	}
	os.Exit(1)
}

//Fatalf форматированный вывод ошибки типа Fatal в консоль. Программа завершает работу
func Fatalf(mes string, v ...interface{}) {
	l := log.New(os.Stderr, fatal, log.Ldate|log.Ltime|log.Lshortfile)
	if err := l.Output(2, fmt.Sprintf(mes, v...)); err != nil {
		fmt.Println("ERROR: не удалось вывести сообщение в лог. Сообщение:", fmt.Sprintf(mes, v...))
	}
	os.Exit(1)
}

//Error вывод лога ошибки в консоль
func Error(v ...interface{}) {
	l := log.New(os.Stderr, er, log.Ldate|log.Ltime|log.Lshortfile)
	if err := l.Output(2, fmt.Sprintln(v...)); err != nil {
		fmt.Println("ERROR: не удалось вывести сообщение в лог. Сообщение:", fmt.Sprintf("%v", v...))
	}
}

//Infof вывод фоматированного информационного лога в консоль
func Infof(mes string, v ...interface{}) {
	l := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	if err := l.Output(2, fmt.Sprintf(mes, v...)); err != nil {
		fmt.Println("ERROR: не удалось вывести сообщение в лог. Сообщение:", fmt.Sprintf(mes, v...))
	}
}

//Errorf вывод форматированного лога ошибки в консоль
func Errorf(mes string, v ...interface{}) {
	l := log.New(os.Stderr, er, log.Ldate|log.Ltime|log.Lshortfile)
	if err := l.Output(2, fmt.Sprintf(mes, v...)); err != nil {
		fmt.Println("ERROR: не удалось вывести сообщение в лог. Сообщение:", fmt.Sprintf(mes, v...))
	}
}