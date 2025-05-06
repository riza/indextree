<h1>indextree</h1>  
<p>Generates the tree of the directory listing page.</p>  
<p>  
  <a href="https://opensource.org/licenses/MIT">  
    <img src="https://img.shields.io/badge/license-MIT-_red.svg">  
  <a href="https://goreportcard.com/badge/github.com/riza/indextree">  
    <img src="https://goreportcard.com/badge/github.com/riza/indextree">  
  </a>  
  <a href="https://github.com/riza/indextree/releases">  
    <img src="https://img.shields.io/github/release/riza/indextree">  
  </a>  
  <a href="https://twitter.com/rizasabuncu">  
    <img src="https://img.shields.io/twitter/follow/rizasabuncu.svg?logo=twitter">  
  </a>
</p>

# What is indextree actually?

*indextree* is a tool for analysing and filtering the "directory listing" pages of web servers, as shown in the screenshot below.

<img src="https://raw.githubusercontent.com/riza/indextree/main/res/dirlist.png">

# Installation
indextree requires **go1.22** to install successfully. Run the following command to get the repo -

```sh
go install github.com/riza/indextree@latest
```

# Usage

```sh
Usage of indextree:
  -b    show banner (default true)
  -d    directory mode: only show matching directories
  -e string
        extensions to filter, example: -e jpg,png,gif
  -f    file mode: only show matching files
  -m string
        match in url, example: -mu admin,login
  -of   show only files
  -t    show tree (default true)
  -u string
        url to parse index
```

```sh
➤ indextree -u http://127.0.0.1/ -e txt,xlsx                                                                                                  git:main
    _           __          __               
   (_)___  ____/ /__  _  __/ /_________  ___ 
  / / __ \/ __  / _ \| |/_/ __/ ___/ _ \/ _ \
 / / / / / /_/ /  __/>  </ /_/ /  /  __/  __/
/_/_/ /_/\__,_/\___/_/|_|\__/_/   \___/\___/ v1.0.5
 
├── http://127.0.0.1/HOME/
├── http://127.0.0.1/secrets/
│   └── http://127.0.0.1/secrets/passwords.xlsx
│   └── http://127.0.0.1/secrets/private_key.txt
```

# Directory Mode & File Mode

You can use directory mode and file mode to focus on specific content:

```sh
# Only show directories that match "admin"
➤ indextree -u http://127.0.0.1/ -m admin -d

# Only show files that match "password"
➤ indextree -u http://127.0.0.1/ -m password -f
```

These modes help you to focus on either directories or files that match your search criteria. The tool also includes protection against infinite loops when scanning recursive directory structures.

## Disclaimer

This tool is developed and shared solely for educational and research purposes. The intention behind its creation is to foster learning and exploration within the field of cybersecurity. The tool is not intended for any malicious or illegal activities.

By accessing and using this tool, you agree to use it responsibly and in compliance with all applicable laws and regulations. The developers of this tool shall not be held liable for any misuse or damage caused by its usage.

Please use this tool ethically and responsibly, and only on systems and networks that you have permission to test. 


## Donate

<a href="https://www.buymeacoffee.com/rizasabuncu" target="_blank"><img src="https://www.buymeacoffee.com/assets/img/custom_images/orange_img.png" alt="Buy Me A Coffee" style="height: 41px !important;width: 174px !important;box-shadow: 0px 3px 2px 0px rgba(190, 190, 190, 0.5) !important;-webkit-box-shadow: 0px 3px 2px 0px rgba(190, 190, 190, 0.5) !important;" ></a>
