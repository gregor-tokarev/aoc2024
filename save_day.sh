read -p "Day number:" day_number

mkdir "day${day_number}"

if [ -f ./main.go ]; then
mv main.go "day${day_number}/main.go"
fi

if [ -f ./data.txt ]; then
mv data.txt "day${day_number}/day.txt"
fi
