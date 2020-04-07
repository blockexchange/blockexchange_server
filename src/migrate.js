
const pgmigrations = require("postgres-migrations");
const pool = require("./pool");

module.exports = () => {
  return new Promise(function(resolve, reject){
    pool.connect()
    .then(client => {
      return pgmigrations.migrate({client: client }, "migrations")
      .then(() => client.release())
      .then(resolve)
      .catch(e => {
        client.release();
        console.error(e.stack);
        reject();
      });
    });
  });
};
