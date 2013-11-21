set xlabel "Liczba klientów wysyłających zapytania"
set ylabel "Liczba obsłużonych zapytań na sekundę"

set term png
set output "concurrencyTest.png"

plot 'newtest.tst' using 8:4 smooth sbezier notitle

set xlabel "Liczba klientów wysyłających zapytania"
set ylabel "Liczba obsłużonych zapytań na sekundę na klienta"

set output

set term png
set output "concurrencyTestPerThread.png"

plot 'newtest.tst' using 8:($4/$8) smooth sbezier notitle
