# CloudflareDDNS
This program is based on [Cloudflare API V4](https://api.cloudflare.com/#dns-records-for-a-zone-update-dns-record).

# Description

This program allows you to manage more than one domains by using Cloudflare ddns services.
If you have more than one ips and you need to ddns all of them but your route is not support DDNS or your can't get router access.

THIS PROGRAM IS THE BEST CHOICE!!

# How to use it?

You just need to prepare a config file and a machine. 
## Write config
You have to write your email, Global Key, Main Domain, Subdomains and Detect IP URL in config file.

+ email: Which is your Cloudflare account email address.

+ apikey: Your Cloudflare Global API Key.
(How to get? See end of README)

+ mainDomain: Your top-level domains.(like "example.com" not "www.example.com")

+ subDomainArray: Your Subdomains array.
Like ["a", "b", "c"]

  - if your maindomain is "example.com" , these subdomains will be update.

    * a.example.com

    * b.example.com

    * c.example.com

+ IPdetectAddress: A url which can return your IP address. In china you can use "http://members.3322.org/dyndns/getip" to get your IP.

    - WARNING: This url must only return ip!



## Example config

```
{
    "email" : "your@Email",
    "apiKey" : "your Global key",
    "mainDomain" : "google.com",
    "subDomainArray" : ["a", "b", "c"],
    "IPdetectAddress" : "http://members.3322.org/dyndns/getip"
}
```
## Run program
0. Download program from release page or find it in ./bin
1. Open a terminal.
2. Into the directory where the program is.
3. Make sure you know where the config file is.
4. Use command line to run the program.
```
# .\cloudflareDDNS-linux64 -c config.json
```
5. If you are first run this program and you are using linux.Please run this command first.
```
# sudo chmod +x .\cloudflareDDNS-linux64
```
Windows is as same as Linux.
## How to get Cloudflare Global API Key

>### View API Key
>
>To retrieve your API key:
>1. Log in to the Cloudflare dashboard.
>2. Under the My Profile dropdown, click My Profile.
>3. Click the API tokens tab.
>4. In the API keys section, choose one of two options: Global API Key or Origin CA Key. Choose the API Key that you would like to view.

[From Cloudflare (see end of page)](https://support.cloudflare.com/hc/en-us/articles/200167836-Where-do-I-find-my-Cloudflare-API-key-)
