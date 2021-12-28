# Benchmark for a 2 out of 2 setup

From the root folder, run


    docker-compose -f build/docker-compose.yml up

In another terminal, do the keygen

    sh benchmark/keygen.sh

Copy the public returned in the response

Run keysign.sh

    sh benchmark/keysign.sh <PUBKEY_FROM_KEYGEN> <MESSAGE_TO_SIGN>

Example

    sh keysign.sh thorpub1addwnpepqd5jn2jnvhrp7vswp53h72y3adtghvyfhgk8ycvvz25j94rym5je65w2mxc MHhlNjJlYjA2ODcwNjJjZWQ3NjNlOTRjNWI0NjkzMGQ1NzFkNDZlYzJlNjY5MWQ1MWIzODI2YTE5ZmIzZGY1OTY4