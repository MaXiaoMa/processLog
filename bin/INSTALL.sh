echo "compile start ...."
cd ../
go build -v -o ./bin/logAnalysis
cd ./bin
echo "compile completed!"
echo "start create project folder..."
mkdir logAnalysis
cd logAnalysis
mkdir bin
mkdir conf
mkdir data
mkdir log
mv ../logAnalysis ../start.sh ../stop.sh bin
mv ../../conf/config.json conf
cd bin
chmod 755 logAnalysis start.sh stop.sh
cd ../
cd ../
echo "create folder completed!"
