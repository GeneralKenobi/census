MONGO_INIT_DIR="/docker-entrypoint-initdb.d"
INIT_SCRIPTS_DIR="init-scripts"

echo "Starting"

if [ ! -d "$MONGO_INIT_DIR" ]; then
  echo "$MONGO_INIT_DIR doesn't exist - creating"
  mkdir "$MONGO_INIT_DIR"
fi

echo "Copying init scripts to $MONGO_INIT_DIR"
cp "$INIT_SCRIPTS_DIR"/* "$MONGO_INIT_DIR/"

echo "Content of $MONGO_INIT_DIR"
echo "-----"
ls -al "$MONGO_INIT_DIR"
echo "-----"
echo "Done"
