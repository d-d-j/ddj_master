set xlabel "Liczba wątków wysyłających INSERT"
set ylabel "Czas wykonania 1000 zapytań SELECT na 1000000 rekordów"

set term png
set output "concurrentSelectsTest.png"

plot 'serverTests/1000select100000values.txt' using 3:1 smooth sbezier notitle

