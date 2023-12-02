package docker_compose_yaml

var orderTemplate = `
version: '3'
services:
  ${containerOrderName}:
    container_name: ${containerOrderName}
    image: ${image}
    restart: always
    labels:
      service: hyperledger-fabric
    environment:
      - FABRIC_LOGGING_SPEC=INFO
      - ORDERER_GENERAL_LISTENADDRESS=0.0.0.0
      - ORDERER_GENERAL_LISTENPORT=7050
      - ORDERER_GENERAL_LOCALMSPID=OrdererMSP
      - ORDERER_GENERAL_LOCALMSPDIR=/etc/hyperledger/fabric/msp
      - ORDERER_GENERAL_TLS_ENABLED=true
      - ORDERER_GENERAL_TLS_PRIVATEKEY=/etc/hyperledger/fabric/tls/server.key
      - ORDERER_GENERAL_TLS_CERTIFICATE=/etc/hyperledger/fabric/tls/server.crt
      - ORDERER_GENERAL_TLS_ROOTCAS=[/etc/hyperledger/fabric/tls/ca.crt]
      - ORDERER_GENERAL_CLUSTER_CLIENTCERTIFICATE=/etc/hyperledger/fabric/tls/server.crt
      - ORDERER_GENERAL_CLUSTER_CLIENTPRIVATEKEY=/etc/hyperledger/fabric/tls/server.key
      - ORDERER_GENERAL_CLUSTER_ROOTCAS=[/etc/hyperledger/fabric/tls/ca.crt]
      - ORDERER_GENERAL_BOOTSTRAPMETHOD=none
      - ORDERER_CHANNELPARTICIPATION_ENABLED=true
      - ORDERER_ADMIN_TLS_ENABLED=true
      - ORDERER_ADMIN_TLS_CERTIFICATE=/etc/hyperledger/fabric/tls/server.crt
      - ORDERER_ADMIN_TLS_PRIVATEKEY=/etc/hyperledger/fabric/tls/server.key
      - ORDERER_ADMIN_TLS_ROOTCAS=[/etc/hyperledger/fabric/tls/ca.crt]
      - ORDERER_ADMIN_TLS_CLIENTROOTCAS=[/etc/hyperledger/fabric/tls/ca.crt]
      - ORDERER_ADMIN_LISTENADDRESS=0.0.0.0:7053
    working_dir: /etc/hyperledger/fabric
    command: orderer
    volumes:
       - ${certs}/certs/ordererOrganizations/${orgName}.${dns}/orderers/${containerOrderName}/msp:/etc/hyperledger/fabric/msp
       - ${certs}/certs/ordererOrganizations/${orgName}.${dns}/orderers/${containerOrderName}/tls:/etc/hyperledger/fabric/tls
       - ${certs}/production:/var/hyperledger/production/
    ports:
      - ${7050}:7050
      - ${7053}:7053
    extra_hosts:
      $(cat $PWD/order_ip.txt)`

// peer模板
var peerTemplate = `
version: '3'
services:
  ${containerPeerName}:
    container_name: ${containerPeerName}
    image: ${image}
    restart: always
    labels:
      service: hyperledger-fabric
    environment:
      - CORE_VM_ENDPOINT=unix:///host/var/run/docker.sock
      - CORE_VM_DOCKER_HOSTCONFIG_NETWORKMODE=peer_default
      - FABRIC_LOGGING_SPEC=INFO
      - CORE_PEER_TLS_ENABLED=true
      - CORE_PEER_PROFILE_ENABLED=true
      - CORE_PEER_ADDRESSAUTODETECT=true
      - CORE_PEER_TLS_CERT_FILE=/etc/hyperledger/fabric/tls/server.crt
      - CORE_PEER_TLS_KEY_FILE=/etc/hyperledger/fabric/tls/server.key
      - CORE_PEER_TLS_ROOTCERT_FILE=/etc/hyperledger/fabric/tls/ca.crt
      - CORE_PEER_ID=${containerPeerName}
      - CORE_PEER_ADDRESS=${containerPeerName}:7051
      - CORE_PEER_LISTENADDRESS=0.0.0.0:7051
      - CORE_PEER_CHAINCODEADDRESS=${containerPeerName}:7052
      - CORE_PEER_CHAINCODELISTENADDRESS=0.0.0.0:7052
      - CORE_PEER_GOSSIP_BOOTSTRAP=$(cat $PWD/${orgName}-peer_bootstrap.txt|paste -s -d ' ')
      - CORE_PEER_GOSSIP_EXTERNALENDPOINT=${containerPeerName}:7051
      - CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/fabric/msp
      - CORE_PEER_LOCALMSPID=${orgName}MSP
    volumes:
        - /var/run/docker.sock:/host/var/run/docker.sock
        - ${certs}/certs/peerOrganizations/${orgName}.${dns}/peers/${containerPeerName}/msp:/etc/hyperledger/fabric/msp
        - ${certs}/certs/peerOrganizations/${orgName}.${dns}/peers/${containerPeerName}/tls:/etc/hyperledger/fabric/tls
        - ${certs}/production:/var/hyperledger/production
    working_dir: /etc/hyperledger/fabric/peer
    command: peer node start
    ports:
      - ${7051}:7051
      - ${7052}:7052
    extra_hosts:
      $(cat $PWD/order_ip.txt)`

