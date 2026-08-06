#!/usr/bin/env ash

sed -i 's/^eula=false$/eula=true/' eula.txt
exec java -jar server.jar --nogui
