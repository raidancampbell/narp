# NARP
No ARP (probe) for you

Address Resolution Protocol, detailed in [RFC 826](https://tools.ietf.org/html/rfc826), provides a simple marriage between the [link layer](https://en.wikipedia.org/wiki/Data_link_layer) and the [network layer](https://en.wikipedia.org/wiki/Network_layer).  For a given subnet, any host can broadcast an ARP request for another host on the same subnet.  All hosts on the network receive the broadcast and answers are unauthenticated.

[RFC 5227](https://tools.ietf.org/html/rfc5227) expands upon the initial ARP specification by providing a new type of ARP frame, known as an ARP Probe.  An ARP Probe is meant to prevent IP address collisions.  When a host first wishes to use an IP address on a given network, RFC 5227 compliant operating systems must send an ARP Probe for the desired IP address.  After a timeout period, the host considers the IP address unused, and is free to claim it.

NARP answers all ARP Probes, preventing any new host from claiming an IP.

![NARP](narp.gif?raw=true)