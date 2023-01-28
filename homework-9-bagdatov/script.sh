
echo "test: check existing user with autothrization and -v --raw flag for trailers"
curl --user admin:admin -v --raw http://localhost:8080/app/user/1 
echo
echo

echo "test: check existing user without autothrization and -v --raw flag for trailers"
curl -v --raw http://localhost:8080/app/user/1 
echo
echo

echo "test: register new user"
curl --user admin:admin --data "username=test1&password=test1" http://localhost:8080/app/user/2
echo
echo

echo "test: create new wallet for this user"
curl --user test1:test1 -X POST http://localhost:8080/app/wallet/testwallet
echo
echo

echo "test: check his wallets"
curl --user admin:admin http://localhost:8080/app/user/2
echo
echo

echo "test: check amount of money in the wallet"
curl --user test1:test1 http://localhost:8080/app/wallet/testwallet
echo
echo

echo "start mining"
curl --user test1:test1 -X OPTIONS http://localhost:8080/app/wallet/testwallet/start
echo
echo

sleep 3

echo "test: check amount of money in the wallet"
curl --user test1:test1 http://localhost:8080/app/wallet/testwallet
echo
echo

echo "test: stop mining"
curl --user test1:test1 -X OPTIONS http://localhost:8080/app/wallet/testwallet/stop
echo
echo

echo "test: check amount of money in the wallet"
curl --user test1:test1 http://localhost:8080/app/wallet/testwallet
echo
echo