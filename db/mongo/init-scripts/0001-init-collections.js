db.createUser({user: "admin", pwd: "admin", roles: ["readWrite"]});

db.createCollection("person");
db.person.createIndex({email: 1}, {name: "person_email", unique: true});
