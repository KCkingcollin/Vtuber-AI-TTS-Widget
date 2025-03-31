# Vtuber-AI-TTS-Widget
VATTS (Vtuber-AI-TTS-Widget)

If you're on Linux download the Linux zip from the [releases](https://github.com/KCkingcollin/Vtuber-AI-TTS-Widget/releases), and simply run VATTS\
If you're on Windows download the Windows zip from the [releases](https://github.com/KCkingcollin/Vtuber-AI-TTS-Widget/releases), and simply run VATTS.exe

When you first start it, it will download all that it needs to run, do not close the windows or shutdown during the process otherwise it may not finish installing everything it needs.\
The first thing you type will take a while for it to process, but it should be faster afterwords.

alt+shift+t will be the default keys to hide the GUI, pressing it again will bring it back, if there is a problem with it popping back up, you can run VATTS again, and it will tell you if its running already and give you the option to kill it.\
Experimental configuration is available via the settings.json file in the config folder, you will need to restart the app for it to take affect.
If you need to revert to default settings simply delete the settings.json file and restart the app.

If you'd like to build the app then you'll need python 3.13 and GO, you can use the requirements.txt to set up a venv for python with all you'll need to build it, then just use pyinstaller and the provided spec file to make the server exe, make sue to create the python venv folder in the project directory as there is no guarantee that i wont add something from there to the spec file, you may also need the voice(s).bin/onnx/json files for the spec file to work, here is where to download them:
https://github.com/thewh1teagle/kokoro-onnx/releases/download/model-files/voices.json
https://github.com/thewh1teagle/kokoro-onnx/releases/download/model-files/voices.bin
https://github.com/thewh1teagle/kokoro-onnx/releases/download/model-files/kokoro-v0_19.onnx

After building the python exe you'll need to run ```go mod tidy``` to get all the modules that are required to build it. Then if you're on linux simply run ```go build -v -o VATTS main.go```, however if you are on windows then you are SOL, you will need to use wsl2 or set up a Linux environment so you can compile it there, I do not know why, but it seems to be far more difficult to build the GO windows exe in windows as crazy as that sounds, anyway,  once you're in a Linux environment you'll need to run ```GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc go build -v -ldflags="-H windowsgui" -o VATTS.exe```, you will obviously need go in the Linux environment, but not as obvious you'll also need mingw-gcc

Finally, just put all the exes and necessary files in the right folders, the python exe and its `_internal` folder will be in the `kokoro/dist` folder, so move both them into the project root or where ever `VATTS` will be. You don't reall need the kokoro folder in the final project, VATTS will download what it needs as long as it at least has the server exe, the internal folder, and itself.

