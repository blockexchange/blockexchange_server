
import { get_colormapping } from '../api/colormapping.js';

var promise;

export default function(){
  if (!promise){
    promise = new Promise((resolve, reject) => {
      get_colormapping()
      .then(resolve)
      .catch(reject);
    });
  }
  return promise;
}
