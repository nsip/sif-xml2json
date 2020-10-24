# Web-Service for Converting SIF XML Format to JSON Format

## Build Prerequisite

0. Except 'config.toml' content, Do NOT change any project structure & file name & content.

1. Make sure current working directory of your command line environment is identical to the directory of this README.md file.
   i.e. under "/sif-xml2json/"

## Native Build

0. It is NOT supported to make building on Windows OS. If you are using Windows OS, please choose 'Docker Build'.

1. Make sure `golang` dev package & `git` are available on your machine.

2. Run `./build.sh` to build service which embedded with SIF Spec 3.4.6 & 3.4.7.

3. Run `./release.sh [linux64|win64|mac] 'dest-path'` to extract minimal executable package on different.
   e.g. `./release.sh win64 ~/Desktop/sif-xml2json/` extracts windows version bin package into "~/Desktop/sif-xml2json/".
   then 'server' executable is available under "~/Desktop/sif-xml2json/".

4. Jump into "~/Desktop/sif-xml2json/", modify 'config.toml' if needed.
   Please set [Service] & [Version] to your own value.

5. Run `server`.
   Default port is `1324`, can be set at config.toml.

## Docker Build
  
0. Make sure `Docker` is available and running on your machine.

1. Run `docker build --rm -t nsip/sif-xml2json:latest .` to make docker image.

2. In order to do configuration before running docker image.
   Copy '/sif-xml2json/config/config.toml' to current directory, modify if needed, and name it like `config_d.toml`.
   Please set [Service] & [Version] to your own value.

3. Run `docker run --rm --mount type=bind,source=$(pwd)/config_d.toml,target=/config.toml -p 0.0.0.0:1324:1324 nsip/sif-xml2json`.
   Default port is `1324`, can be set at config.toml. If not 1324, change above command's '1324' to your own number.

## Test

0. Make sure `curl` is available on your machine.

1. Run `curl IP:Port` to get the list of all available API path of sif-xml2json.
   `IP` : your sif-xml2json server running machine ip.
   `Port`: set in 'config.toml' file, default is 1324, can be changed in 'config.toml'.

2. Run `curl -X POST IP:Port/Service/Version/2json?sv=3.4.7 -d @path/to/your/sif.xml`
   to convert a XML SIF to JSON.

   `IP` : your sif-xml2json server running machine ip.
   `Port`: Get from server's 'config.toml'-[WebService]-[Port], default is 1324.
   `Service`: service name. Get from server's 'config.toml'-[Service].
   `Version`: service version. Get from server's 'config.toml'-[Version].
   `sv`: SIF Spec Version, available 3.4.6 & 3.4.7
