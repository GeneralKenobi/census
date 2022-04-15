MONGO_KEYS_DIR="/mongo-keys"
MONGO_KEY_FILE="mongodb.key"

openssl rand -base64 756 > "$MONGO_KEY_FILE"
chmod 600 "$MONGO_KEY_FILE"
chown 999:999 "$MONGO_KEY_FILE" # 999:999 is the mongodb linux user in the container

mv "$MONGO_KEY_FILE" "$MONGO_KEYS_DIR/$MONGO_KEY_FILE"

echo "Content of $MONGO_KEYS_DIR"
echo "-----"
ls -al "$MONGO_KEYS_DIR"
echo "-----"
echo "Done"
