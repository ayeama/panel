#!/usr/bin/env ash

if [ ! -f server.properties ] || [ ! -f eula.txt ]; then
  java -jar server.jar --initSettings
  sed -i 's/^eula=false$/eula=true/' eula.txt
fi

exec java -jar server.jar --nogui
