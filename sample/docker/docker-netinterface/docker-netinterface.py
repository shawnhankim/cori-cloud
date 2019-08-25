import netifaces as ni
import os

print("\nNetwork Interface Version")
print("-------------------------------------------------")
print(ni.version)

print("\nNetwork Interface Properties")
print("-------------------------------------------------")
print(dir(ni))

print("\nNetwork Interfaces")
print("-------------------------------------------------")
print(ni.interfaces())

print("\nIP")
print("-------------------------------------------------")
print(ni.ifaddresses('en0'))
print("")
print(ni.ifaddresses('utun1'))

print("\nIP Address Family")
print("-------------------------------------------------")
print(ni.address_families)

print("\nIP Address for utun1")
print("-------------------------------------------------")
ip_info = ni.ifaddresses('utun1')
ip_addr = ip_info[ni.AF_INET][0]['addr']
print(ip_addr)

print("\nIP Address for en0")
print("-------------------------------------------------")
ip_info = ni.ifaddresses('en0')
ip_addr = ip_info[ni.AF_INET][0]['addr']
print(ip_addr)

print("\nGateway")
print("-------------------------------------------------")
print(ni.gateways())

print("\nGateway based IP")
print("-------------------------------------------------")
iface = ni.gateways()["default"][ni.AF_INET][1]
print(iface)
ip_addr = ni.ifaddresses(iface)[ni.AF_INET][0]["addr"]
print(ip_addr)

