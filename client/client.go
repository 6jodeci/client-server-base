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
		panic("failed to connect to server")
	}
	//закрыть соединение после обработки
	defer conn.Close()
	//получение строки от клиента и преобразование в байты
	conn.Write([]byte(InputString() + END_BYTES))
	//прочитывание полученных данных
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
	}
	//читаем сообщение, которое обработал сервер
	fmt.Println(message)
}

//фунция получает строку от клиента
func InputString() string {
	//считывает текст с клавиатуры пока не будет введен \n
	msg, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	//удаляем символ \n
	return strings.Replace(msg, "\n", "", -1)
}
//TODO 22:30 ДОДЕЛАТЬ!