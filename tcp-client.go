// tcp-client
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"runtime"
	"time"

	"gopkg.in/robfig/cron.v3"
)

func main() {
	p := fmt.Println
	t := time.Now()
	c := cron.New()

	var lsconn net.Conn
	var err error
	var buf *bufio.Reader

	var maxproc int
	var dlevel int
	var fout string
	var ipaddr string
	var ipport string
	var dout string
	var rotate int
	var trotate time.Duration
	var crotate string

	flag.IntVar(&maxproc, "maxproc", 1, "Максимальное кол-во одновременных потоков.")
	flag.IntVar(&dlevel, "dlevel", 0, "Уровень отладки. 0 - Err, 1 - Info, 2 - All")
	flag.StringVar(&ipaddr, "ipaddr", "127.0.0.1", "IP аддрес сервера.")
	flag.StringVar(&ipport, "ipport", "10113", "IP порт сервера.")
	flag.StringVar(&fout, "fout", "log", "Расширение исходящего файла.")
	flag.StringVar(&dout, "dout", "out", "Путь расположения исходящих файлов.")
	flag.IntVar(&rotate, "rotate", 1, "Метод ротации файлов по времени. 0 - через указанный интервал согласно параметра 'trotate' или 1 - по переменной 'crotate'")
	flag.DurationVar(&trotate, "trotate", 3600000000000, "Период создания иходящего файла. Формат: 10s = 10 секунд, 10m = 10 минут, 10h = 10 часов, 10d = 10 дней и т.д.")
	flag.StringVar(&crotate, "crotate", "0 0-23 * * *", "Период создания исходящего файла в формате cron (https://en.wikipedia.org/wiki/Cron).")

	flag.Parse()

	c.Start()

	p("==========================================================================================")

	fmt.Println("dlevel:", dlevel)
	fmt.Println("ipaddr:", ipaddr)
	fmt.Println("ipport:", ipport)
	fmt.Println("fout:", fout)
	fmt.Println("dout:", dout)
	fmt.Println("rotate:", rotate)
	fmt.Println("trotate:", trotate)
	fmt.Println("crotate:", crotate)

	numcpu := runtime.NumCPU()
	fmt.Println("NumCPU", numcpu)
	fmt.Println("maxproc:", maxproc)
	runtime.GOMAXPROCS(maxproc)

	p("==========================================================================================")

	p(t.Format("2006-01-02 15:04:05") + " - Start.")
	//p(t.Format("200601021504"))
	//p(t.Format("15:04"))

	derr := os.MkdirAll(dout, 0775)
	if derr != nil {
		fmt.Println(derr)
		os.Exit(1)
	}

	c.AddFunc(crotate, func() {
		//fmt.Println("Cron funcion ---- " + time.Now().Format("2006-01-02 15:04:05"))
		if rotate == 1 {
			if len(crotate) != 0 {
				tf1 := time.Now()
				newFile := dout + "/" + tf1.Format("20060102150405") + "." + fout
				err := os.Rename(dout+"/out.tmp", newFile)
				if err != nil {
					fmt.Println(err)
				} else {
					if dlevel > 0 {
						fmt.Println(tf1.Format("2006-01-02 15:04:05") + " - Create NEW file: " + newFile)
						if err := os.Chmod(newFile, 0644); err != nil {
							fmt.Println(err)
						}
					}
				}
			} else {
				os.Exit(1)
			}
		}
	})

	if rotate == 0 {
		if trotate != 0 {
			go func() {
				for {
					select {
					case <-time.After(trotate):
						tf1 := time.Now()
						newFile := dout + "/" + tf1.Format("20060102150405") + "." + fout
						err := os.Rename(dout+"/out.tmp", newFile)
						if err != nil {
							fmt.Println(err)
						} else {
							if dlevel > 0 {
								fmt.Println(tf1.Format("2006-01-02 15:04:05") + " - Create NEW file: " + newFile)
								if err := os.Chmod(newFile, 0644); err != nil {
									fmt.Println(err)
								}
							}

						}
					}
				}
			}()
		}
	}

	connect := func() {
		// непрерывное подключение к сокету
		for {
			lsconn, err = net.Dial("tcp", ipaddr+":"+ipport)
			if err == nil {
				buf = bufio.NewReader(lsconn)
				break
			}
			fmt.Println(err)
			time.Sleep(time.Second * 10)
		}
	}

	connect()

	for {
		//select {
		//case <-time.After(time.Second * trotate):
		//	fmt.Println(time.Now())
		//default:
		message, err1 := buf.ReadString('\n')
		if err1 == io.EOF {
			if dlevel > 0 {
				p(t.Format("2006-01-02 15:04:05") + " - Error connect Server.")
				//fmt.Println(err1)
			}
			lsconn.Close()
			connect()
			time.Sleep(time.Second * 10)
		} else {
			if dlevel > 1 {
				fmt.Print(message)
			}

			f, errf := os.OpenFile(dout+"/out.tmp", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
			if errf != nil {
				fmt.Println(errf)
				os.Exit(1)
			}

			if _, err := f.Write([]byte(message)); err != nil {
				f.Close() // ignore error; Write error takes precedence
				fmt.Println(errf)
				os.Exit(1)
			}

			if errf := f.Close(); errf != nil {
				fmt.Println(errf)
				os.Exit(1)
			}
		}
		//}
	}

}
