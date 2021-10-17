package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

const (
	END_BYTES   = "\000\001\002\003\004\005" //информация об окончании передачи данных для сервера
	ADDR_SERVER = ":8080"
)

func main() {
	//подключение к tcp с портом ADDR_SERVER
	conn, err := net.Dial("tcp", ADDR_SERVER)
	if err != nil {
		panic("failed to connect server")
	}
	//закрыть соединение после обработки
	defer conn.Close()
	//получение строки от клиента и преобразование в байты
	go ClientOutput(conn)
	ClientInput(conn)
}

//ввод данных и отправка клиенту
func ClientInput(conn net.Conn) {
	for {
		conn.Write([]byte(InputString() + END_BYTES))
	}
}

//принимает данные со стороны сервера
func ClientOutput(conn net.Conn) {
	//прочитывание полученных данных
	var (
		buffer  = make([]byte, 512) //сохранение промежуточного результата до 512 байт
		message string              //результат который мы приняли от сервера
	)
close:
	for {
		message = ""
		//цикл при помощи которого мы будем читать данные и заносить их в буфер, а из буфера заносить в сообщение
		for {
			length, err := conn.Read(buffer) //указывает буфер, куда будут заноситься данные
			if err != nil {
				break close
			}
			//если все хорошо, заносим в сообщение содержимое буфера до конца length и переводим в строку
			message += string(buffer[:length])
			if strings.HasSuffix(message, END_BYTES) {
				message = strings.TrimSuffix(message, END_BYTES)
				break
			}
		}
		//читаем сообщение, которое обработал сервер
		fmt.Println(message)
	}
}

//фунция получает строку от клиента
func InputString() string {
	//считывает текст с клавиатуры пока не будет введен \n
	msg, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	//удаляем символ \n
	return strings.Replace(msg, "\n", "", -1)
}

//TODO 22:30 ДОДЕЛАТЬ!
