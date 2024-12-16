read -p "Day number:" day_number

mkdir "day${day_number}"
mv main.go "day${day_number}/main.go"
mv data.txt "day${day_number}/day.txt"
