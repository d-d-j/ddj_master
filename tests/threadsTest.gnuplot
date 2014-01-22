set term eps
set output "threadsTest.eps"
set xlabel "Zapytanie"
set ylabel "Czas odpowiedzi [ms]"
set xtics 20000

plot "threads1.temp" using 9 smooth sbezier with lines title "1 thread", "threads10.temp" using 9 smooth sbezier with lines title "10 threads", "threads20.temp" using 9 smooth sbezier with lines title "20 threads", "threads50.temp" using 9 smooth sbezier with lines title "50 threads", "threads100.temp" using 9 smooth sbezier with lines title "100 threads"
