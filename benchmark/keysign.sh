#!/bin/sh

#How to run
# sh keysign.sh <<pub_key from keygen>> <<base64 message to sign>>
# example:
# sh keysign.sh thorpub1addwnpepqd5jn2jnvhrp7vswp53h72y3adtghvyfhgk8ycvvz25j94rym5je65w2mxc MHhlNjJlYjA2ODcwNjJjZWQ3NjNlOTRjNWI0NjkzMGQ1NzFkNDZlYzJlNjY5MWQ1MWIzODI2YTE5ZmIzZGY1OTY4

keysignRequest=$(cat <<EOF
              {
                  "pool_pub_key": "$1",
                  "messages": ["$2"],
                  "signer_pub_keys": [
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

echo $keysignRequest
do_keysign() {
  response=$(curl -d "$keysignRequest" -H 'Content-Type: application/json' $1 --silent)
  echo "${response}"
  return 1
}

start=`gdate +%s.%N`
{
  res1=$(do_keysign http://localhost:8080/keysign)
  echo "Response from Node 1 ${res1}"
}&
{
  res2=$(do_keysign http://localhost:8081/keysign)
  echo "Response from Node 2 ${res2}"
}&
{
  res3=$(do_keysign http://localhost:8082/keysign)
  echo "Response from Node 3 ${res3}"
}&
{
  res4=$(do_keysign http://localhost:8083/keysign)
  echo "Response from Node 4 ${res4}"
}&
wait
end=`gdate +%s.%N`
runtime=$( echo "$end - $start" | bc -l )

echo "Keysign Complete. Took ${runtime} seconds"