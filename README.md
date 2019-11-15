
<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.0 Transitional//EN">
<html>
<body lang="ru-RU" dir="ltr">

# tcp-client

program build Linux32

tcp-client -help
  -dlevel int
	Уровень отладки. 0 - Err, 1 - Info, 2 - All
  -dout string
	Путь расположения исходящих файлов. (default "out")
  -fout string
	Расширение исходящего файла. (default "log")
  -ipaddr string
	IP аддрес сервера. (default "127.0.0.1")
  -ipport string
	IP порт сервера. (default "10113")
  -maxproc int
	Максимальное кол-во одновременных потоков. (default 1)
  -trotate duration
	Период создания иходящего файла. Формат: 10s = 10 секунд, 10m = 10 минут, 10h = 10 часов, 10d = 10 дней и т.д. (default 1h0m0s)

build go version go1.13.3 linux/amd64

</body>
</html>

