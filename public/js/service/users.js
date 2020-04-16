import { get_all } from '../api/user.js';

var promise;

export default function(){
  if (!promise){
    promise = new Promise((resolve, reject) => {
      get_all()
      .then(resolve)
      .catch(reject);
    });
  }

  return promise;
}
