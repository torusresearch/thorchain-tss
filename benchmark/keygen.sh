#!/bin/sh

keygenRequest=$(cat <<EOF
              {
                  "keys":[
                      "thorpub1addwnpepqtdklw8tf3anjz7nn5fly3uvq2e67w2apn560s4smmrt9e3x52nt2svmmu3",
                      "thorpub1addwnpepqtspqyy6gk22u37ztra4hq3hdakc0w0k60sfy849mlml2vrpfr0wvm6uz09",
                      "thorpub1addwnpepq2ryyje5zr09lq7gqptjwnxqsy2vcdngvwd6z7yt5yjcnyj8c8cn559xe69",
                      "thorpub1addwnpepqfjcw5l4ay5t00c32mmlky7qrppepxzdlkcwfs2fd5u73qrwna0vzag3y4j"
                  ],
                  "tss_version":"0.14.0",
                  "block_height":100
              }
              EOF
              )

do_keygen() {
  response=$(curl -d "$keygenRequest" -H 'Content-Type: application/json' $1 --silent)
  echo "${response}"
  return 1
}

start=`date +%s`
{
  res1=$(do_keygen http://localhost:8080/keygen)
  echo "Response from Node 1 ${res1}"
}&
{
  res2=$(do_keygen http://localhost:8081/keygen)
  echo "Response from Node 2 ${res2}"
}&
{
  res3=$(do_keygen http://localhost:8082/keygen)
  echo "Response from Node 3 ${res3}"
}&
{
  res4=$(do_keygen http://localhost:8083/keygen)
  echo "Response from Node 4 ${res4}"
}&
wait
end=`date +%s`
runtime=$( echo "$end - $start" | bc -l )

echo "Keygen Complete. Took ${runtime} seconds"