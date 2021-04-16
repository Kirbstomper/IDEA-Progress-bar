# IDEA-Progress-bar
Backend for generating a custom progress bar when sent an image

This project is to act as the backend for my website used to generate custom intelij progress bars.
You can send an icon and it will be used in the plugin file generated and returned

Big thanks to https://github.com/batya239/NyanProgressBar for providing most of the base for the plugin.
Also thanks to https://github.com/disintegration/imaging for providing an easy to use go library for manipuilating images



# Usage when running in docker

Sending a post request to localhost:8080/upload with a multipart form containing the following pairs in the form (substituting your color values as you wish)


image : your image file you wish to turn into a loading icon

config : {"R":255, "B":0, "G":255 }
