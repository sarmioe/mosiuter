# MOSIUTER, see the world, free and private

## Introduction

> MOSIUTER is a modern encryption protocol.
>
> This protocol uses multi-mode obfuscation + encryption to transmit your data, ensuring that even if others get your data, they don't know your true intentions, and disguise it like a normal Internet request.
>
> Both the transmission protocol and encryption features are specified by an IV, and this IV changes every few minutes to ensure that the data packet cannot be judged by the DPI rules.
>
> First, your data will be optimized, and then encryption such as snowflake or obfs4 will be used to remove the features, and then encrypted with TLS, so that the data packet looks like normal Internet encrypted data.

## Quick Start

>Download mosiuterclient and mosiuterserver
>
>Put mosiuterserver on the server and start it
>
>At this time, it will generate two key files. You only need to open the file ending with .pem. This key file is used to encrypt the key for transmitting the first few important data packets
>
>Then you need to configure two servers according to the configuration in the first two steps (at least two, used to distract DPI attention, and up to 10. Of course, the more the better)
>
>Then edit the ip.txt file in the local mosiuterclient directory, delete the file content, and then write the server's IP address and port number
>
>Then start mosiuterclient. It will automatically connect to the server randomly and start data transmission. You need to wait for a few seconds for the first few data packets to be transmitted
>
>Then you need to wait for the client's prompt. It will output your IP: port. Finally, set the proxy server to its output. Then directly visit the website with a browser to start encrypting your Internet traffic

## How to compile

> mosiuter is divided into client (mosiuterclient) and server (mosiuterserver), you can find these two folders in the root directory
>
> Step 0 You need to CD to the root folder and Go language environment version 1.23.4
>
> ### First, let's talk about compiling with Go build
>
> Client: `cd ./mosiuterclient` (Enter first) then `go build`
>
> Server: `cd ./mosiuterserver` (Enter first) then `go build`
>
> ### Then talk about compiling with MakeFile
>
> Client: make client
>
> Server: make server
>
> Clean up the compiled binary files: make debuild
>
> Run the build-all script in the directory

## How to it work?

### Transport protocol

> The code indicates which protocol is used for transmission based on IV1

| Code | Protocol             |
| ---- | -------------------- |
| 1    | TCP RAW              |
| 2    | UDP defaults to QUIC |
| 3    | TLS                  |
| 4    | HTTPS                |
| 5    | DoT                  |
| 6    | DoH                  |
| 7    | mKCP                 |
| 8    | gRCP                 |

### Camouflage protocol

> What is the code IV2? It is the camouflage feature
>
> And it will add a certain random value. It will be agreed on where to place it before the data packet starts to be transmitted (multiple random values may be inserted in multiple places in a data packet)
>
> What makes the censors even more angry is that every time there is a random value (an integer with a minimum of 1 and a maximum of 20), this code will randomly generate a new one, so that the proxy server is also disguised as a reverse proxy, and the client pretends to use this reverse proxy to surf the Internet normally
>
> There is also TTL value camouflage, which makes the proxy server look more like a reverse proxy server
>
> There is also automatic rotation of the proxy server IP. The proxy server IP is changed every few minutes, which makes DPI angry. It also supports IPv6, which makes DPI even more angry
>
> What makes DPI even more angry is that this thing will also send several real and fake data packets. It will really create a request to access the website, but it will be deleted after reaching the proxy server
>
> If a certain IP replays the data packet to the proxy server, it will really return a legitimate page (or a legitimate website DNS query result) to this IP, and in order to avoid being exposed, this IP will also be recorded Each replay will return the same data

| Code | Disguise features                | Disguise behavior                                            |
| ---- | -------------------------------- | ------------------------------------------------------------ |
| 9    | Watch video websites             | Simulate CDN behavior Divide traffic into multiple small segments Randomly distribute on different IPs |
| 10   | Listen to music websites         | Simulate buffering and pausing Intermittently transmit traffic instead of continuous transmission |
| 11   | Download large files             | Simulate breakpoint resumption Re-request different parts of the file after a period of time |
| 12   | Log in to the cloudflare website | Small data packets Slightly longer interval Simulate logging in to a website |
| 13   | Play online games                | High-frequency small data packets Simulate UDP traffic Combined with Ping value randomization (but there is also a limit Maximum 500ms Minimum 100ms and slightly smaller fluctuations) |
| 14   | Video call                       | Simulate video conferencing protocols such as WebRTC or Zoom |
| 15   | Random data                      | Pre-disguise with OBFS4+Shandowsocks Then insert some useless data |

### Packet acceleration

> MOSIUTER uses a variety of advanced packet size reduction modes to ensure that the packet size is smaller and more data is transmitted
>
> 1. Preamble deletion
>
> 2. Frame gap deletion
>
> 3. Try frame fragmentation as much as possible
>
> 4. If the protocol is based on TCP, BBR will be enabled
>
> 5. If the protocol is based on UDP, UOT will be enabled
>
> 6. Enable header compression
>
> 7. Enable server-side (to proxy server) Gzip compression and Brotli compression
>
> 8. If https is randomly reached, HTTP/3 will be enabled
>
> 9. Merge similar data
>
> 10. IPv4/v6 automatically selects the best
>
> 11. Enable video or image compression
>
> And more

## Simplified actual Internet operation process

> First, the client needs to get the server key and a server trust list
>
> The trust list can have a maximum of 10 entries and a minimum of 2 entries at a time. It can be a single IPv4\IPv6 address. It is recommended to build an IPv4 server first. After all, compared with IPV6, this thing is more stable in transmission

> 1. Communication proxy server Establish a WS connection with it. If there is a middleman tampering, immediately cut off the network (kill swithy)
>
> 2. Then start the second data packet exchange: notify the server of the random data IV1
>
> 3. After the server gets this value, it sends a received data packet to the client
>
> 4. Start to negotiate the random data position (IV3). Add a few random packets for the number. The value is an integer with a maximum of 35 and a minimum of 15
>
> 5. The client starts random data and uses the configured encryption key as the transmission encryption of the random data position notification data packet
>
> 6. After the server confirms that everything is fine, the client starts random IV2
>
> 7. After IV2 is randomized After encrypting the data to the server, notify the server of IV2.
>
> 8. After the client receives the data packet, the WS connection is cut off. Server 1 notifies the encryption method used to the next server (encryption notification).
>
> 9. Get browser data and start data packet optimization.
>
> 10. (When the time is reached) The client ends the connection with this server and opens a new connection to start connecting to this new address.
>
> 11. When the random time expires, the connection with this server will be cut off and the next server will be opened. This cycle continues.
