set xlabel "Liczba klientów wysyłających zapytania"
set ylabel "Liczba obsłużonych zapytań na sekundę"

set term eps
set output "concurrencyTest.eps"

plot 'newtest.tst' using 8:4 smooth sbezier notitle

set xlabel "Liczba klientów wysyłających zapytania"
set ylabel "Liczba obsłużonych zapytań na sekundę na klienta"

set output

set term eps
set output "concurrencyTestPerThread.eps"

plot 'newtest.tst' using 8:($4/$8) smooth sbezier notitle
