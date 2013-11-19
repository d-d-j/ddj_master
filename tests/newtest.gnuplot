set xlabel "Liczba wątków wysyłających równolegle zapytania"
set ylabel "Liczba obsłużonych zapytań w ciągu sekundy"

set term png
set output "concurrencyTest.png"

plot 'newtest.tst' using 8:4 smooth sbezier notitle

set xlabel "Liczba wątków wysyłających równolegle zapytania"
set ylabel "Liczba zapytań obsłużonych w ciągu sekundy przez każdy wątek"

set output

set term png
set output "concurrencyTestPerThread.png"

plot 'newtest.tst' using 8:($4/$8) smooth sbezier notitle
