if ! getent group plakken > /dev/null; then
  groupadd -r plakken
fi

if ! getent passwd plakken > /dev/null; then
    useradd -r -d /var/lib/plakken -s /sbin/nologin -g plakken -c "Plakken server" plakken
fi
if ! test -d /var/lib/plakken; then
    mkdir -p /var/lib/plakken
    chmod 0750 /var/lib/plakken
    chown -R plakken:plakken /var/lib/plakken
fi
