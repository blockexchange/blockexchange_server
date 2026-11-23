var Module = typeof Module != "undefined" ? Module : {};

Module['print'] = (text) => {
    console.log(text)
};

Module['printErr'] = (text) => {
    console.log(text)
};

importScripts('luanti.js');
