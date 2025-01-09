ssh mark@ubuntu-latitude

works with "config"  file:
$cat config
Host ubuntu-latitude
	HostName 192.168.2.200
	User mark
	IdentityFile ~/.ssh/Id_rsa.latitude


Other way:
ssh -i ~/.ssh/Id_rsa.latitude mark@ubuntu-latitude