// cli模板
var cliTemplate = `
version: '3'
services:
 ${containerCliName}:
    container_name: ${containerCliName}
    image: ${image}
    labels:
      service: hyperledger-fabric
    restart: always
    tty: true
    stdin_open: true
    environment:
      - GOPATH=/etc/hyperledger/fabric
      - CORE_VM_ENDPOINT=unix:///host/var/run/docker.sock
      - FABRIC_LOGGING_SPEC=INFO
      - CORE_PEER_TLS_ENABLED=true
      - CORE_PEER_TLS_CERT_FILE=/etc/hyperledger/fabric/tls/server.crt
      - CORE_PEER_TLS_KEY_FILE=/etc/hyperledger/fabric/tls/server.key
      - CORE_PEER_TLS_ROOTCERT_FILE=/etc/hyperledger/fabric/tls/ca.crt
      - CORE_PEER_ADDRESS=${peerAddress}   
      - CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/fabric/msp
      - CORE_PEER_LOCALMSPID=${orgName}MSP
    working_dir: /etc/hyperledger/fabric
    command: /bin/bash
    volumes:
      - /var/run/:/host/var/run/
      - ${certs}/certs:/certs
      - ${certs}/certs/peerOrganizations/${orgName}.${dns}/users/Admin@${orgName}.${dns}/msp:/etc/hyperledger/fabric/msp
      - ${certs}/certs/peerOrganizations/${orgName}.${dns}/peers/${containerCliName}/tls:/etc/hyperledger/fabric/tls
      - ${certs}/chaincode:/etc/hyperledger/fabric/src
      - ${certs}/blocktx:/etc/hyperledger/fabric/blocktx
      - ${certs}/configtx:/etc/hyperledger/fabric/configtx
    extra_hosts:
      $(cat $PWD/order_ip.txt)
      $(cat $PWD/${orgName}-peer_host.txt)`

var explorerTemplate = `
version: '2.1'
services:
  ${containerExplorerDBName}:
    image: ${exploreDBImage}
    container_name: ${containerExplorerDBName}
    hostname: ${containerExplorerDBName}
    restart: always
    environment:
      - DATABASE_DATABASE=fabricexplorer
      - DATABASE_USERNAME=hppoc
      - DATABASE_PASSWORD=password
    healthcheck:
      test: "pg_isready -h localhost -p 5432 -q -U postgres"
      interval: 30s
      timeout: 10s
      retries: 5
    volumes:
      - ${pgdata}/pgdata:/var/lib/postgresql/data
      - /etc/localtime:/etc/localtime

  ${containerExplorerName}:
    image: ${exploreImage}
    container_name: ${containerExplorerName}
    hostname: ${containerExplorerName}
    restart: always
    environment:
      - DATABASE_HOST=${containerExplorerDBName}
      - DATABASE_DATABASE=fabricexplorer
      - DATABASE_USERNAME=hppoc
      - DATABASE_PASSWD=password
      - LOG_LEVEL_APP=debug
      - LOG_LEVEL_DB=debug
      - LOG_LEVEL_CONSOLE=info
      - LOG_CONSOLE_STDOUT=true
      - DISCOVERY_AS_LOCALHOST=false
    volumes:
      - ${fabricConfig}/config.json:/opt/explorer/app/platform/fabric/config.json
      $(cat $PWD/explorer_peers_volumes.txt)
      - ${certs}/certs:/certs
      - ${walletstore}/walletstore:/opt/explorer/wallet
      - /etc/localtime:/etc/localtime
    ports:
      - ${8080}:8080
    depends_on:
      ${containerExplorerDBName}:
        condition: service_healthy
    extra_hosts:
      $(cat $PWD/explorer_peers_url_tmp.txt)`

var couchDBTemplate = `
version: '2'
services:
  ${containerCouchdbName}:
    container_name: ${containerCouchdbName}
    image: ${image} // couchdb:3.3.2
    restart: always
    environment:
      - COUCHDB_USER=admin
      - COUCHDB_PASSWORD=PL,OKM09*
    volumes:
      - ${db}/db/couchdb:/opt/couchdb/data
    ports:
      - "${5984}:5984"`
