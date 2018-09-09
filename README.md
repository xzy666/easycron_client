## struct

1.cron(command) 定时执行某个命令行 每个都是在独立的goroutine中
2.S端给过来的永远是一个封装的Task
3.C端会存一个JobList，记录所有已经S端交付在Cron运行的任务
4.S端的数据入库
5.C端的数据实时保存在内存中(必须保证内存的使用量)

