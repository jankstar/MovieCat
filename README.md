# MovieCat (moviecat) Chrome App (PWA)

Chrome App (PWA): Record from camera or screen and save as a movie.

use it <br/>
[`moviecat.azurewebsites.net`](https://moviecat.azurewebsites.net)

or clone it<br/>
`git clone https://github.com/jankstar/MovieCat.git`

<br/>

# App

It is a PWA app developed in Quasar and Vue with Tailwindcss.
The application should be used in Chrome, because Chrome uses a certificate in interaction with a Mac, so that the application can access the screen, a camera or the microphone. Some other browsers do not support these functions!

## Resolution, frame rate

Before starting the recording, the resolution and the frame rate should be adjusted. These settings determine the required memory or file size of the recording.

## Compression

The media format depends on the capabilities of the browser. The file format and compression depend on these selections. The selection must be made before starting the recording.

## Recording duration / memory size

Depending on the computer, the memory for the recording is limited, i.e. the recording is kept in blocks in the memory and this memory is limited to 4 GByte on a Mac with 64 bit operating system, for example.
Depending on the resolution, frame rate and compression, the recording time is thus limited.

It would be possible to create several files, which you later join together with ffmpeg - this is not currently supported.

For decoding the WebM data, `arrayBuffer` objects are used; these are limited to 2 GByte - this function is only available for smaller files.
The upload processes slices in 2 GByte blocks (array of blobs).

## File api

If the function `showDirectoryPicker` is available, a card for a directory and file handling is displayed. A file can be read and also written to this directory when saving (download button/function).
Depending on the mime type, only files with the defined extension are displayed and can be loaded.
For WebM, the file is decoded and the metadata, e.g. duration, is determined. This function is asynchronous and depends on the file size, it needs a few seconds.

## Supported codec

The supported codec depends on the browser - Chrome in version 113, for example, only supports WebM as a container and VP8, VP9 or H.264 for video encoding. The media player, however, can play mp4. For this reason, only files with the WebM extension are selected in the file api. The button for native upload also loads any other file. For ts (mpeg), for example, only the sound is decoded and played on. So if you need a player, you should use VLC.

# Server

The application uses a server in go (gin-gonic), but node.js or any other can also be used.
The go server works locally with generated certificates that have to be created in the /key directory. For testing, I generated the certificates on my local IP address, so that the server can start at another IP address without HTTPS and thus use HTTPS when deploying to Azure via the proxy there.
If golang is deployed to the Azure Cloud, version 1.17 must be used.

Of course, node.js can also be used.

Please note: the application must either use HTTPS or `localhost`, which is a prerequisite for every PWA application.

# Build

The application is generated with the following command:

`quasar build -m pwa`

In the configuration, `minify: false ` can be set if the source code is to be analysed. In addition, `useFilenameHashes: false ` can be used to prevent new file names from being generated again and again for the generated assets and, if necessary, to force the server to reload them.

The server is generated under linux (Azure Cloud) with the command `env GOOS=linux GOARCH=amd64 go build -o server main.go `. In the Azure Cloud, the server is then started `/home/site/wwwroot/server`. Of course, node.js can also be used.

(2023/04/01)
