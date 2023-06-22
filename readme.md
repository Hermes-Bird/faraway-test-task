# Faraway test task

To run POW-server
1) ``make create-netork`` to create a docker network for containers to communicate with each other
2) ``make build-docker-server`` create image of POW-server
3) ``make run-docker-server`` run docker server container

To run POW-client
1) ``make build-docker-client``
2) ``make run-docker-client [server_adderss=$]`` where server address has a `docker_network_name:port` format (default is `pow-server:8080`)

So my explanation why POW can help with DDOS-atacks in some way is that because POW requires additional time and computing power to figure out solution so users can not take access for it so easily. The POW algorithm is easy to check but hard to compute so it requires minimal computational power from our side (except generating a puzzle for client which should not be precomputed so we could not just take timeframe).
Usually what came request is some use of the server by that I mean query to a database for example and by this point of view POW makes some form of gate before that like "if u really need it u must proof it". Also, we could dynamically adjust difficulty of puzzle depending on amount of requests (what I didn't implement but as I know bitcoin does).

But we have some kinda problem with TCP itself because SYN flood attack does not require to do any computation at all and will avoid our POW implementation because most of the things with TCP are controlled by operational system itself and should be configured on that layer. So this implementation can not guarantee total security.