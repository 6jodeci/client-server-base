package main

import (
	"fmt"
	"net"
	"strings"
)

const (
	END_BYTES = "\000\001\002\003\004\005" //информация об окончании передачи данных для сервера
	PORT      = ":8080"
)

func main() {
	//сервер прослушивает данные на tcp на своем порту
	listen, err := net.Listen("tcp", PORT)
	if err != nil {
		//сервер не смог подняться
		panic("server error")
	}
	//закрытие соединение после завершения всех действий сервера (в данном случае работает бесконечно)
	defer listen.Close()
	//сервер принимает соединение по указанному порту
	for {
		//принятие соединения
		conn, err := listen.Accept()
		if err != nil {
			break
		}
		go handleConnect(conn) // коннект пареллелен относительно функции main, потому-что соединений может быть множество
	}
}

func handleConnect(conn net.Conn) {
	defer conn.Close() //завершить соединение по окончанию функции
	var (
		buffer  = make([]byte, 512) //сохранение промежуточного результата до 512 байт
		message string              //результат который мы приняли от сервера

	)
	//цикл при помощи которого мы будем читать данные и заносить их в буфер, а из буфера заносить в сообщение
	for {
		length, err := conn.Read(buffer) //указывает буфер, куда будут заноситься данные
		if length == 0 || err != nil {   //проверка длины, чтобы не равнялась нулю и обработка ошибки
			break
		}
		//если все хорошо, заносим в сообщение содержимое буфера до конца length и переводим в строку
		message += string(buffer[:length])
		//проверка находятся ли байты окончания в конце сообщения
		if strings.HasSuffix(message, END_BYTES) {
			message = strings.TrimSuffix(message, END_BYTES)
			break
		}
	}
	//выводит введенные клиентом сообщения на сервер 
	fmt.Println(message)
	//переводит message в  UpperCase
	conn.Write([]byte(strings.ToUpper(message)))
}
//TODO 22:30 ДОДЕЛАТЬ!