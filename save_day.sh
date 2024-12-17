read -p "Day number:" day_number

if [ -d "day${day_number}" ]; then
mkdir "day${day_number}"
fi

if [ -f ./main.go ]; then
mv main.go "day${day_number}/main.go"
fi

if [ -f ./data.txt ]; then
mv data.txt "day${day_number}/day.txt"
fi
