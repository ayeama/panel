#!/usr/bin/env ash

if [ ! -f eula.txt ]; then
  echo "eula=true" > eula.txt
fi

exec java @user_jvm_args.txt @libraries/net/neoforged/neoforge/21.1.249/unix_args.txt nogui
