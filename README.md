#go_crud_app


Be sure to adjust the systemvariables DBUSER and DBPASS to the proper values that you are using for your database. 

Please run "CREATE DATABASE swosh;" and "use swosh;" before running the code in this repo.

The development of this app is documented on www.twitch.tv/viktorvallmark


Installing Golang

Linux:

Run this command as sudo:  rm -rf /usr/local/go && tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz
Do not untar it into any existing /usr/local/go tree

Export Go bin to path:

export PATH=$PATH:/usr/local/go/bin

Windows:

Enter this url in your browser of choice: https/go.dev/dl/go1.21.5.windows-amd64.msi

Run the MSI.

Verify the install by doing go version in the command prompt.
