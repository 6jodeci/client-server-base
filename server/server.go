package main

import (
	"log"
	"net"
	"strings"
)

const (
	END_BYTES = "\000\001\002\003\004\005" //информация об окончании передачи данных для сервера
	PORT      = ":8080"                    //порт
)

var (
	Connections = make(map[net.Conn]bool) //хранит соединения клиетов
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
	Connections[conn] = true //добавляет connection
	var (
		buffer  = make([]byte, 512) //сохранение промежуточного результата до 512 байт
		message string              //результат который мы приняли от сервера
	)
	//бесконечный цикл открытого соединения
close:
	for {
		//чтобы строка не накапливала промежуточный результат
		message = ""
		//цикл при помощи которого мы будем читать данные и заносить их в буфер, а из буфера заносить в сообщение
		for {
			length, err := conn.Read(buffer) //указывает буфер, куда будут заноситься данные
			if err != nil {
				break close
			} //закрывает соединение
			//если все хорошо, заносим в сообщение содержимое буфера до конца length и переводим в строку
			message += string(buffer[:length])
			//проверка находятся ли байты окончания в конце сообщения
			if strings.HasSuffix(message, END_BYTES) {
				message = strings.TrimSuffix(message, END_BYTES)
				break
			}
		}
		//логирует введенные клиентом сообщения на сервер
		log.Println(message)
		//отправка сообщения всем остальным пользователям
		for c := range Connections {
			if c == conn {
				continue
			} //если это отправляюший пользователь, то пропускаем
			conn.Write([]byte(strings.ToUpper(message) + END_BYTES)) //переводит message в UpperCase с условием конца строки
		}
	}
	delete(Connections, conn) //удаление соединения
}

//TODO 22:30 ДОДЕЛАТЬ!
