#!/bin/sh

#How to run
# sh keysign.sh
# copy the pubkey returned and use it in keysign.sh

keygenRequest=$(cat <<EOF
              {
                  "keys":[
                      "thorpub1addwnpepqtdklw8tf3anjz7nn5fly3uvq2e67w2apn560s4smmrt9e3x52nt2svmmu3",
                      "thorpub1addwnpepqtspqyy6gk22u37ztra4hq3hdakc0w0k60sfy849mlml2vrpfr0wvm6uz09"
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
wait
end=`date +%s`
runtime=$( echo "$end - $start" | bc -l )

echo "Keygen Complete. Took ${runtime} seconds"