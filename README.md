# SimSSH v1.0.0


SimSSH is a simple ssh client through which a user can execute a linux command on multiple servers simultaneously. This tool is very handy when you have to set up and configure multiple servers and don't want to manually ssh into each and everyone and write the same commands again and again. 

Key features :-
 - Interactive mode
 - Read commands from text file
 - Password / Key authentication
 - Easy to confiure and use

### Usage :-

 - Clone the repository
 ` git clone https://github.com/rohanchavan1918/simssh.git`
   
- Edit the `hosts.json` file which is present in the simssh file
```
{
   "hosts":[
      {
         "host":"host:port",
         "username":"rohan",
         "use_password":false,
         "password":"",
         "publickey":"publickey path"
      },
      {
         "host":"host:port",
         "username":"rohan",
         "use_password":true,
         "password":"supersecuredpassword",
         "publickey":""
      },
      ... Add as many hosts as you want
   ]
}
```
You can configure the tool to connect to these hosts using the password or through SSH key. If you want the tool to connect using password, the "use_password" should be true and password should be present. If the "use_password" is set to false , then you have to provide the path of the ssh key. Thats it !
### NOTE !
Since this file might contain plain text passwords, it is your responsibility to store it properly, you can store this file any where and the path can be provided to the tool, keep reading....

- Running the tool - 
		You can build from the source if you have installed go in your system or simply use the binary file. You can run this tool in interactive mode or the script can run predefined commands from a text file.

	- Binary  - The binary file is located in simssh/bin/ folder. cd to the folder and run the below command.
			Using  interactive mode
		 `simssh --hosts path_to_hosts.json --mode interactive `
			 ![interactive](https://user-images.githubusercontent.com/25792843/110240882-f434ce00-7f73-11eb-95df-74cff49ce3e5.png)
	- Using batch mode
		`simssh -hosts path_to_hosts.json --mode batch --cmd_file file_to_cmds.txt` 
			![batch](https://user-images.githubusercontent.com/25792843/110241145-4fb38b80-7f75-11eb-93d7-ab318768b2ba.png)

- ### Running from source
	
	`simssh`  has external dependencies, and so they need to be pulled in first:

	`go get && go build`

	This will create a  `simssh`  binary for you. If you want to install it in the  `$GOPATH/bin`  folder you can run:
	`go install`

- ### Refrences - 
	https://medium.com/tarkalabs/ssh-recipes-in-go-part-one-5f5a44417282
	https://gist.githubusercontent.com/boyzhujian/73b5ecd37efd6f8dd38f56e7588f1b58/raw/79dc8598ab51e13986252f68aebcefc0e97c19ee/gistfile1.txt

### TODO 
- Implement SSH
- Sometimes the output from the server is not printed completely.
- configurable verbosity through flags.

### Thanks and feel free to use and contribute !!
