const keythereum = require("keythereum");

/**
 * 
 * @param datadir directory of keystore ; absolute path
 * @param address 
 * @param password 
 * @returns 
 */
const keystoreToPrivate = (datadir, address, password) => {
    const keyObject = keythereum.importFromFile(address, datadir);
    const privateKey = keythereum.recover(password, keyObject);
    console.log(address, ": ", privateKey.toString('hex'));

    return privateKey.toString('hex');
}

module.exports = keystoreToPrivate
