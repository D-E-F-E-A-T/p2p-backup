$nodes = "Bitcoin", "BitcoinCash", "Litecoin", "Ethereum", "Ripple", "Waves"

foreach ($1 in $nodes) {
    "Stop api on $1"
    ssh $1 "systemctl stop nextgen_api"
    "Loading update to $1"
    scp ./node_linux $1":"/root/api/node_linux
    "Start api on $1"
    ssh $1 "systemctl start nextgen_api"
}

"Updating finished"
