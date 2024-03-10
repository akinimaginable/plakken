if getent passwd plakken > /dev/null; then
  userdel -r plakken
fi

if getent group plakken > /dev/null; then
  groupdel plakken
fi