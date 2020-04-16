
import { get_info } from '../api/info.js';

var promise;

export default function(){
  if (!promise){
    promise = new Promise((resolve, reject) => {
      get_info()
      .then(resolve)
      .catch(reject);
    });
  }
  return promise;
}
