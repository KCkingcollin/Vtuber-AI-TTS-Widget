# Vtuber-AI-TTS-Widget
VATTS (Vtuber-AI-TTS-Widget)

## Users
If you're on Linux download the Linux zip from the [releases](https://github.com/KCkingcollin/Vtuber-AI-TTS-Widget/releases), and simply run VATTS\
If you're on Windows download the Windows zip from the [releases](https://github.com/KCkingcollin/Vtuber-AI-TTS-Widget/releases), and simply run VATTS.exe

When you first start it, it will download all that it needs to run, do not close the windows or shutdown during the process otherwise it may not finish installing everything it needs.\
The first thing you type will take a while for it to process, but it should be faster afterwords.

alt+shift+t will be the default keys to hide the GUI, pressing it again will bring it back, if there is a problem with it popping back up, you can run VATTS again, and it will tell you if its running already and give you the option to kill it.\
Experimental configuration is available via the settings.json file in the config folder, you will need to restart the app for it to take affect.
If you need to revert to default settings simply delete the settings.json file and restart the app.


## Developers
If you'd like to build the app then you'll need python 3.13.2, GO, and if you're building it for windows x86_64-w64-mingw32-gcc.\
If you are on windows then you are SOL, you will need to use wsl2 or set up a Linux environment so you can compile it there, I do not know why, but it seems to be far more difficult to build the windows exes in windows as crazy as that sounds. I tried in VM and I almost jumped off a bridge, so just don't.

You may need the voice(s).bin/onnx/json files for the spec file to work, put them into the `kokoro-tts` folder. Here is where to download them:\
https://github.com/thewh1teagle/kokoro-onnx/releases/download/model-files/voices.json\
https://github.com/thewh1teagle/kokoro-onnx/releases/download/model-files/voices.bin\
https://github.com/thewh1teagle/kokoro-onnx/releases/download/model-files/kokoro-v0_19.onnx

First move to the `kokoro-tts` folder and run ```python -m venv venv``` or if you want to make the windows environment ```wine python -m venv winvenv```.\
Make sure the venv is in the `kokoro-tts` folder.

You can use the requirements.txt to install all you'll need to build it, just run ```pip install -r requirements.txt``` or for the windows build run ```wine winvenv/Scripts/pip.exe install -r requirements.txt```, then just run ```pyinstaller tts-server.spec``` or for a windows exe run ```wine ./winvenv/Scripts/pyinstaller.exe tts-server-win.spec```.

After building the python exe you'll need to run ```go mod tidy``` to get all the modules that are required to build it. Then simply run ```go build -v -o VATTS main.go```, for the windows exe run ```GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc go build -v -ldflags="-H windowsgui" -o VATTS.exe```.

Finally, just put all the exes and necessary files in the right folders, the python exe and its `_internal` folder will be in the `kokoro-tts/dist` folder, so move both them into the project root or where ever `VATTS` will be. You don't reall need the kokoro-tts folder in the final project, VATTS will download what it needs as long as it at least has the server exe, the internal folder, and itself.

