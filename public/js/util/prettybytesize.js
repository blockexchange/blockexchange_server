
// https://stackoverflow.com/questions/10420352/converting-file-size-in-bytes-to-human-readable-string/10420404
export default function(fileSizeInBytes) {
    var i = -1;
    var byteUnits = [' kB', ' MB', ' GB', ' TB', 'PB', 'EB', 'ZB', 'YB'];
    do {
        fileSizeInBytes = fileSizeInBytes / 1000;
        i++;
    } while (fileSizeInBytes > 1000);

    return Math.max(fileSizeInBytes, 0.1).toFixed(1) + byteUnits[i];
}
